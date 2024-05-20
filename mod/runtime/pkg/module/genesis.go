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

package beacon

import (
	"context"
	"encoding/json"
	"fmt"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// DefaultGenesis returns default genesis state as raw bytes
// for the beacon module.
func (AppModule) DefaultGenesis() json.RawMessage {
	defaultGenesis, err := deneb.DefaultBeaconState()
	if err != nil {
		panic(err)
	}

	bz, err := json.Marshal(defaultGenesis)
	if err != nil {
		panic(err)
	}
	return bz
}

// ValidateGenesis performs genesis state validation for the evm module.
func (AppModule) ValidateGenesis(
	bz json.RawMessage,
) error {
	data := new(deneb.BeaconState)
	if err := json.Unmarshal(bz, data); err != nil {
		return err
	}

	seenValidators := make(map[[48]byte]struct{})
	for _, validator := range data.Validators {
		if _, ok := seenValidators[validator.Pubkey]; ok {
			return fmt.Errorf(
				"duplicate pubkey in genesis state: %x",
				validator.Pubkey,
			)
		}
		seenValidators[validator.Pubkey] = struct{}{}
	}
	return nil
}

// InitGenesis performs genesis initialization for the beacon module.
func (am AppModule) InitGenesis(
	ctx context.Context,
	bz json.RawMessage,
) ([]appmodulev2.ValidatorUpdate, error) {
	data := new(deneb.BeaconState)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, err
	}

	// Load the store.
	store := am.keeper.BeaconStore().WithContext(ctx)
	sdb := state.NewBeaconStateFromDB(store, am.chainSpec)
	if err := sdb.WriteGenesisStateDeneb(data); err != nil {
		return nil, err
	}

	// Build ValidatorUpdates for CometBFT.
	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	blsType := "bls12_381"
	for _, validator := range data.Validators {
		validatorUpdates = append(validatorUpdates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: blsType,
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		},
		)
	}
	return validatorUpdates, nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(
	_ context.Context,
) (json.RawMessage, error) {
	return json.Marshal(&deneb.BeaconState{})
}
