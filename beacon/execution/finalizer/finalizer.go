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

	"cosmossdk.io/core/header"
)

// TODO: make this not hte blockchain service.
type Enginecaller interface {
	FinalizeBlock(ctx context.Context, beaconBlock header.Info, blockHash []byte) error
}

type FinalizationRequest struct {
	blockHash   []byte
	beaconBlock header.Info
}

func NewFinalizationRequest(blockHash []byte, beaconBlock header.Info) *FinalizationRequest {
	return &FinalizationRequest{
		blockHash:   blockHash,
		beaconBlock: beaconBlock,
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
	req chan *FinalizationRequest

	// cancelFunc is used to cancel the previous finalization request.
	cancelFunc context.CancelFunc
}

// New creates a new Finalizer with the provided Enginecaller.
func New(engineCaller Enginecaller) *Finalizer {
	return &Finalizer{
		engineCaller: engineCaller,
		req:          make(chan *FinalizationRequest),
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
		select {
		case <-ctx.Done():
			return
		case r := <-f.req:
			var cancelCtx context.Context
			cancelCtx, f.cancelFunc = context.WithCancel(ctx)
			// todo: this gorountine might get kind of nasty when we need to sync
			// a really old chain. We should figure out a way to make this nicer
			// when running in replay mode.
			go f.processFinalizationRequest(cancelCtx, r)
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
	if f.cancelFunc != nil {
		f.cancelFunc()
	}

	// Clear the channel since we always want to use the latest request.
	for len(f.req) > 0 {
		<-f.req
	}

	f.req <- &FinalizationRequest{blockHash: blockHash, beaconBlock: beaconBlock}
}

func (f *Finalizer) processFinalizationRequest(ctx context.Context, r *FinalizationRequest) {
	// The finalization process can be cancelled by calling the cancel function of the context.
	_ = f.engineCaller.FinalizeBlock(ctx, r.beaconBlock, r.blockHash)
}
