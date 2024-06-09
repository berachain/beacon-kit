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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

// DBPruner is a struct that holds the prunable interface and a notifier
// channel.
type DBPruner[
	BeaconBlockT BeaconBlock,
	BlockEventT BlockEvent[BeaconBlockT],
	PrunableT Prunable,
	SubscriptionT Subscription,
] struct {
	prunable     Prunable
	logger       log.Logger[any]
	name         string
	feed         BlockFeed[BeaconBlockT, BlockEventT, SubscriptionT]
	pruneRangeFn func(BlockEventT) (uint64, uint64)
}

func NewPruner[
	BeaconBlockT BeaconBlock,
	BlockEventT BlockEvent[BeaconBlockT],
	PrunableT Prunable,
	SubscriptionT Subscription,
](
	logger log.Logger[any],
	prunable Prunable,
	name string,
	feed BlockFeed[BeaconBlockT, BlockEventT, SubscriptionT],
	pruneRangeFn func(BlockEventT) (uint64, uint64),
) *DBPruner[BeaconBlockT, BlockEventT, PrunableT, SubscriptionT] {
	return &DBPruner[BeaconBlockT, BlockEventT, PrunableT, SubscriptionT]{
		logger:       logger,
		prunable:     prunable,
		name:         name,
		feed:         feed,
		pruneRangeFn: pruneRangeFn,
	}
}

// Start starts the Pruner by listening for new indexes to prune.
func (p *DBPruner[
	BeaconBlockT, BlockEventT, PrunableT, SubscriptionT,
]) Start(ctx context.Context) {
	ch := make(chan BlockEventT)
	sub := p.feed.Subscribe(ch)
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:
				if event.Is(events.BeaconBlockFinalized) {
					start, end := p.pruneRangeFn(event)
					if err := p.prunable.Prune(start, end); err != nil {
						p.logger.Error(
							"‼️ error pruning index ‼️",
							"error", err,
						)
					}
				}
			}
		}
	}()
}

// Name returns the name of the Pruner.
func (p *DBPruner[
	BeaconBlockT, BlockEventT, PrunableT, SubscriptionT,
]) Name() string {
	return p.name
}
