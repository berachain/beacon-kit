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

package builder

import (
	"context"
	"io"

	"cosmossdk.io/depinject"
	sdklog "cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/node"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/runtime"
	servertypes "github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/server/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	dbm "github.com/cosmos/cosmos-db"
)

// NodeBuilder is a construction helper for creating nodes that implement
// the types.NodeI interface.
// TODO: #Make nodebuilder build a node. Currently this is just a builder for
// the AppCreator function, which is eventually called by cosmos to build a
// node.
type NodeBuilder[
	NodeT types.Node,
	LoggerT interface {
		log.AdvancedLogger[any, LoggerT]
		log.Configurable[LoggerT, LoggerConfigT]
	},
	LoggerConfigT any,
] struct {
	node NodeT
	// components is a list of components to provide.
	components []any
}

// New returns a new NodeBuilder.
func New[
	NodeT types.Node,
	LoggerT interface {
		log.AdvancedLogger[any, LoggerT]
		log.Configurable[LoggerT, LoggerConfigT]
	},
	LoggerConfigT any,
](
	opts ...Opt[NodeT, LoggerT, LoggerConfigT],
) *NodeBuilder[NodeT, LoggerT, LoggerConfigT] {
	nb := &NodeBuilder[NodeT, LoggerT, LoggerConfigT]{
		node: node.New[NodeT](),
	}
	for _, opt := range opts {
		opt(nb)
	}
	return nb
}

// Build uses the node builder options and runtime parameters to
// build a new instance of the node.
// It is necessary to adhere to the types.AppCreator[T] interface.
func (nb *NodeBuilder[NodeT, LoggerT, LoggerConfigT]) Build(
	logger sdklog.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) NodeT {
	// variables to hold the components needed to set up BeaconApp
	var (
		chainSpec       common.ChainSpec
		abciMiddleware  runtime.Middleware
		serviceRegistry *service.Registry
		consensusEngine components.ConsensusEngine
		apiBackend      interface{ AttachNode(NodeT) }
		storeKey        = new(storetypes.KVStoreKey)
		storeKeyDblPtr  = &storeKey
	)

	// build all node components using depinject
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Provide(
				nb.components...,
			),
			depinject.Supply(
				appOpts,
				logger.Impl().(LoggerT),
			),
			// TODO: cosmos depinject bad project, fixed with dig.
			// depinject.Invoke(
			// 	SetLoggerConfig[LoggerT, LoggerConfigT],
			// ),
		),
		&storeKeyDblPtr,
		&chainSpec,
		&abciMiddleware,
		&serviceRegistry,
		&consensusEngine,
		&apiBackend,
	); err != nil {
		panic(err)
	}

	if consensusEngine == nil || apiBackend == nil {
		panic("consensus engine or api backend is nil")
	}

	// set the application to a new BeaconApp with necessary ABCI handlers
	nb.node.RegisterApp(
		runtime.NewBeaconKitApp(
			*storeKeyDblPtr, logger, db, traceStore, true, abciMiddleware,
			append(
				DefaultBaseappOptions(appOpts),
				WithCometParamStore(chainSpec),
				WithPrepareProposal(consensusEngine.PrepareProposal),
				WithProcessProposal(consensusEngine.ProcessProposal),
			)...,
		),
	)
	// TODO: so hood
	apiBackend.AttachNode(nb.node)
	nb.node.SetServiceRegistry(serviceRegistry)

	// TODO: put this in some post node creation hook/listener.
	if err := nb.node.Start(context.Background()); err != nil {
		logger.Error("failed to start node", "err", err)
		panic(err)
	}
	return nb.node
}
