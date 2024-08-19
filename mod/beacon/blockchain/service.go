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

package blockchain

import (
	"context"
	"sync"

	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is the blockchain service.
type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BeaconBlockHeaderT BeaconBlockHeader,
	BeaconStateT ReadOnlyBeaconState[
		BeaconStateT, BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	],
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	PayloadAttributesT interface {
		IsNil() bool
		Version() uint32
		GetSuggestedFeeRecipient() common.ExecutionAddress
	},
] struct {
	// sb represents the backend storage for beacon states and associated
	// sidecars.
	sb StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
	]
	// logger is used for logging messages in the service.
	logger log.Logger[any]
	// cs holds the chain specifications.
	cs common.ChainSpec
	// ee is the execution engine responsible for processing execution payloads.
	ee ExecutionEngine[PayloadAttributesT]
	// lb is a local builder for constructing new beacon states.
	lb LocalBuilder[BeaconStateT]
	// sp is the state processor for beacon blocks and states.
	sp StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	]
	// metrics is the metrics for the service.
	metrics *chainMetrics
	// genesisBroker is the event feed for genesis data.
	genesisBroker *broker.Broker[*asynctypes.Event[GenesisT]]
	// blkBroker is the event feed for new blocks.
	blkBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]]
	// validatorUpdateBroker is the event feed for validator updates.
	//nolint:lll // annoying formatter.
	validatorUpdateBroker *broker.Broker[*asynctypes.Event[transition.ValidatorUpdates]]
	// optimisticPayloadBuilds is a flag used when the optimistic payload
	// builder is enabled.
	optimisticPayloadBuilds bool
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once
}

// NewService creates a new validator service.
func NewService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BeaconBlockHeaderT BeaconBlockHeader,
	BeaconStateT ReadOnlyBeaconState[
		BeaconStateT, BeaconBlockHeaderT,
		ExecutionPayloadHeaderT,
	],
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	PayloadAttributesT interface {
		IsNil() bool
		Version() uint32
		GetSuggestedFeeRecipient() common.ExecutionAddress
	},
](
	sb StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
	],
	logger log.Logger[any],
	cs common.ChainSpec,
	ee ExecutionEngine[PayloadAttributesT],
	lb LocalBuilder[BeaconStateT],
	sp StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	],
	ts TelemetrySink,
	genesisBroker *broker.Broker[*asynctypes.Event[GenesisT]],
	blkBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]],
	//nolint:lll // annoying formatter.
	validatorUpdateBroker *broker.Broker[*asynctypes.Event[transition.ValidatorUpdates]],
	optimisticPayloadBuilds bool,
) *Service[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	GenesisT, PayloadAttributesT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		GenesisT, PayloadAttributesT,
	]{
		sb:                      sb,
		logger:                  logger,
		cs:                      cs,
		ee:                      ee,
		lb:                      lb,
		sp:                      sp,
		metrics:                 newChainMetrics(ts),
		genesisBroker:           genesisBroker,
		blkBroker:               blkBroker,
		validatorUpdateBroker:   validatorUpdateBroker,
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "blockchain"
}

func (s *Service[
	_, _, _, _, _, _, _, _, _, _,
]) Start(ctx context.Context) error {
	subBlkCh, err := s.blkBroker.Subscribe()
	if err != nil {
		return err
	}
	subGenCh, err := s.genesisBroker.Subscribe()
	if err != nil {
		return err
	}
	go s.start(ctx, subBlkCh, subGenCh)
	return nil
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, GenesisT, _,
]) start(
	ctx context.Context,
	subBlkCh chan *asynctypes.Event[BeaconBlockT],
	subGenCh chan *asynctypes.Event[GenesisT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-subBlkCh:
			switch msg.Type() {
			case events.BeaconBlockReceived:
				s.handleBeaconBlockReceived(msg)
			case events.BeaconBlockFinalizedRequest:
				s.handleBeaconBlockFinalization(msg)
			}
		case msg := <-subGenCh:
			if msg.Type() == events.GenesisDataProcessRequest {
				s.handleProcessGenesisDataRequest(msg)
			}
		}
	}
}

func (s *Service[
	_, _, _, _, _, _, _, _, GenesisT, _,
]) handleProcessGenesisDataRequest(msg *asynctypes.Event[GenesisT]) {
	if msg.Error() != nil {
		s.logger.Error("Error processing genesis data", "error", msg.Error())
		return
	}

	// Process the genesis data.
	valUpdates, err := s.ProcessGenesisData(msg.Context(), msg.Data())
	if err != nil {
		s.logger.Error("Failed to process genesis data", "error", err)
	}

	// Publish the validator set updated event.
	if err = s.validatorUpdateBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.ValidatorSetUpdated,
			valUpdates,
			err,
		),
	); err != nil {
		s.logger.Error(
			"Failed to publish validator set updated event",
			"error",
			err,
		)
	}
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockReceived(
	msg *asynctypes.Event[BeaconBlockT],
) {
	// If the block is nil, exit early.
	if msg.Error() != nil {
		s.logger.Error("Error processing beacon block", "error", msg.Error())
		return
	}

	// Publish the verified block event.
	if err := s.blkBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.BeaconBlockVerified,
			msg.Data(),
			s.VerifyIncomingBlock(msg.Context(), msg.Data()),
		),
	); err != nil {
		s.logger.Error("Failed to publish verified block", "error", err)
	}
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockFinalization(
	msg *asynctypes.Event[BeaconBlockT],
) {
	// If there's an error in the event, log it and return
	if msg.Error() != nil {
		s.logger.Error("Error verifying beacon block", "error", msg.Error())
		return
	}

	// Process the verified block
	valUpdates, err := s.ProcessBeaconBlock(msg.Context(), msg.Data())
	if err != nil {
		s.logger.Error("Failed to process verified beacon block", "error", err)
	}

	// Publish the validator set updated event
	if err = s.validatorUpdateBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(),
			events.ValidatorSetUpdated,
			valUpdates,
			err,
		)); err != nil {
		s.logger.Error(
			"Failed to publish validator set updated event",
			"error",
			err,
		)
	}
}
