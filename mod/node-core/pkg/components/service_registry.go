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
	blockstore "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/middleware"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

// ServiceRegistryInput is the input for the service registry provider.
type ServiceRegistryInput[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, *Deposit,
		*Eth1Data, *ExecutionPayload, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT, BlobSidecarT],
	KVStoreT any,
	LoggerT any,
	NodeAPIContextT NodeAPIContext,
] struct {
	depinject.In
	ABCIService *middleware.ABCIMiddleware[
		BeaconBlockT, BlobSidecarsT, *Genesis, *SlotData,
	]
	BlockStoreService *blockstore.Service[
		BeaconBlockT, BeaconBlockStoreT,
	]
	ChainService *blockchain.Service[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconBlockHeaderT, BeaconStateT, *Deposit, *ExecutionPayload,
		*ExecutionPayloadHeader, *Genesis, *PayloadAttributes,
	]
	DAService      *da.Service[AvailabilityStoreT, BlobSidecarsT]
	DBManager      *DBManager
	DepositService *deposit.Service[
		BeaconBlockT, BeaconBlockBodyT, *Deposit,
		*ExecutionPayload, WithdrawalCredentials,
	]
	Dispatcher       Dispatcher
	EngineClient     *EngineClient
	Logger           LoggerT
	NodeAPIServer    *server.Server[NodeAPIContextT]
	ReportingService *ReportingService
	TelemetrySink    *metrics.TelemetrySink
	ValidatorService *validator.Service[
		*AttestationData, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, *Deposit, *DepositStore,
		*Eth1Data, *ExecutionPayload, *ExecutionPayloadHeader,
		*ForkData, *SlashingInfo, *SlotData,
	]
}

// ProvideServiceRegistry is the depinject provider for the service registry.
func ProvideServiceRegistry[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, *Deposit,
		*Eth1Data, *ExecutionPayload, *SlashingInfo,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarsT, BlobSidecarT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	NodeAPIContextT NodeAPIContext,
](
	in ServiceRegistryInput[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconBlockHeaderT, BeaconBlockStoreT, BeaconStateT,
		BeaconStateMarshallableT, BlobSidecarT, BlobSidecarsT,
		KVStoreT, LoggerT, NodeAPIContextT,
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
	)
}
