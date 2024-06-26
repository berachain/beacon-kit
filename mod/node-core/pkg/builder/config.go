// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package builder

import (
	runtimev2 "cosmossdk.io/api/cosmos/app/runtime/v2"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	beacon "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module"
	modulev2 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v2"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// DefaultDepInjectConfig returns the default configuration for the dependency
// injection framework.
func DefaultDepInjectConfig() depinject.Config {
	addrCdc := addresscodec.NewBech32Codec("bera")
	return depinject.Configs(
		appconfig.Compose(&appv1alpha1.Config{
			Modules: []*appv1alpha1.ModuleConfig{
				{
					Name: runtime.ModuleName,
					Config: appconfig.WrapAny(&runtimev2.Module{
						AppName:       DefaultAppName,
						PreBlockers:   []string{},
						BeginBlockers: []string{},
						EndBlockers:   []string{beacon.ModuleName},
						InitGenesis:   []string{beacon.ModuleName},
						GasConfig: &runtimev2.GasConfig{
							ValidateTxGasLimit: 100_000,
							QueryGasLimit:      100_000,
							SimulationGasLimit: 100_000,
						},
					}),
				},
				{
					Name:   beacon.ModuleName,
					Config: appconfig.WrapAny(&modulev2.Module{}),
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
