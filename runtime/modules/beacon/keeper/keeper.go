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
	"github.com/berachain/beacon-kit/beacon/core/state"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/beacon/forkchoice/ssf"
	filedb "github.com/berachain/beacon-kit/db/file"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	beaconstore "github.com/berachain/beacon-kit/store/beacon"
	"github.com/berachain/beacon-kit/store/blob"
	forkchoicestore "github.com/berachain/beacon-kit/store/forkchoice"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/ethereum/go-ethereum/common"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
type Keeper struct {
	availabilityStore *blob.Store
	beaconStore       *beaconstore.Store
	forkchoiceStore   *forkchoicestore.Store
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
) *Keeper {
	return &Keeper{
		availabilityStore: blob.NewStore(fdb),
		beaconStore:       beaconstore.NewStore(env),
		forkchoiceStore:   forkchoicestore.NewStore(env.KVStoreService),
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
) ([]abci.ValidatorUpdate, error) {
	store := k.beaconStore.WithContext(ctx)
	// Get the public key of the validator
	val, err := k.beaconStore.GetValidatorsByEffectiveBalance(ctx, 100)
	if err != nil {
		panic(err)
	}

	validatorUpdates := make([]abci.ValidatorUpdate, 0)
	for _, validator := range val {
		pk := crypto.PublicKey{
			Sum: &crypto.PublicKey_Bls12381{Bls12381: validator.Pubkey[:]},
		}

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
		validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
			PubKey: pk,
			//#nosec:G701 // will not realistically cause a problem.
			Power: int64(validator.EffectiveBalance),
		})
	}
	return validatorUpdates, nil
}

// AvailabilityStore returns the availability store struct initialized with a.
func (k *Keeper) AvailabilityStore(
	_ context.Context,
) state.AvailabilityStore {
	return k.availabilityStore
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k *Keeper) BeaconState(
	ctx context.Context,
) state.BeaconState {
	return k.beaconStore.WithContext(ctx)
}

// context and the store key.
//
// TODO: Decouple from the Specific SingleSlotFinalityStore Impl.
func (k *Keeper) ForkchoiceStore(
	ctx context.Context,
) ssf.SingleSlotFinalityStore {
	return k.forkchoiceStore.WithContext(ctx)
}

// InitGenesis initializes the genesis state of the module.
func (k *Keeper) InitGenesis(
	ctx context.Context,
	data types.GenesisState,
) ([]abci.ValidatorUpdate, error) {
	// Set the genesis RANDAO mix.
	st := k.BeaconState(ctx)
	if err := st.UpdateRandaoMixAtIndex(0, data.Mix()); err != nil {
		return nil, err
	}

	// Compare this snippet from beacon/keeper/keeper.go:
	if err := st.UpdateStateRootAtIndex(0, [32]byte{}); err != nil {
		return nil, err
	}

	// Set the genesis block root.
	if err := st.UpdateBlockRootAtIndex(0, [32]byte{}); err != nil {
		return nil, err
	}

	// Set the genesis block data.
	fcs := k.ForkchoiceStore(ctx)
	hash := common.HexToHash(data.Eth1GenesisHash)
	fcs.SetGenesisEth1Hash(hash)
	fcs.SetSafeEth1BlockHash(hash)
	fcs.SetFinalizedEth1BlockHash(hash)

	// Set the genesis block header.
	if err := st.SetLatestBlockHeader(
		&beacontypes.BeaconBlockHeader{
			Slot:          0,
			ProposerIndex: 0,
			ParentRoot:    [32]byte{},
			StateRoot:     [32]byte{},
			BodyRoot:      [32]byte{},
		},
	); err != nil {
		return nil, err
	}

	// TODO: don't need to set any validators here if we are setting in
	// EndBlock. TODO: we should only do updates in EndBlock and actually do the
	// full initial update here.
	return nil, nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	return &types.GenesisState{
		Eth1GenesisHash: k.ForkchoiceStore(ctx).GenesisEth1Hash().Hex(),
	}
}
