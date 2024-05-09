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
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/staking"
	"github.com/berachain/beacon-kit/mod/beacon/staking/abi"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	"github.com/berachain/beacon-kit/mod/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	stda "github.com/berachain/beacon-kit/mod/state-transition/pkg/da"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/randao"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/verification"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // bullish.
func ProvideRuntime(
	cfg *config.Config,
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	jwtSecret *jwt.Secret,
	kzgTrustedSetup *gokzg4844.JSONTrustedSetup,
	// TODO: this is really poor coupling, we should fix.
	storageBackend runtime.BeaconStorageBackend[
		*datypes.BlobSidecars,
		*depositdb.KVStore,
		consensus.ReadOnlyBeaconBlockBody,
	],
	logger log.Logger,
) (*runtime.BeaconKitRuntime[
	*datypes.BlobSidecars,
	*depositdb.KVStore,
	consensus.ReadOnlyBeaconBlockBody,
	runtime.BeaconStorageBackend[
		*datypes.BlobSidecars,
		*depositdb.KVStore,
		consensus.ReadOnlyBeaconBlockBody,
	],
], error) {
	// Set the module as beacon-kit to override the cosmos-sdk naming.
	logger = logger.With("module", "beacon-kit")

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
	executionEngine := execution.New(engineClient, logger)

	// Build the staking service.
	stakingService := staking.NewService(
		staking.WithBeaconStorageBackend(storageBackend),
		staking.WithChainSpec(chainSpec),
		staking.WithDepositABI(depositABI),
		staking.WithDepositStore(storageBackend.DepositStore(nil)),
		staking.WithExecutionEngine(executionEngine),
		staking.WithLogger(logger.With("service", "staking")),
	)

	// Build the local builder service.
	localBuilder := payloadbuilder.New[state.BeaconState](
		&cfg.PayloadBuilder,
		chainSpec,
		logger.With("service", "payload-builder"),
		executionEngine,
		cache.NewPayloadIDCache[engineprimitives.PayloadID, [32]byte, math.Slot](),
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
	randaoProcessor := randao.NewProcessor[
		consensus.BeaconBlockBody,
		consensus.BeaconBlock,
		state.BeaconState,
	](
		chainSpec,
		signer,
		logger.With("service", "randao"),
	)

	// Build the builder service.
	validatorService := validator.NewService[*datypes.BlobSidecars](
		&cfg.Validator,
		logger.With("service", "validator"),
		chainSpec,
		signer,
		dablob.NewSidecarFactory[consensus.BeaconBlockBody](
			chainSpec,
			consensus.KZGPositionDeneb,
		),
		randaoProcessor,
		storageBackend.DepositStore(nil),
		localBuilder,
		[]validator.PayloadBuilder[state.BeaconState]{localBuilder},
	)

	// Build the blockchain service.
	chainService := blockchain.NewService[*datypes.BlobSidecars](
		storageBackend,
		logger.With("service", "blockchain"),
		chainSpec,
		executionEngine,
		localBuilder,
		stakingService,
		verification.NewBlockVerifier(chainSpec),
		core.NewStateProcessor[*datypes.BlobSidecars](
			chainSpec,
			stda.NewBlobProcessor[
				consensus.ReadOnlyBeaconBlockBody, *datypes.BlobSidecars,
			](
				logger.With("module", "blob-processor"),
				chainSpec,
				dablob.NewVerifier(blobProofVerifier),
			),
			randaoProcessor,
			signer,
			logger.With("module", "state-processor"),
		),
		verification.NewPayloadVerifier(chainSpec),
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger.With("module", "service-registry")),
		service.WithService(validatorService),
		service.WithService(chainService),
		service.WithService(stakingService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime[
		*datypes.BlobSidecars,
		*depositdb.KVStore,
		consensus.ReadOnlyBeaconBlockBody,
	](
		logger.With(
			"module",
			"beacon-kit.runtime",
		),
		svcRegistry,
		storageBackend,
	)
}
