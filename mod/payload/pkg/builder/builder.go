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

package builder

import (
	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/execution"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/payload"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/state-transition/state"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// PayloadBuilder is used to build payloads on the
// execution client.
type PayloadBuilder[
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, WithdrawalT,
	],
	BeaconBlockHeaderT any,
	BlobsBundleT engineprimitives.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	Eth1DataT any,
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT any,
	KVStoreT any,
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	ValidatorT any,
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
] struct {
	// cfg holds the configuration settings for the PayloadBuilder.
	cfg *Config
	// chainSpec holds the chain specifications for the PayloadBuilder.
	chainSpec common.ChainSpec
	// logger is used for logging within the PayloadBuilder.
	logger log.Logger[any]
	// ee is the execution engine.
	ee execution.Engine[
		BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
		ExecutionPayloadHeaderT, PayloadAttributesT, PayloadIDT, WithdrawalT,
	]
	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc *cache.PayloadIDCache[
		PayloadIDT, [32]byte, math.Slot,
	]
	// attributesFactory is used to create attributes for the
	attributesFactory payload.AttributesFactory[
		BeaconStateT, PayloadAttributesT, WithdrawalT,
	]
}

// New creates a new service.
func New[
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, WithdrawalT,
	],
	BeaconBlockHeaderT any,
	BlobsBundleT engineprimitives.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	Eth1DataT any,
	ExecutionPayloadT types.ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadEnvelopeT engineprimitives.ExecutionPayloadEnvelope[
		ExecutionPayloadEnvelopeT, BlobsBundleT, ExecutionPayloadT,
	],
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT any,
	KVStoreT any,
	PayloadAttributesT engineprimitives.PayloadAttributes[
		PayloadAttributesT, WithdrawalT,
	],
	PayloadIDT ~[8]byte,
	ValidatorT any,
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
](
	cfg *Config,
	chainSpec common.ChainSpec,
	logger log.Logger[any],
	ee execution.Engine[
		BlobsBundleT, ExecutionPayloadT, ExecutionPayloadEnvelopeT,
		ExecutionPayloadHeaderT, PayloadAttributesT, PayloadIDT, WithdrawalT,
	],
	pc *cache.PayloadIDCache[
		PayloadIDT, [32]byte, math.Slot,
	],
	attributesFactory payload.AttributesFactory[
		BeaconStateT, PayloadAttributesT, WithdrawalT,
	],
) *PayloadBuilder[
	BeaconStateT, BeaconBlockHeaderT, BlobsBundleT, Eth1DataT,
	ExecutionPayloadT, ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
	ForkT, KVStoreT, PayloadAttributesT, PayloadIDT, ValidatorT, WithdrawalT,
] {
	return &PayloadBuilder[
		BeaconStateT, BeaconBlockHeaderT, BlobsBundleT, Eth1DataT,
		ExecutionPayloadT, ExecutionPayloadEnvelopeT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT, PayloadAttributesT,
		PayloadIDT, ValidatorT, WithdrawalT,
	]{
		cfg:               cfg,
		chainSpec:         chainSpec,
		logger:            logger,
		ee:                ee,
		pc:                pc,
		attributesFactory: attributesFactory,
	}
}

// Enabled returns true if the payload builder is enabled.
func (pb *PayloadBuilder[
	BeaconStateT, BeaconBlockHeaderT, BlobsBundleT, Eth1DataT,
	ExecutionPayloadT, ExecutionPayloadEnvelopeT, ExecutionPayloadHeaderT,
	ForkT, KVStoreT, PayloadAttributesT, PayloadIDT, ValidatorT, WithdrawalT,
]) Enabled() bool {
	return pb.cfg.Enabled
}
