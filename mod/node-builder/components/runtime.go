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

package components

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/randao"
	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/da"
	"github.com/berachain/beacon-kit/mod/da/kzg"
	"github.com/berachain/beacon-kit/mod/execution"
	engineclient "github.com/berachain/beacon-kit/mod/execution/client"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/node-builder/utils/jwt"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/builder"
	"github.com/berachain/beacon-kit/mod/payload/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/services/blockchain"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking"
	"github.com/berachain/beacon-kit/mod/runtime/services/staking/abi"
	"github.com/berachain/beacon-kit/mod/validator"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // bullish.
func ProvideRuntime(
	cfg *config.Config,
	chainSpec primitives.ChainSpec,
	signer runtime.BLSSigner,
	jwtSecret *jwt.Secret,
	kzgTrustedSetup *gokzg4844.JSONTrustedSetup,
	// TODO: this is really poor coupling, we should fix.
	bsb runtime.BeaconStorageBackend[types.ReadOnlyBeaconBlock],
	logger log.Logger,
) (*runtime.BeaconKitRuntime, error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit")

	// Create the base service, we will the create shallow copies for each
	// service.
	baseService := service.NewBaseService(
		cfg, bsb, chainSpec, logger,
	)

	// Build the client to interact with the Engine API.
	engineClient := engineclient.New(
		engineclient.WithEngineConfig(&cfg.Engine),
		engineclient.WithJWTSecret(jwtSecret),
		engineclient.WithLogger(
			logger.With("module", "beacon-kit.engine.client"),
		),
	)

	// TODO: move.
	engineClient.Start(context.Background())

	// Extrac the staking ABI.
	depositABI, err := abi.BeaconDepositContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Build the execution engine.
	executionEngine := execution.NewEngine(engineClient, logger)

	// Build the staking service.
	stakingService := service.New[staking.Service](
		staking.WithBaseService(baseService.ShallowCopy("staking")),
		staking.WithDepositABI(depositABI),
		staking.WithDepositStore(bsb.DepositStore(nil)),
		staking.WithExecutionEngine(executionEngine),
	)

	// Build the local builder service.
	localBuilder := service.New[payloadbuilder.PayloadBuilder](
		payloadbuilder.WithLogger(
			logger.With("service", "payload-builder"),
		),
		payloadbuilder.WithChainSpec(chainSpec),
		payloadbuilder.WithConfig(&cfg.PayloadBuilder),
		payloadbuilder.WithExecutionEngine(executionEngine),
		payloadbuilder.WithPayloadCache(
			cache.NewPayloadIDCache[engineprimitives.PayloadID, [32]byte, math.Slot](),
		),
	)

	// Build the Blobs Verifier
	blobProofVerifier, err := kzg.NewBlobProofVerifier(
		cfg.KZG.Implementation, kzgTrustedSetup,
	)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"successfully loaded blob verifier",
		"impl",
		cfg.KZG.Implementation,
	)

	// Build the Randao Processor.
	randaoProcessor := randao.NewProcessor(
		randao.WithSigner(signer),
		randao.WithLogger(logger.With("service", "randao")),
		randao.WithConfig(chainSpec),
	)

	// Build the builder service.
	blobFactory := da.NewSidecarFactory[types.BeaconBlockBody](
		chainSpec,
		types.KZGPositionDeneb,
	)
	validatorService := validator.NewService(
		validator.WithBlobFactory(blobFactory),
		validator.WithChainSpec(chainSpec),
		validator.WithConfig(&cfg.Validator),
		validator.WithDepositStore(bsb.DepositStore(nil)),
		validator.WithLocalBuilder(localBuilder),
		validator.WithLogger(logger.With("service", "validator")),
		validator.WithRandaoProcessor(randaoProcessor),
		validator.WithSigner(signer),
	)

	// Build the blockchain service.
	chainService := service.New[blockchain.Service](
		blockchain.WithBaseService(baseService.ShallowCopy("blockchain")),
		blockchain.WithBlockValidator(core.NewBlockValidator(chainSpec)),
		blockchain.WithExecutionEngine(executionEngine),
		blockchain.WithLocalBuilder(localBuilder),
		blockchain.WithPayloadValidator(core.NewPayloadValidator(chainSpec)),
		blockchain.WithStakingService(stakingService),
		blockchain.WithStateProcessor(
			core.NewStateProcessor(
				chainSpec,
				da.NewBlobVerifier(blobProofVerifier),
				randaoProcessor,
				logger.With("module", "state-processor"),
			)),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger),
		service.WithService(validatorService),
		service.WithService(chainService),
		service.WithService(stakingService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime(
		runtime.WithBeaconStorageBackend(bsb),
		runtime.WithConfig(cfg),
		runtime.WithLogger(logger),
		runtime.WithServiceRegistry(svcRegistry),
	)
}
