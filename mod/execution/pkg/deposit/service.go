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

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Service represents the deposit service that processes deposit events.
type Service[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	DepositT Deposit[DepositT, LogT, WithdrawalCredentialsT],
	ExecutionPayloadT ExecutionPayload,
	LogT any,
	WithdrawalCredentialsT any,
] struct {
	// logger is used for logging information and errors.
	logger log.Logger
	// eth1FollowDistance is the follow distance for Ethereum 1.0 blocks.
	eth1FollowDistance math.U64
	// dc is the contract interface for interacting with the deposit contract.
	dc Contract[DepositT]
	// ds is the deposit store that stores deposits.
	ds Store[DepositT]
	// dispatcher is the dispatcher for the service.
	dispatcher asynctypes.EventDispatcher
	// subFinalizedBlockEvents is the channel holding BeaconBlockFinalized
	// events.
	subFinalizedBlockEvents chan async.Event[BeaconBlockT]
	// metrics is the metrics for the deposit service.
	metrics *metrics
	// failedBlocks is a map of blocks that failed to be processed to be
	// retried.
	failedBlocks map[math.U64]struct{}
}

// NewService creates a new instance of the Service struct.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[DepositT, ExecutionPayloadT],
	DepositT Deposit[DepositT, LogT, WithdrawalCredentialsT],
	ExecutionPayloadT ExecutionPayload,
	LogT any,
	WithdrawalCredentialsT any,
](
	logger log.Logger,
	eth1FollowDistance math.U64,
	telemetrySink TelemetrySink,
	ds Store[DepositT],
	dc Contract[DepositT],
	dispatcher asynctypes.EventDispatcher,
) *Service[
	BeaconBlockT, BeaconBlockBodyT, DepositT,
	ExecutionPayloadT, LogT, WithdrawalCredentialsT,
] {
	return &Service[
		BeaconBlockT, BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, LogT, WithdrawalCredentialsT,
	]{
		dc:                      dc,
		dispatcher:              dispatcher,
		ds:                      ds,
		eth1FollowDistance:      eth1FollowDistance,
		failedBlocks:            make(map[math.Slot]struct{}),
		subFinalizedBlockEvents: make(chan async.Event[BeaconBlockT]),
		logger:                  logger,
		metrics:                 newMetrics(telemetrySink),
	}
}

// Start subscribes the Deposit service to BeaconBlockFinalized events and
// begins the main event loop to handle them accordingly.
func (s *Service[
	_, _, _, _, _, _,
]) Start(ctx context.Context) error {
	if err := s.dispatcher.Subscribe(
		async.BeaconBlockFinalized, s.subFinalizedBlockEvents,
	); err != nil {
		s.logger.Error("failed to subscribe to event", "event",
			async.BeaconBlockFinalized, "err", err)
		return err
	}

	// Listen for finalized block events and fetch deposits for the block.
	go s.eventLoop(ctx)

	// Catchup deposits for failed blocks.
	go s.depositCatchupFetcher(ctx)
	return nil
}

// eventLoop starts the main event loop to listen and handle
// BeaconBlockFinalized events.
func (s *Service[
	_, _, _, _, _, _,
]) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.subFinalizedBlockEvents:
			s.depositFetcher(ctx, event)
		}
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _,
]) Name() string {
	return "deposit-handler"
}
