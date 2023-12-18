package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/cosmos/x/evm/store"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	if err := genesisStore.Store(data.Eth1GenesisHash); err != nil {
		panic(err)
	}
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	return &types.GenesisState{
		Eth1GenesisHash: genesisStore.Retrieve().Hex(),
	}
}
