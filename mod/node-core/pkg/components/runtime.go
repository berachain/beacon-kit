// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	"cosmossdk.io/core/log"
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
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/services/version"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/event"
)

type BeaconState = core.BeaconState[
	*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
	*types.Validator, *engineprimitives.Withdrawal,
]

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[types.BeaconBlockBody],
	*types.BeaconBlock,
	types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
]

// NewDefaultBeaconKitRuntime creates a new BeaconKitRuntime with the default
// services.
//
//nolint:funlen // bullish.
func ProvideRuntime(
	cfg *config.Config,
	blobProofVerifier kzg.BlobProofVerifier,
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	engineClient *engineclient.EngineClient[*types.ExecutionPayload],
	storageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
	telemetrySink *metrics.TelemetrySink,
	logger log.Logger,
) (*BeaconKitRuntime, error) {
	// Build the execution engine.
	executionEngine := execution.New[*types.ExecutionPayload](
		engineClient,
		logger.With("service", "execution-engine"),
		telemetrySink,
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
	localBuilder := payloadbuilder.New[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	](
		&cfg.PayloadBuilder,
		chainSpec,
		logger.With("service", "payload-builder"),
		executionEngine,
		cache.NewPayloadIDCache[engineprimitives.PayloadID, [32]byte, math.Slot](),
	)

	stateProcessor := core.NewStateProcessor[
		*types.BeaconBlock,
		types.BeaconBlockBody,
		*types.BeaconBlockHeader,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
		*types.ExecutionPayload,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.ForkData,
		*types.Validator,
		*engineprimitives.Withdrawal,
		types.WithdrawalCredentials,
	](
		chainSpec,
		executionEngine,
		signer,
	)

	// Build the event feed.
	blockFeed := event.FeedOf[events.Block[*types.BeaconBlock]]{}

	// slice of pruners to pass to the DBManager.
	pruners := []*pruner.Pruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription]{}

	// Build the deposit pruner.\
	depositPruner := pruner.NewPruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription,
	](
		logger.With("service", manager.DepositPrunerName),
		storageBackend.DepositStore(nil),
		manager.DepositPrunerName,
		&blockFeed,
		deposit.BuildPruneRangeFn[
			types.BeaconBlockBody,
			*types.BeaconBlock,
			events.Block[*types.BeaconBlock],
			*types.Deposit,
			*types.ExecutionPayload,
			types.WithdrawalCredentials,
		](chainSpec),
	)
	pruners = append(pruners, depositPruner)

	avs := storageBackend.AvailabilityStore(nil).IndexDB
	if avs != nil {
		// build the availability pruner if IndexDB is available.
		availabilityPruner := pruner.NewPruner[
			*types.BeaconBlock,
			events.Block[*types.BeaconBlock],
			event.Subscription,
		](
			logger.With("service", manager.AvailabilityPrunerName),
			avs.(*filedb.RangeDB),
			manager.AvailabilityPrunerName,
			&blockFeed,
			dastore.BuildPruneRangeFn[
				*types.BeaconBlock,
				events.Block[*types.BeaconBlock],
			](chainSpec),
		)
		pruners = append(pruners, availabilityPruner)
	}

	// Build the DBManager service.
	dbManagerService, err := manager.NewDBManager[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription,
	](
		logger.With("service", "db-manager"),
		pruners...,
	)
	if err != nil {
		return nil, err
	}

	blobProcessor := dablob.NewProcessor[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
	](
		logger.With("service", "blob-processor"),
		chainSpec,
		dablob.NewVerifier(blobProofVerifier, telemetrySink),
		types.BlockBodyKZGOffset,
		telemetrySink,
	)

	// Build the builder service.
	validatorService := validator.NewService[
		*types.BeaconBlock,
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		*types.ForkData,
	](
		&cfg.Validator,
		logger.With("service", "validator"),
		chainSpec,
		storageBackend,
		blobProcessor,
		stateProcessor,
		signer,
		dablob.NewSidecarFactory[
			*types.BeaconBlock,
			types.BeaconBlockBody,
		](
			chainSpec,
			types.KZGPositionDeneb,
			telemetrySink,
		),
		localBuilder,
		[]validator.PayloadBuilder[BeaconState, *types.ExecutionPayload]{
			localBuilder,
		},
		telemetrySink,
	)

	// Build the blockchain service.
	chainService := blockchain.NewService[
		*dastore.Store[types.BeaconBlockBody],
		*types.BeaconBlock,
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
	](
		storageBackend,
		logger.With("service", "blockchain"),
		chainSpec,
		executionEngine,
		localBuilder,
		blobProcessor,
		stateProcessor,
		telemetrySink,
		&blockFeed,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		cfg.Validator.EnableOptimisticPayloadBuilds,
	)

	// Build the deposit service.
	depositService := deposit.NewService[
		types.BeaconBlockBody,
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		*depositdb.KVStore[*types.Deposit],
		*types.ExecutionPayload,
		event.Subscription,
	](
		logger.With("service", "deposit"),
		math.U64(chainSpec.Eth1FollowDistance()),
		engineClient,
		telemetrySink,
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
		service.WithService(version.NewReportingService(
			logger,
			telemetrySink,
			sdkversion.Version,
		)),
		service.WithService(dbManagerService),
	)

	// Pass all the services and options into the BeaconKitRuntime.
	return runtime.NewBeaconKitRuntime[
		*dastore.Store[types.BeaconBlockBody],
		*types.BeaconBlock,
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		blockchain.StorageBackend[
			*dastore.Store[types.BeaconBlockBody],
			types.BeaconBlockBody,
			BeaconState,
			*datypes.BlobSidecars,
			*types.Deposit,
			*depositdb.KVStore[*types.Deposit],
		],
	](
		chainSpec,
		logger,
		svcRegistry,
		storageBackend,
		telemetrySink,
	)
}
