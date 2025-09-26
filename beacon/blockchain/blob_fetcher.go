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
	Slot        math.Slot                                    `json:"slot"`
	Header      *ctypes.BeaconBlockHeader                    `json:"header"`
	Commitments eip4844.KZGCommitments[common.ExecutionHash] `json:"commitments"`
}

// blobFetcher handles asynchronous fetching of blobs in the background.
type blobFetcher struct {
	logger         log.Logger
	blobProcessor  BlobProcessor
	blobRequester  BlobRequester
	storageBackend StorageBackend

	// Directory for persistent queue
	queueDir string
	// Channel to signal new requests available
	notifyChan chan struct{}
	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// NewBlobFetcher creates a new background blob fetcher.
func NewBlobFetcher(
	dataDir string,
	logger log.Logger,
	blobProcessor BlobProcessor,
	blobRequester BlobRequester,
	storageBackend StorageBackend,
) (BlobFetcher, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create queue directory
	queueDir := filepath.Join(dataDir, "blob_fetcher_queue")
	if err := os.MkdirAll(queueDir, 0755); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create blob fetcher queue directory: %w", err)
	}

	return &blobFetcher{
		logger:         logger,
		blobProcessor:  blobProcessor,
		blobRequester:  blobRequester,
		storageBackend: storageBackend,
		queueDir:       queueDir,
		notifyChan:     make(chan struct{}, 1), // Buffered to avoid blocking
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start begins the background blob fetching process.
func (bf *blobFetcher) Start() {
	// In case node crashed or was restarted, process any pending requests
	select {
	case bf.notifyChan <- struct{}{}:
	default:
	}

	go bf.run()
}

// Stop gracefully shuts down the blob fetcher.
func (bf *blobFetcher) Stop() {
	bf.cancel()
	close(bf.notifyChan)
}

// SetHeadSlot updates the head slot for blob fetching.
func (bf *blobFetcher) SetHeadSlot(slot math.Slot) {
	bf.blobRequester.SetHeadSlot(slot.Unwrap())
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
		Slot:        slot,
		Header:      block.GetHeader(),
		Commitments: commitments,
	}

	// Serialize to JSON file with slot as filename
	filename := filepath.Join(bf.queueDir, fmt.Sprintf("%020d.json", slot))
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal blob fetch request: %w", err)
	}

	if writeErr := os.WriteFile(filename, data, 0600); writeErr != nil {
		return fmt.Errorf("failed to queue blob fetch request: %w", writeErr)
	}

	bf.logger.Info("Queued blob fetch request", "slot", slot, "expected_blobs", len(commitments))

	// Signal that a new request is available
	select {
	case bf.notifyChan <- struct{}{}:
	default:
		// Already signaled
	}

	return nil
}

func (bf *blobFetcher) run() {
	for {
		select {
		case <-bf.ctx.Done():
			bf.logger.Info("Blob fetcher shutting down")
			return

		case <-bf.notifyChan:
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

		request, cleanup, err := bf.getNextRequest()
		if err != nil {
			if errors.Is(err, errNoMoreRequests) {
				return
			}
			bf.logger.Error("Failed to get next request", "error", err)
			cleanup()
			continue
		}

		err = bf.processFetchRequest(request)
		if err != nil {
			bf.logger.Error("Failed to process blob fetch request", "slot", request.Slot, "error", err)
		}

		cleanup()
	}
}

// getNextRequest reads the next request from disk queue.
// Returns the request and a cleanup function to remove the file after processing.
func (bf *blobFetcher) getNextRequest() (BlobFetchRequest, func(), error) {
	var request BlobFetchRequest
	noopCleanup := func() {}

	// List all request files (already sorted by name)
	files, err := os.ReadDir(bf.queueDir)
	if err != nil {
		bf.logger.Error("Failed to read queue directory", "dir", bf.queueDir, "error", err)
		return request, noopCleanup, fmt.Errorf("failed to read queue directory: %w", err)
	}

	var filename string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			filename = filepath.Join(bf.queueDir, file.Name())
			break
		}
	}

	if len(filename) == 0 {
		return request, noopCleanup, errNoMoreRequests
	}

	// Create cleanup function that will always try to remove the file
	cleanup := func() {
		removeErr := os.Remove(filename)
		if removeErr != nil && !os.IsNotExist(removeErr) {
			bf.logger.Error("Failed to delete request file", "file", filename, "error", removeErr)
		}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return request, cleanup, fmt.Errorf("failed to read request file: %w", err)
	}

	if err = json.Unmarshal(data, &request); err != nil {
		return request, cleanup, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	return request, cleanup, nil
}

// processFetchRequest handles a single blob fetch request.
func (bf *blobFetcher) processFetchRequest(req BlobFetchRequest) error {
	bf.logger.Info("Fetching blobs from peers", "slot", req.Slot, "expected_blobs", len(req.Commitments))

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
	fetchedBlobs, err := bf.blobRequester.RequestBlobs(req.Slot.Unwrap(), verifier)
	if err != nil {
		return fmt.Errorf("failed to request valid blobs for slot %d: %w", req.Slot, err)
	}

	// Process and store the validated blobs
	err = bf.blobProcessor.ProcessSidecars(bf.storageBackend.AvailabilityStore(), fetchedBlobs)
	if err != nil {
		return fmt.Errorf("failed to process blobs for slot %d: %w", req.Slot, err)
	}

	bf.logger.Info("Successfully fetched and stored blobs", "slot", req.Slot, "count", len(fetchedBlobs))
	return nil
}
