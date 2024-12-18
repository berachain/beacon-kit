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
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/state-transition/core"
)

// StateProcessorInput is the input for the state processor for the depinject
// framework.
type StateProcessorInput[
	LoggerT log.AdvancedLogger[LoggerT],
] struct {
	depinject.In
	Logger          LoggerT
	ChainSpec       chain.ChainSpec
	ExecutionEngine *engine.Engine
	DepositStore    DepositStore
	Signer          crypto.BLSSigner
	TelemetrySink   *metrics.TelemetrySink
}

// ProvideStateProcessor provides the state processor to the depinject
// framework.
func ProvideStateProcessor[
	LoggerT log.AdvancedLogger[LoggerT],
	DepositStoreT DepositStore,
	KVStoreT BeaconStore[KVStoreT],
](
	in StateProcessorInput[LoggerT],
) *core.StateProcessor[
	*Context,
	KVStoreT,
] {
	return core.NewStateProcessor[
		*Context,
		KVStoreT,
	](
		in.Logger.With("service", "state-processor"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.DepositStore,
		in.Signer,
		crypto.GetAddressFromPubKey,
		in.TelemetrySink,
	)
}
