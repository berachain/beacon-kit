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
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	cmtcfg "github.com/cometbft/cometbft/config"
)

// TODO: we could make engine type configurable
func ProvideNodeAPIEngine() *echo.Engine {
	return echo.NewDefaultEngine()
}

type NodeAPIBackendInput struct {
	depinject.In

	ChainSpec      chain.Spec
	StorageBackend *storage.Backend
	CometConfig    *cmtcfg.Config
}

func ProvideNodeAPIBackend(
	in NodeAPIBackendInput,
) (*backend.Backend, error) {
	return backend.New(
		in.StorageBackend,
		in.ChainSpec,
		in.CometConfig,
	)
}

type NodeAPIServerInput struct {
	depinject.In

	Engine   NodeAPIEngine
	Config   *config.Config
	Handlers []handlers.Handlers
	Logger   *phuslu.Logger
}

func ProvideNodeAPIServer(in NodeAPIServerInput) *server.Server {
	in.Logger.AddKeyValColor(
		"service",
		"node-api-server",
		log.Blue,
	)
	return server.New(
		in.Config.NodeAPI,
		in.Engine,
		in.Logger.With("service", "node-api-server"),
		in.Handlers...,
	)
}
