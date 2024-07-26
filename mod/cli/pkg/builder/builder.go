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
	"cosmossdk.io/runtime/v2"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/api/grpc"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/cli/pkg/utils/context"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/server/pkg/components/cometbft"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder[
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
] struct {
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
	nodeBuilderFunc serverv2.AppCreator[NodeT, T]
	// rootCmdSetup is a function that sets up the root command.
	rootCmdSetup rootCmdSetup[NodeT, T]
	// server is the server to be used by the commands.
	server *serverv2.Server[NodeT, T]
}

// New returns a new CLIBuilder with the given options.
func New[
	NodeT types.Node[T], T transaction.Tx any,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	ValidatorUpdateT any,
](
	opts ...Opt[NodeT, T, ValidatorUpdateT],
) *CLIBuilder[NodeT, T, ValidatorUpdateT] {
	cb := &CLIBuilder[NodeT, T, ValidatorUpdateT]{
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
func (
	cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT],
) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		autoCliOpts autocli.AppOptions
		mm          *runtime.MM[T]
		clientCtx   client.Context
		chainSpec   common.ChainSpec
		logger      log.Logger
		cmtServer   *cometbft.Server[NodeT, T, ValidatorUpdateT]
	)
	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			cb.depInjectCfg,
			depinject.Supply(
				append(
					cb.suppliers, &components.StorageBackend{})...,
			),
			depinject.Provide(
				cb.components...,
			),
		),
		&mm,
		&logger,
		&clientCtx,
		&chainSpec,
		&autoCliOpts,
		&cmtServer,
	); err != nil {
		return nil, err
	}

	// build the server
	// TOOD: move into server once depinject gets sorted
	cb.server = serverv2.NewServer(
		logger,
		cmtServer,
		grpc.New[NodeT, T](),
	)

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
	// get list of custom commands
	cmdList := cmdlib.Commands[NodeT](
		rootCmd,
		mm,
		cb.nodeBuilderFunc,
		chainSpec,
	)

	// get the default command config with the server
	cmdConfig, err := cmdlib.DefaultCommandConfig(
		rootCmd.Command(),
		cb.nodeBuilderFunc,
		logger,
		cmdList,
		cb.server,
	)
	if err != nil {
		return nil, err
	}
	// add the commands to the root command
	cmdlib.AddCommands[NodeT, T](
		rootCmd.Command(),
		cb.nodeBuilderFunc,
		logger,
		cmdConfig,
		cb.server,
	)

	return rootCmd, nil
}

// defaultRunHandler returns the default run handler for the CLIBuilder.
func (
	cb *CLIBuilder[NodeT, T, ExecutionPayloadT, ValidatorUpdateT],
) defaultRunHandler(logger log.Logger) func(
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
func (cb *CLIBuilder[
	NodeT, T, ExecutionPayloadT, ValidatorUpdateT,
]) InterceptConfigsPreRunHandler(
	cmd *cobra.Command, logger log.Logger, customAppConfigTemplate string,
	customAppConfig interface{}, cmtConfig *cmtcfg.Config,
) error {
	serverCtx, err := context.InterceptConfigsAndCreateContext(
		cmd, customAppConfigTemplate, customAppConfig, cmtConfig, logger)
	if err != nil {
		return err
	}

	// set server context
	return server.SetCmdServerContext(cmd, serverCtx)
}
