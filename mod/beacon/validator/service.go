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

package validator

import (
	"context"

	async "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	async1 "github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is responsible for building beacon blocks and sidecars.
type Service[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	BlobSidecarsT any,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkDataT ForkData[ForkDataT],
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT],
] struct {
	// cfg is the validator config.
	cfg *Config
	// logger is a logger.
	logger log.Logger[any]
	// chainSpec is the chain spec.
	chainSpec common.ChainSpec
	// signer is used to retrieve the public key of this node.
	signer crypto.BLSSigner
	// blobFactory is used to create blob sidecars for blocks.
	blobFactory BlobFactory[BeaconBlockT, BlobSidecarsT]
	// sb is the beacon state backend.
	sb StorageBackend[BeaconStateT, DepositStoreT]
	// dispatcher is the dispatcher.
	dispatcher async.EventDispatcher
	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		ExecutionPayloadHeaderT,
	]
	// localPayloadBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks are done by submitting forkchoice updates through.
	// The local Builder.
	localPayloadBuilder PayloadBuilder[BeaconStateT, ExecutionPayloadT]
	// remotePayloadBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remotePayloadBuilders []PayloadBuilder[BeaconStateT, ExecutionPayloadT]
	// metrics is a metrics collector.
	metrics *validatorMetrics
	// subNewSlot is a channel for new slot events.
	subNewSlot chan async1.Event[SlotDataT]
}

// NewService creates a new validator service.
func NewService[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	BlobSidecarsT any,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkDataT ForkData[ForkDataT],
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT],
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	sb StorageBackend[BeaconStateT, DepositStoreT],
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		ExecutionPayloadHeaderT,
	],
	signer crypto.BLSSigner,
	blobFactory BlobFactory[BeaconBlockT, BlobSidecarsT],
	localPayloadBuilder PayloadBuilder[BeaconStateT, ExecutionPayloadT],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT, ExecutionPayloadT],
	ts TelemetrySink,
	dispatcher async.EventDispatcher,
) *Service[
	AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, ForkDataT, SlashingInfoT, SlotDataT,
] {
	return &Service[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, ForkDataT, SlashingInfoT,
		SlotDataT,
	]{
		cfg:                   cfg,
		logger:                logger,
		sb:                    sb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
		metrics:               newValidatorMetrics(ts),
		dispatcher:            dispatcher,
		subNewSlot:            make(chan async1.Event[SlotDataT]),
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "validator"
}

// Start listens for NewSlot events and builds a block and sidecars for the
// requested slot data.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, SlotDataT,
]) Start(
	ctx context.Context,
) error {
	// subscribe to new slot events
	err := s.dispatcher.Subscribe(async1.NewSlot, s.subNewSlot)
	if err != nil {
		return err
	}
	// set the handler for new slot events
	go s.eventLoop(ctx)
	return nil
}

func (s *Service[_, _, _, _, _, _, _, _, _, _, _, _, SlotDataT]) eventLoop(
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.subNewSlot:
			s.handleNewSlot(event)
		}
	}
}

// handleNewSlot builds a block and sidecars for the requested
// slot data and emits an event containing the built block and sidecars.
func (s *Service[
	_, BeaconBlockT, _, _, BlobSidecarsT, _, _, _, _, _, _, _, SlotDataT,
]) handleNewSlot(req async1.Event[SlotDataT]) {
	var (
		blk      BeaconBlockT
		sidecars BlobSidecarsT
		err      error
	)
	// build the block and sidecars for the requested slot data
	blk, sidecars, err = s.buildBlockAndSidecars(
		req.Context(), req.Data(),
	)
	if err != nil {
		s.logger.Error("failed to build block", "err", err)
	}

	// emit a built block event with the built block and the error
	if bbErr := s.dispatcher.Publish(
		async1.NewEvent(req.Context(), async1.BuiltBeaconBlock, blk, err),
	); bbErr != nil {
		s.logger.Error("failed to dispatch built block", "err", err)
	}

	// emit a built sidecars event with the built sidecars and the error
	if scErr := s.dispatcher.Publish(
		async1.NewEvent(req.Context(), async1.BuiltSidecars, sidecars, err),
	); scErr != nil {
		s.logger.Error("failed to dispatch built sidecars", "err", err)
	}
}
