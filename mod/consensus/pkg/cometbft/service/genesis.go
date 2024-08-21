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

package cometbft

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	cmttypes "github.com/cometbft/cometbft/types"
)

// LoadHeight loads a particular height.
func (app *Service) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// DefaultGenesis returns the default genesis state for the application.
func (app *Service) DefaultGenesis() map[string]json.RawMessage {
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
func (app *Service) ValidateGenesis(
	_ map[string]json.RawMessage,
) error {
	// Implement the validation logic for the provided genesis state.
	// This should validate the genesis state for each module in the
	// application.
	return nil
}

// ExportAppStateAndValidators exports the state of the application for a
// genesis
// file.
func (app *Service) ExportAppStateAndValidators(
	forZeroHeight bool,
	_, _ []string,
) (servertypes.ExportedApp, error) {
	// We export at last height + 1, because that's the height at which
	// CometBFT will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		// height = 0
		panic("not supported, just look at the genesis file u goofy")
	}

	// genState, err := app.ModuleManager.ExportGenesisForModules(
	// 	ctx,
	// 	modulesToExport,
	// )
	// if err != nil {
	// 	return servertypes.ExportedApp{}, err
	// }

	// appState, err := json.MarshalIndent(genState, "", "  ")
	// if err != nil {
	// 	return servertypes.ExportedApp{}, err
	// }

	// TODO: Pull these in from the BeaconKeeper, should be easy.
	validators := []cmttypes.GenesisValidator(nil)

	return servertypes.ExportedApp{
		AppState:        nil,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.GetConsensusParams(context.TODO()),
	}, nil
}
