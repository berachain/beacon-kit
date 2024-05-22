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
	bkcomponents "github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	beacon "github.com/berachain/beacon-kit/mod/node-builder/pkg/components/module"
	"github.com/berachain/beacon-kit/mod/primitives"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
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

	// TODO: Deprecate.
	ConsensusParamsKeeper consensuskeeper.Keeper
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	dCfg depinject.Config,
	chainSpec primitives.ChainSpec,
	baseAppOptions ...func(*baseapp.BaseApp),
) *BeaconApp {
	app := &BeaconApp{}
	appBuilder := &runtime.AppBuilder{}
	if err := depinject.Inject(
		depinject.Configs(
			dCfg,
			depinject.Provide(
				bkcomponents.ProvideAvailibilityStore,
				bkcomponents.ProvideBlsSigner,
				bkcomponents.ProvideTrustedSetup,
				bkcomponents.ProvideDepositStore,
				bkcomponents.ProvideConfig,
				bkcomponents.ProvideEngineClient,
				bkcomponents.ProvideJWTSecret,
				bkcomponents.ProvideTelemetrySink,
			),
			depinject.Supply(
				appOpts,
				logger,
				chainSpec,
			),
		),
		&appBuilder,
		&app.ConsensusParamsKeeper,
	); err != nil {
		panic(err)
	}

	// Build the runtime.App using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Get the beacon module.
	//
	// TODO: Cleanup.
	beaconModule, ok := app.ModuleManager.
		Modules[beacon.ModuleName].(beacon.AppModule)
	if !ok {
		panic("beacon module not found")
	}

	app.SetPrepareProposal(beaconModule.ABCIHandler().PrepareProposalHandler)
	app.SetProcessProposal(beaconModule.ABCIHandler().ProcessProposalHandler)
	app.SetPreBlocker(beaconModule.ABCIHandler().FinalizeBlock)

	// Check for goleveldb cause bad project.
	if appOpts.Get("app-db-backend") == "goleveldb" {
		panic("goleveldb is not supported")
	}

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	// TODO: this needs to be made un-hood.
	if err := beaconModule.StartServices(
		context.Background(),
	); err != nil {
		panic(err)
	}

	return app
}
