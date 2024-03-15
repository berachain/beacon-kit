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
	"github.com/berachain/beacon-kit/beacon/forkchoice/ssf"
	"github.com/berachain/beacon-kit/runtime"
	"github.com/berachain/beacon-kit/runtime/modules/beacon/types"
	beaconstore "github.com/berachain/beacon-kit/store/beacon"
	forkchoicestore "github.com/berachain/beacon-kit/store/forkchoice"
	"github.com/ethereum/go-ethereum/common"
)

// Keeper maintains the link to data storage and exposes access to the
// underlying `BeaconState` methods for the x/beacon module.
type Keeper struct {
	beaconStore     *beaconstore.Store
	forkchoiceStore *forkchoicestore.Store
	vsu             runtime.ValsetUpdater
}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	env appmodule.Environment,
	vsu runtime.ValsetUpdater,
) *Keeper {
	return &Keeper{
		beaconStore:     beaconstore.NewStore(env.KVStoreService),
		forkchoiceStore: forkchoicestore.NewStore(env.KVStoreService),
		vsu:             vsu,
	}
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
) error {
	// Set the genesis RANDAO mix.
	st := k.BeaconState(ctx)
	if err := st.SetRandaoMix(data.Mix()); err != nil {
		return err
	}

	// Set the genesis block data.
	fcs := k.ForkchoiceStore(ctx)
	hash := common.HexToHash(data.Eth1GenesisHash)
	fcs.SetGenesisEth1Hash(hash)
	fcs.SetSafeEth1BlockHash(hash)
	fcs.SetFinalizedEth1BlockHash(hash)
	return nil
}

// ExportGenesis exports the current state of the module as genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	return &types.GenesisState{
		Eth1GenesisHash: k.ForkchoiceStore(ctx).GenesisEth1Hash().Hex(),
	}
}
