package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/store"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/types"
)

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) error {
	genesisStore := store.NewGenesis(ctx.KVStore(k.storeKey))
	genesisStore.Store(data.Eth1GenesisHash)
	return nil
}
