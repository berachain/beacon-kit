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
	"slices"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// appModule defines runtime as an AppModule.
type appModule struct {
	app *App
}

func (m appModule) IsOnePerModuleType() {}
func (m appModule) IsAppModule()        {}

var (
	_ appmodule.AppModule = appModule{}
)

func init() {
	appconfig.RegisterModule(&runtimev1alpha1.Module{},
		appconfig.Provide(
			ProvideApp,
			codec.ProvideAddressCodec,
			ProvideKVStoreKey,
			ProvideEnvironment,
			ProvideModuleManager,
		),
		appconfig.Invoke(SetupAppBuilder),
	)
}

func ProvideApp(middleware Middleware) (
	*AppBuilder,
	appmodule.AppModule,
	error,
) {
	app := &App{Middleware: middleware}
	return &AppBuilder{app: app}, appModule{app}, nil
}

type AppInputs struct {
	depinject.In

	Logger        log.Logger
	Config        *runtimev1alpha1.Module
	AppBuilder    *AppBuilder
	ModuleManager *module.Manager
}

func SetupAppBuilder(inputs AppInputs) {
	app := inputs.AppBuilder.app
	app.config = inputs.Config
	app.logger = inputs.Logger
	// app.ModuleManager = inputs.ModuleManager
}

func registerStoreKey(wrapper *AppBuilder, key storetypes.StoreKey) {
	wrapper.app.storeKeys = append(wrapper.app.storeKeys, key)
}

func ProvideKVStoreKey(
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) *storetypes.KVStoreKey {
	if slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		return nil
	}

	storeKey := storetypes.NewKVStoreKey(key.Name())
	registerStoreKey(app, storeKey)
	return storeKey
}

func ProvideModuleManager(
	modules map[string]appmodule.AppModule,
) *module.Manager {
	return module.NewManagerFromMap(modules)
}

func ProvideEnvironment(
	logger log.Logger,
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) (store.KVStoreService, appmodule.Environment) {
	var (
		kvService store.KVStoreService = failingStoreService{}
	)

	// skips modules that have no store
	if !slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		storeKey := ProvideKVStoreKey(config, key, app)
		kvService = kvStoreService{key: storeKey}

	}

	return kvService, appmodule.Environment{
		Logger: logger.With(
			log.ModuleKey,
			fmt.Sprintf("x/%s", key.Name()),
		),
		KVStoreService: kvService,
	}

}
