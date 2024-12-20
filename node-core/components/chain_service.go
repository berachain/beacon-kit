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
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/da/da"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput[
	StorageBackendT any,
	LoggerT any,
	BlockStoreT BlockStore,
	DepositStoreT any,
	DepositContractT any,
	AvailabilityStoreT any,
	ConsensusSidecarsT any,
] struct {
	depinject.In

	AppOpts         config.AppOptions
	ChainSpec       chain.ChainSpec
	Cfg             *config.Config
	EngineClient    *client.EngineClient
	ExecutionEngine *engine.Engine
	LocalBuilder    LocalBuilder
	Logger          LoggerT
	Signer          crypto.BLSSigner
	StateProcessor  StateProcessor[*Context]
	StorageBackend  StorageBackendT
	BlobProcessor   BlobProcessor[
		AvailabilityStoreT, ConsensusSidecarsT,
	]
	TelemetrySink         *metrics.TelemetrySink
	BeaconDepositContract DepositContractT
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService[
	AvailabilityStoreT AvailabilityStore,
	ConsensusBlockT ConsensusBlock,
	ConsensusSidecarsT da.ConsensusSidecars,
	DepositStoreT DepositStore,
	DepositContractT deposit.Contract,
	GenesisT Genesis,
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
	StorageBackendT StorageBackend[AvailabilityStoreT, BlockStoreT, DepositStoreT],
	BlockStoreT BlockStore,
](
	in ChainServiceInput[
		StorageBackendT, LoggerT, BlockStoreT, DepositStoreT,
		DepositContractT, AvailabilityStoreT, ConsensusSidecarsT,
	],
) *blockchain.Service[
	AvailabilityStoreT, DepositStoreT, ConsensusBlockT,
	BlockStoreT, GenesisT, ConsensusSidecarsT,
] {
	return blockchain.NewService[
		AvailabilityStoreT,
		DepositStoreT,
		ConsensusBlockT,
		BlockStoreT,
		GenesisT,
		ConsensusSidecarsT,
	](
		cast.ToString(in.AppOpts.Get(flags.FlagHome)),
		in.StorageBackend,
		in.BlobProcessor,
		in.BeaconDepositContract,
		math.U64(in.ChainSpec.Eth1FollowDistance()),
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
