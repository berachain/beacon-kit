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

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/async/event"
	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/da/blob"
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/node-core/storage"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/payload"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/state-transition/core"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/state-transition/state"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/storage/deposit"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/telemetry"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is responsible for building beacon blocks.
type Service[
	AvailabilityStoreT any,
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, WithdrawalT,
	],
	BlobsBundleT engineprimitives.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	BlobSidecarsT any,
	DepositT types.Deposit[DepositT, ForkDataT, WithdrawalCredentialsT],
	DepositStoreT deposit.Store[DepositT],
	Eth1DataT types.Eth1Data[Eth1DataT],
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	//nolint:lll // annoying formatter
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT types.Fork[ForkT],
	ForkDataT types.ForkData[ForkDataT],
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	KVStoreT any,
	ValidatorT types.Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT types.WithdrawalCredentials[WithdrawalCredentialsT],
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
	blobFactory blob.SidecarFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobsBundleT, BlobSidecarsT,
	]
	// bsb is the beacon state backend.
	bsb storage.Backend[
		AvailabilityStoreT, BeaconStateT, DepositStoreT,
	]
	// stateProcessor is responsible for processing the state.
	stateProcessor core.StateProcessor[
		BeaconBlockT, BeaconStateT, BlobSidecarsT,
		*transition.Context, DepositT, ExecutionPayloadHeaderT,
	]
	// localPayloadBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks are done by submitting forkchoice updates through.
	// The local Builder.
	localPayloadBuilder payload.Builder[
		BeaconStateT, BlobsBundleT, ExecutionPayloadT,
		ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalT,
	]
	// remotePayloadBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remotePayloadBuilders []payload.Builder[
		BeaconStateT, BlobsBundleT, ExecutionPayloadT,
		ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalT,
	]
	// metrics is a metrics collector.
	metrics *validatorMetrics
	// blkBroker is a publisher for blocks.
	blkBroker event.Feed[*asynctypes.Event[BeaconBlockT]]
	// sidecarBroker is a publisher for sidecars.
	sidecarBroker event.Feed[*asynctypes.Event[BlobSidecarsT]]
	// newSlotSub is a feed for slots.
	newSlotSub chan *asynctypes.Event[math.Slot]
}

// NewService creates a new validator service.
func NewService[
	AvailabilityStoreT any,
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, WithdrawalT,
	],
	BlobsBundleT engineprimitives.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	BlobSidecarsT any,
	DepositT types.Deposit[DepositT, ForkDataT, WithdrawalCredentialsT],
	DepositStoreT deposit.Store[DepositT],
	Eth1DataT types.Eth1Data[Eth1DataT],
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	//nolint:lll // annoying formatter
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT types.Fork[ForkT],
	ForkDataT types.ForkData[ForkDataT],
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	KVStoreT any,
	ValidatorT types.Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT types.WithdrawalCredentials[WithdrawalCredentialsT],
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	bsb storage.Backend[
		AvailabilityStoreT, BeaconStateT, DepositStoreT,
	],
	stateProcessor core.StateProcessor[
		BeaconBlockT, BeaconStateT, BlobSidecarsT,
		*transition.Context, DepositT, ExecutionPayloadHeaderT,
	],
	signer crypto.BLSSigner,
	blobFactory blob.SidecarFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobsBundleT, BlobSidecarsT,
	],
	localPayloadBuilder payload.Builder[
		BeaconStateT, BlobsBundleT, ExecutionPayloadT,
		ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalT,
	],
	remotePayloadBuilders []payload.Builder[
		BeaconStateT, BlobsBundleT, ExecutionPayloadT,
		ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
		PayloadAttributesT, PayloadIDT, WithdrawalT,
	],
	ts telemetry.Sink,
	blkBroker event.Feed[*asynctypes.Event[BeaconBlockT]],
	sidecarBroker event.Feed[*asynctypes.Event[BlobSidecarsT]],
	newSlotSub chan *asynctypes.Event[math.Slot],
) *Service[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobsBundleT, BlobSidecarsT, DepositT, DepositStoreT,
	Eth1DataT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
	ExecutionPayloadHeaderT, ForkT, ForkDataT, PayloadAttributesT,
	PayloadIDT, KVStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
] {
	return &Service[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobsBundleT, BlobSidecarsT, DepositT, DepositStoreT,
		Eth1DataT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
		ExecutionPayloadHeaderT, ForkT, ForkDataT, PayloadAttributesT,
		PayloadIDT, KVStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]{
		cfg:                   cfg,
		logger:                logger,
		bsb:                   bsb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
		metrics:               newValidatorMetrics(ts),
		blkBroker:             blkBroker,
		sidecarBroker:         sidecarBroker,
		newSlotSub:            newSlotSub,
	}
}

// Name returns the name of the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) Name() string {
	return "validator"
}

// Start starts the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) Start(
	ctx context.Context,
) error {
	go s.start(ctx)
	return nil
}

// start starts the service.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) start(
	ctx context.Context,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-s.newSlotSub:
			if req.Type() == events.NewSlot {
				s.handleNewSlot(req)
			}
		}
	}
}

// handleBlockRequest handles a block request.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) handleNewSlot(msg *asynctypes.Event[math.Slot]) {
	blk, sidecars, err := s.buildBlockAndSidecars(
		msg.Context(), msg.Data(),
	)
	if err != nil {
		s.logger.Error("failed to build block", "err", err)
	}

	// Publish our built block to the broker.
	if blkErr := s.blkBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			msg.Context(), events.BeaconBlockBuilt, blk, err,
		)); blkErr != nil {
		// Propagate the error from buildBlockAndSidecars
		s.logger.Error("failed to publish block", "err", err)
	}

	// Publish our built blobs to the broker.
	if sidecarsErr := s.sidecarBroker.Publish(
		msg.Context(),
		asynctypes.NewEvent(
			// Propagate the error from buildBlockAndSidecars
			msg.Context(), events.BlobSidecarsBuilt, sidecars, err,
		),
	); sidecarsErr != nil {
		s.logger.Error("failed to publish sidecars", "err", err)
	}
}
