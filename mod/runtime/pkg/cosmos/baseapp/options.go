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

package baseapp

import (
	"context"
	"errors"
	"fmt"

	pruningtypes "cosmossdk.io/store/pruning/types"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// File for storing in-package BaseApp optional functions,
// for options that need access to non-exported fields of the BaseApp

// SetPruning sets a pruning option on the multistore associated with the app.
func SetPruning(opts pruningtypes.PruningOptions) func(*BaseApp) {
	return func(bapp *BaseApp) { bapp.cms.SetPruning(opts) }
}

// SetMinRetainBlocks returns a BaseApp option function that sets the minimum
// block retention height value when determining which heights to prune during
// ABCI Commit.
func SetMinRetainBlocks(minRetainBlocks uint64) func(*BaseApp) {
	return func(bapp *BaseApp) { bapp.setMinRetainBlocks(minRetainBlocks) }
}

// SetIAVLCacheSize provides a BaseApp option function that sets the size of
// IAVL cache.
func SetIAVLCacheSize(size int) func(*BaseApp) {
	return func(bapp *BaseApp) { bapp.cms.SetIAVLCacheSize(size) }
}

// SetIAVLDisableFastNode enables(false)/disables(true) fast node usage from the
// IAVL store.
func SetIAVLDisableFastNode(disable bool) func(*BaseApp) {
	return func(bapp *BaseApp) { bapp.cms.SetIAVLDisableFastNode(disable) }
}

// SetInterBlockCache provides a BaseApp option function that sets the
// inter-block cache.
func SetInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) func(*BaseApp) {
	return func(app *BaseApp) { app.setInterBlockCache(cache) }
}

// SetChainID sets the chain ID in BaseApp.
func SetChainID(chainID string) func(*BaseApp) {
	return func(app *BaseApp) { app.chainID = chainID }
}

// SetStoreLoader allows customization of the rootMultiStore initialization.
func SetStoreLoader(loader StoreLoader) func(*BaseApp) {
	return func(app *BaseApp) { app.SetStoreLoader(loader) }
}

func (app *BaseApp) SetName(name string) {
	app.name = name
}

// SetParamStore sets a parameter store on the BaseApp.
func (app *BaseApp) SetParamStore(ps ParamStore) {
	app.paramStore = ps
}

// SetVersion sets the application's version string.
func (app *BaseApp) SetVersion(v string) {
	app.version = v
}

// SetAppVersion sets the application's version this is used as part of the
// header in blocks and is returned to the consensus engine in EndBlock.
func (app *BaseApp) SetAppVersion(ctx context.Context, v uint64) error {
	if app.paramStore == nil {
		return errors.New("param store must be set to set app version")
	}

	cp, err := app.paramStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get consensus params: %w", err)
	}
	if cp.Version == nil {
		return errors.New("version is not set in param store")
	}
	cp.Version.App = v
	if err := app.paramStore.Set(ctx, cp); err != nil {
		return err
	}
	return nil
}

func (app *BaseApp) SetDB(db dbm.DB) {
	app.db = db
}

func (app *BaseApp) SetCMS(cms storetypes.CommitMultiStore) {
	app.cms = cms
}

func (app *BaseApp) SetInitChainer(initChainer sdk.InitChainer) {
	app.initChainer = initChainer
}

func (app *BaseApp) PreBlocker() sdk.PreBlocker {
	return app.preBlocker
}

func (app *BaseApp) SetPreBlocker(preBlocker sdk.PreBlocker) {
	app.preBlocker = preBlocker
}

func (app *BaseApp) SetBeginBlocker(beginBlocker sdk.BeginBlocker) {
	app.beginBlocker = beginBlocker
}

func (app *BaseApp) SetEndBlocker(endBlocker sdk.EndBlocker) {
	app.endBlocker = endBlocker
}

func (app *BaseApp) SetPrepareCheckStater(
	prepareCheckStater sdk.PrepareCheckStater,
) {
	app.prepareCheckStater = prepareCheckStater
}

func (app *BaseApp) SetPrecommiter(precommiter sdk.Precommiter) {
	app.precommiter = precommiter
}

// SetStoreLoader allows us to customize the rootMultiStore initialization.
func (app *BaseApp) SetStoreLoader(loader StoreLoader) {
	app.storeLoader = loader
}

// SetProcessProposal sets the process proposal function for the BaseApp.
func (app *BaseApp) SetProcessProposal(handler sdk.ProcessProposalHandler) {
	app.processProposal = handler
}

// SetPrepareProposal sets the prepare proposal function for the BaseApp.
func (app *BaseApp) SetPrepareProposal(handler sdk.PrepareProposalHandler) {
	app.prepareProposal = handler
}
