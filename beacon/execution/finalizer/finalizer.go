// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package finalizer

import (
	"context"
	"sync"

	"cosmossdk.io/core/header"
)

const defaultBufferSize = 256

// TODO: maybe we want to create some sort of filter that ensures that we call the sync every
// X blocks or so?
// I'm thinking the issue that could arise from the current optimizations, is that by the time
// CometBFT kicks into Consensus Mode from Replay mode, the execution client might not be
// fully sync'd yet
// and then if we just so happen to be chosen to propose a block shortly there after, we may run
// into problems.
// so TLDR, maybe we only discard sync requests if they are within X blocks (maybe 64 blocks)
// of the previous one?

// NewFinalizationRequest creates a new FinalizationRequest with the provided
// blockHash and beaconBlock.
// blockHash is the hash of the block to be finalized.
// beaconBlock is the header information of the beacon block.
func NewFinalizationRequest(blockHash []byte, beaconBlock header.Info) *FinalizationRequest {
	return &FinalizationRequest{
		blockHash: blockHash,
		slot:      uint64(beaconBlock.Height),
	}
}

// Finalizer is responsible for finalizing execution layer blocks. It's goal is to allow us to
// decouple the finalization of execution layer blocks from the ABCI lifecycle.
type Finalizer struct {
	// engineCaller is the caller of the execution engine.
	engineCaller Enginecaller

	// req is a queue of finalization requests, this allows us
	// to decouple the finalization of execution layer blocks
	// from the ABCI lifecycle.
	signal        chan struct{}
	latestRequest *FinalizationRequest
	lrMu          sync.Mutex
}

// New creates a new Finalizer with the provided Enginecaller.
func New(engineCaller Enginecaller) *Finalizer {
	return &Finalizer{
		engineCaller: engineCaller,
		signal:       make(chan struct{}, defaultBufferSize),
	}
}

// Start begins the Finalizer's main loop in a new goroutine.
func (f *Finalizer) Start(ctx context.Context) {
	go f.run(ctx)
}

// run is the main loop for the Finalizer. It listens for finalization requests
// and processes them.
func (f *Finalizer) run(ctx context.Context) {
	for {
		// We acquire the lock to prevent messinesses between emptying the queue in
		// `RequestionFinalization` and the above `select` statement. This lock ensures
		// that if we are about to push a new request onto the channel, we will wait
		// to ensure we are going to get the latest request, opposed to consuming a
		// slightly older one that is going to be replaced momentarily.
		select {
		case <-ctx.Done():
			return
		case <-f.signal:
			f.lrMu.Lock()
			// This is safe, since there will never be a signal on the channel
			// unless there is a request in the queue.
			copiedLr := *f.latestRequest
			f.latestRequest = nil
			f.lrMu.Unlock()
			f.processFinalizationRequest(ctx, &copiedLr)
		}
	}
}

// RequestFinalization requests the finalization of an execution layer block.
func (f *Finalizer) RequestFinalization(blockHash []byte, beaconBlock header.Info) {
	// When we push on a new request, we want to attempt to abort the previous.
	// TODO: we actually only want to abort the previous if it is trying to sync
	// a block that is in the FUTURE of the block we are trying to finalize.
	// We should probably add some safety checks to this finalization requestor
	// to prevent making calls to the execution client that would result in
	// bad things happening.
	f.lrMu.Lock()
	f.latestRequest = NewFinalizationRequest(blockHash, beaconBlock)
	// We clear out the channel to ensure that we are only ever processing the
	// latest request.
	for len(f.signal) > 0 {
		<-f.signal
	}
	f.lrMu.Unlock()
	f.signal <- struct{}{}
}

func (f *Finalizer) processFinalizationRequest(ctx context.Context, r *FinalizationRequest) {
	// The finalization process can be cancelled by calling the cancel function of the context.
	_ = f.engineCaller.FinalizeBlock(ctx, r.slot, r.blockHash)
}
