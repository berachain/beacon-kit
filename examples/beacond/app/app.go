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
	"io"

	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"

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
	beaconkitconfig "github.com/berachain/beacon-kit/config"
	beaconkitruntime "github.com/berachain/beacon-kit/runtime"
	beaconkeeper "github.com/berachain/beacon-kit/runtime/modules/beacon/keeper"
	stakingwrapper "github.com/berachain/beacon-kit/runtime/modules/staking"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
)

var (
	_ runtime.AppI            = (*BeaconApp)(nil)
	_ servertypes.Application = (*BeaconApp)(nil)
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

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Provide(
				beaconkitruntime.ProvideRuntime,
				bls12381.ProvideBlsSigner,
			),
			depinject.Supply(
				// supply the application options
				appOpts,
				// supply the logger
				logger,
				// supply beaconkit options
				beaconkitconfig.MustReadConfigFromAppOpts(appOpts),
				// supply our custom staking wrapper.
				stakingwrapper.NewKeeper(app.StakingKeeper), // StakingKeeper is nil here.
			),
		),
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
		&app.BeaconKitRuntime,
	); err != nil {
		panic(err)
	}
	// Build the app using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)
	// Build all the ABCI Componenets.
	prepare, process, preBlocker := app.BeaconKitRuntime.BuildABCIComponents(
		baseapp.NewDefaultProposalHandler(app.Mempool(), app).
			PrepareProposalHandler(),
		baseapp.NewDefaultProposalHandler(app.Mempool(), app).
			ProcessProposalHandler(),
		nil,
		stakingwrapper.NewKeeper(app.StakingKeeper),
	)

	// Set all the newly built ABCI Componenets on the App.
	app.SetPrepareProposal(prepare)
	app.SetProcessProposal(process)
	app.SetPreBlocker(preBlocker)

	// TODO: Fix Depinject.
	app.BeaconKeeper.SetValsetChangeProvider(
		// TODO add in dep inject
		stakingwrapper.NewKeeper(app.StakingKeeper))

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

	return app
}

// PostStartup is called after the app has started up and CometBFT is connected.
func (app *BeaconApp) PostStartup(
	ctx context.Context,
	clientCtx client.Context,
) error {
	// Initial check for execution client sync.
	app.BeaconKitRuntime.StartServices(
		ctx,
		clientCtx,
	)
	return nil
}

// kvStoreKeys returns the KVStoreKeys for the app.
func (app *BeaconApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}
