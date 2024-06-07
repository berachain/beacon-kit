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

package app

import (
	"context"
	"io"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	bkcomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	beacon "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module"
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
				bkcomponents.ProvideBlobProofVerifier,
				bkcomponents.ProvideTelemetrySink,
			),
			depinject.Supply(
				appOpts,
				logger,
				chainSpec,
			),
		),
		&appBuilder,
	); err != nil {
		panic(err)
	}

	// Build the runtime.App using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)
	app.SetTxDecoder(bkcomponents.NoOpTxConfig{}.TxDecoder())
	app.setupBeaconModule()

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// TODO: Unhack this.
func (app *BeaconApp) setupBeaconModule() {
	// Get the beacon module.
	//
	// TODO: Cleanup.
	beaconModule, ok := app.ModuleManager.
		Modules[beacon.ModuleName].(beacon.AppModule)
	if !ok {
		panic("beacon module not found")
	}

	// Set the beacon module's handlers.
	app.SetPrepareProposal(
		beaconModule.ABCIValidatorMiddleware().
			PrepareProposalHandler,
	)
	app.SetProcessProposal(
		beaconModule.
			ABCIValidatorMiddleware().
			ProcessProposalHandler,
	)
	app.SetPreBlocker(beaconModule.ABCIFinalizeBlockMiddleware().PreBlock)

	// TODO: this needs to be made un-hood.
	if err := beaconModule.StartServices(
		context.Background(),
	); err != nil {
		panic(err)
	}
}
