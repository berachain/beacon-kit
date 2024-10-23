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

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is the blockchain service.
type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	ConsensusBlockT ConsensusBlock[BeaconBlockT],
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
	PayloadAttributesT PayloadAttributes,
] struct {
	// storageBackend represents the backend storage for beacon states and
	// associated sidecars.
	storageBackend StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
	]
	// logger is used for logging messages in the service.
	logger log.Logger
	// chainSpec holds the chain specifications.
	chainSpec common.ChainSpec
	// dispatcher is the dispatcher for the service.
	dispatcher asynctypes.Dispatcher
	// executionEngine is the execution engine responsible for processing
	//
	// execution payloads.
	executionEngine ExecutionEngine[PayloadAttributesT]
	// localBuilder is a local builder for constructing new beacon states.
	localBuilder LocalBuilder[BeaconStateT]
	// stateProcessor is the state processor for beacon blocks and states.
	stateProcessor StateProcessor[
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

	// subFinalBlkReceived is a channel holding FinalBeaconBlockReceived events.
	subFinalBlkReceived chan async.Event[BeaconBlockT]
	// subBlockReceived is a channel holding BeaconBlockReceived events.
	subBlockReceived chan async.Event[BeaconBlockT]
	// subGenDataReceived is a channel holding GenesisDataReceived events.
	subGenDataReceived chan async.Event[GenesisT]
}

// NewService creates a new validator service.
func NewService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	ConsensusBlockT ConsensusBlock[BeaconBlockT],
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
	PayloadAttributesT PayloadAttributes,
](
	storageBackend StorageBackend[
		AvailabilityStoreT,
		BeaconStateT,
	],
	logger log.Logger,
	chainSpec common.ChainSpec,
	dispatcher asynctypes.Dispatcher,
	executionEngine ExecutionEngine[PayloadAttributesT],
	localBuilder LocalBuilder[BeaconStateT],
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	],
	telemetrySink TelemetrySink,
	optimisticPayloadBuilds bool,
) *Service[
	AvailabilityStoreT,
	ConsensusBlockT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	GenesisT, PayloadAttributesT,
] {
	return &Service[
		AvailabilityStoreT,
		ConsensusBlockT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		GenesisT, PayloadAttributesT,
	]{
		storageBackend:          storageBackend,
		logger:                  logger,
		chainSpec:               chainSpec,
		dispatcher:              dispatcher,
		executionEngine:         executionEngine,
		localBuilder:            localBuilder,
		stateProcessor:          stateProcessor,
		metrics:                 newChainMetrics(telemetrySink),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
		subFinalBlkReceived:     make(chan async.Event[BeaconBlockT]),
		subBlockReceived:        make(chan async.Event[BeaconBlockT]),
		subGenDataReceived:      make(chan async.Event[GenesisT]),
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "blockchain"
}

// Start subscribes the Blockchain service to GenesisDataReceived,
// BeaconBlockReceived, and FinalBeaconBlockReceived events, and begins
// the main event loop to handle them accordingly.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _,
]) Start(ctx context.Context) error {
	if err := s.dispatcher.Subscribe(
		async.GenesisDataReceived, s.subGenDataReceived,
	); err != nil {
		return err
	}

	if err := s.dispatcher.Subscribe(
		async.BeaconBlockReceived, s.subBlockReceived,
	); err != nil {
		return err
	}

	if err := s.dispatcher.Subscribe(
		async.FinalBeaconBlockReceived, s.subFinalBlkReceived,
	); err != nil {
		return err
	}

	// start the main event loop to listen and handle events.
	go s.eventLoop(ctx)
	return nil
}

// eventLoop listens for events and handles them accordingly.
func (s *Service[
	_, _, BeaconBlockT, _, _, _, _, _, _, GenesisT, _,
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
/*                                Event Handlers                              */
/* -------------------------------------------------------------------------- */

// handleGenDataReceived processes the genesis data received and emits a
// GenesisDataProcessed event containing the resulting validator updates.
func (s *Service[
	_, _, _, _, _, _, _, _, _, GenesisT, _,
]) handleGenDataReceived(msg async.Event[GenesisT]) {
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
		async.NewEvent(
			msg.Context(),
			async.GenesisDataProcessed,
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

// handleBeaconBlockReceived emits a BeaconBlockVerified event with the error
// result from VerifyIncomingBlock.
func (s *Service[
	_, _, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockReceived(
	msg async.Event[BeaconBlockT],
) {
	// If the block is nil, exit early.
	if msg.Error() != nil {
		s.logger.Error("Error processing beacon block", "error", msg.Error())
		return
	}

	// emit a BeaconBlockVerified event with
	// the error result from VerifyIncomingBlock
	if err := s.dispatcher.Publish(
		async.NewEvent(
			msg.Context(),
			async.BeaconBlockVerified,
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

// handleBeaconBlockFinalization processes the finalized beacon block and emits
// a FinalValidatorUpdatesProcessed event containing the resulting validator
// updates.
func (s *Service[
	_, _, BeaconBlockT, _, _, _, _, _, _, _, _,
]) handleBeaconBlockFinalization(
	msg async.Event[BeaconBlockT],
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
		async.NewEvent(
			msg.Context(),
			async.FinalValidatorUpdatesProcessed,
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
