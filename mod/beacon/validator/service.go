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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Service is responsible for building beacon blocks.
type Service[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
	],
	BeaconStateT BeaconState[
		*types.BeaconBlockHeader,
		BeaconStateT,
		*types.ExecutionPayloadHeader,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore[*types.Deposit],
	ForkDataT interface {
		New(
			primitives.Version,
			primitives.Root,
		) ForkDataT
		ComputeRandaoSigningRoot(
			primitives.DomainType,
			math.Epoch,
		) (primitives.Root, error)
	},
] struct {
	// cfg is the validator config.
	cfg *Config
	// logger is a logger.
	logger log.Logger[any]
	// chainSpec is the chain spec.
	chainSpec primitives.ChainSpec
	// signer is used to retrieve the public key of this node.
	signer crypto.BLSSigner
	// blobFactory is used to create blob sidecars for blocks.
	blobFactory BlobFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
	]
	// bsb is the beacon state backend.
	bsb StorageBackend[BeaconStateT, *types.Deposit, DepositStoreT]
	// blobProcessor is used to process blobs.
	blobProcessor BlobProcessor[BlobSidecarsT]
	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
	]
	// localPayloadBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks are done by submitting forkchoice updates through.
	// The local Builder.
	localPayloadBuilder PayloadBuilder[BeaconStateT, *types.ExecutionPayload]
	// remotePayloadBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remotePayloadBuilders []PayloadBuilder[BeaconStateT, *types.ExecutionPayload]
	// metrics is a metrics collector.
	metrics *validatorMetrics
}

// NewService creates a new validator service.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload],
	BeaconStateT BeaconState[
		*types.BeaconBlockHeader,
		BeaconStateT,
		*types.ExecutionPayloadHeader,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore[*types.Deposit],
	ForkDataT interface {
		New(
			primitives.Version,
			primitives.Root,
		) ForkDataT
		ComputeRandaoSigningRoot(
			primitives.DomainType,
			math.Epoch,
		) (primitives.Root, error)
	},
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	bsb StorageBackend[BeaconStateT, *types.Deposit, DepositStoreT],
	blobProcessor BlobProcessor[BlobSidecarsT],
	stateProcessor StateProcessor[BeaconBlockT, BeaconStateT, *transition.Context],
	signer crypto.BLSSigner,
	blobFactory BlobFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
	],
	localPayloadBuilder PayloadBuilder[BeaconStateT, *types.ExecutionPayload],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT, *types.ExecutionPayload],
	ts TelemetrySink,
) *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
] {
	return &Service[
		BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositStoreT, ForkDataT,
	]{
		cfg:                   cfg,
		logger:                logger,
		blobProcessor:         blobProcessor,
		bsb:                   bsb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
		metrics:               newValidatorMetrics(ts),
	}
}

// Name returns the name of the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) Name() string {
	return "validator"
}

// Start starts the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) Start(
	context.Context,
) error {
	s.logger.Info(
		"starting validator service 🛜 ",
		"optimistic_payload_builds", s.cfg.EnableOptimisticPayloadBuilds,
	)
	return nil
}
