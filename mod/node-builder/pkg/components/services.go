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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/services/version"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/event"
)

// TODO: MOVE THESE INTO SEPARATE FILES

// ==============================dbmanager=============================

// DBManagerInput is the input for the DBManager service through the dep inject
// framework.
type DBManagerInput struct {
	depinject.In
	Logger  log.Logger
	Pruners []*pruner.Pruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription]
}

// ProvideDBManager provides the DBManager service through the dep inject
// framework.
func ProvideDBManager(
	in DBManagerInput,
) (*manager.DBManager[*types.BeaconBlock,
	events.Block[*types.BeaconBlock],
	event.Subscription], error) {
	// Build the DBManager service.
	dbm, err := manager.NewDBManager[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription,
	](
		in.Logger.With("service", "db-manager"),
		in.Pruners...,
	)
	if err != nil {
		return nil, err
	}
	return dbm, nil
}

// ==============================blockchain=============================

type ChainServiceInput struct {
	depinject.In
	Logger         log.Logger
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	BlockFeed      *event.FeedOf[events.Block[*types.BeaconBlock]]
	StateProcessor blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	]
	ChainSpec       primitives.ChainSpec
	ExecutionEngine *execution.Engine[*types.ExecutionPayload]
	LocalBuilder    *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	]
	BlobProcessor *dablob.Processor[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
	]
	TelemetrySink *metrics.TelemetrySink
	Cfg           *config.Config
}

func ProvideChainService(
	in ChainServiceInput,
) *blockchain.Service[
	*dastore.Store[types.BeaconBlockBody],
	*types.BeaconBlock,
	types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*types.Deposit,
	*depositdb.KVStore[*types.Deposit],
] {
	// Build the blockchain service.
	return blockchain.NewService[
		*dastore.Store[types.BeaconBlockBody],
		*types.BeaconBlock,
		types.BeaconBlockBody,
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
}

// ==============================validator=============================

type ValidatorServiceInput struct {
	depinject.In
	Cfg            *config.Config
	Logger         log.Logger
	ChainSpec      primitives.ChainSpec
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	StateProcessor blockchain.StateProcessor[
		*types.BeaconBlock,
		BeaconState,
		*datypes.BlobSidecars,
		*transition.Context,
		*types.Deposit,
	]
	Signer       crypto.BLSSigner
	LocalBuilder *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	]
	BlobProcessor *dablob.Processor[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
	]
	TelemetrySink *metrics.TelemetrySink
}

func ProvideValidatorService(
	in ValidatorServiceInput,
) *validator.Service[
	*types.BeaconBlock,
	types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	*types.ForkData,
] {
	return validator.NewService[
		*types.BeaconBlock,
		types.BeaconBlockBody,
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
			types.BeaconBlockBody,
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
}

// ==============================deposit=============================

type DepositServiceInput struct {
	depinject.In
	Logger         log.Logger
	ChainSpec      primitives.ChainSpec
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	BeaconDepositContract *deposit.WrappedBeaconDepositContract[
		*types.Deposit, types.WithdrawalCredentials,
	]
	BlockFeed     *event.FeedOf[events.Block[*types.BeaconBlock]]
	EngineClient  *engineclient.EngineClient[*types.ExecutionPayload]
	TelemetrySink *metrics.TelemetrySink
}

func ProvideDepositService(
	in DepositServiceInput,
) *deposit.Service[
	*types.BeaconBlock,
	types.BeaconBlockBody,
	events.Block[*types.BeaconBlock],
	*types.Deposit,
	*types.ExecutionPayload,
	event.Subscription,
	types.WithdrawalCredentials,
] {
	// Build the deposit service.
	return deposit.NewService[
		types.BeaconBlockBody,
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		*depositdb.KVStore[*types.Deposit],
		*types.ExecutionPayload,
		event.Subscription,
	](
		in.Logger.With("service", "deposit"),
		math.U64(in.ChainSpec.Eth1FollowDistance()),
		in.EngineClient,
		in.TelemetrySink,
		in.StorageBackend.DepositStore(nil),
		in.BeaconDepositContract,
		in.BlockFeed,
	)
}

// ==============================reporting=============================.
type ReportingServiceInput struct {
	depinject.In
	Logger        log.Logger
	TelemetrySink *metrics.TelemetrySink
}

func ProvideReportingService(
	in ReportingServiceInput,
) *version.ReportingService {
	return version.NewReportingService(
		in.Logger,
		in.TelemetrySink,
		sdkversion.Version,
	)
}
