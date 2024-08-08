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
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/genesis"
	"cosmossdk.io/core/legacy"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// appModule defines runtime as an AppModule.
type appModule struct {
	app *App
}

func (m appModule) IsOnePerModuleType() {}
func (m appModule) IsAppModule()        {}

func (m appModule) RegisterServices(
	configurator module.Configurator,
) { // nolint:staticcheck // SA1019: Configurator is deprecated but still used in runtime v1.

}

func (m appModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"autocli": {
					Service: autocliv1.Query_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "AppOptions",
							Short:     "Query the custom autocli options",
						},
					},
				},
				"reflection": {
					Service: reflectionv1.ReflectionService_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "FileDescriptors",
							Short:     "Query the app's protobuf file descriptors",
						},
					},
				},
			},
		},
	}
}

var (
	_ appmodule.AppModule = appModule{}
	_ module.HasServices  = appModule{}
)

// BaseAppOption is a depinject.AutoGroupType which can be used to pass
// BaseApp options into the depinject. It should be used carefully.
type BaseAppOption func(*baseapp.BaseApp)

// IsManyPerContainerType indicates that this is a
// depinject.ManyPerContainerType.
func (b BaseAppOption) IsManyPerContainerType() {}

func init() {
	appconfig.RegisterModule(&runtimev1alpha1.Module{},
		appconfig.Provide(
			ProvideApp,
			// to decouple runtime from sdk/codec ProvideInterfaceReistry can be
			// registered from the app
			// i.e. in the call to depinject.Inject(...)
			codec.ProvideInterfaceRegistry,
			codec.ProvideLegacyAmino,
			codec.ProvideProtoCodec,
			codec.ProvideAddressCodec,
			ProvideKVStoreKey,
			ProvideTransientStoreKey,
			ProvideMemoryStoreKey,
			ProvideGenesisTxHandler,
			ProvideEnvironment,
			ProvideTransientStoreService,
			ProvideModuleManager,
		),
		appconfig.Invoke(SetupAppBuilder),
	)
}

func ProvideApp(
	interfaceRegistry codectypes.InterfaceRegistry,
) (
	*AppBuilder,
	appmodule.AppModule,
	error,
) {

	std.RegisterInterfaces(interfaceRegistry)

	app := &App{}
	appBuilder := &AppBuilder{app}

	return appBuilder, appModule{app}, nil
}

type AppInputs struct {
	depinject.In

	Logger            log.Logger
	Config            *runtimev1alpha1.Module
	AppBuilder        *AppBuilder
	ModuleManager     *module.Manager
	BaseAppOptions    []BaseAppOption
	InterfaceRegistry codectypes.InterfaceRegistry
	LegacyAmino       legacy.Amino
}

func SetupAppBuilder(inputs AppInputs) {
	app := inputs.AppBuilder.app
	app.baseAppOptions = inputs.BaseAppOptions
	app.config = inputs.Config
	app.logger = inputs.Logger
	app.ModuleManager = inputs.ModuleManager
	app.ModuleManager.RegisterInterfaces(inputs.InterfaceRegistry)
	app.ModuleManager.RegisterLegacyAminoCodec(inputs.LegacyAmino)
}

func registerStoreKey(wrapper *AppBuilder, key storetypes.StoreKey) {
	wrapper.app.storeKeys = append(wrapper.app.storeKeys, key)
}

func storeKeyOverride(
	config *runtimev1alpha1.Module,
	moduleName string,
) *runtimev1alpha1.StoreKeyConfig {
	for _, cfg := range config.GetOverrideStoreKeys() {
		if cfg.GetModuleName() == moduleName {
			return cfg
		}
	}
	return nil
}

func ProvideKVStoreKey(
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) *storetypes.KVStoreKey {
	if slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		return nil
	}

	override := storeKeyOverride(config, key.Name())

	var storeKeyName string
	if override != nil {
		storeKeyName = override.GetKvStoreKey()
	} else {
		storeKeyName = key.Name()
	}

	storeKey := storetypes.NewKVStoreKey(storeKeyName)
	registerStoreKey(app, storeKey)
	return storeKey
}

func ProvideTransientStoreKey(
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) *storetypes.TransientStoreKey {
	if slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		return nil
	}

	storeKey := storetypes.NewTransientStoreKey(
		fmt.Sprintf("transient:%s", key.Name()),
	)
	registerStoreKey(app, storeKey)
	return storeKey
}

func ProvideMemoryStoreKey(
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) *storetypes.MemoryStoreKey {
	if slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		return nil
	}

	storeKey := storetypes.NewMemoryStoreKey(
		fmt.Sprintf("memory:%s", key.Name()),
	)
	registerStoreKey(app, storeKey)
	return storeKey
}

func ProvideModuleManager(
	modules map[string]appmodule.AppModule,
) *module.Manager {
	return module.NewManagerFromMap(modules)
}

func ProvideGenesisTxHandler(appBuilder *AppBuilder) genesis.TxHandler {
	return NoopGenTxHandler{}
}

type NoopGenTxHandler struct{}

func (NoopGenTxHandler) ExecuteGenesisTx([]byte) error {
	return nil
}

func ProvideEnvironment(
	logger log.Logger,
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) (store.KVStoreService, store.MemoryStoreService, appmodule.Environment) {
	var (
		kvService    store.KVStoreService     = failingStoreService{}
		memKvService store.MemoryStoreService = failingStoreService{}
	)

	// skips modules that have no store
	if !slices.Contains(config.GetSkipStoreKeys(), key.Name()) {
		storeKey := ProvideKVStoreKey(config, key, app)
		kvService = kvStoreService{key: storeKey}

		memStoreKey := ProvideMemoryStoreKey(config, key, app)
		memKvService = memStoreService{key: memStoreKey}
	}

	return kvService, memKvService, NewEnvironment(
		kvService,
		logger.With(log.ModuleKey, fmt.Sprintf("x/%s", key.Name())),
	)
}

func ProvideTransientStoreService(
	config *runtimev1alpha1.Module,
	key depinject.ModuleKey,
	app *AppBuilder,
) store.TransientStoreService {
	storeKey := ProvideTransientStoreKey(config, key, app)
	if storeKey == nil {
		return failingStoreService{}
	}

	return transientStoreService{key: storeKey}
}
