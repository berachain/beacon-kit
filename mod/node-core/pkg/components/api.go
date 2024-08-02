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

// import (
// 	"github.com/berachain/beacon-kit/mod/config"
// 	"github.com/berachain/beacon-kit/mod/depinject"
// 	"github.com/berachain/beacon-kit/mod/log"
// 	"github.com/berachain/beacon-kit/mod/node-api/backend"
// 	"github.com/berachain/beacon-kit/mod/node-api/engines/echo"
// 	"github.com/berachain/beacon-kit/mod/node-api/handlers"
// 	"github.com/berachain/beacon-kit/mod/node-api/server"
// 	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// // TODO: we could make engine type configurable
// func ProvideNodeAPIEngine() *NodeAPIEngine {
// 	return echo.NewDefaultEngine()
// }

// type NodeAPIBackendInput struct {
// 	depinject.In

// 	StorageBackend *StorageBackend
// 	ChainSpec      common.ChainSpec
// }

// func ProvideNodeAPIBackend(in NodeAPIBackendInput) *NodeAPIBackend {
// 	return backend.New[
// 		*AvailabilityStore,
// 		*BeaconBlock,
// 		*BeaconBlockBody,
// 		*BeaconBlockHeader,
// 		*BeaconState,
// 		*BeaconStateMarshallable,
// 		*BlobSidecars,
// 		*BlockStore,
// 		sdk.Context,
// 		*Deposit,
// 		*DepositStore,
// 		*Eth1Data,
// 		*ExecutionPayloadHeader,
// 		*Fork,
// 		nodetypes.Node,
// 		*StateManager,
// 		*StorageBackend,
// 		*Validator,
// 		*Withdrawal,
// 		WithdrawalCredentials,
// 	](
// 		in.StorageBackend,
// 		in.ChainSpec,
// 	)
// }

// type NodeAPIServerInput struct {
// 	depinject.In

// 	Engine   *NodeAPIEngine
// 	Config   *config.Config
// 	Logger   *Logger
// 	Handlers []handlers.Handlers[NodeAPIContext]
// }

// func ProvideNodeAPIServer(in NodeAPIServerInput) *NodeAPIServer {
// 	in.Logger.AddKeyValColor("service", "node-api-server",
// 		log.Blue)
// 	return server.New[
// 		NodeAPIContext,
// 		*NodeAPIEngine,
// 	](
// 		in.Config.NodeAPI,
// 		in.Engine,
// 		in.Logger.With("service", "node-api-server"),
// 		in.Handlers...,
// 	)
// }
