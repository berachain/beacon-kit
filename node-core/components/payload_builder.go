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
	"github.com/berachain/beacon-kit/config"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// LocalBuilderInput is an input for the dep inject framework.
type LocalBuilderInput[
	BeaconStateT any,
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	LoggerT log.AdvancedLogger[LoggerT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
] struct {
	depinject.In
	AttributesFactory AttributesFactory[
		BeaconStateT, *engineprimitives.PayloadAttributes[WithdrawalT],
	]
	Cfg             *config.Config
	ChainSpec       common.ChainSpec
	ExecutionEngine *engine.Engine[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
		PayloadID,
		WithdrawalsT,
	]
	Logger LoggerT
}

// ProvideLocalBuilder provides a local payload builder for the
// depinject framework.
func ProvideLocalBuilder[
	BeaconStateT BeaconState[
		BeaconStateT, BeaconStateMarshallableT,
		ExecutionPayloadHeaderT, *Fork, KVStoreT, *Validator,
		Validators, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
](
	in LocalBuilderInput[
		BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT, LoggerT,
		WithdrawalT, WithdrawalsT,
	],
) *payloadbuilder.PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	*engineprimitives.PayloadAttributes[WithdrawalT], PayloadID, WithdrawalT,
] {
	return payloadbuilder.New[
		BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		*engineprimitives.PayloadAttributes[WithdrawalT], PayloadID, WithdrawalT,
	](
		&in.Cfg.PayloadBuilder,
		in.ChainSpec,
		in.Logger.With("service", "payload-builder"),
		in.ExecutionEngine,
		cache.NewPayloadIDCache[
			PayloadID,
			[32]byte, math.Slot,
		](),
		in.AttributesFactory,
	)
}
