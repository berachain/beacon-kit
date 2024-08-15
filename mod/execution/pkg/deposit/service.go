// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package deposit

import (
	"context"

	async "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Service represents the deposit service that processes deposit events.
type Service[
	BeaconBlockT BeaconBlock[DepositT, BeaconBlockBodyT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BlockEventT BlockEvent[
		DepositT, BeaconBlockBodyT, BeaconBlockT, ExecutionPayloadT,
	],
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
	ExecutionPayloadT ExecutionPayload,
	WithdrawalCredentialsT any,
] struct {
	// logger is used for logging information and errors.
	logger log.Logger[any]
	// eth1FollowDistance is the follow distance for Ethereum 1.0 blocks.
	eth1FollowDistance math.U64
	// dc is the contract interface for interacting with the deposit contract.
	dc Contract[DepositT]
	// ds is the deposit store that stores deposits.
	ds Store[DepositT]
	// dispatcher is the dispatcher for the service.
	dispatcher async.EventDispatcher
	// finalizedBlockEventID is the event ID for the finalized block event.
	finalizedBlockEventID async.EventID
	// subFinalizedBlockEvents is the channel that provides finalized block events.
	subFinalizedBlockEvents async.Subscription[BlockEventT]
	// metrics is the metrics for the deposit service.
	metrics *metrics
	// failedBlocks is a map of blocks that failed to be processed to be
	// retried.
	failedBlocks map[math.U64]struct{}
}

// NewService creates a new instance of the Service struct.
func NewService[
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	BeaconBlockT BeaconBlock[DepositT, BeaconBlockBodyT, ExecutionPayloadT],
	BlockEventT BlockEvent[
		DepositT, BeaconBlockBodyT,
		BeaconBlockT, ExecutionPayloadT,
	],
	DepositStoreT Store[DepositT],
	ExecutionPayloadT ExecutionPayload,
	WithdrawalCredentialsT any,
	DepositT Deposit[DepositT, WithdrawalCredentialsT],
](
	logger log.Logger[any],
	eth1FollowDistance math.U64,
	telemetrySink TelemetrySink,
	ds Store[DepositT],
	dc Contract[DepositT],
	finalizedBlockEventID async.EventID,
	dispatcher async.EventDispatcher,
) *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT, DepositT,
	ExecutionPayloadT, WithdrawalCredentialsT,
] {
	return &Service[
		BeaconBlockT, BeaconBlockBodyT, BlockEventT, DepositT,
		ExecutionPayloadT,
		WithdrawalCredentialsT,
	]{
		dc:                      dc,
		dispatcher:              dispatcher,
		ds:                      ds,
		eth1FollowDistance:      eth1FollowDistance,
		failedBlocks:            make(map[math.Slot]struct{}),
		finalizedBlockEventID:   finalizedBlockEventID,
		subFinalizedBlockEvents: async.NewSubscription[BlockEventT](),
		logger:                  logger,
		metrics:                 newMetrics(telemetrySink),
	}
}

// Start starts the service and begins processing block events.
func (s *Service[
	_, _, _, _, _, _,
]) Start(ctx context.Context) error {
	if err := s.dispatcher.Subscribe(
		s.finalizedBlockEventID, s.subFinalizedBlockEvents,
	); err != nil {
		s.logger.Error("failed to subscribe to event", "event",
			s.finalizedBlockEventID, "err", err)
		return err
	}

	// Listen for finalized block events and fetch deposits for the block.
	s.subFinalizedBlockEvents.Listen(ctx, s.getDepositFetcher(ctx))

	// Catchup deposits for failed blocks.
	go s.depositCatchupFetcher(ctx)
	return nil
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _,
]) Name() string {
	return "deposit-handler"
}
