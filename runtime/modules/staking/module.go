package staking

import (
	"cosmossdk.io/x/staking"
	"cosmossdk.io/x/staking/keeper"
	"cosmossdk.io/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// AppModule implements an application module for the staking module.
type AppModule struct {
	staking.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) AppModule {
	return AppModule{
		AppModule: staking.NewAppModule(cdc, keeper, ak, bk),
	}
}
