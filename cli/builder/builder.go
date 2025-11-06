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

package builder

import (
	"os"

	"cosmossdk.io/depinject"
	cmdlib "github.com/berachain/beacon-kit/cli/commands"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log/phuslu"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder struct {
	name        string
	description string
	// components is a list of component providers for depinject.
	components []any
	// suppliers is a list of suppliers for depinject.
	suppliers []any
	// nodeBuilderFunc is a function that builds the Node,
	// eventually called by the cosmos-sdk.
	// TODO: CLI should not know about the AppCreator
	nodeBuilderFunc      servertypes.AppCreator
	chainSpecBuilderFunc servertypes.ChainSpecCreator
}

// New returns a new CLIBuilder with the given options.
func New(opts ...Opt) *CLIBuilder {
	cb := &CLIBuilder{
		suppliers: []any{
			os.Stdout, // supply io.Writer for logger
		},
	}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

// Build builds the CLI commands.
func (cb *CLIBuilder) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		clientCtx client.Context
		logger    *phuslu.Logger
	)

	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Supply(cb.suppliers...),
			depinject.Provide(
				cb.components...,
			),
		),
		&logger,
		&clientCtx,
	); err != nil {
		return nil, err
	}

	// pass in deps to build the root command
	rootCmd := cmdlib.New(
		cb.name,
		cb.description,
		cb.defaultRunHandler(logger),
		clientCtx,
	)

	// apply default root command setup
	cmdlib.DefaultRootCommandSetup(
		rootCmd,
		&cometbft.Service{},
		cb.nodeBuilderFunc,
		cb.chainSpecBuilderFunc,
	)

	return rootCmd, nil
}

// defaultRunHandler returns the default run handler for the CLIBuilder.
func (cb *CLIBuilder) defaultRunHandler(logger *phuslu.Logger) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {
		return cb.InterceptConfigsPreRunHandler(
			cmd,
			logger,
			DefaultAppConfigTemplate(),
			DefaultAppConfig(),
			cometbft.DefaultConfig(),
		)
	}
}

func (cb *CLIBuilder) InterceptConfigsPreRunHandler(
	cmd *cobra.Command,
	logger *phuslu.Logger,
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
) error {
	return config.SetupCommand(
		cmd,
		customAppConfigTemplate,
		customAppConfig,
		cmtConfig,
		logger,
	)
}
