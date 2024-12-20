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
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: we could make engine type configurable
func ProvideNodeAPIEngine() *echo.Engine {
	return echo.NewDefaultEngine()
}

type NodeAPIBackendInput[
	StorageBackendT any,
] struct {
	depinject.In

	ChainSpec      chain.ChainSpec
	StateProcessor StateProcessor[*Context]
	StorageBackend StorageBackendT
}

func ProvideNodeAPIBackend[
	AvailabilityStoreT AvailabilityStore,
	BeaconBlockStoreT BlockStore,
	DepositStoreT DepositStore,
	KVStoreT any,
	NodeT interface {
		CreateQueryContext(height int64, prove bool) (sdk.Context, error)
	},
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconBlockStoreT, DepositStoreT,
	],
](
	in NodeAPIBackendInput[StorageBackendT],
) *backend.Backend[
	AvailabilityStoreT,
	BeaconBlockStoreT,
	sdk.Context, DepositStoreT,
	NodeT, KVStoreT, StorageBackendT,
] {
	return backend.New[
		AvailabilityStoreT,
		BeaconBlockStoreT,
		sdk.Context,
		DepositStoreT,
		NodeT,
		KVStoreT,
		StorageBackendT,
	](
		in.StorageBackend,
		in.ChainSpec,
		in.StateProcessor,
	)
}

type NodeAPIServerInput[
	LoggerT log.AdvancedLogger[LoggerT],
	NodeAPIContextT NodeAPIContext,
] struct {
	depinject.In

	Engine   NodeAPIEngine[NodeAPIContextT]
	Config   *config.Config
	Handlers []handlers.Handlers[NodeAPIContextT]
	Logger   LoggerT
}

func ProvideNodeAPIServer[
	LoggerT log.AdvancedLogger[LoggerT],
	NodeAPIContextT NodeAPIContext,
](
	in NodeAPIServerInput[LoggerT, NodeAPIContextT],
) *server.Server[NodeAPIContextT] {
	in.Logger.AddKeyValColor("service", "node-api-server",
		log.Blue)
	return server.New[NodeAPIContextT](
		in.Config.NodeAPI,
		in.Engine,
		in.Logger.With("service", "node-api-server"),
		in.Handlers...,
	)
}
