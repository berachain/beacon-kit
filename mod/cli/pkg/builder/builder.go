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

	"cosmossdk.io/depinject"
	sdklog "cosmossdk.io/log"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/cli/pkg/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/runtime"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
] struct {
	name        string
	description string
	// components is a list of component providers for depinject.
	components []any
	// suppliers is a list of suppliers for depinject.
	suppliers []any
	// runHandler is a function to set up run handlers for the command.
	runHandler runHandler
	// nodeBuilderFunc is a function that builds the Node,
	// eventually called by the cosmos-sdk.
	// TODO: CLI should not know about the AppCreator
	nodeBuilderFunc servertypes.AppCreator[T]
	// rootCmdSetup is a function that sets up the root command.
	rootCmdSetup rootCmdSetup[T]
}

// New returns a new CLIBuilder with the given options.
func New[
	T types.Node,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	opts ...Opt[T, ExecutionPayloadT],
) *CLIBuilder[T, ExecutionPayloadT] {
	cb := &CLIBuilder[T, ExecutionPayloadT]{
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
func (cb *CLIBuilder[T, ExecutionPayloadT]) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		clientCtx client.Context
		chainSpec common.ChainSpec
		logger    log.AdvancedLogger[any, sdklog.Logger]
	)

	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Supply(
				append(
					cb.suppliers, []any{&runtime.App{}, &components.StorageBackend{}})...,
			),
			depinject.Provide(
				cb.components...,
			),
		),
		&logger,
		&clientCtx,
		&chainSpec,
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
	cmdlib.DefaultRootCommandSetup[T, ExecutionPayloadT](
		rootCmd,
		&runtime.App{},
		cb.nodeBuilderFunc,
		chainSpec,
	)

	return rootCmd, nil
}

// defaultRunHandler returns the default run handler for the CLIBuilder.
func (cb *CLIBuilder[T, ExecutionPayloadT]) defaultRunHandler(
	logger log.AdvancedLogger[any, sdklog.Logger],
) func(cmd *cobra.Command) error {
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

func (cb *CLIBuilder[T, ExecutionPayloadT]) InterceptConfigsPreRunHandler(
	cmd *cobra.Command,
	logger log.AdvancedLogger[any, sdklog.Logger],
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
) error {
	serverCtx, err := config.SetupConfigAndContext(
		cmd,
		customAppConfigTemplate,
		customAppConfig,
		cmtConfig,
		logger,
	)
	if err != nil {
		return err
	}

	// set server context
	return server.SetCmdServerContext(cmd, serverCtx)
}
