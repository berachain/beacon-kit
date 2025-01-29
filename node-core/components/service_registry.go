// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	service "github.com/berachain/beacon-kit/node-core/services/registry"
	"github.com/berachain/beacon-kit/node-core/services/shutdown"
	"github.com/berachain/beacon-kit/node-core/services/version"
	"github.com/berachain/beacon-kit/observability/telemetry"
)

// ServiceRegistryInput is the input for the service registry provider.
type ServiceRegistryInput[
	LoggerT log.AdvancedLogger[LoggerT],
] struct {
	depinject.In
	ChainService     *blockchain.Service
	EngineClient     *client.EngineClient
	Logger           LoggerT
	NodeAPIServer    *server.Server[echo.Context]
	ReportingService *version.ReportingService
	TelemetrySink    *metrics.TelemetrySink
	TelemetryService *telemetry.Service
	ValidatorService *validator.Service
	CometBFTService  *cometbft.Service[LoggerT]
	ShutdownService  *shutdown.Service
}

// ProvideServiceRegistry is the depinject provider for the service registry.
func ProvideServiceRegistry[
	LoggerT log.AdvancedLogger[LoggerT],
](
	in ServiceRegistryInput[LoggerT],
) *service.Registry {
	// Note: the order of opts matters since the registry will start these services
	// in the order they are  declared in this slice, and in reverse order
	// during shutdown.
	opts := []service.RegistryOption{
		// we want shutdownservice to be the first service to start and the last to stop
		service.WithService(in.ShutdownService),

		service.WithService(in.ValidatorService),
		service.WithService(in.NodeAPIServer),
		service.WithService(in.ReportingService),
		service.WithService(in.TelemetryService),

		// engineClient will block until it connects to the execution layer
		service.WithService(in.EngineClient),

		// only once we connect to an execution client will we start the
		// chain service and cometbft service
		service.WithService(in.ChainService),
		service.WithService(in.CometBFTService),
	}

	return service.NewRegistry(in.Logger, opts...)
}
