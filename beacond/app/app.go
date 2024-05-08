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
	"io"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	consensuskeeper "cosmossdk.io/x/consensus/keeper"
	"github.com/berachain/beacon-kit/beacond/x/beacon/keeper"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	bkcomponents "github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	beaconkitruntime "github.com/berachain/beacon-kit/mod/runtime"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

var (
	_ runtime.AppI            = (*BeaconApp)(nil)
	_ servertypes.Application = (*BeaconApp)(nil)
)

// BeaconApp extends an ABCI application, but with most of its parameters
// exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type BeaconApp struct {
	*runtime.App
	BeaconKeeper     *keeper.Keeper
	BeaconKitRuntime *beaconkitruntime.BeaconKitRuntime[
		*datypes.BlobSidecars,
		*deposit.KVStore,
		consensus.ReadOnlyBeaconBlockBody,
		beaconkitruntime.BeaconStorageBackend[
			*datypes.BlobSidecars,
			*deposit.KVStore,
			consensus.ReadOnlyBeaconBlockBody,
		]]
	ConsensusParamsKeeper consensuskeeper.Keeper
}

// NewBeaconKitAppWithDefaultBaseAppOptions returns a reference to an
// initialized BeaconApp.
func NewBeaconKitAppWithDefaultBaseAppOptions(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) BeaconApp {
	return *NewBeaconKitApp(
		logger, db, traceStore, true,
		appOpts,
		server.DefaultBaseappOptions(appOpts)...,
	)
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *BeaconApp {
	app := &BeaconApp{}
	appBuilder := &runtime.AppBuilder{}
	if err := depinject.Inject(
		depinject.Configs(
			Config(),
			depinject.Provide(
				bkcomponents.ProvideRuntime,
				bkcomponents.ProvideBlsSigner,
				bkcomponents.ProvideTrustedSetup,
				bkcomponents.ProvideJWTSecret,
				bkcomponents.ProvideDepositStore,
				bkcomponents.ProvideConfig,
			),
			depinject.Supply(
				// supply the application options
				appOpts,
				// supply the logger
				logger,
				// TODO: allow nodebuilder to inject.
				spec.LocalnetChainSpec(),
			),
		),
		&appBuilder,
		&app.ConsensusParamsKeeper,
		&app.BeaconKeeper,
		&app.BeaconKitRuntime,
	); err != nil {
		panic(err)
	}
	// Build the runtime.App using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Build all the ABCI Components.
	prepare, process, preBlocker := app.BeaconKitRuntime.BuildABCIComponents()

	// Set all the newly built ABCI Components on the App.
	app.SetPrepareProposal(prepare)
	app.SetProcessProposal(process)
	app.SetPreBlocker(preBlocker)

	/**** End of BeaconKit Configuration ****/

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
func (app BeaconApp) PostStartup(
	ctx context.Context,
	_ client.Context,
) error {
	// Initial check for execution client sync.
	app.BeaconKitRuntime.StartServices(
		ctx,
	)
	return nil
}
