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

package processor

import (
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Processor is responsible for building beacon blocks.
type Processor[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	BlobSidecarsT,
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
	cfg *validator.Config
	// logger is a logger.
	logger log.Logger[any]
	// chainSpec is the chain spec.
	chainSpec common.ChainSpec
	// signer is used to retrieve the public key of this node.
	signer crypto.BLSSigner
	// blobFactory is used to create blob sidecars for blocks.
	blobFactory BlobFactory[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	]
	// bsb is the beacon state backend.
	bsb StorageBackend[
		BeaconStateT, DepositT, DepositStoreT, ExecutionPayloadHeaderT,
	]
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
}

// New creates a new validator service.
func New[
	AttestationDataT any,
	BeaconBlockT BeaconBlock[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		AttestationDataT, DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconStateT BeaconState[ExecutionPayloadHeaderT],
	BlobSidecarsT,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT Eth1Data[Eth1DataT],
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkDataT ForkData[ForkDataT],
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT],
](
	cfg *validator.Config,
	logger log.Logger[any],
	chainSpec common.ChainSpec,
	bsb StorageBackend[
		BeaconStateT, DepositT, DepositStoreT, ExecutionPayloadHeaderT,
	],
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
		ExecutionPayloadHeaderT,
	],
	signer crypto.BLSSigner,
	blobFactory BlobFactory[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	localPayloadBuilder PayloadBuilder[BeaconStateT, ExecutionPayloadT],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT, ExecutionPayloadT],
	ts TelemetrySink,
) *Processor[
	AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, ForkDataT, SlashingInfoT, SlotDataT,
] {
	return &Processor[
		AttestationDataT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, ForkDataT, SlashingInfoT, SlotDataT,
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
	}
}
