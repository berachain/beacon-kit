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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput[
	BeaconBlockT any,
	BeaconStateT any,
	DepositT any,
	StorageBackendT any,
	LoggerT any,
] struct {
	depinject.In

	ChainSpec       common.ChainSpec
	Cfg             *config.Config
	EngineClient    *EngineClient
	ExecutionEngine *ExecutionEngine
	Dispatcher      *Dispatcher
	LocalBuilder    LocalBuilder[BeaconStateT, *ExecutionPayload]
	Logger          LoggerT
	Signer          crypto.BLSSigner
	StateProcessor  StateProcessor[
		BeaconBlockT, BeaconStateT, *Context,
		DepositT, *ExecutionPayloadHeader,
	]
	StorageBackend StorageBackendT
	TelemetrySink  *metrics.TelemetrySink
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		*Eth1Data, *ExecutionPayload, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT any,
	DepositT any,
	DepositStoreT any,
	GenesisT Genesis[DepositT, *ExecutionPayloadHeader],
	KVStoreT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
](
	in ChainServiceInput[
		BeaconBlockT, BeaconStateT, DepositT, StorageBackendT, LoggerT,
	],
) *blockchain.Service[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconBlockHeaderT, BeaconStateT, DepositT, *ExecutionPayload,
	*ExecutionPayloadHeader, GenesisT, *PayloadAttributes,
] {
	return blockchain.NewService[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		DepositT,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		GenesisT,
		*PayloadAttributes,
	](
		in.StorageBackend,
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
