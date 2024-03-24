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
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	"github.com/berachain/beacon-kit/beacon/core/state"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/beacon/forkchoice/ssf"
	filedb "github.com/berachain/beacon-kit/db/file"
	"github.com/berachain/beacon-kit/runtime"
	beaconstore "github.com/berachain/beacon-kit/store/beacon"
	"github.com/berachain/beacon-kit/store/blob"
	forkchoicestore "github.com/berachain/beacon-kit/store/forkchoice"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
type Keeper struct {
	availabilityStore *blob.Store
	beaconStore       *beaconstore.Store
	forkchoiceStore   *forkchoicestore.Store
	vsu               runtime.ValsetUpdater
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	fdb *filedb.DB,
	env appmodule.Environment,
	vsu runtime.ValsetUpdater,
) *Keeper {
	return &Keeper{
		availabilityStore: blob.NewStore(fdb),
		beaconStore:       beaconstore.NewStore(env),
		forkchoiceStore:   forkchoicestore.NewStore(env.KVStoreService),
		vsu:               vsu,
	}
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
	data state.BeaconStateDeneb,
) error {
	// Set the genesis RANDAO mix.
	st := k.BeaconState(ctx)
	if err := st.UpdateRandaoMixAtIndex(0, types.Mix(data.RandaoMix)); err != nil {
		return err
	}

	// Compare this snippet from beacon/keeper/keeper.go:
	if err := st.UpdateStateRootAtIndex(0, [32]byte{}); err != nil {
		return err
	}

	// Set the genesis block root.
	if err := st.UpdateBlockRootAtIndex(0, [32]byte{}); err != nil {
		return err
	}

	// Set the genesis block data.
	fcs := k.ForkchoiceStore(ctx)
	fcs.SetGenesisEth1Hash(data.Eth1GenesisHash)
	fcs.SetSafeEth1BlockHash(data.Eth1GenesisHash)
	fcs.SetFinalizedEth1BlockHash(data.Eth1GenesisHash)

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
		return err
	}

	return nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) *state.BeaconStateDeneb {
	return &state.BeaconStateDeneb{
		Eth1GenesisHash: k.ForkchoiceStore(ctx).GenesisEth1Hash(),
	}
}
