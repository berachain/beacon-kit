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

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"cosmossdk.io/core/registry"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	// ConsensusVersion defines the current x/beacon module consensus version.
	ConsensusVersion = 1
	// ModuleName is the module name constant used in many places.
	ModuleName = "beacon"
)

var (
	_ appmodulev2.AppModule  = AppModule{}
	_ module.HasABCIGenesis  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}
)

// AppModule implements an application module for the evm module.
type AppModule struct {
	keeper *storage.Backend[
		*dastore.Store[types.BeaconBlockBody],
		state.BeaconState,
	]
	chainSpec primitives.ChainSpec
}

// NewAppModule creates a new AppModule object.
func NewAppModule(
	keeper *storage.Backend[
		*dastore.Store[types.BeaconBlockBody], state.BeaconState,
	],
	chainSpec primitives.ChainSpec,
) AppModule {
	return AppModule{
		keeper:    keeper,
		chainSpec: chainSpec,
	}
}

// Name is the name of this module.
func (am AppModule) Name() string {
	return ModuleName
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterInterfaces registers the module's interface types.
func (am AppModule) RegisterInterfaces(registry.InterfaceRegistrar) {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// EndBlock returns the validator set updates from the beacon state.
func (am AppModule) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	store := am.keeper.BeaconStore().WithContext(ctx)

	// Get the public key of the validator
	val, err := store.GetValidatorsByEffectiveBalance()
	if err != nil {
		panic(err)
	}

	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	for _, validator := range val {
		// TODO: Config
		// Max 100 validators in the active set.
		// TODO: this is kinda hood.
		if validator.EffectiveBalance == 0 {
			var idx math.ValidatorIndex
			idx, err = store.WithContext(ctx).
				ValidatorIndexByPubkey(validator.Pubkey)
			if err != nil {
				return nil, err
			}
			if err = store.WithContext(ctx).
				RemoveValidatorAtIndex(idx); err != nil {
				return nil, err
			}
		}

		// TODO: this works, but there is a bug where if we send a validator to
		// 0 voting power, it can somehow still propose the next block? This
		// feels big bad.
		validatorUpdates = append(validatorUpdates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: "bls12_381",
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}

	// Save the store.
	store.Save()
	return validatorUpdates, nil
}
