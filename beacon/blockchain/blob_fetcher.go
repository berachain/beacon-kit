// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
)

var (
	errNoMoreRequests = errors.New("no more requests in queue")
)

// BlobFetchRequest contains the minimal data needed to fetch and validate blobs.
type BlobFetchRequest struct {
	Header        *ctypes.BeaconBlockHeader                    `json:"header"`
	Commitments   eip4844.KZGCommitments[common.ExecutionHash] `json:"commitments"`
	LastRetryTime time.Time                                    `json:"last_retry_time"`
	FailureCount  int                                          `json:"failure_count"`
}

// blobFetcher handles asynchronous fetching of blobs in the background.
type blobFetcher struct {
	logger         log.Logger
	blobProcessor  BlobProcessor
	blobRequester  BlobRequester
	storageBackend StorageBackend
	chainSpec      BlobFetcherChainSpec

	// Directory for persistent queue
	queueDir string
	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc

	stopOnce sync.Once
}

// NewBlobFetcher creates a new background blob fetcher.
func NewBlobFetcher(
	dataDir string,
	logger log.Logger,
	blobProcessor BlobProcessor,
	blobRequester BlobRequester,
	storageBackend StorageBackend,
	chainSpec BlobFetcherChainSpec,
) (BlobFetcher, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create queue directory
	queueDir := filepath.Join(dataDir, "blob_fetcher_queue")
	if err := os.MkdirAll(queueDir, 0700); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create blob fetcher queue directory: %w", err)
	}

	return &blobFetcher{
		logger:         logger,
		blobProcessor:  blobProcessor,
		blobRequester:  blobRequester,
		storageBackend: storageBackend,
		chainSpec:      chainSpec,
		queueDir:       queueDir,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start begins the background blob fetching process.
func (bf *blobFetcher) Start() {
	go bf.run()
}

// Stop gracefully shuts down the blob fetcher.
func (bf *blobFetcher) Stop() {
	bf.stopOnce.Do(func() {
		bf.cancel()
	})
}

// SetHeadSlot updates the head slot for blob fetching.
func (bf *blobFetcher) SetHeadSlot(slot math.Slot) {
	bf.blobRequester.SetHeadSlot(slot)
}

// QueueBlobRequest queues a request to fetch blobs for a specific slot.
func (bf *blobFetcher) QueueBlobRequest(slot math.Slot, block *ctypes.BeaconBlock) error {
	// Don't queue if no blobs expected
	commitments := block.GetBody().GetBlobKzgCommitments()
	if len(commitments) == 0 {
		return nil
	}

	// Create request with header and commitments needed for validation
	request := BlobFetchRequest{
		Header:      block.GetHeader(),
		Commitments: commitments,
	}

	// Serialize to JSON file with slot as filename
	filename := filepath.Join(bf.queueDir, fmt.Sprintf("%010d.json", slot.Unwrap()))
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal blob fetch request: %w", err)
	}

	// Write the request to a tmp file first, then rename atomically. This prevents the
	// main run loop from seeing a partially written file causing JSON unmarshal errors.
	tempFile := filename + ".tmp"
	if writeErr := os.WriteFile(tempFile, data, 0600); writeErr != nil {
		return fmt.Errorf("failed to write temp blob fetch request: %w", writeErr)
	}
	if renameErr := os.Rename(tempFile, filename); renameErr != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to rename blob fetch request: %w", renameErr)
	}

	bf.logger.Info("Queued blob fetch request", "slot", slot.Unwrap(), "expected_blobs", len(commitments))
	return nil
}

func (bf *blobFetcher) run() {
	// Ticker to periodically check for requests (both new and ready to retry)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-bf.ctx.Done():
			bf.logger.Info("Blob fetcher shutting down")
			return

		case <-ticker.C:
			// Process all pending requests from disk
			bf.processAllPendingRequests()
		}
	}
}

