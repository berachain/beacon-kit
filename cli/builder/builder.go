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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	cmdlib "github.com/berachain/beacon-kit/cli/commands"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-core/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder[
	T types.Node,
	LoggerT log.AdvancedLogger[LoggerT],
] struct {
	name        string
	description string
	// components is a list of component providers for depinject.
	components []any
	// suppliers is a list of suppliers for depinject.
	suppliers []any
	// nodeBuilderFunc is a function that builds the Node,
	// eventually called by the cosmos-sdk.
	// TODO: CLI should not know about the AppCreator
	nodeBuilderFunc servertypes.AppCreator[T, LoggerT]
}

// New returns a new CLIBuilder with the given options.
func New[
	T types.Node,
	LoggerT log.AdvancedLogger[LoggerT],
](
	opts ...Opt[T, LoggerT],
) *CLIBuilder[T, LoggerT] {
	cb := &CLIBuilder[T, LoggerT]{
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
func (cb *CLIBuilder[
	T, LoggerT,
]) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		clientCtx client.Context
		chainSpec chain.ChainSpec
		logger    LoggerT
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
	cmdlib.DefaultRootCommandSetup[T](
		rootCmd,
		&cometbft.Service[LoggerT]{},
		cb.nodeBuilderFunc,
		chainSpec,
	)

	return rootCmd, nil
}

// defaultRunHandler returns the default run handler for the CLIBuilder.
func (cb *CLIBuilder[_, LoggerT]) defaultRunHandler(
	logger LoggerT,
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

func (cb *CLIBuilder[_, LoggerT]) InterceptConfigsPreRunHandler(
	cmd *cobra.Command,
	logger LoggerT,
	customAppConfigTemplate string,
	customAppConfig interface{},
	cmtConfig *cmtcfg.Config,
) error {
	return config.SetupCommand[LoggerT](
		cmd,
		customAppConfigTemplate,
		customAppConfig,
		cmtConfig,
		logger,
	)
}
