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
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	_ "cosmossdk.io/x/auth/tx/config" // import for side-effects
	"github.com/berachain/beacon-kit/beacond/x/beacon"
	beaconv1alpha1 "github.com/berachain/beacon-kit/beacond/x/beacon/api/module/v1alpha1"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
)

const AppName = "BeaconKitApp"

// Config returns the default app config.
func Config() depinject.Config {
	addrCdc := addresscodec.NewBech32Codec("bera")
	return depinject.Configs(
		appconfig.Compose(&appv1alpha1.Config{
			Modules: []*appv1alpha1.ModuleConfig{
				{
					Name: runtime.ModuleName,
					Config: appconfig.WrapAny(&runtimev1alpha1.Module{
						AppName:       AppName,
						PreBlockers:   []string{},
						BeginBlockers: []string{},
						EndBlockers:   []string{beacon.ModuleName},
						InitGenesis:   []string{beacon.ModuleName},
					}),
				},
				{
					Name:   beacon.ModuleName,
					Config: appconfig.WrapAny(&beaconv1alpha1.Module{}),
				},
				{
					Name:   consensustypes.ModuleName,
					Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
				},
				{
					Name: "tx",
					Config: appconfig.WrapAny(&txconfigv1.Config{
						SkipAnteHandler: true,
					}),
				},
			},
		}),
		depinject.Supply(
			func() address.Codec { return addrCdc },
			func() address.ValidatorAddressCodec { return addrCdc },
			func() address.ConsensusAddressCodec { return addrCdc },
		),
	)
}
