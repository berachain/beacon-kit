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
	"cosmossdk.io/core/appmodule/v2"
	broker "github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	blockstore "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-api/engines/echo"
	beaconapi "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon"
	builderapi "github.com/berachain/beacon-kit/mod/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/mod/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/mod/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/mod/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/mod/node-api/handlers/node"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components/signer"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components/storage"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types/services/version"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/service"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	runtime "github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
)

type (
	RuntimeApp = runtime.App[
		*AttestationData,
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*Genesis,
		*SlashingInfo,
		*SlotData,
		*StorageBackend,
	]

	// AttestationData is a type alias for the attestation data.
	AttestationData = types.AttestationData

	// AttributesFactory is a type alias for the attributes factory.
	AttributesFactory = attributes.Factory[
		*BeaconState,
		*PayloadAttributes,
		*Withdrawal,
	]

	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = dastore.Store[*BeaconBlockBody]

	// BeaconBlock type aliases.
	BeaconBlock       = types.BeaconBlock
	BeaconBlockBody   = types.BeaconBlockBody
	BeaconBlockHeader = types.BeaconBlockHeader

	// BeaconState is a type alias for the BeaconState.
	BeaconState = statedb.StateDB[
		*BeaconBlockHeader,
		*BeaconStateMarshallable,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*BeaconStore,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// BeaconStateMarshallable is a type alias for the BeaconStateMarshallable.
	BeaconStateMarshallable = types.BeaconState[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		BeaconBlockHeader,
		Eth1Data,
		ExecutionPayloadHeader,
		Fork,
		Validator,
	]

	BeaconStore = beacondb.Store[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		Validators,
	]

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = dablob.Processor[
		*AvailabilityStore,
		*BeaconBlockBody,
	]

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = datypes.BlobSidecars

	// BlobVerifier is a type alias for the blob verifier.
	BlobVerifier = dablob.Verifier

	// BlockStoreService is a type alias for the block store service.
	BlockStoreService = blockstore.Service[*BeaconBlock, *BlockStore]

	// BlockStore is a type alias for the block store.
	BlockStore = block.KVStore[*BeaconBlock]

	// ChainService is a type alias for the chain service.
	ChainService = blockchain.Service[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Genesis,
		*PayloadAttributes,
		*Withdrawal,
	]

	// Context is a type alias for the transition context.
	Context = transition.Context

	// DAService is a type alias for the DA service.
	DAService = da.Service[
		*AvailabilityStore,
		*BeaconBlockBody,
		*BlobSidecars,
		*SidecarsBroker,
		*ExecutionPayload,
	]

	// DBManager is a type alias for the database manager.
	DBManager = manager.DBManager

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

	// DepositContract is a type alias for the deposit contract.
	DepositContract = deposit.WrappedBeaconDepositContract[
		*Deposit,
		WithdrawalCredentials,
	]

	// DepositService is a type alias for the deposit service.
	DepositService = deposit.Service[
		*BeaconBlock,
		*BeaconBlockBody,
		*BlockEvent,
		*Deposit,
		*ExecutionPayload,
		WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// Eth1Data is a type alias for the eth1 data.
	Eth1Data = types.Eth1Data

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient[
		*ExecutionPayload,
		*PayloadAttributes,
	]

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine[
		*ExecutionPayload,
		*PayloadAttributes,
		PayloadID,
		*Withdrawal,
	]

	// ExecutionPayload type aliases.
	ExecutionPayload       = types.ExecutionPayload
	ExecutionPayloadHeader = types.ExecutionPayloadHeader

	// Fork is a type alias for the fork.
	Fork = types.Fork

	// ForkData is a type alias for the fork data.
	ForkData = types.ForkData

	// Genesis is a type alias for the genesis.
	Genesis = types.Genesis[
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*types.AttestationData,
		*types.SlashingInfo,
	]

	// IndexDB is a type alias for the range DB.
	IndexDB = filedb.RangeDB

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = payloadbuilder.PayloadBuilder[
		*BeaconState,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*PayloadAttributes,
		PayloadID,
		*Withdrawal,
	]

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// NodeAPIContext is a type alias for the node API context.
	NodeAPIContext = echo.Context

	// NodeAPIEngine is a type alias for the node API engine.
	NodeAPIEngine = echo.Engine

	// NodeAPIServer is a type alias for the node API server.
	NodeAPIServer = server.Server[
		NodeAPIContext,
		*NodeAPIEngine,
	]

	// PayloadAttributes is a type alias for the payload attributes.
	PayloadAttributes = engineprimitives.PayloadAttributes[*Withdrawal]

	// PayloadID is a type alias for the payload ID.
	PayloadID = engineprimitives.PayloadID

	// ReportingService is a type alias for the reporting service.
	ReportingService = version.ReportingService

	// SidecarFactory is a type alias for the sidecar factory.
	SidecarFactory = dablob.SidecarFactory[
		*BeaconBlock,
		*BeaconBlockBody,
	]

	// SlashingInfo is a type alias for the slashing info.
	SlashingInfo = types.SlashingInfo

	// StateProcessor is the type alias for the state processor interface.
	StateProcessor = core.StateProcessor[
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*Context,
		*Deposit,
		*Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Fork,
		*ForkData,
		*BeaconStore,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = storage.Backend[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BeaconStore,
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// Validator is a type alias for the validator.
	Validator = types.Validator

	Validators = types.Validators

	// ValidatorService is a type alias for the validator service.
	ValidatorService = validator.Service[
		*AttestationData,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconState,
		*BlobSidecars,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*ForkData,
		*SlashingInfo,
		*SlotData,
	]

	// ValidatorUpdate is a type alias for the validator update.
	ValidatorUpdate = appmodule.ValidatorUpdate

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal

	// WithdrawalCredentials is a type alias for the withdrawal credentials.
	WithdrawalCredentials = types.WithdrawalCredentials
)

/* -------------------------------------------------------------------------- */
/*                                   Events                                   */
/* -------------------------------------------------------------------------- */

type (
	// BlockEvent is a type alias for the block event.
	BlockEvent = asynctypes.Event[*BeaconBlock]

	// GenesisEvent is a type alias for the genesis event.
	GenesisEvent = asynctypes.Event[*Genesis]

	// SidecarEvent is a type alias for the sidecar event.
	SidecarEvent = asynctypes.Event[*BlobSidecars]

	// SlotEvent is a type alias for the slot event.
	SlotEvent = asynctypes.Event[*SlotData]

	// StatusEvent is a type alias for the status event.
	StatusEvent = asynctypes.Event[*service.StatusEvent]

	// ValidatorUpdateEvent is a type alias for the validator update event.
	ValidatorUpdateEvent = asynctypes.Event[transition.ValidatorUpdates]
)

/* -------------------------------------------------------------------------- */
/*                                   Brokers                                  */
/* -------------------------------------------------------------------------- */

type (
	// GenesisBroker is a type alias for the genesis feed.
	GenesisBroker = broker.Broker[*GenesisEvent]

	// SidecarsBroker is a type alias for the blob feed.
	SidecarsBroker = broker.Broker[*SidecarEvent]

	// BlockBroker is a type alias for the block feed.
	BlockBroker = broker.Broker[*BlockEvent]

	// SlotBroker is a type alias for the slot feed.
	SlotBroker = broker.Broker[*SlotEvent]

	// StatusBroker is a type alias for the status feed.
	StatusBroker = broker.Broker[*StatusEvent]

	// ValidatorUpdateBroker is a type alias for the validator update feed.
	ValidatorUpdateBroker = broker.Broker[*ValidatorUpdateEvent]
)

/* -------------------------------------------------------------------------- */
/*                                  Pruners                                   */
/* -------------------------------------------------------------------------- */

type (
	// DAPruner is a type alias for the DA pruner.
	DAPruner = pruner.Pruner[*IndexDB]

	// DepositPruner is a type alias for the deposit pruner.
	DepositPruner = pruner.Pruner[*DepositStore]

	// BlockPruner is a type alias for the block pruner.
	BlockPruner = pruner.Pruner[*BlockStore]
)

/* -------------------------------------------------------------------------- */
/*                                API Handlers                                */
/* -------------------------------------------------------------------------- */

type (
	// BeaconAPIHandler is a type alias for the beacon handler.
	BeaconAPIHandler = beaconapi.Handler[
		*BeaconBlockHeader, NodeAPIContext, *Fork, *Validator,
	]

	// BuilderAPIHandler is a type alias for the builder handler.
	BuilderAPIHandler = builderapi.Handler[NodeAPIContext]

	// ConfigAPIHandler is a type alias for the config handler.
	ConfigAPIHandler = configapi.Handler[NodeAPIContext]

	// DebugAPIHandler is a type alias for the debug handler.
	DebugAPIHandler = debugapi.Handler[NodeAPIContext]

	// EventsAPIHandler is a type alias for the events handler.
	EventsAPIHandler = eventsapi.Handler[NodeAPIContext]

	// NodeAPIHandler is a type alias for the node handler.
	NodeAPIHandler = nodeapi.Handler[NodeAPIContext]
)
