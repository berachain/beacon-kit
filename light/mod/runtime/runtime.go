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
	lightcore "github.com/berachain/beacon-kit/light/mod/core"
	"github.com/berachain/beacon-kit/light/mod/provider"
	"github.com/berachain/beacon-kit/light/mod/runtime/services/blockchain"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/execution"
	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"

	"github.com/berachain/beacon-kit/mod/runtime"
)

// NewDefaultBeaconLightRuntime creates a new BeaconKitRuntime with the default
// services.
func NewDefaultBeaconLightRuntime(
	cfg *config.Config,
	// signer core.BLSSigner,
	jwtSecret *jwt.Secret,
	// kzgTrustedSetup *gokzg4844.JSONTrustedSetup,
	bsb runtime.BeaconStorageBackend,
	provider *provider.Provider,
	logger log.Logger,
) (*runtime.BeaconKitRuntime, error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit")

	// Create the base service, we will the create shallow copies for each
	// service.
	baseService := service.NewBaseService(
		cfg, bsb, logger,
	)

	// Build the client to interact with the Engine API.
	engineClient := engineclient.New(
		engineclient.WithEngineConfig(&cfg.Engine),
		engineclient.WithJWTSecret(jwtSecret),
		engineclient.WithLogger(logger),
	)

	// TODO: move.
	engineClient.Start(context.Background())

	// // Extrac the staking ABI.
	// depositABI, err := abi.BeaconDepositContractMetaData.GetAbi()
	// if err != nil {
	// 	return nil, err
	// }

	// Build the execution engine.
	executionEngine := execution.NewEngine(engineClient, logger)

	// // Build the staking service.
	// stakingService := service.New[staking.Service](
	// 	staking.WithBaseService(baseService.ShallowCopy("staking")),
	// 	staking.WithDepositABI(depositABI),
	// 	staking.WithExecutionEngine(executionEngine),
	// )

	// // Build the local builder service.
	// localBuilder := service.New[localbuilder.Service](
	// 	localbuilder.WithBaseService(baseService.ShallowCopy("local-builder")),
	// 	localbuilder.WithBuilderConfig(&cfg.Builder),
	// 	localbuilder.WithExecutionEngine(executionEngine),
	// 	localbuilder.WithPayloadCache(cache.NewPayloadIDCache()),
	// )

	// // Build the Blobs Vierifer
	// blobProofVerifier, err := da.NewBlobProofVerifier(
	// 	cfg.KZG.Implementation, kzgTrustedSetup,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// logger.Info(
	// 	"successfully loaded blob verifier",
	// 	"impl",
	// 	cfg.KZG.Implementation,
	// )

	// // Build the Blobs Processor.
	// blobsProcessor := blobs.NewProcessor(
	// 	da.NewBlobVerifier(blobProofVerifier), logger)

	// // Build the Randao Processor.
	// randaoProcessor := randao.NewProcessor(
	// 	randao.WithSigner(signer),
	// 	randao.WithLogger(logger.With("service", "randao")),
	// 	randao.WithConfig(cfg),
	// )

	// // Build the builder service.
	// builderService := service.New[builder.Service](
	// 	builder.WithBaseService(baseService.ShallowCopy("builder")),
	// 	builder.WithBuilderConfig(&cfg.Builder),
	// 	builder.WithLocalBuilder(localBuilder),
	// 	builder.WithRandaoProcessor(randaoProcessor),
	// 	builder.WithSigner(signer),
	// )

	// Build the blockchain service.
	chainService := service.New[blockchain.Service](
		blockchain.WithBaseService(baseService.ShallowCopy("blockchain")),
		blockchain.WithBlockValidator(core.NewBlockValidator(&cfg.Beacon)),
		blockchain.WithExecutionEngine(executionEngine),
		// blockchain.WithLocalBuilder(localBuilder),
		blockchain.WithPayloadValidator(lightcore.NewPayloadValidator(&cfg.Beacon)),
		// blockchain.WithStakingService(stakingService),
		blockchain.WithStateProcessor(
			lightcore.NewStateProcessor(
				&cfg.Beacon,
				logger,
			)),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger),
		// service.WithService(builderService),
		service.WithService(chainService),
		// service.WithService(stakingService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime(
		runtime.WithBeaconStorageBackend(bsb),
		runtime.WithConfig(cfg),
		runtime.WithLogger(logger),
		runtime.WithServiceRegistry(svcRegistry),
	)
}
