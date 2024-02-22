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

package app

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	authkeeper "cosmossdk.io/x/auth/keeper"
	bankkeeper "cosmossdk.io/x/bank/keeper"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	"cosmossdk.io/x/gov"
	govclient "cosmossdk.io/x/gov/client"
	govtypes "cosmossdk.io/x/gov/types"
	mintkeeper "cosmossdk.io/x/mint/keeper"
	_ "cosmossdk.io/x/protocolpool"
	slashingkeeper "cosmossdk.io/x/slashing/keeper"
	stakingkeeper "cosmossdk.io/x/staking/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	beaconkitconfig "github.com/itsdevbear/bolaris/config"
	beaconkitruntime "github.com/itsdevbear/bolaris/runtime"
	beaconkeeper "github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper"
	stakingwrapper "github.com/itsdevbear/bolaris/runtime/modules/staking"
	"github.com/itsdevbear/bolaris/types/cosmos"
)

//nolint:gochecknoinits // from sdk.
func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".beacond")
}

const TermsOfServiceURL = "https://github.com/berachain/beacon-kit/blob/main/TERMS_OF_SERVICE.md"

var (
	_ runtime.AppI            = (*BeaconApp)(nil)
	_ servertypes.Application = (*BeaconApp)(nil)
	// DefaultNodeHome default home directories for the application daemon.
	DefaultNodeHome string //nolint:gochecknoglobals // from sdk.
)

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		// appconfig.LoadYAML(AppConfigYAML),
		BeaconAppConfig,
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
				govtypes.ModuleName: gov.NewAppModuleBasic(
					[]govclient.ProposalHandler{},
				),
			},
		),
	)
}

// BeaconApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type BeaconApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper

	// beacon-kit required keepers
	BeaconKeeper    *beaconkeeper.Keeper
	BeaconKitRunner *beaconkitruntime.BeaconKitRuntime
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	bech32Prefix string,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *BeaconApp {
	var (
		app        = &BeaconApp{}
		appBuilder *runtime.AppBuilder
		bkCfg      = beaconkitconfig.MustReadConfigFromAppOpts(appOpts)

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				// supply the application options
				appOpts,
				// supply the logger
				logger,
				// supply the beacon config
				&(bkCfg.Beacon),
			),
		)
	)

	if err := depinject.Inject(
		appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.UpgradeKeeper,
		&app.EvidenceKeeper,
		&app.ConsensusParamsKeeper,
		&app.BeaconKeeper,
	); err != nil {
		panic(err)
	}

	// Build the app using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	/**** Start of BeaconKit Configuration ****/
	var err error
	if app.BeaconKitRunner, err = beaconkitruntime.NewDefaultBeaconKitRuntime(
		bkCfg, app.BeaconKeeper,
		stakingwrapper.NewKeeper(app.StakingKeeper),
		app.Logger(),
	); err != nil {
		panic(err)
	}

	if err = app.BeaconKitRunner.RegisterApp(app.BaseApp); err != nil {
		panic(err)
	}

	/**** End of BeaconKit Configuration ****/

	// register streaming services
	if err = app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		panic(err)
	}

	// Load the app.
	if err = app.Load(loadLatest); err != nil {
		panic(err)
	}

	ctx := cosmos.NewEmptyContextWithMS(context.Background(), app.CommitMultiStore())
	app.BeaconKitRunner.StartServices(ctx)

	// Initial check for execution client sync.
	if err := app.BeaconKitRunner.InitialSyncCheck(ctx); err != nil {
		panic(err)
	}

	return app
}

// Name returns the name of the App.
func (app *BeaconApp) Name() string { return app.BaseApp.Name() }

// LegacyAmino returns BeaconApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *BeaconApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *BeaconApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// SimulationManager implements the SimulationApp interface.
func (app *BeaconApp) SimulationManager() *module.SimulationManager {
	return nil
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *BeaconApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	ctx := cosmos.NewEmptyContextWithMS(apiSvr.ClientCtx.CmdContext, app.CommitMultiStore())
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(
		apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger,
	); err != nil {
		panic(err)
	}

	v, ok := ctx.Value(server.ServerContextKey).(*server.Context)
	if !ok {
		panic(fmt.Errorf("unexpected server context type: %T", v))
	}
	app.BeaconKitRunner.SetCometCfg(v.Config)
}
