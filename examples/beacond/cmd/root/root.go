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

//nolint:govet,gomnd,lll // from sdk.
package root

import (
	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	cmdconfig "github.com/itsdevbear/bolaris/examples/beacond/cmd/config"

	"github.com/itsdevbear/bolaris/examples/beacond/app"
	"github.com/itsdevbear/bolaris/io/cli/tos"
)

// NewRootCmd creates a new root command for simd. It is called once in the main
// function.
func NewRootCmd() *cobra.Command {
	var (
		autoCliOpts        autocli.AppOptions
		moduleBasicManager *module.Manager
		clientCtx          client.Context
	)
	if err := depinject.Inject(
		depinject.Configs(
			app.AppConfig(),
			depinject.Supply(
				log.NewNopLogger(),
				simtestutil.NewAppOptionsWithFlagHome(tempDir()),
			),
			depinject.Provide(
				cmdconfig.ProvideClientContext,
				cmdconfig.ProvideKeyring,
			),
		),
		&autoCliOpts,
		&moduleBasicManager,
		&clientCtx,
	); err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   "beacond",
		Short: "beacon-kit sample app",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			clientCtx = clientCtx.WithCmdContext(cmd.Context()).WithViper("")
			clientCtx, err := client.ReadPersistentCommandFlags(
				clientCtx,
				cmd.Flags(),
			)
			if err != nil {
				return err
			}

			if err = tos.VerifyTosAcceptedOrPrompt(
				app.AppName, app.TermsOfServiceURL, clientCtx, cmd,
			); err != nil {
				return err
			}

			customClientTemplate, customClientConfig := cmdconfig.InitClientConfig()
			clientCtx, err = config.CreateClientConfig(
				clientCtx,
				customClientTemplate,
				customClientConfig,
			)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := cmdconfig.InitAppConfig()
			customCMTConfig := cmdconfig.InitCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(
				cmd,
				customAppTemplate,
				customAppConfig,
				customCMTConfig,
			)
		},
	}

	initRootCmd(
		rootCmd,
		clientCtx.TxConfig,
		clientCtx.InterfaceRegistry,
		clientCtx.Codec,
		moduleBasicManager,
	)

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}
