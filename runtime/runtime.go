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

//nolint:revive // blank import to register uber maxprocs.
package runtime

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/async/dispatch"
	"github.com/berachain/beacon-kit/async/notify"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/builder"
	localbuilder "github.com/berachain/beacon-kit/beacon/builder/local"
	"github.com/berachain/beacon-kit/beacon/core/randao"
	"github.com/berachain/beacon-kit/beacon/execution"
	"github.com/berachain/beacon-kit/beacon/staking"
	"github.com/berachain/beacon-kit/beacon/sync"
	"github.com/berachain/beacon-kit/cache"
	"github.com/berachain/beacon-kit/config"
	stakingabi "github.com/berachain/beacon-kit/contracts/abi"
	"github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	engineclient "github.com/berachain/beacon-kit/engine/client"
	"github.com/berachain/beacon-kit/health"
	"github.com/berachain/beacon-kit/lib/abi"
	_ "github.com/berachain/beacon-kit/lib/maxprocs"
	"github.com/berachain/beacon-kit/runtime/service"
	"github.com/cosmos/cosmos-sdk/client"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime struct {
	cfg      *config.Config
	logger   log.Logger
	fscp     BeaconStorageBackend
	services *service.Registry
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

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // This function is long because it sets up the services.
func NewDefaultBeaconKitRuntime(
	appOpts AppOptions,
	signer crypto.Signer[[bls12381.SignatureLength]byte],
	logger log.Logger,
	bsb BeaconStorageBackend,
	vcp ValsetChangeProvider,
) (*BeaconKitRuntime, error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit")

	// Read the configuration from the application options.
	cfg, err := config.ReadConfigFromAppOpts(appOpts)
	if err != nil {
		return nil, err
	}

	// Build the service dispatcher.
	gcd, err := dispatch.NewGrandCentralDispatch(
		dispatch.WithLogger(logger),
		dispatch.WithDispatchQueue(
			"dispatch.forkchoice",
			dispatch.QueueTypeSerial,
		),
	)
	if err != nil {
		return nil, err
	}

	// Create the base service, we will the create shallow copies for each
	// service.
	baseService := service.NewBaseService(
		cfg, bsb, gcd, logger,
	)

	// Build the client to interact with the Engine API.
	engineClient := engineclient.New(
		engineclient.WithEngineConfig(&cfg.Engine.Config),
		engineclient.WithLogger(logger),
	)

	// Build the Notification Service.
	notificationService := service.New(
		notify.WithBaseService(baseService.ShallowCopy("notify")),
		notify.WithGCD(gcd),
	)

	// Extrac the staking ABI.
	depositABI, err := stakingabi.BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Build the staking service.
	stakingService := service.New[staking.Service](
		staking.WithBaseService(baseService.ShallowCopy("staking")),
		staking.WithDepositABI(abi.NewWrappedABI(depositABI)),
		staking.WithValsetChangeProvider(vcp),
	)

	// Build the execution service.
	executionService := service.New[execution.Service](
		execution.WithBaseService(baseService.ShallowCopy("execution")),
		execution.WithEngineCaller(engineClient),
		execution.WithStakingService(stakingService),
	)

	// Build the local builder service.
	localBuilder := service.New[localbuilder.Service](
		localbuilder.WithBaseService(baseService.ShallowCopy("local-builder")),
		localbuilder.WithBuilderConfig(&cfg.Builder),
		localbuilder.WithExecutionService(executionService),
		localbuilder.WithPayloadCache(cache.NewPayloadIDCache()),
	)

	// Build the Randao Processor.
	randaoProcessor := randao.NewProcessor(
		randao.WithSigner(signer),
		randao.WithLogger(logger.With("service", "randao")),
	)

	// Build the builder service.
	builderService := service.New[builder.Service](
		builder.WithBaseService(baseService.ShallowCopy("builder")),
		builder.WithBuilderConfig(&cfg.Builder),
		builder.WithLocalBuilder(localBuilder),
		builder.WithRandaoProcessor(randaoProcessor),
	)

	// Build the sync service.
	syncService := service.New[sync.Service](
		sync.WithBaseService(baseService.ShallowCopy("sync")),
		sync.WithEngineClient(engineClient),
		sync.WithConfig(sync.DefaultConfig()),
	)

	// Build the blockchain service.
	chainService := service.New[blockchain.Service](
		blockchain.WithBaseService(baseService.ShallowCopy("blockchain")),
		blockchain.WithExecutionService(executionService),
		blockchain.WithLocalBuilder(localBuilder),
		blockchain.WithRandaoProcessor(randaoProcessor),
		blockchain.WithSyncService(syncService),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger),
		service.WithService(builderService),
		service.WithService(chainService),
		service.WithService(executionService),
		service.WithService(notificationService),
		service.WithService(stakingService),
		service.WithService(syncService),
	)

	// Build the health service.
	healthService := service.New[health.Service](
		health.WithBaseService(baseService.ShallowCopy("health")),
		health.WithServiceRegistry(svcRegistry),
	)

	if err = svcRegistry.RegisterService(healthService); err != nil {
		return nil, err
	}

	// Pass all the services and options into the BeaconKitRuntime.
	return NewBeaconKitRuntime(
		WithBeaconStorageBackend(bsb),
		WithConfig(cfg),
		WithLogger(logger),
		WithServiceRegistry(svcRegistry),
	)
}

// StartServices starts the services.
func (r *BeaconKitRuntime) StartServices(
	ctx context.Context,
	clientCtx client.Context,
) {
	var syncService *sync.Service
	if err := r.services.FetchService(&syncService); err != nil {
		panic(err)
	}
	syncService.SetClientContext(clientCtx)
	r.services.StartAll(ctx)
}
