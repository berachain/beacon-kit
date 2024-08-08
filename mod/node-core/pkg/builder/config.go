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
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	beacon "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module"
	beaconv1alpha1 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v1alpha1"
)

// DefaultDepInjectConfig returns the default configuration for the dependency
// injection framework.
func DefaultDepInjectConfig() depinject.Config {
	return depinject.Configs(
		appconfig.Compose(&appv1alpha1.Config{
			Modules: []*appv1alpha1.ModuleConfig{
				{
					Name: "runtime",
					Config: appconfig.WrapAny(&runtimev1alpha1.Module{
						AppName:       "BeaconKit",
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
			},
		}),
	)
}
