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
	"github.com/berachain/beacon-kit/config"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput[
	BeaconBlockT any,
	BeaconStateT any,
	DepositT any,
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	StorageBackendT any,
	LoggerT any,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	DepositStoreT any,
	DepositContractT any,
] struct {
	depinject.In

	ChainSpec    common.ChainSpec
	Cfg          *config.Config
	EngineClient *client.EngineClient[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	]
	ExecutionEngine *engine.Engine[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
		PayloadID,
		WithdrawalsT,
	]
	Dispatcher     Dispatcher
	LocalBuilder   LocalBuilder[BeaconStateT, ExecutionPayloadT]
	Logger         LoggerT
	Signer         crypto.BLSSigner
	StateProcessor StateProcessor[
		BeaconBlockT, BeaconStateT, *Context,
		DepositT, ExecutionPayloadHeaderT,
	]
	StorageBackend        StorageBackendT
	TelemetrySink         *metrics.TelemetrySink
	BlockStore            BeaconBlockStoreT
	DepositStore          DepositStoreT
	BeaconDepositContract DepositContractT
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	ConsensusBlockT ConsensusBlock[BeaconBlockT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		*Eth1Data, ExecutionPayloadT, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, ExecutionPayloadHeaderT, *Fork, KVStoreT,
		*Validator, Validators, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT any,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	DepositContractT deposit.Contract[DepositT],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
](
	in ChainServiceInput[
		BeaconBlockT, BeaconStateT, DepositT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, StorageBackendT, LoggerT,
		WithdrawalT, WithdrawalsT, BeaconBlockStoreT, DepositStoreT, DepositContractT,
	],
) *blockchain.Service[
	AvailabilityStoreT,
	ConsensusBlockT, BeaconBlockT, BeaconBlockBodyT,
	BeaconBlockHeaderT, BeaconStateT, BeaconBlockStoreT, DepositT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, GenesisT,
	*engineprimitives.PayloadAttributes[WithdrawalT],
] {
	return blockchain.NewService[
		AvailabilityStoreT,
		ConsensusBlockT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		BeaconBlockStoreT,
		DepositT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		GenesisT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	](
		in.StorageBackend,
		in.BlockStore,
		in.DepositStore,
		in.BeaconDepositContract,
		math.U64(in.ChainSpec.Eth1FollowDistance()),
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.Dispatcher,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.StateProcessor,
		in.TelemetrySink,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
}
