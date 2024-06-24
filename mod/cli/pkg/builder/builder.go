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
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/cometbft"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder[NodeT types.Node[T], T transaction.Tx] struct {
	depInjectCfg depinject.Config
	name         string
	description  string
	// components is a list of component providers for depinject.
	components []any
	// suppliers is a list of suppliers for depinject.
	suppliers []any
	// runHandler is a function to set up run handlers for the command.
	runHandler runHandler
	// nodeBuilderFunc is a function that builds the Node,
	// eventually called by the cosmos-sdk.
	// TODO: CLI should not know about the AppCreator
	nodeBuilderFunc serverv2.AppCreator[T]
	// rootCmdSetup is a function that sets up the root command.
	rootCmdSetup rootCmdSetup[NodeT, T]
}

// New returns a new CLIBuilder with the given options.
func New[NodeT types.Node[T], T transaction.Tx](opts ...Opt[NodeT, T]) *CLIBuilder[NodeT, T] {
	cb := &CLIBuilder[NodeT, T]{
		suppliers: []any{
			os.Stdout, // supply io.Writer for logger
			viper.GetViper(),
		},
	}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

// Build builds the CLI commands.
func (cb *CLIBuilder[NodeT, T]) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
		chainSpec   common.ChainSpec
		logger      log.Logger
		cmtServer   *cometbft.CometBFTServer[transaction.Tx]
	)
	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			cb.depInjectCfg,
			depinject.Supply(
				cb.suppliers...,
			),
			depinject.Provide(
				cb.components...,
			),
		),
		&mm,
		&logger,
		&clientCtx,
		&cmtServer,
		&chainSpec,
		&autoCliOpts,
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

	// enhance the root command with the autoCliOpts
	if err := rootCmd.Enhance(autoCliOpts.EnhanceRootCommand); err != nil {
		return nil, err
	}

	// hood for now
	cmdList := cmdlib.Commands[NodeT](
		rootCmd,
		mm,
		cb.nodeBuilderFunc,
		chainSpec,
	)

	cmdlib.
		DefaultCommandConfig(
			rootCmd,
			cb.nodeBuilderFunc,
			logger,
			[]*serverv2.ServerComponent[transaction.Tx]{cmtServer},
		)

	return rootCmd, nil
}

// defaultRunHandler returns the default run handler for the CLIBuilder.
func (cb *CLIBuilder[NodeT, T]) defaultRunHandler(logger log.Logger) func(
	cmd *cobra.Command,
) error {
	return func(cmd *cobra.Command) error {
		return cb.InterceptConfigsPreRunHandler(
			cmd,
			logger,
			DefaultAppConfigTemplate(),
			DefaultAppConfig(),
			DefaultCometConfig(),
		)
	}
}

// InterceptConfigsPreRunHandler is identical to
// InterceptConfigsAndCreateContext except it also sets the server context on
// the command and the server logger.
func (cb *CLIBuilder[NodeT, T]) InterceptConfigsPreRunHandler(
	cmd *cobra.Command, logger log.Logger, customAppConfigTemplate string,
	customAppConfig interface{}, cmtConfig *cmtcfg.Config,
) error {
	serverCtx, err := server.InterceptConfigsAndCreateContext(
		cmd, customAppConfigTemplate, customAppConfig, cmtConfig)
	if err != nil {
		return err
	}

	serverCtx.Logger = logger

	// set server context
	return server.SetCmdServerContext(cmd, serverCtx)
}
