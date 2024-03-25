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

package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	"cosmossdk.io/depinject/appconfig"
	authtypes "cosmossdk.io/x/auth/types"
	banktypes "cosmossdk.io/x/bank/types"
	"github.com/berachain/beacon-kit/runtime/modules/beacon"
	beaconv1alpha1 "github.com/berachain/beacon-kit/runtime/modules/beacon/api/module/v1alpha1"
	"github.com/cosmos/cosmos-sdk/runtime"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	_ "cosmossdk.io/x/auth"
	_ "cosmossdk.io/x/auth/tx/config"            // import for side-effects
	_ "cosmossdk.io/x/bank"                      // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
)

const AppName = "BeaconKitApp"

var (
	// module account permissions
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
	}

	// application configuration (used by depinject)
	BeaconAppConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: runtime.ModuleName,
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName:       AppName,
					PreBlockers:   []string{},
					BeginBlockers: []string{},
					EndBlockers: []string{
						beacon.ModuleName,
					},
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					InitGenesis: []string{
						authtypes.ModuleName,
						banktypes.ModuleName,
						beacon.ModuleName,
					},
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             "cosmos",
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module.
					// This is configurable with the following: Authority:
					// "group", // A custom module authority can be set using a
					// module name Authority:
					// "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a
					// specific address
				}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   beacon.ModuleName,
				Config: appconfig.WrapAny(&beaconv1alpha1.Module{}),
			},
		},
	})
)
