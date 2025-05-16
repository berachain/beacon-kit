// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput struct {
	depinject.In

	ChainSpec             chain.Spec
	Cfg                   *config.Config
	ExecutionEngine       *engine.Engine
	LocalBuilder          LocalBuilder
	Logger                *phuslu.Logger
	Signer                crypto.BLSSigner
	StateProcessor        StateProcessor
	StorageBackend        *storage.Backend
	BlobProcessor         BlobProcessor
	TelemetrySink         *metrics.TelemetrySink
	BeaconDepositContract deposit.Contract
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService(in ChainServiceInput) *blockchain.Service {
	return blockchain.NewService(
		in.StorageBackend,
		in.BlobProcessor,
		in.BeaconDepositContract,
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.StateProcessor,
		in.TelemetrySink,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
}
