// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
)

// EngineClientInputs is the input for the EngineClient.
type EngineClientInputs struct {
	depinject.In
	// ChainSpec is the chain spec.
	ChainSpec primitives.ChainSpec
	// Config is the BeaconKit configuration.
	Config *config.Config
	// Logger is the logger.
	Logger log.Logger
	// TelemetrySink is the telemetry sink.
	TelemetrySink *metrics.TelemetrySink
	// JWTSecret is the jwt secret. It is optional, since
	// it is not required when connecting to the execution client
	// over IPC.
	JWTSecret *jwt.Secret `optional:"true"`
}

// ProvideEngineClient creates a new EngineClient.
func ProvideEngineClient(
	in EngineClientInputs,
) *engineclient.EngineClient[*types.ExecutionPayload] {
	return engineclient.New[*types.ExecutionPayload](
		&in.Config.Engine,
		in.Logger.With("service", "engine.client"),
		in.JWTSecret,
		in.TelemetrySink,
		new(big.Int).SetUint64(in.ChainSpec.DepositEth1ChainID()),
	)
}
