// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/cosmos/x/beacon/store"
	"github.com/itsdevbear/bolaris/cosmos/x/beacon/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	if err := genesisStore.Store(data.Eth1GenesisHash); err != nil {
		panic(err)
	}

	hash := common.HexToHash(data.Eth1GenesisHash)

	fcs := store.NewForkchoice(ctx.KVStore(k.storeKey))
	fcs.SetSafeEth1BlockHash(hash)
	fcs.SetFinalizedEth1BlockHash(hash)
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	return &types.GenesisState{
		Eth1GenesisHash: genesisStore.Retrieve().Hex(),
	}
}
