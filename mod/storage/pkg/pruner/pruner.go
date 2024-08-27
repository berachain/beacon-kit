// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package pruner

import (
	"context"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// Compile-time check to ensure pruner implements the Pruner interface.
var _ Pruner[Prunable] = (*pruner[BeaconBlock, Prunable])(nil)

// pruner is a struct that holds the prunable interface and a notifier
// channel.
type pruner[
	BeaconBlockT BeaconBlock,
	PrunableT Prunable,
] struct {
	prunable                Prunable
	logger                  log.Logger
	name                    string
	subBeaconBlockFinalized chan async.Event[BeaconBlockT]
	pruneRangeFn            func(async.Event[BeaconBlockT]) (uint64, uint64)
}

// NewPruner creates a new Pruner.
func NewPruner[
	BeaconBlockT BeaconBlock,
	PrunableT Prunable,
](
	logger log.Logger,
	prunable Prunable,
	name string,
	subBeaconBlockFinalized chan async.Event[BeaconBlockT],
	pruneRangeFn func(async.Event[BeaconBlockT]) (uint64, uint64),
) Pruner[PrunableT] {
	return &pruner[BeaconBlockT, PrunableT]{
		logger:                  logger,
		prunable:                prunable,
		name:                    name,
		pruneRangeFn:            pruneRangeFn,
		subBeaconBlockFinalized: subBeaconBlockFinalized,
	}
}

// Start starts the Pruner by listening for new indexes to prune.
func (p *pruner[BeaconBlockT, PrunableT]) Start(ctx context.Context) {
	go p.listen(ctx)
}

// listen listens for new finalized blocks and prunes the prunable store based
// on the received finalized block event.
func (p *pruner[_, PrunableT]) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-p.subBeaconBlockFinalized:
			p.onFinalizeBlock(event)
		}
	}
}

// onFinalizeBlock will prune the prunable store based on the received
// finalized block event.
func (p *pruner[BeaconBlockT, PrunableT]) onFinalizeBlock(
	event async.Event[BeaconBlockT],
) {
	start, end := p.pruneRangeFn(event)
	if err := p.prunable.Prune(start, end); err != nil {
		p.logger.Error("‼️ error pruning index ‼️", "error", err)
	}
}

// Name returns the name of the Pruner.
func (p *pruner[_, _]) Name() string {
	return p.name
}
