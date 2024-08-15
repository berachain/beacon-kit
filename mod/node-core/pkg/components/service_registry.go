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

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

// ServiceRegistryInput is the input for the service registry provider.
type ServiceRegistryInput[
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In
	BlockStoreService *BlockStoreService
	ChainService      *ChainService
	DAService         *DAService
	DBManager         *DBManager
	DepositService    *DepositService
	EngineClient      *EngineClient
	Logger            LoggerT
	NodeAPIServer     *NodeAPIServer
	ReportingService  *ReportingService
	TelemetrySink     *metrics.TelemetrySink
	ValidatorService  *ValidatorService
	Dispatcher        *Dispatcher
}

// ProvideServiceRegistry is the depinject provider for the service registry.
func ProvideServiceRegistry[
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	in ServiceRegistryInput[LoggerT],
) *service.Registry {
	return service.NewRegistry(
		service.WithLogger(in.Logger),
		service.WithService(in.Dispatcher),
		service.WithService(in.ValidatorService),
		service.WithService(in.BlockStoreService),
		service.WithService(in.ChainService),
		service.WithService(in.DAService),
		service.WithService(in.DepositService),
		service.WithService(in.NodeAPIServer),
		service.WithService(in.ReportingService),
		service.WithService(in.DBManager),
		service.WithService(in.EngineClient),
	)
}
