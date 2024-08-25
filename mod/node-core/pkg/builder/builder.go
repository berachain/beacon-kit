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
	"io"

	"cosmossdk.io/depinject"
	servertypes "github.com/berachain/beacon-kit/mod/cli/pkg/commands/server/types"
	"github.com/berachain/beacon-kit/mod/config"
	cometbft "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	cmtcfg "github.com/cometbft/cometbft/config"
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
		log.AdvancedLogger[LoggerT]
		log.Configurable[LoggerT, LoggerConfigT]
	},
	LoggerConfigT any,
] struct {
	// components is a list of components to provide.
	components []any
}

// New returns a new NodeBuilder.
func New[
	NodeT types.Node,
	LoggerT interface {
		log.AdvancedLogger[LoggerT]
		log.Configurable[LoggerT, LoggerConfigT]
	},
	LoggerConfigT any,
](
	opts ...Opt[NodeT, LoggerT, LoggerConfigT],
) *NodeBuilder[NodeT, LoggerT, LoggerConfigT] {
	nb := &NodeBuilder[NodeT, LoggerT, LoggerConfigT]{}
	for _, opt := range opts {
		opt(nb)
	}
	return nb
}

// Build uses the node builder options and runtime parameters to
// build a new instance of the node.
// It is necessary to adhere to the types.AppCreator[T] interface.
func (nb *NodeBuilder[NodeT, LoggerT, LoggerConfigT]) Build(
	logger LoggerT,
	db dbm.DB,
	_ io.Writer,
	cmtCfg *cmtcfg.Config,
	appOpts servertypes.AppOptions,
) NodeT {
	// variables to hold the components needed to set up BeaconApp
	var (
		apiBackend interface {
			AttachQueryBackend(*cometbft.Service[LoggerT])
		}
		beaconNode NodeT
		cmtService *cometbft.Service[LoggerT]
		config     *config.Config
	)

	// build all node components using depinject
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Provide(
				nb.components...,
			),
			depinject.Supply(
				appOpts,
				logger,
				db,
				cmtCfg,
			),
		),
		&apiBackend,
		&beaconNode,
		&cmtService,
		&config,
	); err != nil {
		panic(err)
	}
	// attach the logger config
	logger.WithConfig(any(config.GetLogger()).(LoggerConfigT))

	if apiBackend == nil {
		panic("node or api backend is nil")
	}

	// TODO: so hood
	apiBackend.AttachQueryBackend(cmtService)
	return beaconNode
}
