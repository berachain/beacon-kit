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
	beaconstore "github.com/berachain/beacon-kit/beacond/store/beacon"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/da"
	"github.com/berachain/beacon-kit/mod/primitives"
	filedb "github.com/berachain/beacon-kit/mod/storage/filedb"
	bls12381 "github.com/cosmos/cosmos-sdk/crypto/keys/bls12_381"
	"github.com/ethereum/go-ethereum/common"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
type Keeper struct {
	availabilityStore *da.Store
	beaconStore       *beaconstore.Store
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
) *Keeper {
	return &Keeper{
		availabilityStore: da.NewStore(fdb),
		beaconStore:       beaconstore.NewStore(env),
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
			var idx primitives.ValidatorIndex
			idx, err = store.WithContext(ctx).
				ValidatorIndexByPubkey(validator.Pubkey[:])
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
	return k.beaconStore.WithContext(ctx)
}

// InitGenesis initializes the genesis state of the module.
//
// TODO: This whole thing needs to be abstracted into mod/core/state
func (k *Keeper) InitGenesis(
	ctx context.Context,
	data state.BeaconStateDeneb,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Set the genesis RANDAO mix.
	st := k.BeaconState(ctx)
	for i, mix := range data.RandaoMixes {
		if err := st.UpdateRandaoMixAtIndex(
			//#nosec:G701 // will not cause a problem.
			uint64(i), mix,
		); err != nil {
			return nil, err
		}
	}

	// Compare this snippet from beacon/keeper/keeper.go:
	if err := st.UpdateStateRootAtIndex(0, [32]byte{}); err != nil {
		return nil, err
	}

	if err := st.SetSlot(0); err != nil {
		return nil, err
	}

	// Set the genesis block header.
	if err := st.UpdateEth1BlockHash(data.Eth1BlockHash); err != nil {
		return nil, err
	}

	emptyHeader := &beacontypes.BeaconBlockHeader{
		Slot:          0,
		ProposerIndex: 0,
		ParentRoot:    [32]byte{},
		StateRoot:     [32]byte{},
		BodyRoot:      [32]byte{},
	}

	// Set the genesis block header.
	if err := st.SetLatestBlockHeader(
		emptyHeader,
	); err != nil {
		return nil, err
	}

	headerHtr, err := emptyHeader.HashTreeRoot()
	if err != nil {
		return nil, err
	}

	// Set the genesis block root.
	if err = st.UpdateBlockRootAtIndex(0, headerHtr); err != nil {
		return nil, err
	}

	if err = st.SetEth1DepositIndex(0); err != nil {
		return nil, err
	}

	// TODO: don't need to set any validators here if we are setting in
	// EndBlock. TODO: we should only do updates in EndBlock and actually do the
	// full initial update here.

	store := k.beaconStore.WithContext(ctx)
	validatorUpdates := make([]appmodulev2.ValidatorUpdate, 0)
	for i, validator := range data.Validators {
		if err = store.AddValidator(validator); err != nil {
			return nil, err
		}

		if err = store.IncreaseBalance(
			primitives.ValidatorIndex(i), validator.EffectiveBalance,
		); err != nil {
			return nil, err
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

	// Set the genesis slashing data.
	for i := range data.Slashings {
		//#nosec:G701 // will not realistically cause a problem.
		if err = store.UpdateSlashingAtIndex(uint64(i), 0); err != nil {
			return nil, err
		}
	}

	if err = store.SetGenesisValidatorsRoot(
		data.GenesisValidatorsRoot,
	); err != nil {
		return nil, err
	}
	return validatorUpdates, nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(_ context.Context) *state.BeaconStateDeneb {
	return &state.BeaconStateDeneb{
		Eth1BlockHash: common.Hash{},
	}
}
