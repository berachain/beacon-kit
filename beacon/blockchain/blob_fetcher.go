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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
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
	"github.com/berachain/beacon-kit/da/blobreactor"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BlobFetcherConfig configures the background blob fetcher.
type BlobFetcherConfig struct {
	// CheckInterval is how often the queue is scanned for pending requests.
	CheckInterval time.Duration
	// RetryInterval is the minimum time between attempts for one request. There is no retry cap: requests are retried until they succeed or
	// their slot leaves the DA window.
	RetryInterval time.Duration
}

//nolint:mnd // defaults
func DefaultBlobFetcherConfig() BlobFetcherConfig {
	return BlobFetcherConfig{
		CheckInterval: 5 * time.Second,
		RetryInterval: 30 * time.Second,
	}
}

// Metric keys for the background blob fetcher and its queue.
const (
	metricsBlobQueueDepth        = "beacon_kit.blockchain.blob_fetcher.queue_depth"
	metricsBlobRequestsQueued    = "beacon_kit.blockchain.blob_fetcher.requests_queued"
	metricsBlobRequestsCompleted = "beacon_kit.blockchain.blob_fetcher.requests_completed"
	metricsBlobRequestsExpired   = "beacon_kit.blockchain.blob_fetcher.requests_expired"
	metricsBlobRetries           = "beacon_kit.blockchain.blob_fetcher.retries"
)

// blobFetcher is the background half of blob distribution. It covers every finalized block whose sidecars did
// not arrive through the tip-of-chain lanes: blocks replayed during catch-up (which carry no sidecars once
// blob consensus is enabled) and tip blocks whose synchronous fetch failed. FinalizeSidecars queues one
// request per block, and the fetcher drains the queue by asking peers for whole slot windows over the blob
// reactor.
//
// Its guarantees are the strict half of the design. Every returned slot is verified against the header,
// commitments and block signature recorded at queue time. An empty or short response is a failure, never a
// success. A request is retried until it succeeds or its slot leaves the DA window, and the queue is a
// crash-safe directory of JSON files, so pending fetches survive restarts. The queue depth gates the node's
// synced status, so "synced" keeps implying "holds all in-window blobs".
type blobFetcher struct {
	logger        log.Logger
	chainSpec     BlobFetcherChainSpec
	queue         *blobQueue
	blobProcessor BlobProcessor
	requester     BlobRequester
	storage       StorageBackend
	config        BlobFetcherConfig
	sink          TelemetrySink

	// Head slot gates DA-window expiry of queued requests.
	headSlotMu sync.RWMutex
	headSlot   math.Slot

	// kick wakes the run loop right after new work is queued.
	kick chan struct{}

	ctx      context.Context
	cancel   context.CancelFunc
	stopOnce sync.Once
	// done is closed when the run loop exits, so Stop can wait for in-flight work.
	done chan struct{}
}

// BlobFetcherChainSpec is the part of the chain spec the fetcher needs.
type BlobFetcherChainSpec interface {
	WithinDAPeriod(block, current math.Slot) bool
}

// NewBlobFetcher creates the background blob fetcher; its queue lives under dataDir.
func NewBlobFetcher(
	dataDir string,
	logger log.Logger,
	blobProcessor BlobProcessor,
	requester BlobRequester,
	storageBackend StorageBackend,
	chainSpec BlobFetcherChainSpec,
	config BlobFetcherConfig,
	telemetrySink TelemetrySink,
) (BlobFetcher, error) {
	queue, err := newBlobQueue(filepath.Join(dataDir, "blobs", "download_queue"), logger, telemetrySink)
	if err != nil {
		return nil, err
	}

	return &blobFetcher{
		logger:        logger,
		chainSpec:     chainSpec,
		queue:         queue,
		blobProcessor: blobProcessor,
		requester:     requester,
		storage:       storageBackend,
		config:        config,
		sink:          telemetrySink,
		kick:          make(chan struct{}, 1),
		done:          make(chan struct{}),
	}, nil
}

func (bf *blobFetcher) Start(ctx context.Context) {
	bf.ctx, bf.cancel = context.WithCancel(ctx)
	go bf.run()
}

// Stop cancels the run loop and waits for it to exit, so no fetch is still touching the queue or store when
// shutdown proceeds to the storage backends.
func (bf *blobFetcher) Stop() {
	bf.stopOnce.Do(func() {
		if bf.cancel != nil {
			bf.cancel()
			<-bf.done
		}
	})
}

// SetHeadSlot updates the fetcher's (and the reactor's) view of the chain head.
func (bf *blobFetcher) SetHeadSlot(slot math.Slot) {
	bf.headSlotMu.Lock()
	bf.headSlot = slot
	bf.headSlotMu.Unlock()

	bf.requester.SetHeadSlot(slot)
}

func (bf *blobFetcher) getHeadSlot() math.Slot {
	bf.headSlotMu.RLock()
	defer bf.headSlotMu.RUnlock()
	return bf.headSlot
}

// PendingRequests returns the number of queued blob-fetch requests. A node with pending in-window requests must not report itself as
// synced.
func (bf *blobFetcher) PendingRequests() int {
	return bf.queue.PendingCount()
}

