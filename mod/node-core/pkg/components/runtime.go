// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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
	"cosmossdk.io/depinject"
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

// RuntimeInput is the input for the runtime provider.
type RuntimeInput struct {
	depinject.In
	Cfg           *config.Config
	BlobProcessor *dablob.Processor[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
	]
	BlockFeed        *event.FeedOf[*feed.Event[*types.BeaconBlock]]
	ChainSpec        primitives.ChainSpec
	DBManagerService *manager.DBManager[
		*types.BeaconBlock,
		*feed.Event[*types.BeaconBlock],
		event.Subscription,
	]
	DepositService *deposit.Service[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		*feed.Event[*types.BeaconBlock],
		*types.Deposit,
		*types.ExecutionPayload,
		event.Subscription,
		types.WithdrawalCredentials,
	]
	Signer          crypto.BLSSigner
	EngineClient    *engineclient.EngineClient[*types.ExecutionPayload]
	ExecutionEngine *execution.Engine[*types.ExecutionPayload]
	StateProcessor  blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	]
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	LocalBuilder *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	]
	TelemetrySink *metrics.TelemetrySink
	Logger        log.Logger
}

// ProvideRuntime is a depinject provider that returns a BeaconKitRuntime.
func ProvideRuntime(
	in RuntimeInput,
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
		&in.Cfg.Validator,
		in.Logger.With("service", "validator"),
		in.ChainSpec,
		in.StorageBackend,
		in.BlobProcessor,
		in.StateProcessor,
		in.Signer,
		dablob.NewSidecarFactory[
			*types.BeaconBlock,
			*types.BeaconBlockBody,
		](
			in.ChainSpec,
			types.KZGPositionDeneb,
			in.TelemetrySink,
		),
		in.LocalBuilder,
		[]validator.PayloadBuilder[BeaconState, *types.ExecutionPayload]{
			in.LocalBuilder,
		},
		in.TelemetrySink,
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
		in.StorageBackend,
		in.Logger.With("service", "blockchain"),
		in.ChainSpec,
		in.ExecutionEngine,
		in.LocalBuilder,
		in.BlobProcessor,
		in.StateProcessor,
		in.TelemetrySink,
		in.BlockFeed,
		// If optimistic is enabled, we want to skip post finalization FCUs.
		in.Cfg.Validator.EnableOptimisticPayloadBuilds,
	)
	// Build the service registry.
	svcRegistry := service.NewRegistry(
		service.WithLogger(in.Logger.With("service", "service-registry")),
		service.WithService(validatorService),
		service.WithService(chainService),
		service.WithService(in.DepositService),
		service.WithService(in.EngineClient),
		service.WithService(version.NewReportingService(
			in.Logger,
			in.TelemetrySink,
			sdkversion.Version,
		)),
		service.WithService(in.DBManagerService),
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
		in.ChainSpec,
		in.Logger,
		svcRegistry,
		in.StorageBackend,
		in.TelemetrySink,
	)
}
