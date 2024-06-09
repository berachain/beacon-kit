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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
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
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/event"
)

type BeaconState = core.BeaconState[
	*types.BeaconBlockHeader, *types.Eth1Data,
	*types.ExecutionPayloadHeader, *types.Fork,
	*types.Validator, *engineprimitives.Withdrawal,
]

// BeaconKitRuntime is a type alias for the BeaconKitRuntime.
type BeaconKitRuntime = runtime.BeaconKitRuntime[
	*dastore.Store[*types.BeaconBlockBody],
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
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
	blobProcessor *dablob.Processor[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
	],
	blockFeed *event.FeedOf[*feed.Event[*types.BeaconBlock]],
	chainSpec primitives.ChainSpec,
	dbManagerService *manager.DBManager[
		*types.BeaconBlock,
		*feed.Event[*types.BeaconBlock],
		event.Subscription,
	],
	depositService *deposit.Service[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*feed.Event[*types.BeaconBlock],
		*types.Deposit,
		*types.ExecutionPayload,
		event.Subscription,
		types.WithdrawalCredentials,
	],
	signer crypto.BLSSigner,
	engineClient *engineclient.EngineClient[*types.ExecutionPayload],
	executionEngine *execution.Engine[*types.ExecutionPayload],
	stateProcessor blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	],
	storageBackend blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	],
	localBuilder *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	],
	telemetrySink *metrics.TelemetrySink,
	logger log.Logger,
) (*BeaconKitRuntime, error) {
	// Build the builder service.
	validatorService := validator.NewService[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
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
			*types.BeaconBlockBody,
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
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
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
		blockFeed,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		cfg.Validator.EnableOptimisticPayloadBuilds,
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
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		blockchain.StorageBackend[
			*dastore.Store[*types.BeaconBlockBody],
			*types.BeaconBlockBody,
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