// QueueBlobRequest queues an asynchronous fetch of the given finalized block's sidecars. It records the header, commitments and block
// signature so fetched sidecars can be fully verified later without the block.
func (bf *blobFetcher) QueueBlobRequest(signedBlk *ctypes.SignedBeaconBlock) error {
	blk := signedBlk.GetBeaconBlock()
	commitments := blk.GetBody().GetBlobKzgCommitments()
	if len(commitments) == 0 {
		return nil
	}

	slot := blk.GetHeader().GetSlot()
	if err := bf.queue.Add(slot, BlobFetchRequest{
		Header:      blk.GetHeader(),
		Commitments: commitments,
		Signature:   signedBlk.GetSignature(),
	}); err != nil {
		return err
	}

	bf.sink.IncrementCounter(metricsBlobRequestsQueued)
	bf.logger.Info("Queued blob fetch request", "slot", slot.Unwrap(), "expected_blobs", len(commitments))

	// Wake the run loop; no-op if a kick is already pending.
	select {
	case bf.kick <- struct{}{}:
	default:
	}
	return nil
}

func (bf *blobFetcher) run() {
	defer close(bf.done)
	ticker := time.NewTicker(bf.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-bf.ctx.Done():
			bf.logger.Info("Blob fetcher shutting down")
			return
		case <-ticker.C:
		case <-bf.kick:
		}
		bf.processPending()
	}
}

// processPending drains all currently-ready queue entries, batching contiguous slot windows into by-range requests.
func (bf *blobFetcher) processPending() {
	ready, err := bf.queue.LoadReady(bf.getHeadSlot(), bf.config.RetryInterval, bf.chainSpec.WithinDAPeriod)
	if err != nil {
		bf.logger.Error("Failed to load blob fetch queue", "error", err)
		return
	}
	if len(ready) == 0 {
		return
	}

	// A crash between persisting and removing the queue entry (or a replay of an already-fetched block) leaves an entry the store
	// already satisfies. Complete those without a fetch; refetching would hold the node in "syncing" for as long as no peer serves
	// the slot, despite it holding the data.
	pending := ready[:0]
	for _, entry := range ready {
		req := entry.request
		if !sidecarsAlreadyStored(bf.storage.AvailabilityStore(), req.Header, req.Commitments) {
			pending = append(pending, entry)
			continue
		}
		bf.queue.Remove(entry.filename)
		bf.sink.IncrementCounter(metricsBlobRequestsCompleted)
		bf.logger.Info("Blob sidecars already in store, completing queued fetch",
			"slot", entry.request.Header.GetSlot().Unwrap())
	}
	ready = pending

	// LoadReady returns entries in slot order (files are zero-padded).
	for start := 0; start < len(ready); {
		if bf.ctx.Err() != nil {
			return
		}

		// Batch every entry that fits in one by-range window.
		windowStart := ready[start].request.Header.GetSlot()
		end := start + 1
		for end < len(ready) &&
			ready[end].request.Header.GetSlot()-windowStart < blobreactor.MaxRequestedSlots {
			end++
		}
		bf.fetchWindow(ready[start:end])
		start = end
	}
}

// fetchWindow issues one by-range request covering the given queue entries and persists every slot that verifies; the rest are scheduled
// for retry.
func (bf *blobFetcher) fetchWindow(window []queuedRequest) {
	var (
		firstSlot = window[0].request.Header.GetSlot()
		lastSlot  = window[len(window)-1].request.Header.GetSlot()
		count     = lastSlot.Unwrap() - firstSlot.Unwrap() + 1
		expected  = make(map[math.Slot]*queuedRequest, len(window))
	)
	for i := range window {
		expected[window[i].request.Header.GetSlot()] = &window[i]
	}

	// Peers may return in-range slots we did not queue (they cannot know which slots we still need); those are skipped, not failures.
	verify := func(slot math.Slot, sidecars datypes.BlobSidecars) error {
		entry, ok := expected[slot]
		if !ok {
			return blobreactor.ErrSlotNotRequested
		}
		req := entry.request
		return verifySidecarsBinding(bf.ctx, bf.blobProcessor, req.Header, req.Commitments, req.Signature, sidecars)
	}

	// RequestSidecarsByRange returns an error only when no peer yielded anything usable.
	verified, err := bf.requester.RequestSidecarsByRange(bf.ctx, firstSlot, count, verify)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		bf.logger.Warn("Blob fetch attempt failed, will retry",
			"start_slot", firstSlot.Unwrap(), "slots", len(window), "error", err)
		for _, entry := range window {
			bf.sink.IncrementCounter(metricsBlobRetries)
			bf.queue.UpdateRetry(entry)
		}
		return
	}

	headSlot := bf.getHeadSlot()
	for slot, entry := range expected {
		sidecars, ok := verified[slot]
		if !ok {
			bf.sink.IncrementCounter(metricsBlobRetries)
			bf.queue.UpdateRetry(*entry)
			continue
		}
		// The slot may have left the DA window during the fetch round trip. Pruning would then have advanced
		// the store's lower bound past it, and persisting now would be dropped by the store (or, without the
		// store-level guard, clamped onto another slot). Drop the request instead of persisting stale data.
		if headSlot > 0 && !bf.chainSpec.WithinDAPeriod(slot, headSlot) {
			bf.sink.IncrementCounter(metricsBlobRequestsExpired)
			bf.queue.Remove(entry.filename)
			continue
		}
		if processErr := bf.blobProcessor.ProcessSidecars(bf.storage.AvailabilityStore(), sidecars); processErr != nil {
			bf.logger.Error("Failed to persist fetched blob sidecars",
				"slot", slot.Unwrap(), "error", processErr)
			bf.sink.IncrementCounter(metricsBlobRetries)
			bf.queue.UpdateRetry(*entry)
			continue
		}
		bf.queue.Remove(entry.filename)
		bf.sink.IncrementCounter(metricsBlobRequestsCompleted)
		bf.logger.Info("Fetched and stored blob sidecars",
			"slot", slot.Unwrap(), "count", len(sidecars))
	}
}
