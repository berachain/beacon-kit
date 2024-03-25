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

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
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
		panic("not supported, just look at the genesis file u goofy")
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

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      nil,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}
