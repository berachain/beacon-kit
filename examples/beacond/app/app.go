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
	_ "embed"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	authkeeper "cosmossdk.io/x/auth/keeper"
	bankkeeper "cosmossdk.io/x/bank/keeper"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
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
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	beaconkitconfig "github.com/itsdevbear/bolaris/config"
	cmdconfig "github.com/itsdevbear/bolaris/lib/cmd/config"
	beaconkitruntime "github.com/itsdevbear/bolaris/runtime"
	beaconkeeper "github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper"
	stakingwrapper "github.com/itsdevbear/bolaris/runtime/modules/staking"
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
	)
}

// BeaconApp extends an ABCI application, but with most of its parameters
// exported.
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
	BeaconKeeper     *beaconkeeper.Keeper
	BeaconKitRuntime *beaconkitruntime.BeaconKitRuntime
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
	app := &BeaconApp{}
	appBuilder := &runtime.AppBuilder{}
	clientCtx := client.Context{}
	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Provide(
				beaconkitruntime.ProvideRuntime,
				cmdconfig.ProvideClientContext,
			),
			depinject.Supply(
				// supply the application options
				appOpts,
				// supply the logger
				logger,
				// supply beaconkit options
				beaconkitconfig.MustReadConfigFromAppOpts(appOpts),
				// supply our custom staking wrapper.
				stakingwrapper.NewKeeper(app.StakingKeeper),
			),
		),
		&appBuilder,
		&clientCtx,
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
		&app.BeaconKitRuntime,
	); err != nil {
		panic(err)
	}

	// Build the app using the app builder.
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Build all the ABCI Componenets.
	prepare, process, preBlocker, streamingMgr := app.BeaconKitRuntime.BuildABCIComponents(
		baseapp.NewDefaultProposalHandler(app.Mempool(), app).
			PrepareProposalHandler(),
		baseapp.NewDefaultProposalHandler(app.Mempool(), app).
			ProcessProposalHandler(),
		nil,
		app.Logger(),
	)

	// Set all the newly built ABCI Componenets on the App.
	app.SetPrepareProposal(prepare)
	app.SetProcessProposal(process)
	app.SetPreBlocker(preBlocker)
	app.SetStreamingManager(streamingMgr)

	/**** End of BeaconKit Configuration ****/

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		panic(err)
	}

	// Check for goleveldb cause bad project.
	if appOpts.Get("app-db-backend") == "goleveldb" {
		panic("goleveldb is not supported")
	}

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	// Initial check for execution client sync.
	app.BeaconKitRuntime.StartServices(
		app.NewContext(true),
		clientCtx,
	)

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

func (app *BeaconApp) Close() error {
	app.BeaconKitRuntime.Close()
	return app.App.Close()
}
