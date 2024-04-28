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

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/core"
)

// DefaultGenesis returns default genesis state as raw bytes
// for the beacon module.
func (AppModule) DefaultGenesis() json.RawMessage {
	bz, err := json.Marshal(core.DefaultGenesis())
	if err != nil {
		panic(err)
	}
	return bz
}

// ValidateGenesis performs genesis state validation for the evm module.
func (AppModule) ValidateGenesis(
	_ json.RawMessage,
) error {
	// TODO: implement.

	// IsValidGenesisState gets called whenever there's a deposit event,
	// // it checks whether there's enough effective balance to trigger and
	// // if the minimum genesis time arrived already.
	// //
	// // Spec pseudocode definition:
	// //
	// //	def is_valid_genesis_state(state: BeaconState) -> bool:
	// //	   if state.genesis_time < MIN_GENESIS_TIME:
	// //	       return False
	// //	   if len(get_active_validator_indices(state, GENESIS_EPOCH)) <
	// MIN_GENESIS_ACTIVE_VALIDATOR_COUNT:
	// //	       return False
	// //	   return True
	// //
	// // This method has been modified from the spec to allow whole states not
	// to be saved
	// // but instead only cache the relevant information.
	// func IsValidGenesisState(chainStartDepositCount, currentTime uint64) bool
	// {
	// 	if currentTime < params.BeaconConfig().MinGenesisTime {
	// 		return false
	// 	}
	// 	if chainStartDepositCount <
	// params.BeaconConfig().MinGenesisActiveValidatorCount {
	// 		return false
	// 	}
	// 	return true
	// }

	return nil
}

// InitGenesis performs genesis initialization for the beacon module.
func (am AppModule) InitGenesis(
	ctx context.Context,
	bz json.RawMessage,
) ([]appmodulev2.ValidatorUpdate, error) {
	gs := new(core.Genesis)
	if err := json.Unmarshal(bz, gs); err != nil {
		return nil, err
	}
	return am.keeper.InitGenesis(ctx, gs)
}

// ExportGenesis returns the exported genesis state as raw bytes for the evm
// module.
func (am AppModule) ExportGenesis(
	ctx context.Context,
) (json.RawMessage, error) {
	return json.Marshal(am.keeper.ExportGenesis(ctx))
}
