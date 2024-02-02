// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package runtime

import (
	"context"
	"sync"

	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/async/dispatch"
	"github.com/itsdevbear/bolaris/async/notify"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	initialsync "github.com/itsdevbear/bolaris/beacon/initial-sync"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/state"
	"github.com/prysmaticlabs/prysm/v4/runtime"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime struct {
	cfg        *config.Config
	mu         sync.Mutex
	logger     log.Logger
	fscp       BeaconStateProvider
	services   *runtime.ServiceRegistry
	dispatcher *dispatch.GrandCentralDispatch
}

type BeaconStateProvider interface {
	BeaconState(ctx context.Context) state.BeaconState
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime(
	opts ...Option,
) (*BeaconKitRuntime, error) {
	bkr := &BeaconKitRuntime{
		services: runtime.NewServiceRegistry(),
	}

	for _, opt := range opts {
		if err := opt(bkr); err != nil {
			return nil, err
		}
	}

	return bkr, nil
}

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default services.
func NewDefaultBeaconKitRuntime(
	ctx context.Context, cfg *config.Config, bsp BeaconStateProvider, logger log.Logger,
) (*BeaconKitRuntime, error) {
	// Get JWT Secret for eth1 connection.
	jwtSecret, err := eth.LoadJWTSecret(cfg.ExecutionClient.JWTSecretPath, logger)
	if err != nil {
		return nil, err
	}

	// Build the service dispatcher.
	gcd, err := dispatch.NewGrandCentralDispatch(
		dispatch.WithLogger(logger),
		dispatch.WithDispatchQueue("dispatch.forkchoice", dispatch.QueueTypeSerial),
	)
	if err != nil {
		return nil, err
	}

	// Create the base service, we will the  create shallow copies for each service.
	baseService := service.NewBaseService(&cfg.Beacon, gcd, logger)

	// Create the eth1 client that will be used to interact with the execution client.
	eth1Client, err := eth.NewEth1Client(
		ctx,
		eth.WithHTTPEndpointAndJWTSecret(cfg.ExecutionClient.RPCDialURL, jwtSecret),
		eth.WithLogger(logger),
		eth.WithRequiredChainID(cfg.ExecutionClient.RequiredChainID),
	)
	if err != nil {
		return nil, err
	}

	// Build the Notification Service.
	notificationService := notify.NewService(
		notify.WithGCD(gcd),
		notify.WithLogger(logger),
	)

	// Engine Caller wraps the eth1 client and provides the interface for the
	// blockchain service to interact with the execution client.
	engineCaller := engine.NewCaller(engine.WithEth1Client(eth1Client),
		engine.WithBeaconConfig(&cfg.Beacon),
		engine.WithLogger(logger),
		engine.WithEngineTimeout(cfg.ExecutionClient.RPCTimeout))

	// Build the execution service.
	executionService := execution.New(
		baseService.WithName("execution"),
		execution.WithBeaconStateProvider(bsp),
		execution.WithEngineCaller(engineCaller),
	)

	// Build the blockchain service
	chainService := blockchain.NewService(
		baseService.WithName("blockchain"),
		blockchain.WithBeaconStateProvider(bsp),
		blockchain.WithExecutionService(executionService),
	)

	// Build the sync service.
	syncService := initialsync.NewService(
		baseService.WithName("initial-sync"),
		initialsync.WithEthClient(eth1Client),
		initialsync.WithBeaconStateProvider(bsp),
		initialsync.WithExecutionService(executionService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return NewBeaconKitRuntime(
		WithConfig(cfg),
		WithService(syncService),
		WithService(executionService),
		WithService(chainService),
		WithService(notificationService),
		WithLogger(logger),
		WithBeaconStateProvider(bsp),
		WithDispatcher(gcd),
	)
}

// StartServices starts all services in the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) StartServices() {
	r.services.StartAll()
}

// StopServices stops all services in the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) StopServices() {
	r.logger.Info("stopping all services")
	r.services.StopAll()
}

// FetchService retrieves a service from the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) FetchService(service interface{}) error {
	return r.services.FetchService(service)
}

// InitialSyncCheck.
func (r *BeaconKitRuntime) InitialSyncCheck(ctx context.Context) error {
	var syncService *initialsync.Service

	if err := r.services.FetchService(&syncService); err != nil {
		return err
	}

	return syncService.CheckSyncStatusAndForkchoice(ctx)
}
