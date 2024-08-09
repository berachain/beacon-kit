// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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

package runtime

import (
	"context"
	"encoding/json"
	"io"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/gogoproto/proto"
)

var _ servertypes.Application = &App{}

// App is a wrapper around BaseApp and ModuleManager that can be used in hybrid
// app.go/app config scenarios or directly as a servertypes.Application
// instance.
// To get an instance of *App, *AppBuilder must be requested as a dependency
// in a container which declares the runtime module and the AppBuilder.Build()
// method must be called.
//
// App can be used to create a hybrid app.go setup where some configuration is
// done declaratively with an app config and the rest of it is done the old way.
// See simapp/app.go for an example of this setup.
type App struct {
	*baseapp.BaseApp

	ModuleManager *module.Manager
	Middleware    Middleware
	config        *runtimev1alpha1.Module
	storeKeys     []storetypes.StoreKey
	logger        log.Logger
	// initChainer is the init chainer function defined by the app config.
	// this is only required if the chain wants to add special InitChainer
	// logic.
	initChainer sdk.InitChainer
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp(
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appBuilder *AppBuilder,
	middleware Middleware,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	app := &App{
		Middleware: middleware,
	}

	// Build the runtime.App using the app builder.
	app = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// Load finishes all initialization operations and loads the app.
func (a *App) Load(loadLatest bool) error {
	if len(a.config.GetInitGenesis()) != 0 {
		a.ModuleManager.SetOrderInitGenesis(a.config.GetInitGenesis()...)
		if a.initChainer == nil {
			a.SetInitChainer(a.InitChainer)
		}
	}

	a.SetEndBlocker(a.EndBlocker)

	if loadLatest {
		if err := a.LoadLatestVersion(); err != nil {
			return err
		}
	}

	return nil
}

// EndBlocker application updates every end block.
func (a *App) EndBlocker(ctx context.Context, req proto.Message) (transition.ValidatorUpdates, error) {
	return a.Middleware.EndBlock(ctx, req)
}

// InitChainer initializes the chain.
func (a *App) InitChainer(
	ctx sdk.Context,
	req *abci.InitChainRequest,
) (*abci.InitChainResponse, error) {
	var genesisState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}
	return a.ModuleManager.InitGenesis(ctx, genesisState)
}

// LoadHeight loads a particular height.
func (a *App) LoadHeight(height int64) error {
	return a.LoadVersion(height)
}
