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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/berachain/beacon-kit/primitives/math"
)

// LocalBuilderInput is an input for the dep inject framework.
type LocalBuilderInput[
	LoggerT log.AdvancedLogger[LoggerT],
] struct {
	depinject.In
	AttributesFactory AttributesFactory
	Cfg               *config.Config
	ChainSpec         chain.ChainSpec
	ExecutionEngine   *engine.Engine
	Logger            LoggerT
}

// ProvideLocalBuilder provides a local payload builder for the
// depinject framework.
func ProvideLocalBuilder[
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
](
	in LocalBuilderInput[LoggerT],
) *payloadbuilder.PayloadBuilder {
	return payloadbuilder.New(
		&in.Cfg.PayloadBuilder,
		in.ChainSpec,
		in.Logger.With("service", "payload-builder"),
		in.ExecutionEngine,
		cache.NewPayloadIDCache[
			[32]byte, math.Slot,
		](),
		in.AttributesFactory,
	)
}
