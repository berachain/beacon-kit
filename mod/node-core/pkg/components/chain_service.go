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
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// ChainServiceInput is the input for the chain service provider.
type ChainServiceInput[
	AvailabilityStoreT any,
	BeaconBlockT any,
	BeaconStateT any,
	BlobSidecarsT any,
	BlobFactoryT BlobFactory[BeaconBlockT, BlobSidecarsT],
	BlockStoreT any,
	ContextT any,
	DepositT any,
	DepositStoreT any,
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT any,
	PayloadAttributesT any,
	PayloadIDT ~[8]byte,
	SlotDataT any,
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, ContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	LocalBuilderT LocalBuilder[BeaconStateT, ExecutionPayloadT],
	LoggerT log.AdvancedLogger[any, LoggerT],
	WithdrawalsT Withdrawals,
] struct {
	depinject.In

	BlockBroker           *broker.Broker[*asynctypes.Event[BeaconBlockT]]
	ChainSpec             common.ChainSpec
	Cfg                   *config.Config
	ExecutionEngine       ExecutionEngineT
	GenesisBrocker        *broker.Broker[*asynctypes.Event[GenesisT]]
	LocalBuilder          LocalBuilderT
	Logger                LoggerT
	StateProcessor        StateProcessorT
	StorageBackend        StorageBackendT
	TelemetrySink         *metrics.TelemetrySink
	ValidatorUpdateBroker *broker.Broker[*asynctypes.Event[transition.ValidatorUpdates]]
}

// ProvideChainService is a depinject provider for the blockchain service.
func ProvideChainService[
	AttestationDataT any,
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[
		BeaconBlockT, AttestationDataT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, AttestationDataT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, ValidatorsT, WithdrawalsT,
	],
	BlobSidecarsT any,
	BlobFactoryT BlobFactory[BeaconBlockT, BlobSidecarsT],
	BlockStoreT any,
	ContextT Context[ContextT],
	DepositT any,
	DepositStoreT any,
	Eth1DataT any,
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT any,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	KVStoreT any,
	LocalBuilderT LocalBuilder[BeaconStateT, ExecutionPayloadT],
	LoggerT log.AdvancedLogger[any, LoggerT],
	PayloadAttributesT PayloadAttributes[ExecutionPayloadHeaderT, WithdrawalT],
	PayloadIDT ~[8]byte,
	SlashingInfoT any,
	SlotDataT any,
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, ContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT any,
	ValidatorsT any,
	WithdrawalT any,
	WithdrawalsT Withdrawals,
](
	in ChainServiceInput[
		AvailabilityStoreT, BeaconBlockT, BeaconStateT, BlobSidecarsT,
		BlobFactoryT, BlockStoreT, ContextT, DepositT, DepositStoreT,
		ExecutionEngineT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		GenesisT, PayloadAttributesT, PayloadIDT, SlotDataT,
		StateProcessorT, StorageBackendT, LocalBuilderT, LoggerT, WithdrawalsT,
	],
) *blockchain.Service[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconBlockHeaderT, BeaconStateT, ContextT, DepositT,
	ExecutionEngineT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	GenesisT, LoggerT, PayloadAttributesT, LocalBuilderT,
	PayloadIDT, StateProcessorT, StorageBackendT, WithdrawalsT,
] {
	return blockchain.NewService[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		ContextT,
		DepositT,
		ExecutionEngineT,
		ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		GenesisT,
		LoggerT,
		PayloadAttributesT,
		LocalBuilderT,
		PayloadIDT,
		StateProcessorT,
		StorageBackendT,
		WithdrawalsT,
	](
		in.StorageBackend,
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.StateProcessor,
		in.TelemetrySink,
		in.GenesisBrocker,
		in.BlockBroker,
		in.ValidatorUpdateBroker,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
}
