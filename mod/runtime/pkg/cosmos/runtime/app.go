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
	"errors"
	"io"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gogoproto/proto"
	"github.com/sourcegraph/conc/iter"
)

var _ servertypes.Application = &App{}

// App can be used to create a hybrid app.go setup where some configuration is
// done declaratively with an app config and the rest of it is done the old way.
// See simapp/app.go for an example of this setup.
type App struct {
	*baseapp.BaseApp

	Middleware Middleware
	StoreKeys  []storetypes.StoreKey
	// initChainer is the init chainer function defined by the app config.
	// this is only required if the chain wants to add special InitChainer
	// logic.
	initChainer sdk.InitChainer
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func a(
	storeKey *storetypes.KVStoreKey,
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	middleware Middleware,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	app := &App{
		BaseApp: baseapp.NewBaseApp(
			"BeaconKit",
			logger,
			db,
			baseAppOptions...),
		Middleware: middleware,
	}

	app.SetVersion(version.Version)
	app.MountStore(storeKey, storetypes.StoreTypeIAVL)

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// Load finishes all initialization operations and loads the app.
func (a *App) Load(loadLatest bool) error {
	a.SetInitChainer(a.InitChainer)
	a.SetFinalizeBlocker(a.FinalizeBlocker)

	if loadLatest {
		if err := a.LoadLatestVersion(); err != nil {
			return err
		}
	}

	return nil
}

// FinalizeBlocker application updates every end block.
func (a *App) FinalizeBlocker(
	ctx context.Context,
	req proto.Message,
) (transition.ValidatorUpdates, error) {
	return a.Middleware.FinalizeBlock(ctx, req)
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
	valUpdates, err := a.Middleware.InitGenesis(
		ctx,
		[]byte(genesisState["beacon"]),
	)
	if err != nil {
		return nil, err
	}

	convertedValUpdates, err := iter.MapErr(
		valUpdates,
		convertValidatorUpdate[abci.ValidatorUpdate],
	)
	if err != nil {
		return nil, err
	}

	return &abci.InitChainResponse{
		Validators: convertedValUpdates,
	}, nil

}

// LoadHeight loads a particular height.
func (a *App) LoadHeight(height int64) error {
	return a.LoadVersion(height)
}

// DefaultGenesis returns the default genesis state for the application.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	// Implement the default genesis state for the application.
	// This should return a map of module names to their respective default
	// genesis states.
	gen := make(map[string]json.RawMessage)
	s := types.DefaultGenesisDeneb()
	var err error
	gen["beacon"], err = json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return gen
}

// ValidateGenesis validates the provided genesis state.
func (app *App) ValidateGenesis(genesisData map[string]json.RawMessage) error {
	// Implement the validation logic for the provided genesis state.
	// This should validate the genesis state for each module in the
	// application.
	return nil
}

// convertValidatorUpdate abstracts the conversion of a
// transition.ValidatorUpdate to an appmodulev2.ValidatorUpdate.
// TODO: this is so hood, bktypes -> sdktypes -> generic is crazy
// maybe make this some kind of codec/func that can be passed in?
func convertValidatorUpdate[ValidatorUpdateT any](
	u **transition.ValidatorUpdate,
) (ValidatorUpdateT, error) {
	var valUpdate ValidatorUpdateT
	update := *u
	if update == nil {
		return valUpdate, errors.New("undefined validator update")
	}
	return any(abci.ValidatorUpdate{
		PubKeyBytes: update.Pubkey[:],
		PubKeyType:  crypto.CometBLSType,
		//#nosec:G701 // this is safe.
		Power: int64(update.EffectiveBalance.Unwrap()),
	}).(ValidatorUpdateT), nil
}
