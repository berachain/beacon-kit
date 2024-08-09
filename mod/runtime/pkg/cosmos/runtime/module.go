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

package runtime

import (
	"fmt"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
)

func init() {
	appconfig.RegisterModule(&runtimev1alpha1.Module{},
		appconfig.Provide(
			ProvideApp,
			ProvideKVStoreKey,
			ProvideEnvironment,
		),
		appconfig.Invoke(SetupAppBuilder),
	)
}

func ProvideApp(middleware Middleware) (
	*AppBuilder,
	error,
) {
	app := &App{Middleware: middleware}
	return &AppBuilder{app: app}, nil
}

type AppInputs struct {
	depinject.In

	Logger     log.Logger
	AppBuilder *AppBuilder
}

func SetupAppBuilder(inputs AppInputs) {
	app := inputs.AppBuilder.app
	app.logger = inputs.Logger
}

func ProvideKVStoreKey(
	app *AppBuilder,
) *storetypes.KVStoreKey {
	storeKey := storetypes.NewKVStoreKey("beacon")
	app.app.storeKeys = append(app.app.storeKeys, storeKey)
	return storeKey
}

func ProvideEnvironment(
	logger log.Logger,
	app *AppBuilder,
) (store.KVStoreService, appmodule.Environment) {
	var (
		kvService store.KVStoreService = failingStoreService{}
	)

	// skips modules that have no store
	storeKey := ProvideKVStoreKey(app)
	kvService = kvStoreService{key: storeKey}

	return kvService, appmodule.Environment{
		Logger: logger.With(
			log.ModuleKey,
			fmt.Sprintf("x/%s", "beacon"),
		),
		KVStoreService: kvService,
	}

}
