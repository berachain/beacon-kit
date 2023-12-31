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

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	initialsync "github.com/itsdevbear/bolaris/beacon/initial-sync"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/types"
	"github.com/prysmaticlabs/prysm/v4/runtime"
)

// BeaconKitRuntime is a struct that holds the service registry.
type BeaconKitRuntime struct {
	mu       sync.Mutex
	logger   log.Logger
	fscp     ForkChoiceStoreProvider
	services *runtime.ServiceRegistry
}

type ForkChoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime and applies the provided options.
func NewBeaconKitRuntime(opts ...Option) (*BeaconKitRuntime, error) {
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
	ctx context.Context, cfg *config.Config, fcsp ForkChoiceStoreProvider, logger log.Logger,
) (*BeaconKitRuntime, error) {
	jwtSecret, err := eth.LoadJWTSecret(cfg.ExecutionClient.JWTSecretPath, logger)
	if err != nil {
		return nil, err
	}

	// Create the eth1 client that will be used to interact with the execution client.
	opts := []eth.Option{
		eth.WithHTTPEndpointAndJWTSecret(cfg.ExecutionClient.RPCDialURL, jwtSecret),
		eth.WithLogger(logger),
		eth.WithRequiredChainID(cfg.ExecutionClient.RequiredChainID),
	}

	eth1Client, err := eth.NewEth1Client(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// Engine Caller wraps the eth1 client and provides the interface for the
	// blockchain service to interact with the execution client.
	engineCallerOpts := []engine.Option{
		engine.WithEth1Client(eth1Client),
		engine.WithBeaconConfig(&cfg.BeaconConfig),
		engine.WithLogger(logger),
		engine.WithEngineTimeout(cfg.ExecutionClient.RPCTimeout),
	}
	engineCaller := engine.NewCaller(engineCallerOpts...)

	executionOpts := []execution.Option{
		execution.WithBeaconConfig(&cfg.BeaconConfig),
		execution.WithLogger(logger),
		execution.WithForkChoiceStoreProvider(fcsp),
		execution.WithEngineCaller(engineCaller),
	}
	executionService := execution.New(executionOpts...)

	// Build the blockchain service
	chainOpts := []blockchain.Option{
		blockchain.WithBeaconConfig(&cfg.BeaconConfig),
		blockchain.WithLogger(logger),
		blockchain.WithForkChoiceStoreProvider(fcsp),
		blockchain.WithExecutionService(executionService),
	}

	chainService := blockchain.NewService(chainOpts...)

	syncServiceOpts := []initialsync.Option{
		initialsync.WithLogger(logger),
		initialsync.WithEthClient(eth1Client),
		initialsync.WithForkChoiceStoreProvider(fcsp),
	}

	// Build Sync Service
	syncService := initialsync.NewService(syncServiceOpts...)

	return NewBeaconKitRuntime([]Option{
		WithService(syncService),
		WithService(executionService),
		WithService(chainService),
		WithLogger(logger),
		WithForkChoiceStoreProvider(fcsp),
	}...)
}

// StartServices starts all services in the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) StartServices() {
	r.services.StartAll()
}

// StopServices stops all services in the BeaconKitRuntime's service registry.
func (r *BeaconKitRuntime) StopServices() {
	r.services.StopAll()
}

func (r *BeaconKitRuntime) FetchService(service interface{}) error {
	return r.services.FetchService(service)
}
