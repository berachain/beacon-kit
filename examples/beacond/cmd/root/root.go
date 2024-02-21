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
	"os"

	"github.com/spf13/cobra"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"

	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/itsdevbear/bolaris/examples/beacond/app"
	"github.com/itsdevbear/bolaris/io/cli/tos"

	beaconconfig "github.com/itsdevbear/bolaris/config"
)

// NewRootCmd creates a new root command for simd. It is called once in the main function.
//

func NewRootCmd() *cobra.Command {
	var (
		interfaceRegistry  codectypes.InterfaceRegistry
		appCodec           codec.Codec
		txConfig           client.TxConfig
		legacyAmino        *codec.LegacyAmino
		autoCliOpts        autocli.AppOptions
		moduleBasicManager module.BasicManager
	)

	if err := depinject.Inject(
		depinject.Configs(
			app.MakeAppConfig(""),
			depinject.Supply(
				log.NewNopLogger(),
				simtestutil.NewAppOptionsWithFlagHome(tempDir()),
				&beaconconfig.Beacon{},
			),
			depinject.Provide(),
		),
		&interfaceRegistry,
		&appCodec,
		&txConfig,
		&legacyAmino,
		&autoCliOpts,
		&moduleBasicManager,
	); err != nil {
		panic(err)
	}

	initClientCtx := client.Context{}.
		WithCodec(appCodec).
		WithInterfaceRegistry(interfaceRegistry).
		WithLegacyAmino(legacyAmino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("") // In BeaconApp, we don't use any prefix for env variables.

	rootCmd := &cobra.Command{
		Use:   "beacond",
		Short: "beacon-kit sample app",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = tos.VerifyTosAcceptedOrPrompt(
				app.AppName, app.TermsOfServiceURL, initClientCtx, cmd,
			); err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL.
			txConfigWithTextual, err := tx.NewTxConfigWithOptions(
				codec.NewProtoCodec(interfaceRegistry),
				tx.ConfigOptions{
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
				},
			)
			if err != nil {
				return err
			}

			initClientCtx = initClientCtx.WithTxConfig(txConfigWithTextual)
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, txConfig, interfaceRegistry, appCodec, moduleBasicManager)

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

// add server commands.
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string,
	appCreator servertypes.AppCreator, appExport servertypes.AppExporter,
	addStartFlags servertypes.ModuleInitFlags, opts server.StartCmdOptions) {
	cometCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		server.ShowNodeIDCmd(),
		server.ShowValidatorCmd(),
		server.ShowAddressCmd(),
		server.VersionCmd(),
		cmtcmd.ResetAllCmd,
		cmtcmd.ResetStateCmd,
		server.BootstrapStateCmd(appCreator),
	)

	startCmd := server.StartCmdWithOptions(appCreator, defaultNodeHome, opts)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		cometCmd,
		server.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		server.NewRollbackCmd(appCreator, defaultNodeHome),
	)
}
