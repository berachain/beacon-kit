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

package keeper

import (
	"context"

	"cosmossdk.io/core/appmodule"
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/state/deneb"
	"github.com/berachain/beacon-kit/mod/da"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/storage/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/deposit"
	filedb "github.com/berachain/beacon-kit/mod/storage/filedb"
	bls12381 "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
type Keeper struct {
	availabilityStore *da.Store
	beaconStore       *beacondb.KVStore[
		*primitives.Fork,
		*primitives.BeaconBlockHeader,
		engineprimitives.ExecutionPayload,
		*primitives.Eth1Data,
		*primitives.Validator,
	]
	depositStore *deposit.KVStore
	cfg          primitives.ChainSpec
}

// TODO: move this.
func DenebPayloadFactory() engineprimitives.ExecutionPayload {
	return &engineprimitives.ExecutableDataDeneb{}
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
	cfg primitives.ChainSpec,
	ddb *deposit.KVStore,
) *Keeper {
	return &Keeper{
		availabilityStore: da.NewStore(cfg, fdb),
		beaconStore: beacondb.New[
			*primitives.Fork,
			*primitives.BeaconBlockHeader,
			engineprimitives.ExecutionPayload,
			*primitives.Eth1Data,
			*primitives.Validator,
		](env.KVStoreService, DenebPayloadFactory),
		cfg:          cfg,
		depositStore: ddb,
	}
}

// ApplyAndReturnValidatorSetUpdates returns the validator set updates from
// the beacon state.
//
// TODO: this function is horribly inefficient and should be replaced with a
// more efficient implementation, that does not update the entire
// valset every block.
func (k *Keeper) ApplyAndReturnValidatorSetUpdates(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	store := k.beaconStore.WithContext(ctx)
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
			PubKeyType: (&bls12381.PubKey{}).Type(),
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}

	// Save the store.
	store.Save()
	return validatorUpdates, nil
}

// AvailabilityStore returns the availability store struct initialized with a.
func (k *Keeper) AvailabilityStore(
	_ context.Context,
) core.AvailabilityStore {
	return k.availabilityStore
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k *Keeper) BeaconState(
	ctx context.Context,
) state.BeaconState {
	return state.NewBeaconStateFromDB(k.beaconStore.WithContext(ctx), k.cfg)
}

// DepositStore returns the deposit store struct initialized with a.
func (k *Keeper) DepositStore(
	_ context.Context,
) *deposit.KVStore {
	return k.depositStore
}

// InitGenesis initializes the genesis state of the module.
func (k *Keeper) InitGenesis(
	ctx context.Context,
	data *deneb.BeaconState,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Load the store.
	store := k.beaconStore.WithContext(ctx)
	sdb := state.NewBeaconStateFromDB(store, k.cfg)
	if err := sdb.WriteGenesisStateDeneb(data); err != nil {
		return nil, err
	}

	// Build ValidatorUpdates for CometBFT.
	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	blsType := (&bls12381.PubKey{}).Type()
	for _, validator := range data.Validators {
		validatorUpdates = append(validatorUpdates, appmodulev2.ValidatorUpdate{
			PubKey:     validator.Pubkey[:],
			PubKeyType: blsType,
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}
	return validatorUpdates, nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(_ context.Context) *deneb.BeaconState {
	return &deneb.BeaconState{}
}
