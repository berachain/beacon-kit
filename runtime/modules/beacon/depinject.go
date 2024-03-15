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

package evm

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	stakingtypes "cosmossdk.io/x/staking/types"
	modulev1alpha1 "github.com/berachain/beacon-kit/runtime/modules/beacon/api/module/v1alpha1"
	"github.com/berachain/beacon-kit/runtime/modules/beacon/keeper"
	stakingwrapper "github.com/berachain/beacon-kit/runtime/modules/staking/keeper"
)

//nolint:gochecknoinits // GRRRR fix later.
func init() {
	appconfig.RegisterModule(&modulev1alpha1.Module{},
		appconfig.Provide(ProvideModule),
		appconfig.Provide(stakingwrapper.ProvideStakingKeeper),
	)
}

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	Config        *modulev1alpha1.Module
	Environment   appmodule.Environment
	ValsetUpdater *stakingwrapper.Keeper
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
	Hooks  stakingtypes.StakingHooksWrapper
}

// ProvideModule is a function that provides the module to the application.
func ProvideModule(in DepInjectInput) DepInjectOutput {
	k := keeper.NewKeeper(
		in.Environment,
		in.ValsetUpdater,
	)

	return DepInjectOutput{
		Keeper: k,
		Module: NewAppModule(k),
		Hooks: stakingtypes.StakingHooksWrapper{
			StakingHooks: k.Hooks(),
		},
	}
}
