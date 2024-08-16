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
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/engines/echo"
	beaconapi "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon"
	builderapi "github.com/berachain/beacon-kit/mod/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/mod/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/mod/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/mod/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/mod/node-api/handlers/node"
	proofapi "github.com/berachain/beacon-kit/mod/node-api/handlers/proof"
	"github.com/berachain/beacon-kit/mod/node-api/server"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/services/version"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/service"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/middleware"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* -------------------------------------------------------------------------- */
/*                                  Services                                  */
/* -------------------------------------------------------------------------- */

type (
	// ABCIMiddleware is a type alias for the ABCIMiddleware.
	ABCIMiddleware = middleware.ABCIMiddleware[
		*AvailabilityStore,
		*BeaconBlock,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*Genesis,
		*SlotData,
	]

	// AttributesFactory is a type alias for the attributes factory.
	AttributesFactory = attributes.Factory[
		*BeaconState,
		*PayloadAttributes,
		*Withdrawal,
	]

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = dablob.Processor[
		*AvailabilityStore,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BlobSidecar,
		*BlobSidecars,
	]

	// BlobVerifier is a type alias for the blob verifier.
	BlobVerifier = dablob.Verifier[
		*BeaconBlockHeader,
		*BlobSidecar,
		*BlobSidecars,
	]

	// BlockStoreService is a type alias for the block store service.
	BlockStoreService = blockstore.Service[*BeaconBlock, *BlockStore]

	// ChainService is a type alias for the chain service.
	ChainService = blockchain.Service[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*Deposit,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Genesis,
		*PayloadAttributes,
		*Withdrawal,
	]

	// ConsensusEngine is a type alias for the consensus engine.
	ConsensusEngine = cometbft.ConsensusEngine[
		*AttestationData,
		*BeaconState,
		*SlashingInfo,
		*SlotData,
		*StorageBackend,
		*ValidatorUpdate,
	]

	// ConsensusMiddleware is a type alias for the consensus middleware.
	ConsensusMiddleware = cometbft.Middleware[
		*AttestationData,
		*SlashingInfo,
		*SlotData,
	]

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

	// DepositService is a type alias for the deposit service.
	DepositService = deposit.Service[
		*BeaconBlock,
		*BeaconBlockBody,
		*BlockEvent,
		*Deposit,
		*ExecutionPayload,
		WithdrawalCredentials,
	]

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
		engineprimitives.Withdrawals,
	]

	// KVStore is a type alias for the KV store.
	KVStore = beacondb.KVStore[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		Validators,
	]

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = payloadbuilder.PayloadBuilder[
		*BeaconState,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*PayloadAttributes,
		PayloadID,
		*Withdrawal,
	]

	// NodeAPIEngine is a type alias for the node API engine.
	NodeAPIEngine = echo.Engine

	// NodeAPIServer is a type alias for the node API server.
	NodeAPIServer = server.Server[
		NodeAPIContext,
		*NodeAPIEngine,
	]

	// ReportingService is a type alias for the reporting service.
	ReportingService = version.ReportingService

	// SidecarFactory is a type alias for the sidecar factory.
	SidecarFactory = dablob.SidecarFactory[
		*BeaconBlock,
		*BeaconBlockBody,
	]

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
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		Withdrawals,
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
		*BlobSidecars,
		*BlockStore,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

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
)

/* -------------------------------------------------------------------------- */
/*                                    Types                                   */
/* -------------------------------------------------------------------------- */

type (
	// AttestationData is a type alias for the attestation data.
	AttestationData = types.AttestationData

	// BeaconState is a type alias for the BeaconState.
	BeaconState = statedb.StateDB[
		*BeaconBlockHeader,
		*BeaconStateMarshallable,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// BeaconStateMarshallable is a type alias for the BeaconState.
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

	// BlobSidecar is a type alias for the blob sidecar.
	BlobSidecar = datypes.BlobSidecar

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = datypes.BlobSidecars

	// // BlockStore is a type alias for the block store.
	// BlockStore = block.KVStore[*BeaconBlock]

	// Context is a type alias for the transition context.
	Context = transition.Context

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

	// DepositContract is a type alias for the deposit contract.
	DepositContract = deposit.WrappedBeaconDepositContract[
		*Deposit,
		WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// Eth1Data is a type alias for the eth1 data.
	Eth1Data = types.Eth1Data

	// ExecutionPayload type aliases.
	ExecutionPayload       = types.ExecutionPayload
	ExecutionPayloadHeader = types.ExecutionPayloadHeader

	// Fork is a type alias for the fork.
	Fork = types.Fork

	// ForkData is a type alias for the fork data.
	ForkData = types.ForkData

	// Genesis is a type alias for the Genesis type.
	Genesis = types.Genesis[
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// LoggerConfig is a type alias for the logger config.
	LoggerConfig = phuslu.Config

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*AttestationData,
		*SlashingInfo,
	]

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// NodeAPIBackend is a type alias for the node API backend.
	NodeAPIBackend = backend.Backend[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		*BeaconState,
		*BeaconStateMarshallable,
		*BlobSidecars,
		*BlockStore,
		sdk.Context,
		*Deposit,
		*DepositStore,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		nodetypes.Node,
		*KVStore,
		*StorageBackend,
		*Validator,
		Validators,
		*Withdrawal,
		WithdrawalCredentials,
	]

	// NodeAPIContext is a type alias for the node API context.
	NodeAPIContext = echo.Context

	// PayloadAttributes is a type alias for the payload attributes.
	PayloadAttributes = engineprimitives.PayloadAttributes[*Withdrawal]

	// PayloadID is a type alias for the payload ID.
	PayloadID = engineprimitives.PayloadID

	// SlashingInfo is a type alias for the slashing info.
	SlashingInfo = types.SlashingInfo

	// Validator is a type alias for the validator.
	Validator = types.Validator

	// Validators is a type alias for the validators.
	Validators = types.Validators

	// ValidatorUpdate is a type alias for the validator update.
	ABCIValidatorUpdate = appmodule.ValidatorUpdate

	// ValidatorUpdate is a type alias for the validator update.
	ValidatorUpdate = transition.ValidatorUpdate

	// ValidatorUpdates is a type alias for the validator updates.
	ValidatorUpdates = transition.ValidatorUpdates

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal

	// Withdrawals is a type alias for the engineprimitives withdrawals.
	Withdrawals = engineprimitives.Withdrawals

	// WithdrawalCredentials is a type alias for the withdrawal credentials.
	WithdrawalCredentials = types.WithdrawalCredentials
)

/* -------------------------------------------------------------------------- */
/*                                   Events                                   */
/* -------------------------------------------------------------------------- */

type (
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

	// ProofAPIHandler is a type alias for the proof handler.
	ProofAPIHandler = proofapi.Handler[
		NodeAPIContext, *BeaconBlockHeader, *BeaconState,
		*BeaconStateMarshallable, *ExecutionPayloadHeader, *Validator,
	]
)
