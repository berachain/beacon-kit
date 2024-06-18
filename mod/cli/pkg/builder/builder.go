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
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	cmdlib "github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is the builder for the commands.Root (root command).
type CLIBuilder[T servertypes.Application] struct {
	depInjectCfg depinject.Config
	name         string
	description  string
	// components is a list of component providers for depinject.
	components []any
	// supplies is a list of suppliers for depinject.
	supplies   []any
	runHandler runHandler
	// appCreator is a function that builds the Node, eventually called by the
	// cosmos-sdk.
	// TODO: CLI should not know about the AppCreator
	appCreator servertypes.AppCreator[T]
	// rootCmdSetup is a function that sets up the root command.
	rootCmdSetup rootCmdSetup[T]
}

// New returns a new CLIBuilder with the given options.
func New[T servertypes.Application](opts ...Opt[T]) *CLIBuilder[T] {
	cb := &CLIBuilder[T]{
		supplies: []any{log.NewLogger(os.Stdout), viper.GetViper()},
	}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

// Build builds the CLI commands.
func (cb *CLIBuilder[T]) Build() (*cmdlib.Root, error) {
	// allocate memory to hold the dependencies
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
		chainSpec   primitives.ChainSpec
	)
	// build dependencies for the root command
	if err := depinject.Inject(
		depinject.Configs(
			cb.depInjectCfg,
			depinject.Supply(
				cb.supplies...,
			),
			depinject.Provide(
				cb.components...,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
		&chainSpec,
	); err != nil {
		return nil, err
	}

	// pass in deps to build the root command
	rootCmd := cmdlib.New(
		cb.name,
		cb.description,
		defaultRunHandler(cb.runHandler),
		clientCtx,
	)

	// enhance the root command with the autoCliOpts
	if err := rootCmd.Enhance(autoCliOpts.EnhanceRootCommand); err != nil {
		return nil, err
	}

	// apply default root command setup
	cmdlib.DefaultRootCommandSetup(
		rootCmd,
		mm,
		cb.appCreator,
		chainSpec,
	)

	return rootCmd, nil
}

// defaultRunHandler returns a runHandler that uses the default configuration.
func defaultRunHandler(runHandler runHandler) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {
		return runHandler(
			cmd,
			DefaultAppConfigTemplate(),
			DefaultAppConfig(),
			DefaultCometConfig(),
		)
	}
}
