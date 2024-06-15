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
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLIBuilder is a builder that incrementally constructs a CLI command.
type CLIBuilder[T servertypes.Application] struct {
	depInjectCfg depinject.Config
	name         string
	description  string
	components   []any
	moduleDeps   []any
	runHandler   func(cmd *cobra.Command,
		customAppConfigTemplate string,
		customAppConfig interface{},
		cmtConfig *cmtcfg.Config,
	) error
	appCreator   servertypes.AppCreator[T]
	rootCmdSetup func(cmd *cmdlib.Root,
		mm *module.Manager,
		appCreator servertypes.AppCreator[T],
		chainSpec primitives.ChainSpec,
	)
}

// New returns a new CLIBuilder with the given options.
func New[T servertypes.Application](opts ...Opt[T]) *CLIBuilder[T] {
	cb := &CLIBuilder[T]{}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

// Build builds the CLI commands.
func (cb *CLIBuilder[T]) Build() (*cmdlib.Root, error) {
	// dependencies for the root command
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
				append([]any{
					log.NewLogger(os.Stdout),
					viper.GetViper(),
				},
					cb.moduleDeps...)...,
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

	cmd := &cobra.Command{
		Use:   cb.name,
		Short: cb.description,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			var err error
			clientCtx, err = client.ReadPersistentCommandFlags(
				clientCtx,
				cmd.Flags(),
			)
			if err != nil {
				return err
			}

			customClientTemplate, customClientConfig := InitClientConfig()
			clientCtx, err = config.CreateClientConfig(
				clientCtx,
				customClientTemplate,
				customClientConfig,
			)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(
				clientCtx, cmd,
			); err != nil {
				return err
			}

			return cb.runHandler(
				cmd,
				DefaultAppConfigTemplate(),
				DefaultAppConfig(),
				DefaultCometConfig(),
			)
		},
	}

	cmdlib.DefaultRootCommandSetup(
		cmdlib.New(cmd),
		mm,
		cb.appCreator,
		chainSpec,
	)

	if err := autoCliOpts.EnhanceRootCommand(cmd); err != nil {
		return nil, err
	}

	return cmdlib.New(cmd), nil
}

// InitClientConfig sets up the default client configuration, allowing for
// overrides.
func InitClientConfig() (string, interface{}) {
	clientCfg := config.DefaultConfig()
	clientCfg.KeyringBackend = "test"
	return config.DefaultClientConfigTemplate, clientCfg
}
