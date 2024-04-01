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

//nolint:govet,revive // from sdk.
package nodebuilder

import (
	"context"
	"os"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	modclient "github.com/berachain/beacon-kit/mod/node-builder/client"
	cmdlib "github.com/berachain/beacon-kit/mod/node-builder/commands"
	"github.com/berachain/beacon-kit/mod/node-builder/commands/utils/tos"
	_ "github.com/berachain/beacon-kit/mod/node-builder/utils/maxprocs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// NewRootCmd creates a new root command for simd. It is called once in the main
// function.
func NewRootCmd[T servertypes.Application](
	appName, appDescription string,
	appConfig depinject.Config,
	appCreator servertypes.AppCreator[T],
) *cobra.Command {
	var (
		autoCliOpts autocli.AppOptions
		mm          *module.Manager
		clientCtx   client.Context
	)
	if err := depinject.Inject(
		depinject.Configs(
			appConfig,
			depinject.Supply(
				log.NewNopLogger(),
				simtestutil.NewAppOptionsWithFlagHome(tempDir()),
			),
			depinject.Provide(
				modclient.ProvideClientContext,
				modclient.ProvideKeyring,
			),
		),
		&autoCliOpts,
		&mm,
		&clientCtx,
	); err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   appName,
		Short: appDescription,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			clientCtx, err := client.ReadPersistentCommandFlags(
				clientCtx,
				cmd.Flags(),
			)
			if err != nil {
				return err
			}

			if err = tos.VerifyTosAcceptedOrPrompt(
				appName, modclient.TermsOfServiceURL, clientCtx, cmd,
			); err != nil {
				return err
			}

			customClientTemplate, customClientConfig := modclient.InitClientConfig()
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

			customAppTemplate, customAppConfig := modclient.InitAppConfig()
			customCMTConfig := modclient.InitCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(
				cmd,
				customAppTemplate,
				customAppConfig,
				customCMTConfig,
			)
		},
	}

	cmdlib.DefaultRootCommandSetup(
		rootCmd,
		mm,
		appCreator,
		func(
			_app T,
			_ *server.Context,
			clientCtx client.Context,
			ctx context.Context,
			_ *errgroup.Group,
		) error {
			return interface{}(_app).(BeaconApp).PostStartup(
				ctx,
				clientCtx,
			)
		},
	)

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd
}

type BeaconApp interface {
	PostStartup(
		ctx context.Context,
		clientCtx client.Context,
	) error
}

// // appExport creates a new BeaconApp (optionally at a given height) and
// exports
// // state.
// func appExport(
// 	logger log.Logger,
// 	db dbm.DB,
// 	traceStore io.Writer,
// 	height int64,
// 	forZeroHeight bool,
// 	jailAllowedAddrs []string,
// 	appOpts servertypes.AppOptions,
// 	modulesToExport []string,
// ) (servertypes.ExportedApp, error) {
// 	// this check is necessary as we use the flag in x/upgrade.
// 	// we can exit more gracefully by checking the flag here.
// 	homePath, ok := appOpts.Get(flags.FlagHome).(string)
// 	if !ok || homePath == "" {
// 		return servertypes.ExportedApp{}, errors.New("application home not set")
// 	}

// 	viperAppOpts, ok := appOpts.(*viper.Viper)
// 	if !ok {
// 		return servertypes.ExportedApp{}, errors.New(
// 			"appOpts is not viper.Viper",
// 		)
// 	}

// 	// overwrite the FlagInvCheckPeriod
// 	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
// 	appOpts = viperAppOpts

// 	var beaconApp *app.BeaconApp
// 	if height != -1 {
// 		beaconApp = app.NewBeaconKitApp(
// 			logger,
// 			db,
// 			traceStore,
// 			false,
// 			appOpts,
// 		)

// 		if err := beaconApp.LoadHeight(height); err != nil {
// 			return servertypes.ExportedApp{}, err
// 		}
// 	} else {
// 		beaconApp = app.NewBeaconKitApp(logger, db, traceStore, true, appOpts)
// 	}

// 	return beaconApp.ExportAppStateAndValidators(
// 		forZeroHeight,
// 		jailAllowedAddrs,
// 		modulesToExport,
// 	)
// }

var tempDir = func() string { //nolint:gochecknoglobals // from sdk.
	dir, err := os.MkdirTemp("", "beacond")
	if err != nil {
		dir = modclient.DefaultNodeHome
	}
	defer os.RemoveAll(dir)

	return dir
}
