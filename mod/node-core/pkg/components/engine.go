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
	"math/big"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/config"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
)

// EngineClientInputs is the input for the EngineClient.
type EngineClientInputs[LoggerT any] struct {
	depinject.In
	ChainSpec common.ChainSpec
	Config    *config.Config
	// TODO: this feels like a hood way to handle it.
	JWTSecret     *jwt.Secret `optional:"true"`
	Logger        LoggerT
	TelemetrySink *metrics.TelemetrySink
}

// ProvideEngineClient creates a new EngineClient.
func ProvideEngineClient[
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	LoggerT log.AdvancedLogger[any, LoggerT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
](
	in EngineClientInputs[LoggerT],
) *client.EngineClient[
	ExecutionPayloadT,
	*engineprimitives.PayloadAttributes[WithdrawalT],
] {
	return client.New[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	](
		in.Config.GetEngine(),
		in.Logger.With("service", "engine.client"),
		in.JWTSecret,
		in.TelemetrySink,
		new(big.Int).SetUint64(in.ChainSpec.DepositEth1ChainID()),
	)
}

// EngineClientInputs is the input for the EngineClient.
type ExecutionEngineInputs[
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	LoggerT any,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
] struct {
	depinject.In
	EngineClient *client.EngineClient[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	]
	Logger        LoggerT
	TelemetrySink *metrics.TelemetrySink
}

// ProvideExecutionEngine provides the execution engine to the depinject
// framework.
func ProvideExecutionEngine[
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	LoggerT log.AdvancedLogger[any, LoggerT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
](
	in ExecutionEngineInputs[
		ExecutionPayloadT, ExecutionPayloadHeaderT, LoggerT, WithdrawalT,
		WithdrawalsT,
	],
) *engine.Engine[
	ExecutionPayloadT,
	*engineprimitives.PayloadAttributes[WithdrawalT],
	PayloadID,
	WithdrawalsT,
] {
	return engine.New[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
		PayloadID,
		WithdrawalsT,
	](
		in.EngineClient,
		in.Logger.With("service", "execution-engine"),
		in.TelemetrySink,
	)
}
