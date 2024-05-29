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
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/randao"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/ethereum/go-ethereum/event"
)

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[types.BeaconBlockBody],
	types.BeaconBlockBody,
	core.BeaconState[*types.Validator],
	*datypes.BlobSidecars,
	*depositdb.KVStore,
	runtime.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		core.BeaconState[*types.Validator],
		*datypes.BlobSidecars,
		*depositdb.KVStore,
	],
]

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // bullish.
func ProvideRuntime(
	cfg *config.Config,
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	engineClient *engineclient.EngineClient[*types.ExecutableDataDeneb],
	kzgTrustedSetup *gokzg4844.JSONTrustedSetup,
	storageBackend runtime.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		core.BeaconState[*types.Validator],
		*datypes.BlobSidecars,
		*depositdb.KVStore,
	],
	ts *metrics.TelemetrySink,
	logger log.Logger,
) (*BeaconKitRuntime, error) {
	// Build the execution engine.
	executionEngine := execution.New[engineprimitives.ExecutionPayload](
		engineClient,
		logger.With("service", "execution-engine"),
		ts,
	)

	// Build the deposit contract.
	beaconDepositContract, err := deposit.
		NewWrappedBeaconDepositContract[
		*types.Deposit, types.WithdrawalCredentials,
	](
		chainSpec.DepositContractAddress(),
		engineClient,
	)
	if err != nil {
		return nil, err
	}

	// Build the local builder service.
	localBuilder := payloadbuilder.New[core.BeaconState[*types.Validator]](
		&cfg.PayloadBuilder,
		chainSpec,
		logger.With("service", "payload-builder"),
		executionEngine,
		cache.NewPayloadIDCache[engineprimitives.PayloadID, [32]byte, math.Slot](),
	)

	// Build the Blobs Verifier
	//#nosec:G703 // todo fix depinject stuff.
	blobProofVerifier, _ := kzg.NewBlobProofVerifier(
		cfg.KZG.Implementation, kzgTrustedSetup,
	)

	// // TODO: we need to handle this in the depinject case when the trusted
	// setup
	// // is not ready yet nicer.
	// if err != nil {
	// 	logger.Warn("failed to load blob verifier", "err", err)
	// }

	// logger.Info(
	// 	"successfully loaded blob verifier",
	// 	"impl",
	// 	cfg.KZG.Implementation,
	// )

	// Build the Randao Processor.
	randaoProcessor := randao.NewProcessor[
		types.BeaconBlockBody,
		types.BeaconBlock,
		core.BeaconState[*types.Validator],
	](
		chainSpec,
		signer,
	)

	stateProcessor := core.NewStateProcessor[
		types.BeaconBlock,
		types.BeaconBlockBody,
		core.BeaconState[*types.Validator],
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
		*types.ForkData,
		*types.Validator,
		types.WithdrawalCredentials,
	](
		chainSpec,
		randaoProcessor,
		executionEngine,
		signer,
	)

	// Build the event feed.
	blockFeed := event.FeedOf[events.Block[types.BeaconBlock]]{}

	// Build the builder service.
	validatorService := validator.NewService[
		types.BeaconBlock,
		types.BeaconBlockBody,
		core.BeaconState[*types.Validator], *datypes.BlobSidecars,
	](
		&cfg.Validator,
		logger.With("service", "validator"),
		chainSpec,
		storageBackend,
		stateProcessor,
		signer,
		dablob.NewSidecarFactory[
			types.BeaconBlock,
			types.BeaconBlockBody,
		](
			chainSpec,
			types.KZGPositionDeneb,
			ts,
		),
		randaoProcessor,
		storageBackend.DepositStore(nil),
		localBuilder,
		[]validator.PayloadBuilder[core.BeaconState[*types.Validator]]{
			localBuilder,
		},
		ts,
	)

	// Build the blockchain service.
	chainService := blockchain.NewService[
		*dastore.Store[types.BeaconBlockBody],
		core.BeaconState[*types.Validator],
		*datypes.BlobSidecars,
	](
		storageBackend,
		logger.With("service", "blockchain"),
		chainSpec,
		executionEngine,
		localBuilder,
		dablob.NewProcessor[
			*dastore.Store[types.BeaconBlockBody],
			types.BeaconBlockBody](
			logger.With("service", "blob-processor"),
			chainSpec,
			dablob.NewVerifier(blobProofVerifier, ts),
			types.BlockBodyKZGOffset,
			ts,
		),
		stateProcessor,
		ts,
		&blockFeed,
	)

	// Build the deposit service.
	depositService := deposit.NewService[
		types.BeaconBlock,
		events.Block[types.BeaconBlock],
		*depositdb.KVStore,
		event.Subscription,
	](
		logger.With("service", "deposit"),
		math.U64(chainSpec.Eth1FollowDistance()),
		storageBackend.DepositStore(nil),
		beaconDepositContract,
		&blockFeed,
	)

	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(logger.With("service", "service-registry")),
		service.WithService(validatorService),
		service.WithService(chainService),
		service.WithService(depositService),
		service.WithService(engineClient),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		core.BeaconState[*types.Validator],
		*datypes.BlobSidecars,
		*depositdb.KVStore,
	](
		chainSpec,
		logger,
		svcRegistry,
		storageBackend,
	)
}
