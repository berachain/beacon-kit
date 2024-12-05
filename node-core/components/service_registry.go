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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/middleware"
	"github.com/berachain/beacon-kit/da/da"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/deposit"
	"github.com/berachain/beacon-kit/log"
	blockstore "github.com/berachain/beacon-kit/node-api/block_store"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	service "github.com/berachain/beacon-kit/node-core/services/registry"
	"github.com/berachain/beacon-kit/node-core/services/version"
	"github.com/berachain/beacon-kit/observability/telemetry"
)

// ServiceRegistryInput is the input for the service registry provider.
type ServiceRegistryInput[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	ConsensusBlockT ConsensusBlock[BeaconBlockT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		*Eth1Data, ExecutionPayloadT, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, ExecutionPayloadHeaderT, *Fork, KVStoreT,
		*Validator, Validators, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT, BlobSidecarT],
	DepositT Deposit[DepositT, *ForkData, WithdrawalCredentials],
	DepositStoreT DepositStore[DepositT],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
	NodeAPIContextT NodeAPIContext,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
] struct {
	depinject.In
	ABCIService *middleware.ABCIMiddleware[
		BeaconBlockT, BeaconBlockHeaderT, BlobSidecarsT, GenesisT, *SlotData,
	]
	BlockStoreService *blockstore.Service[
		BeaconBlockT, BeaconBlockStoreT,
	]
	ChainService *blockchain.Service[
		AvailabilityStoreT,
		ConsensusBlockT, BeaconBlockT, BeaconBlockBodyT,
		BeaconBlockHeaderT, BeaconStateT, DepositT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, GenesisT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	]
	DAService *da.Service[
		AvailabilityStoreT,
		ConsensusSidecarsT, BlobSidecarsT, BeaconBlockHeaderT,
	]
	DBManager      *DBManager
	DepositService *deposit.Service[
		BeaconBlockT, BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, WithdrawalCredentials,
	]
	Dispatcher   Dispatcher
	EngineClient *client.EngineClient[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	]
	Logger           LoggerT
	NodeAPIServer    *server.Server[NodeAPIContextT]
	ReportingService *version.ReportingService[
		ExecutionPayloadT,
		*engineprimitives.PayloadAttributes[WithdrawalT],
	]
	TelemetrySink    *metrics.TelemetrySink
	TelemetryService *telemetry.Service
	ValidatorService *validator.Service[
		*AttestationData, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, DepositT, DepositStoreT,
		*Eth1Data, ExecutionPayloadT, ExecutionPayloadHeaderT,
		*ForkData, *SlashingInfo, *SlotData,
	]
	CometBFTService *cometbft.Service[LoggerT]
}

// ProvideServiceRegistry is the depinject provider for the service registry.
func ProvideServiceRegistry[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	ConsensusBlockT ConsensusBlock[BeaconBlockT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		*Eth1Data, ExecutionPayloadT, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, ExecutionPayloadHeaderT, *Fork, KVStoreT,
		*Validator, Validators, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	ConsensusSidecarsT ConsensusSidecars[BlobSidecarsT, BeaconBlockHeaderT],
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT, BlobSidecarT],
	DepositT Deposit[DepositT, *ForkData, WithdrawalCredentials],
	DepositStoreT DepositStore[DepositT],
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT,
		ExecutionPayloadHeaderT, WithdrawalsT],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[LoggerT],
	NodeAPIContextT NodeAPIContext,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
](
	in ServiceRegistryInput[
		AvailabilityStoreT,
		ConsensusBlockT, BeaconBlockT, BeaconBlockBodyT,
		BeaconBlockHeaderT, BeaconBlockStoreT, BeaconStateT,
		BeaconStateMarshallableT,
		ConsensusSidecarsT, BlobSidecarT, BlobSidecarsT,
		DepositT, DepositStoreT, ExecutionPayloadT, ExecutionPayloadHeaderT,
		GenesisT, KVStoreT, LoggerT, NodeAPIContextT, WithdrawalT, WithdrawalsT,
	],
) *service.Registry {
	return service.NewRegistry(
		service.WithLogger(in.Logger),
		service.WithService(in.ABCIService),
		service.WithService(in.Dispatcher),
		service.WithService(in.ValidatorService),
		service.WithService(in.BlockStoreService),
		service.WithService(in.ChainService),
		service.WithService(in.DAService),
		service.WithService(in.DepositService),
		service.WithService(in.NodeAPIServer),
		service.WithService(in.ReportingService),
		service.WithService(in.DBManager),
		service.WithService(in.EngineClient),
		service.WithService(in.TelemetryService),
		service.WithService(in.CometBFTService),
	)
}
