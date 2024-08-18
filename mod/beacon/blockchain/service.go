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

	async "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	async1 "github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
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
		BeaconBlockBodyT,
		BeaconStateT,
	]
	// logger is used for logging messages in the service.
	logger log.Logger[any]
	// cs holds the chain specifications.
	cs common.ChainSpec
	// dispatcher is the dispatcher for the service.
	dispatcher async.Dispatcher
	// ee is the execution engine responsible for processing
	// execution payloads.
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
	// optimisticPayloadBuilds is a flag used when the optimistic payload
	// builder is enabled.
	optimisticPayloadBuilds bool
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once

	// subFinalBlkReceived is a channel for receiving finalize beacon block
	// requests.
	subFinalBlkReceived chan async1.Event[BeaconBlockT]
	// subBlockReceived is a channel for receiving verify beacon block requests.
	subBlockReceived chan async1.Event[BeaconBlockT]
	// subGenDataReceived is a subscription for receiving genesis data
	// received events.
	subGenDataReceived chan async1.Event[GenesisT]
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
		BeaconBlockBodyT,
		BeaconStateT,
	],
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	dispatcher async.Dispatcher,
	executionEngine ExecutionEngine[PayloadAttributesT],
	localBuilder LocalBuilder[BeaconStateT],
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	],
	ts TelemetrySink,
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
		cs:                      chainSpec,
		dispatcher:              dispatcher,
		ee:                      executionEngine,
		lb:                      localBuilder,
		sp:                      stateProcessor,
		metrics:                 newChainMetrics(ts),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
		subFinalBlkReceived:     make(chan async1.Event[BeaconBlockT]),
		subBlockReceived:        make(chan async1.Event[BeaconBlockT]),
		subGenDataReceived:      make(chan async1.Event[GenesisT]),
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "blockchain"
}

// Start sets up the service to listen for FinalizeBeaconBlock,
// VerifyBeaconBlock, and ProcessGenesisData requests, and handles them
// accordingly.
func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, GenesisT, _,
]) Start(ctx context.Context) error {
	if err := s.dispatcher.Subscribe(
		async1.GenesisDataReceived, s.subGenDataReceived,
	); err != nil {
		return err
	}

	if err := s.dispatcher.Subscribe(
		async1.BeaconBlockReceived, s.subBlockReceived,
	); err != nil {
		return err
	}

	if err := s.dispatcher.Subscribe(
		async1.FinalBeaconBlockReceived, s.subFinalBlkReceived,
	); err != nil {
		return err
	}

	// start a goroutine to listen for requests and handle accordingly
	go s.eventLoop(ctx)
	return nil
}

// eventLoop listens for events and handles them accordingly.
func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, GenesisT, _,
]) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.subGenDataReceived:
			s.handleGenDataReceived(event)
		case event := <-s.subBlockReceived:
			s.handleBeaconBlockReceived(event)
		case event := <-s.subFinalBlkReceived:
			s.handleBeaconBlockFinalization(event)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                              Message Handlers                              */
/* -------------------------------------------------------------------------- */

func (s *Service[
	_, _, _, _, _, _, _, _, GenesisT, _,
]) handleGenDataReceived(msg async1.Event[GenesisT]) {
	var (
		valUpdates transition.ValidatorUpdates
		genesisErr error
	)
	if msg.Error() != nil {
		s.logger.Error("Error processing genesis data", "error", msg.Error())
	}

	// Process the genesis data.
	valUpdates, genesisErr = s.ProcessGenesisData(msg.Context(), msg.Data())
	if genesisErr != nil {
		s.logger.Error("Failed to process genesis data", "error", genesisErr)
	}

	// Emit the event containing the validator updates.
	if err := s.dispatcher.Publish(
		async1.NewEvent(
			msg.Context(),
			async1.GenesisDataProcessed,
			valUpdates,
			genesisErr,
		),
	); err != nil {
		s.logger.Error(
			"Failed to emit event in process genesis data",
			"error", err,
		)
		panic(err)
	}
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockReceived(
	msg async1.Event[BeaconBlockT],
) {
	// If the block is nil, exit early.
	if msg.Error() != nil {
		s.logger.Error("Error processing beacon block", "error", msg.Error())
		return
	}

	// emit a BeaconBlockVerified event with the error result from \
	// VerifyIncomingBlock
	if err := s.dispatcher.Publish(
		async1.NewEvent(
			msg.Context(),
			async1.BeaconBlockVerified,
			msg.Data(),
			s.VerifyIncomingBlock(msg.Context(), msg.Data()),
		),
	); err != nil {
		s.logger.Error(
			"Failed to emit event in verify beacon block",
			"error", err,
		)
	}
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockFinalization(
	msg async1.Event[BeaconBlockT],
) {
	var (
		valUpdates  transition.ValidatorUpdates
		finalizeErr error
	)
	// If there's an error in the event, log it and return
	if msg.Error() != nil {
		s.logger.Error("Error verifying beacon block", "error", msg.Error())
		return
	}

	// process the verified block and get the validator updates
	valUpdates, finalizeErr = s.ProcessBeaconBlock(msg.Context(), msg.Data())
	if finalizeErr != nil {
		s.logger.Error("Failed to process verified beacon block",
			"error", finalizeErr,
		)
	}

	// Emit the event containing the validator updates.
	if err := s.dispatcher.Publish(
		async1.NewEvent(
			msg.Context(),
			async1.FinalValidatorUpdatesProcessed,
			valUpdates,
			finalizeErr,
		),
	); err != nil {
		s.logger.Error(
			"Failed to emit event in finalize beacon block",
			"error", err,
		)
	}
}
