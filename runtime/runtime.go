// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/async/dispatch"
	"github.com/itsdevbear/bolaris/async/notify"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/beacon/state"
	"github.com/itsdevbear/bolaris/beacon/sync"
	"github.com/itsdevbear/bolaris/cache"
	"github.com/itsdevbear/bolaris/config"
	builder "github.com/itsdevbear/bolaris/execution/builder/local"
	"github.com/itsdevbear/bolaris/execution/engine"
	eth "github.com/itsdevbear/bolaris/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/io/jwt"
	"github.com/itsdevbear/bolaris/runtime/service"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime struct {
	cfg       *config.Config
	cometCfg  CometBFTConfig
	logger    log.Logger
	ethclient *eth.Eth1Client
	fscp      BeaconStateProvider
	services  *service.Registry
}

// BeaconStateProvider is an interface that provides the
// beacon state to the runtime.
type BeaconStateProvider interface {
	BeaconState(ctx context.Context) state.BeaconState
}

type ValsetChangeProvider interface {
	ApplyChanges(context.Context, []*consensusv1.Deposit, []*enginev1.Withdrawal) error
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime(
	opts ...Option,
) (*BeaconKitRuntime, error) {
	bkr := &BeaconKitRuntime{}
	for _, opt := range opts {
		if err := opt(bkr); err != nil {
			return nil, err
		}
	}

	return bkr, nil
}

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default services.
func NewDefaultBeaconKitRuntime(
	cfg *config.Config, bsp BeaconStateProvider, vcp ValsetChangeProvider, logger log.Logger,
) (*BeaconKitRuntime, error) {
	// Get JWT Secret for eth1 connection.
	jwtSecret, err := jwt.NewFromFile(cfg.Engine.JWTSecretPath)
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
	baseService := service.NewBaseService(
		&cfg.Beacon, bsp, vcp, gcd, logger,
	)

	// Create a payloadCache for the execution service and validator service to share.
	payloadCache := cache.NewPayloadIDCache()

	// Create the eth1 client that will be used to interact with the execution client.
	eth1Client, err := eth.NewEth1Client(
		eth.WithStartupRetryInterval(cfg.Engine.RPCStartupCheckInterval),
		eth.WithHealthCheckInterval(cfg.Engine.RPCHealthCheckInterval),
		eth.WithJWTRefreshInterval(cfg.Engine.RPCJWTRefreshInterval),
		eth.WithEndpointDialURL(cfg.Engine.RPCDialURL),
		eth.WithJWTSecret(jwtSecret),
		eth.WithLogger(logger),
		eth.WithRequiredChainID(cfg.Engine.RequiredChainID),
	)
	if err != nil {
		return nil, err
	}

	// Build the Notification Service.
	notificationService := notify.NewService(
		notify.WithGCD(gcd),
		notify.WithLogger(logger),
	)

	// NewClient wraps the eth1 client and provides the interface for the
	// blockchain service to interact with the execution client.
	engineClient := engine.NewClient(
		engine.WithEth1Client(eth1Client),
		engine.WithBeaconConfig(&cfg.Beacon),
		engine.WithLogger(logger),
		engine.WithEngineTimeout(cfg.Engine.RPCTimeout))

	// Build the execution service.
	executionService := execution.New(
		baseService.WithName("execution"),
		execution.WithEngineCaller(engineClient),
	)

	// Build the local builder service.
	builderService := builder.NewService(
		baseService.WithName("local-builder"),
		builder.WithBuilderConfig(&cfg.Builder),
		builder.WithExecutionService(executionService),
		builder.WithPayloadCache(payloadCache),
	)

	// Build the blockchain service.
	chainService := blockchain.NewService(
		baseService.WithName("blockchain"),
		blockchain.WithBuilderService(builderService),
		blockchain.WithExecutionService(executionService),
	)

	// Build the sync service.
	syncService := sync.NewService(
		baseService.WithName("sync"),
		sync.WithEthClient(eth1Client),
		sync.WithExecutionService(executionService),
	)

	// Create the service registry.
	serviceRegistry := service.NewRegistry(
		service.WithLogger(logger),
		service.WithService(syncService),
		service.WithService(executionService),
		service.WithService(chainService),
		service.WithService(notificationService),
		service.WithService(builderService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return NewBeaconKitRuntime(
		WithConfig(cfg),
		WithLogger(logger),
		WithServiceRegistry(serviceRegistry),
		WithBeaconStateProvider(bsp),
		// We put the eth1 client in the BeaconKitRuntime so we can attach the cmd.Context to it.
		// This is necessary for the eth1 client to be able to shut down gracefully.
		WithEth1Client(eth1Client),
	)
}

// StartServices starts all services in the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) StartServices(cmdCtx context.Context) {
	// First try to start the eth1 client.
	r.ethclient.Start(cmdCtx)

	// Then start all the other services.
	r.services.StartAll(cmdCtx)
}

// SetCometCfg sets the CometBFTConfig for the runtime.
func (r *BeaconKitRuntime) SetCometCfg(cometCfg CometBFTConfig) {
	r.cometCfg = cometCfg
}

// FetchService retrieves a service from the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) FetchService(service interface{}) error {
	return r.services.FetchService(service)
}

// InitialSyncCheck.
func (r *BeaconKitRuntime) InitialSyncCheck(ctx context.Context) error {
	var syncService *sync.Service
	if err := r.services.FetchService(&syncService); err != nil {
		return err
	}

	return syncService.CheckSyncStatusAndForkchoice(ctx)
}
