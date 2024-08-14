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
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/engines/echo"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: we could make engine type configurable
func ProvideNodeAPIEngine() *echo.Engine {
	return echo.NewDefaultEngine()
}

type NodeAPIBackendInput[
	AvailabilityStoreT any,
	BeaconBlockT any,
	BeaconStateT any,
	BlobSidecarsT any,
	BlockStoreT any,
	ContextT any,
	DepositT any,
	DepositStoreT any,
	ExecutionPayloadHeaderT any,
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, ContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
] struct {
	depinject.In

	ChainSpec      common.ChainSpec
	StateProcessor StateProcessorT
	StorageBackend StorageBackendT
}

func ProvideNodeAPIBackend[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	KVStoreT any,
	NodeT nodetypes.Node,
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, TransitionContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	TransitionContextT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT Validators[ValidatorT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
](in NodeAPIBackendInput[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT, BlobSidecarsT,
	BlockStoreT, TransitionContextT, DepositT, DepositStoreT,
	ExecutionPayloadHeaderT, StateProcessorT, StorageBackendT,
]) *backend.Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
	BlobSidecarsT, BlockStoreT, sdk.Context, DepositT, DepositStoreT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, KVStoreT,
	StorageBackendT, ValidatorT, ValidatorsT, WithdrawalT,
	WithdrawalCredentialsT,
] {
	return backend.New[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		BeaconStateMarshallableT,
		BlobSidecarsT,
		BlockStoreT,
		sdk.Context,
		DepositT,
		DepositStoreT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		NodeT,
		KVStoreT,
		StorageBackendT,
		ValidatorT,
		ValidatorsT,
		WithdrawalT,
		WithdrawalCredentialsT,
	](
		in.StorageBackend,
		in.ChainSpec,
		in.StateProcessor,
	)
}

type NodeAPIServerInput[
	ContextT NodeAPIContext,
	EngineT NodeAPIEngine[ContextT],
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In

	Engine   EngineT
	Config   *config.Config
	Handlers []handlers.Handlers[ContextT]
	Logger   LoggerT
}

func ProvideNodeAPIServer[
	ContextT NodeAPIContext,
	EngineT NodeAPIEngine[ContextT],
	LoggerT log.AdvancedLogger[any, LoggerT],
](in NodeAPIServerInput[
	ContextT, EngineT, LoggerT,
]) *server.Server[ContextT, EngineT] {
	in.Logger.AddKeyValColor("service", "node-api-server",
		log.Blue)
	return server.New[ContextT, EngineT](
		in.Config.NodeAPI,
		in.Engine,
		in.Logger.With("service", "node-api-server"),
		in.Handlers...,
	)
}

func DefaultNodeAPIComponents[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT NodeAPIContext,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	EngineT NodeAPIEngine[ContextT],
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	KVStoreT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	NodeT nodetypes.Node,
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, TransitionContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	TransitionContextT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT Validators[ValidatorT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
]() []any {
	return []any{
		ProvideNodeAPIServer[ContextT, EngineT, LoggerT],
		ProvideNodeAPIEngine,
		ProvideNodeAPIBackend[
			AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
			BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
			Eth1DataT, ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT,
			KVStoreT, NodeT, StateProcessorT, StorageBackendT,
			TransitionContextT, ValidatorT, ValidatorsT, WithdrawalT,
			WithdrawalCredentialsT,
		],
	}
}
