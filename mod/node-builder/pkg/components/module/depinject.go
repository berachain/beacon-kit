// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package beacon

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components"
	modulev1alpha1 "github.com/berachain/beacon-kit/mod/node-builder/pkg/components/module/api/module/v1alpha1"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
)

// TODO: we don't allow generics here? Why? Is it fixable?
//
//nolint:gochecknoinits // required by sdk.
func init() {
	appconfig.RegisterModule(&modulev1alpha1.Module{},
		appconfig.Provide(
			ProvideModule,
		),
	)
}

// TODO: DEPINJECT RUNTIME AND REMOVE UNNECESSARY FIELDS
// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In
	BeaconConfig *config.Config
	BKRuntime    *components.BeaconKitRuntime
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out
	Module appmodule.AppModule
}

// ProvideModule is a function that provides the module to the application.
func ProvideModule(in DepInjectInput) (DepInjectOutput, error) {
	// TODO: this is hood as fuck.
	// if in.BeaconConfig.KZG.Implementation == "" {
	// 	in.BeaconConfig.KZG.Implementation = "crate-crypto/go-kzg-4844"
	// }

	// runtime, err := components.ProvideRuntime(
	// 	in.ChainSpec,
	// 	in.StorageBackend,
	// 	in.TelemetrySink,
	// 	in.Environment.Logger.With("module", "beacon-kit"),
	// 	in.ServiceRegistry,
	// )
	// if err != nil {
	// 	return DepInjectOutput{}, err
	// }

	return DepInjectOutput{
		Module: NewAppModule(in.BKRuntime),
	}, nil
}
