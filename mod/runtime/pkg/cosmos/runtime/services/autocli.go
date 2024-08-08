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

package services

import (
	"context"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	cosmosmsg "cosmossdk.io/api/cosmos/msg/v1"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/types/module"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// AutoCLIQueryService implements the cosmos.autocli.v1.Query service.
type AutoCLIQueryService struct {
	autocliv1.UnimplementedQueryServer

	moduleOptions map[string]*autocliv1.ModuleOptions
}

// NewAutoCLIQueryService returns a AutoCLIQueryService for the provided
// modules.
func NewAutoCLIQueryService(
	appModules map[string]appmodule.AppModule,
) *AutoCLIQueryService {
	return &AutoCLIQueryService{
		moduleOptions: ExtractAutoCLIOptions(appModules),
	}
}

// ExtractAutoCLIOptions extracts autocli ModuleOptions from the provided app
// modules.
//
// Example Usage:
//
//	ExtractAutoCLIOptions(ModuleManager.Modules)
func ExtractAutoCLIOptions(
	appModules map[string]appmodule.AppModule,
) map[string]*autocliv1.ModuleOptions {
	moduleOptions := map[string]*autocliv1.ModuleOptions{}
	for modName, mod := range appModules {
		if autoCliMod, ok := mod.(interface {
			AutoCLIOptions() *autocliv1.ModuleOptions
		}); ok {
			moduleOptions[modName] = autoCliMod.AutoCLIOptions()
			continue
		}

		cfg := &autocliConfigurator{}

		// try to auto-discover options based on the last msg and query
		// services registered for the module
		if mod, ok := mod.(module.HasServices); ok {
			mod.RegisterServices(cfg)
		}

		if mod, ok := mod.(appmodule.HasServices); ok {
			err := mod.RegisterServices(cfg)
			if err != nil {
				panic(err)
			}
		}

		// check for errors in the configurator
		if cfg.Error() != nil {
			panic(cfg.Error())
		}

		haveServices := false
		modOptions := &autocliv1.ModuleOptions{}
		if cfg.msgServer.serviceName != "" {
			haveServices = true
			modOptions.Tx = &autocliv1.ServiceCommandDescriptor{
				Service: cfg.msgServer.serviceName,
			}
		}

		if cfg.queryServer.serviceName != "" {
			haveServices = true
			modOptions.Query = &autocliv1.ServiceCommandDescriptor{
				Service: cfg.queryServer.serviceName,
			}
		}

		if haveServices {
			moduleOptions[modName] = modOptions
		}
	}
	return moduleOptions
}

func (a AutoCLIQueryService) AppOptions(
	context.Context,
	*autocliv1.AppOptionsRequest,
) (*autocliv1.AppOptionsResponse, error) {
	return &autocliv1.AppOptionsResponse{
		ModuleOptions: a.moduleOptions,
	}, nil
}

// autocliConfigurator allows us to call RegisterServices and introspect the
// services.
type autocliConfigurator struct {
	msgServer     autocliServiceRegistrar
	queryServer   autocliServiceRegistrar
	registryCache *protoregistry.Files
	err           error
}

var _ module.Configurator = &autocliConfigurator{} //nolint:staticcheck // SA1019: Configurator is deprecated but still used in runtime v1.

func (a *autocliConfigurator) MsgServer() gogogrpc.Server { return &a.msgServer }

func (a *autocliConfigurator) QueryServer() gogogrpc.Server { return &a.queryServer }

func (a *autocliConfigurator) RegisterMigration(
	string,
	uint64,
	module.MigrationHandler,
) error {
	return nil
}

func (a *autocliConfigurator) Register(
	string,
	uint64,
	appmodule.MigrationHandler,
) error {
	return nil
}

func (a *autocliConfigurator) RegisterService(
	sd *grpc.ServiceDesc,
	ss interface{},
) {
	if a.registryCache == nil {
		a.registryCache, a.err = proto.MergedRegistry()
	}

	desc, err := a.registryCache.FindDescriptorByName(
		protoreflect.FullName(sd.ServiceName),
	)
	if err != nil {
		a.err = err
		return
	}

	if protobuf.HasExtension(desc.Options(), cosmosmsg.E_Service) {
		a.msgServer.RegisterService(sd, ss)
	} else {
		a.queryServer.RegisterService(sd, ss)
	}
}
func (a *autocliConfigurator) Error() error { return nil }

// autocliServiceRegistrar is used to capture the service name for registered
// services.
type autocliServiceRegistrar struct {
	serviceName string
}

func (a *autocliServiceRegistrar) RegisterService(
	sd *grpc.ServiceDesc,
	_ interface{},
) {
	a.serviceName = sd.ServiceName
}

var _ autocliv1.QueryServer = &AutoCLIQueryService{}
