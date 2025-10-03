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
	"path/filepath"
	"sync"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BlobFetcherConfig contains configuration for the blob fetcher.
type BlobFetcherConfig struct {
	// CheckInterval is how often to check the queue for pending requests
	CheckInterval time.Duration
	// RetryInterval is the minimum time between retry attempts per blob request
	RetryInterval time.Duration
	// MaxRetries is the maximum number of retry attempts per blob request before giving up and deleting it
	MaxRetries int
}

// DefaultBlobFetcherConfig returns the default configuration.
//
//nolint:mnd // Just defaults
func DefaultBlobFetcherConfig() BlobFetcherConfig {
	return BlobFetcherConfig{
		CheckInterval: 1 * time.Minute,
		RetryInterval: 5 * time.Minute,
		MaxRetries:    72, // 6 hours at 5 minute intervals
	}
}

// blobFetcher handles asynchronous fetching of blobs in the background.
type blobFetcher struct {
	logger    log.Logger
	chainSpec BlobFetcherChainSpec
	queue     *blobQueue         // Queue for persistent requests
	executor  *blobFetchExecutor // Executor for fetch logic
	config    BlobFetcherConfig  // Configuration

	// We need to track current head slot so we know when blob download requests need to be pruned as they are outside the WithinDAPeriod
	headSlotMu sync.RWMutex
	headSlot   math.Slot

	ctx      context.Context
	cancel   context.CancelFunc
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
	config BlobFetcherConfig,
) (BlobFetcher, error) {
	queue, err := newBlobQueue(filepath.Join(dataDir, "blobs", "download_queue"), logger)
	if err != nil {
		return nil, err
	}

	return &blobFetcher{
		logger:    logger,
		chainSpec: chainSpec,
		queue:     queue,
		config:    config,
		executor: &blobFetchExecutor{
			blobProcessor:  blobProcessor,
			blobRequester:  blobRequester,
			storageBackend: storageBackend,
			logger:         logger,
		},
	}, nil
}

// Start begins the background blob fetching process.
func (bf *blobFetcher) Start(ctx context.Context) {
	bf.ctx, bf.cancel = context.WithCancel(ctx)
	go bf.run()
}

// Stop gracefully shuts down the blob fetcher.
func (bf *blobFetcher) Stop() {
	bf.stopOnce.Do(func() { bf.cancel() })
}

// SetHeadSlot updates the head slot for blob fetching.
func (bf *blobFetcher) SetHeadSlot(slot math.Slot) {
	bf.headSlotMu.Lock()
	bf.headSlot = slot
	bf.headSlotMu.Unlock()

	// Also update the reactor's head slot so it can respond correctly to peers
	bf.executor.blobRequester.SetHeadSlot(slot)
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

	if err := bf.queue.Add(slot, request); err != nil {
		return err
	}

	bf.logger.Info("Queued blob fetch request", "slot", slot.Unwrap(), "expected_blobs", len(commitments))
	return nil
}

func (bf *blobFetcher) run() {
	// Ticker to periodically check for requests (both new and ready to retry)
	ticker := time.NewTicker(bf.config.CheckInterval)
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

		bf.headSlotMu.RLock()
		headSlot := bf.headSlot
		bf.headSlotMu.RUnlock()

		request, filename, err := bf.queue.GetNext(headSlot, bf.config.RetryInterval, bf.config.MaxRetries, bf.chainSpec.WithinDAPeriod)
		if err != nil {
			if errors.Is(err, errNoMoreRequests) {
				return
			}
			bf.logger.Error("Failed to get next request", "error", err)
			if filename != "" {
				_ = bf.queue.Remove(filename)
			}
			continue
		}

		err = bf.executor.FetchBlobsAndVerify(bf.ctx, request)
		if err == nil {
			// Successfully processed, remove the request file
			_ = bf.queue.Remove(filename)
			continue
		}

		bf.logger.Error("Failed to process blob fetch request", "slot", request.Header.Slot.Unwrap(), "error", err)

		// Update retry metadata and save back to file
		if updateErr := bf.queue.UpdateRetry(filename, request); updateErr != nil {
			bf.logger.Error("Failed to update retry metadata", "error", updateErr)
			continue
		}

		bf.logger.Warn("Blob fetch failed, will retry later",
			"slot", request.Header.Slot.Unwrap(),
			"failure_count", request.FailureCount+1)
	}
}
