// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package nodebuilder

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/noop"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	cmdlib "github.com/berachain/beacon-kit/mod/node-builder/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppInfo is a struct that holds the application information.
type AppInfo[T servertypes.Application] struct {
	// Name is the name of the application.
	Name string
	// Description is a short description of the application.
	Description string
	// DepInjectConfig is the configuration for the application.
	DepInjectConfig depinject.Config
}

// NodeBuilder is a struct that holds the application information.
type NodeBuilder[T servertypes.Application] struct {
	// Every node has some application it is running.
	appInfo *AppInfo[T]

	// chainSpec is the chain specification for the application.
	chainSpec primitives.ChainSpec

	// rootCmd is the root command for the application.
	rootCmd *cobra.Command
}

// NewNodeBuilder creates a new NodeBuilder.
func NewNodeBuilder[T servertypes.Application]() *NodeBuilder[T] {
	return &NodeBuilder[T]{}
}

// Run runs the application.
func (nb *NodeBuilder[T]) RunNode() error {
	if err := nb.BuildRootCmd(); err != nil {
		return err
	}

	// Run the root command.
	if err := svrcmd.Execute(
		nb.rootCmd, "", components.DefaultNodeHome,
	); err != nil {
		log.NewLogger(nb.rootCmd.OutOrStderr()).
			Error("failure when running app", "error", err)
		return err
	}
	return nil
}

// BuildRootCmd builds the root command for the application.
func (nb *NodeBuilder[T]) BuildRootCmd() error {
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
	)
	if err := depinject.Inject(
		depinject.Configs(
			nb.appInfo.DepInjectConfig,
			depinject.Supply(
				log.NewLogger(os.Stdout),
				viper.GetViper(),
				nb.chainSpec,
				&depositdb.KVStore[*types.Deposit]{},
				&engineclient.EngineClient[*types.ExecutionPayload]{},
				&gokzg4844.JSONTrustedSetup{},
				&noop.Verifier{},
				&dastore.Store[types.BeaconBlockBody]{},
				&signer.BLSSigner{},
			),
			depinject.Provide(
				components.ProvideNoopTxConfig,
				components.ProvideClientContext,
				components.ProvideKeyring,
				components.ProvideConfig,
				components.ProvideTelemetrySink,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
	); err != nil {
		return err
	}

	nb.rootCmd = &cobra.Command{
		Use:   nb.appInfo.Name,
		Short: nb.appInfo.Description,
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

			customClientTemplate, customClientConfig := components.InitClientConfig()
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

			return server.InterceptConfigsPreRunHandler(
				cmd,
				nb.DefaultAppConfigTemplate(),
				nb.DefaultAppConfig(),
				nb.DefaultCometConfig(),
			)
		},
	}

	cmdlib.DefaultRootCommandSetup(
		nb.rootCmd,
		mm,
		nb.AppCreator,
		nb.chainSpec,
	)

	return autoCliOpts.EnhanceRootCommand(nb.rootCmd)
}
