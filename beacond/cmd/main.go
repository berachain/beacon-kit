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

package main

import (
	"log/slog"
	"os"

	clibuilder "github.com/berachain/beacon-kit/mod/cli/pkg/builder"
	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-core/pkg/builder"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	beacon "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/cosmos/cosmos-sdk/server"
	"go.uber.org/automaxprocs/maxprocs"
)

// run runs the beacon node.
func run() error {
	// Set the uber max procs
	if _, err := maxprocs.Set(); err != nil {
		return err
	}

	// Build the node using the node-core.
	nb := nodebuilder.New(
		// Set the DepInject Configuration to the Default.
		nodebuilder.WithDepInjectConfig[types.NodeI](
			nodebuilder.DefaultDepInjectConfig()),
		// Set the Runtime Components to the Default.
		nodebuilder.WithComponents[types.NodeI](
			nodecomponents.DefaultComponentsWithStandardTypes(),
		),
	)

	// Build the root command using the builder
	cb := clibuilder.New(
		// Set the Name to the Default.
		clibuilder.WithName[types.NodeI](nodebuilder.DefaultAppName),
		// Set the Description to the Default.
		clibuilder.WithDescription[types.NodeI](nodebuilder.DefaultDescription),
		// Set the DepInject Configuration to the Default.
		clibuilder.WithDepInjectConfig[types.NodeI](
			nodebuilder.DefaultDepInjectConfig(),
		),
		// Set the Runtime Components to the Default.
		clibuilder.WithComponents[types.NodeI](
			append(
				clicomponents.DefaultClientComponents(),
				// TODO: remove these, and eventually pull cfg and chainspec
				// from built node
				nodecomponents.ProvideNoopTxConfig,
				nodecomponents.ProvideConfig,
				nodecomponents.ProvideChainSpec,
			),
		),
		clibuilder.SupplyModuleDeps[types.NodeI](
			beacon.SupplyModuleDependencies(),
		),
		// Set the Run Handler to the Default.
		clibuilder.WithRunHandler[types.NodeI](
			server.InterceptConfigsPreRunHandler,
		),
		// Set the AppCreator to the NodeBuilder AppCreator.
		clibuilder.WithAppCreator[types.NodeI](nb.AppCreator),
	)

	// we never have to call nb.build() because this function is passed
	// to the cli through the clibuilder.WithAppCreator option, eventually to be
	// called by the cosmos-sdk

	cmd, err := cb.Build()
	if err != nil {
		return err
	}

	// eventually we want to decouple from cosmos cli, and just pass in a built
	// Node and Cmd to a runner

	// for now, running the cmd will start the node
	return cmd.Run(clicomponents.DefaultNodeHome)
}

// main is the entry point.
func main() {
	if err := run(); err != nil {
		//nolint:sloglint // todo fix.
		slog.Error("startup failure", "error", err)
		os.Exit(1)
	}
}
