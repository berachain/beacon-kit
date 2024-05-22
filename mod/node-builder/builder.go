// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package nodebuilder

import (
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	cmdlib "github.com/berachain/beacon-kit/mod/node-builder/pkg/commands"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/utils/tos"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
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
				&depositdb.KVStore{},
				&engineclient.EngineClient[*types.ExecutableDataDeneb]{},
				&gokzg4844.JSONTrustedSetup{},
				&dastore.Store[types.BeaconBlockBody]{},
			),
			depinject.Provide(
				components.ProvideClientContext,
				components.ProvideKeyring,
				components.ProvideConfig,
				components.ProvideBlsSigner,
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

			if err = tos.VerifyTosAcceptedOrPrompt(
				nb.appInfo.Name, components.TermsOfServiceURL, clientCtx, cmd,
			); err != nil {
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
