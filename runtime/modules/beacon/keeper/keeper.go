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

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper/store"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/types"
)

// Keeper maintains the link to data storage and exposes access to the underlying
// `BeaconState` methods for the x/beacon module.
type Keeper struct {
	storeKey  storetypes.StoreKey
	beaconCfg *config.Beacon
}

// Assert Keeper implements BeaconStateProvider interface.
var _ state.BeaconStateProvider = &Keeper{}

// NewKeeper creates new instances of the Beacon Keeper.
func NewKeeper(
	storeKey storetypes.StoreKey,
	beaconCfg *config.Beacon,
) *Keeper {
	return &Keeper{
		storeKey:  storeKey,
		beaconCfg: beaconCfg,
	}
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key for the x/beacon module.
func (k *Keeper) BeaconState(ctx context.Context) state.BeaconState {
	return store.NewBeaconStore(
		ctx,
		k.storeKey,
		k.beaconCfg,
	)
}

// InitGenesis initializes the genesis state of the beacon module.
func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	beaconState := k.BeaconState(ctx)
	hash := common.HexToHash(data.Eth1GenesisHash)

	// At genesis, we assume that the genesis block is also safe and final.
	beaconState.SetGenesisEth1Hash(hash)
	beaconState.SetSafeEth1BlockHash(hash)
	beaconState.SetFinalizedEth1BlockHash(hash)
}

// ExportGenesis exports the current state of the beacon module as genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Eth1GenesisHash: k.BeaconState(ctx).GenesisEth1Hash().Hex(),
	}
}
