package evm

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	store "cosmossdk.io/store/types"

	"github.com/itsdevbear/bolaris/beacon/execution"
	modulev1alpha1 "github.com/itsdevbear/bolaris/cosmos/api/polaris/evm/module/v1alpha1"
	"github.com/itsdevbear/bolaris/cosmos/x/evm/keeper"
)

//nolint:gochecknoinits // GRRRR fix later.
func init() {
	appmodule.Register(&modulev1alpha1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	ModuleKey depinject.OwnModuleKey
	Config    *modulev1alpha1.Module
	Key       *store.KVStoreKey

	ExecutionClient execution.EngineCaller
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

// ProvideModule is a function that provides the module to the application.
func ProvideModule(in DepInjectInput) DepInjectOutput {
	k := keeper.NewKeeper(
		in.ExecutionClient,
		in.Key,
	)
	m := NewAppModule(k)

	return DepInjectOutput{
		Keeper: k,
		Module: m,
	}
}
