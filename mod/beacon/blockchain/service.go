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

	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is the blockchain service.
type Service[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
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
	WithdrawalT any,
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
	// chainSpec holds the chain specifications.
	chainSpec common.ChainSpec
	// dispatcher is the dispatcher for the service.
	dispatcher *dispatcher.Dispatcher
	// executionEngine is the execution engine responsible for processing execution payloads.
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
}

// NewService creates a new validator service.
func NewService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
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
	WithdrawalT any,
](
	sb StorageBackend[
		AvailabilityStoreT,
		BeaconBlockBodyT,
		BeaconStateT,
	],
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	dispatcher *dispatcher.Dispatcher,
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
	GenesisT, PayloadAttributesT, WithdrawalT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		GenesisT, PayloadAttributesT, WithdrawalT,
	]{
		sb:                      sb,
		logger:                  logger,
		chainSpec:               chainSpec,
		dispatcher:              dispatcher,
		executionEngine:         executionEngine,
		localBuilder:            localBuilder,
		stateProcessor:          stateProcessor,
		metrics:                 newChainMetrics(ts),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "blockchain"
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, GenesisT, _, _,
]) Start(ctx context.Context) error {
	finalizeBeaconBlockRequests := make(chan *asynctypes.Message[BeaconBlockT])
	s.dispatcher.RegisterMsgReceiver(messages.FinalizeBeaconBlock, finalizeBeaconBlockRequests)

	verifyBeaconBlockRequests := make(chan *asynctypes.Message[BeaconBlockT])
	s.dispatcher.RegisterMsgReceiver(messages.VerifyBeaconBlock, verifyBeaconBlockRequests)

	processGenDataRequests := make(chan *asynctypes.Message[GenesisT])
	s.dispatcher.RegisterMsgReceiver(messages.ProcessGenesisData, processGenDataRequests)

	go s.start(ctx, finalizeBeaconBlockRequests, verifyBeaconBlockRequests, processGenDataRequests)
	return nil
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, GenesisT, _, _,
]) start(
	ctx context.Context,
	finalizeBeaconBlockRequests chan *asynctypes.Message[BeaconBlockT],
	verifyBeaconBlockRequests chan *asynctypes.Message[BeaconBlockT],
	processGenDataRequests chan *asynctypes.Message[GenesisT],
) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-finalizeBeaconBlockRequests:
			s.handleFinalizeBeaconBlockRequest(msg)
		case msg := <-verifyBeaconBlockRequests:
			s.handleVerifyBeaconBlockRequest(msg)
		case msg := <-processGenDataRequests:
			s.handleProcessGenesisDataRequest(msg)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                              Message Handlers                              */
/* -------------------------------------------------------------------------- */

// handleProcessGenesisDataRequest processes the given genesis data and
// dispatches a response.
func (s *Service[
	_, _, _, _, _, _, _, _, GenesisT, _, _,
]) handleProcessGenesisDataRequest(msg *asynctypes.Message[GenesisT]) {
	if msg.Error() != nil {
		s.logger.Error("Error processing genesis data", "error", msg.Error())
		return
	}

	// Process the genesis data.
	valUpdates, err := s.ProcessGenesisData(msg.Context(), msg.Data())
	if err != nil {
		s.logger.Error("Failed to process genesis data", "error", err)
	}

	// dispatch a response containing the validator updates
	if err := s.dispatcher.Respond(
		asynctypes.NewMessage(
			msg.Context(),
			messages.ProcessGenesisData,
			valUpdates,
			nil,
		),
	); err != nil {
		s.logger.Error("Failed to dispatch response in handleProcessGenesisDataRequest", "error", err)
	}
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _, _,
]) handleVerifyBeaconBlockRequest(
	msg *asynctypes.Message[BeaconBlockT],
) {
	// If the block is nil, exit early.
	if msg.Error() != nil {
		s.logger.Error("Error processing beacon block", "error", msg.Error())
		return
	}

	// dispatch a response with the error result from VerifyIncomingBlock
	s.dispatcher.Respond(
		asynctypes.NewMessage(
			msg.Context(),
			messages.VerifyBeaconBlock,
			msg.Data(),
			s.VerifyIncomingBlock(msg.Context(), msg.Data()),
		),
	)
}

func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _, _,
]) handleFinalizeBeaconBlockRequest(
	msg *asynctypes.Message[BeaconBlockT],
) {
	// If there's an error in the event, log it and return
	if msg.Error() != nil {
		s.logger.Error("Error verifying beacon block", "error", msg.Error())
		return
	}

	// process the verified block and get the validator updates
	valUpdates, err := s.ProcessBeaconBlock(msg.Context(), msg.Data())
	if err != nil {
		s.logger.Error("Failed to process verified beacon block", "error", err)
	}

	// dispatch a response with the validator updates
	s.dispatcher.Respond(
		asynctypes.NewMessage(
			msg.Context(),
			messages.FinalizeBeaconBlock,
			valUpdates,
			err,
		),
	)
}
