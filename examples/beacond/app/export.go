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
	"encoding/json"
	"log"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/staking"
	stakingtypes "cosmossdk.io/x/staking/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExportAppStateAndValidators exports the state of the application for a
// genesis
// file.
func (app *BeaconApp) ExportAppStateAndValidators(
	forZeroHeight bool,
	jailAllowedAddrs, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContextLegacy(
		true,
		cmtproto.Header{Height: app.LastBlockHeight()},
	)

	// We export at last height + 1, because that's the height at which
	// CometBFT will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState, err := app.ModuleManager.ExportGenesisForModules(
		ctx,
		modulesToExport,
	)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//
//	in favor of export at a block height
func (app *BeaconApp) prepForZeroHeightGenesis(
	ctx sdk.Context,
	jailAllowedAddrs []string,
) {
	applyAllowedAddrs := false

	// check if there is a allowed address list
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		allowedAddrsMap[addr] = true
	}

	/* Handle fee distribution state. */

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	err := app.StakingKeeper.IterateRedelegations(
		ctx,
		func(_ int64, red stakingtypes.Redelegation) (stop bool) {
			for i := range red.Entries {
				red.Entries[i].CreationHeight = 0
			}
			err := app.StakingKeeper.SetRedelegation(ctx, red)
			if err != nil {
				panic(err)
			}
			return false
		},
	)
	if err != nil {
		panic(err)
	}

	// iterate through unbonding delegations, reset creation height
	err = app.StakingKeeper.UnbondingDelegations.Walk(
		ctx,
		nil,
		func(key collections.Pair[[]byte, []byte], ubd stakingtypes.UnbondingDelegation) (stop bool, err error) {
			for i := range ubd.Entries {
				ubd.Entries[i].CreationHeight = 0
			}
			err = app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
			if err != nil {
				return true, err
			}
			return false, err
		},
	)
	if err != nil {
		panic(err)
	}
	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.kvStoreKeys()[stakingtypes.StoreKey])
	iter := storetypes.KVStoreReversePrefixIterator(
		store,
		stakingtypes.ValidatorsKey,
	)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(
			stakingtypes.AddressFromValidatorsKey(iter.Key()),
		)
		validator, err := app.StakingKeeper.GetValidator(ctx, addr)
		if err != nil {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		if err = app.StakingKeeper.SetValidator(ctx, validator); err != nil {
			panic(err)
		}
		counter++
	}

	if err := iter.Close(); err != nil {
		app.Logger().
			Error("error while closing the key-value store reverse prefix iterator: ", err)
		return
	}

	_, err = app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