// processAllPendingRequests reads and processes all queued requests from disk.
func (bf *blobFetcher) processAllPendingRequests() {
	for {
		select {
		case <-bf.ctx.Done():
			return
		default:
		}

		request, filename, err := bf.getNextRequest()
		if err != nil {
			if errors.Is(err, errNoMoreRequests) {
				return
			}
			bf.logger.Error("Failed to get next request", "error", err)
			if filename != "" {
				bf.removeRequestFile(filename)
			}
			continue
		}

		err = bf.processFetchRequest(request)
		if err == nil {
			// Successfully processed, remove the request file
			bf.removeRequestFile(filename)
			continue
		}

		bf.logger.Error("Failed to process blob fetch request", "slot", request.Header.Slot.Unwrap(), "error", err)

		// Update retry metadata and save back to file
		request.FailureCount++
		request.LastRetryTime = time.Now()
		var data []byte
		data, err = json.Marshal(request)
		if err != nil {
			bf.logger.Error("Failed to marshal request", "error", err)
			continue
		}

		err = os.WriteFile(filename, data, 0600)
		if err != nil {
			bf.logger.Error("Failed to update request", "error", err)
			continue
		}

		bf.logger.Warn("Blob fetch failed, will retry in 5 minutes",
			"slot", request.Header.Slot.Unwrap(),
			"failure_count", request.FailureCount)
	}
}

// removeRequestFile removes a request file and logs any errors.
func (bf *blobFetcher) removeRequestFile(filename string) {
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		bf.logger.Error("Failed to delete request file", "file", filename, "error", err)
	}
}

// getNextRequest reads the next request from disk queue.
// Returns the request, filename, and error.
func (bf *blobFetcher) getNextRequest() (BlobFetchRequest, string, error) {
	files, err := os.ReadDir(bf.queueDir)
	if err != nil {
		bf.logger.Error("Failed to read queue directory", "dir", bf.queueDir, "error", err)
		return BlobFetchRequest{}, "", fmt.Errorf("failed to read queue directory: %w", err)
	}

	headSlot := bf.blobRequester.HeadSlot()

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := filepath.Join(bf.queueDir, file.Name())
		var data []byte
		data, err = os.ReadFile(filename) // #nosec G304 // filename is constructed from queueDir
		if err != nil {
			return BlobFetchRequest{}, filename, fmt.Errorf("failed to read request file: %w", err)
		}

		var request BlobFetchRequest
		if err = json.Unmarshal(data, &request); err != nil {
			return BlobFetchRequest{}, filename, fmt.Errorf("failed to unmarshal request: %w", err)
		}

		// Check if request is outside availability window
		if headSlot > 0 && !bf.chainSpec.WithinDAPeriod(request.Header.Slot, headSlot) {
			bf.logger.Warn("Request is outside availability window, deleting",
				"slot", request.Header.Slot.Unwrap(),
				"head_slot", headSlot.Unwrap(),
				"failure_count", request.FailureCount)
			bf.removeRequestFile(filename)
			continue
		}

		// Check if this request needs to wait before retry
		if !request.LastRetryTime.IsZero() && time.Since(request.LastRetryTime) < 5*time.Minute {
			continue // Skip, not ready to retry yet
		}

		return request, filename, nil
	}

	return BlobFetchRequest{}, "", errNoMoreRequests
}

// processFetchRequest handles a single blob fetch request.
func (bf *blobFetcher) processFetchRequest(req BlobFetchRequest) error {
	bf.logger.Info("Fetching blobs from peers", "slot", req.Header.Slot.Unwrap(), "expected_blobs", len(req.Commitments))

	select {
	case <-bf.ctx.Done():
		return bf.ctx.Err()
	default:
	}

	// Create a verifier function that validates blobs against the stored header and commitments
	verifier := func(sidecars datypes.BlobSidecars) error {
		return bf.blobProcessor.VerifySidecars(bf.ctx, sidecars, req.Header, req.Commitments)
	}

	// Request blobs with verification - will try multiple peers if verification fails
	fetchedBlobs, err := bf.blobRequester.RequestBlobs(req.Header.Slot, len(req.Commitments), verifier)
	if err != nil {
		return fmt.Errorf("failed to request valid blobs for slot %d: %w", req.Header.Slot.Unwrap(), err)
	}

	// Process and store the validated blobs
	err = bf.blobProcessor.ProcessSidecars(bf.storageBackend.AvailabilityStore(), fetchedBlobs)
	if err != nil {
		return fmt.Errorf("failed to process blobs for slot %d: %w", req.Header.Slot.Unwrap(), err)
	}

	bf.logger.Info("Successfully fetched and stored blobs", "slot", req.Header.Slot.Unwrap(), "count", len(fetchedBlobs))
	return nil
}
