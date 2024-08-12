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

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// LocalBuilderInput is an input for the dep inject framework.
type LocalBuilderInput[
	AttributesFactoryT AttributesFactory[BeaconStateT, PayloadAttributesT],
	BeaconStateT any,
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	LoggerT log.AdvancedLogger[any, LoggerT],
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	WithdrawalsT Withdrawals,
] struct {
	depinject.In
	AttributesFactory AttributesFactoryT
	Cfg               *config.Config
	ChainSpec         common.ChainSpec
	ExecutionEngine   ExecutionEngineT
	Logger            LoggerT
}

// ProvideLocalBuilder provides a local payload builder for the
// depinject framework.
func ProvideLocalBuilder[
	AttributesFactoryT AttributesFactory[BeaconStateT, PayloadAttributesT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, ValidatorsT, WithdrawalsT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	Eth1DataT any,
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkT any,
	KVStoreT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalsT],
	PayloadIDT ~[8]byte,
	ValidatorT any,
	ValidatorsT any,
	WithdrawalsT Withdrawals,
](
	in LocalBuilderInput[
		AttributesFactoryT,
		BeaconStateT,
		ExecutionEngineT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		LoggerT,
		PayloadAttributesT,
		PayloadIDT,
		WithdrawalsT,
	],
) *payloadbuilder.PayloadBuilder[
	AttributesFactoryT,
	BeaconStateT,
	ExecutionEngineT,
	ExecutionPayloadT,
	ExecutionPayloadHeaderT,
	LoggerT,
	PayloadAttributesT,
	PayloadIDT,
	WithdrawalsT,
] {
	return payloadbuilder.New[
		AttributesFactoryT,
		BeaconStateT,
		ExecutionEngineT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		LoggerT,
		PayloadAttributesT,
		PayloadIDT,
		WithdrawalsT,
	](
		&in.Cfg.PayloadBuilder,
		in.ChainSpec,
		in.Logger.With("service", "payload-builder"),
		in.ExecutionEngine,
		cache.NewPayloadIDCache[
			PayloadIDT,
			[32]byte, math.Slot,
		](),
		in.AttributesFactory,
	)
}
