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

package deposit

import (
	"context"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Service represenst the deposit service that processes deposit events.
type Service[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BlockEventT BlockEvent[
		BeaconBlockT, BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	ExecutionPayloadT interface{ GetNumber() math.U64 },
	SubscriptionT interface {
		Unsubscribe()
	},
	DepositT interface{ GetIndex() uint64 },
] struct {
	// logger is used for logging information and errors.
	logger log.Logger[any]
	// eth1FollowDistance is the follow distance for Ethereum 1.0 blocks.
	eth1FollowDistance math.U64
	// ethclient is the Ethereum 1.0 client.
	ethclient EthClient
	// dc is the contract interface for interacting with the deposit contract.
	dc Contract[DepositT]
	// ds is the deposit store that stores deposits.
	ds Store[DepositT]
	// feed is the block feed that provides block events.
	feed BlockFeed[BeaconBlockT, BeaconBlockBodyT, BlockEventT,
		DepositT, ExecutionPayloadT, SubscriptionT]
	// newBlock is the channel for new blocks.
	newBlock chan BeaconBlockT
	// failedBlocks
	failedBlocks map[math.U64]struct{}
}

// NewService creates a new instance of the Service struct.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BlockEventT BlockEvent[
		BeaconBlockT, BeaconBlockBodyT, DepositT, ExecutionPayloadT],
	DepositStoreT Store[DepositT],
	ExecutionPayloadT interface{ GetNumber() math.U64 },
	SubscriptionT interface {
		Unsubscribe()
	},
	DepositT interface{ GetIndex() uint64 },
](
	logger log.Logger[any],
	eth1FollowDistance math.U64,
	ethclient EthClient,
	ds Store[DepositT],
	dc Contract[DepositT],
	feed BlockFeed[
		BeaconBlockT, BeaconBlockBodyT, BlockEventT,
		DepositT, ExecutionPayloadT, SubscriptionT,
	],
) *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
] {
	return &Service[
		BeaconBlockT, BeaconBlockBodyT, BlockEventT,
		ExecutionPayloadT, SubscriptionT, DepositT,
	]{
		feed:               feed,
		logger:             logger,
		ethclient:          ethclient,
		eth1FollowDistance: eth1FollowDistance,
		dc:                 dc,
		ds:                 ds,
		newBlock:           make(chan BeaconBlockT),
		failedBlocks:       make(map[math.U64]struct{}),
	}
}

// Start starts the service and begins processing block events.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) Start(
	ctx context.Context,
) error {
	go s.blockFeedListener(ctx)
	go s.depositFetcher(ctx)
	go s.depositCatchupFetcher(ctx)
	return nil
}

func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) blockFeedListener(ctx context.Context) {
	ch := make(chan BlockEventT)
	sub := s.feed.Subscribe(ch)
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			s.newBlock <- event.Block()
		}
	}
}

// Name returns the name of the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) Name() string {
	return "deposit-handler"
}

// Status returns the current status of the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) Status() error {
	return nil
}

// WaitForHealthy waits for the service to become healthy.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) WaitForHealthy(
	_ context.Context,
) {
}
