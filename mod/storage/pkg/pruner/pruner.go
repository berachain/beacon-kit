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

	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
)

// Compile-time check to ensure pruner implements the Pruner interface.
var _ Pruner[Prunable] = (*pruner[
	BeaconBlock, BlockEvent[BeaconBlock], Prunable,
])(nil)

// pruner is a struct that holds the prunable interface and a notifier
// channel.
type pruner[
	BeaconBlockT BeaconBlock,
	BlockEventT BlockEvent[BeaconBlockT],
	PrunableT Prunable,
] struct {
	prunable              Prunable
	logger                log.Logger[any]
	name                  string
	finalizedBlkEvents    chan BlockEventT
	finalizedBlockEventID types.EventID
	dispatcher            *dispatcher.Dispatcher
	pruneRangeFn          func(BlockEventT) (uint64, uint64)
}

// NewPruner creates a new Pruner.
func NewPruner[
	BeaconBlockT BeaconBlock,
	BlockEventT BlockEvent[BeaconBlockT],
	PrunableT Prunable,
](
	logger log.Logger[any],
	prunable Prunable,
	name string,
	finalizedBlockEventID types.EventID,
	dispatcher *dispatcher.Dispatcher,
	pruneRangeFn func(BlockEventT) (uint64, uint64),
) Pruner[PrunableT] {
	return &pruner[BeaconBlockT, BlockEventT, PrunableT]{
		logger:                logger,
		prunable:              prunable,
		name:                  name,
		finalizedBlkEvents:    make(chan BlockEventT),
		finalizedBlockEventID: finalizedBlockEventID,
		dispatcher:            dispatcher,
		pruneRangeFn:          pruneRangeFn,
	}
}

// Start starts the Pruner by listening for new indexes to prune.
func (p *pruner[_, BlockEventT, _]) Start(ctx context.Context) {
	if err := p.dispatcher.Subscribe(
		p.finalizedBlockEventID, p.finalizedBlkEvents,
	); err != nil {
		p.logger.Error("failed to subscribe to event", "event",
			p.finalizedBlockEventID, "err", err)
		return
	}
	go p.start(ctx)
}

// start listens for new indexes to prune.
func (p *pruner[_, BlockEventT, _]) start(
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-p.finalizedBlkEvents:
			start, end := p.pruneRangeFn(event)
			if err := p.prunable.Prune(start, end); err != nil {
				p.logger.Error("‼️ error pruning index ‼️", "error", err)
			}
		}
	}
}

// Name returns the name of the Pruner.
func (p *pruner[_, _, _]) Name() string {
	return p.name
}
