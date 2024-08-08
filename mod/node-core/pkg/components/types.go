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
	"github.com/berachain/beacon-kit/mod/async/pkg/dispatcher"
	"github.com/berachain/beacon-kit/mod/async/pkg/messaging"
	asyncserver "github.com/berachain/beacon-kit/mod/async/pkg/server"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	blockstore "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
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
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// ABCIMiddleware is a type alias for the ABCIMiddleware.
	ABCIMiddleware = middleware.ABCIMiddleware[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBundle,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*Genesis,
		*SlotData,
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
	BeaconBlockBundle = datypes.BlockBundle
	BeaconBlockBody   = types.BeaconBlockBody
	BeaconBlockHeader = types.BeaconBlockHeader

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

	// Context is a type alias for the transition context.
	Context = transition.Context

	// DAService is a type alias for the DA service.
	DAService = da.Service[
		*AvailabilityStore,
		*BeaconBlockBody,
		*BlobSidecars,
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
		*FinalizedBlockEvent,
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
		engineprimitives.Withdrawals,
	]

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

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*types.AttestationData,
		*types.SlashingInfo,
	]

	// IndexDB is a type alias for the range DB.
	IndexDB = filedb.RangeDB

	// KVStore is a type alias for the KV store.
	KVStore = beacondb.KVStore[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		Validators,
	]

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
		*KVStore,
		*Validator,
		Validators,
		*Withdrawal,
		engineprimitives.Withdrawals,
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

	// Validator is a type alias for the validator.
	Validator = types.Validator

	// Validators is a type alias for the validators.
	Validators = types.Validators

	// ValidatorService is a type alias for the validator service.
	ValidatorService = validator.Service[
		*AttestationData,
		*BeaconBlock,
		*BeaconBlockBundle,
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
/*                                   Messages                                 */
/* -------------------------------------------------------------------------- */

// events
type (
	// FinalizedBlockEvent is a type alias for the block event.
	FinalizedBlockEvent = asynctypes.Event[*BeaconBlock]
)

// messages
type (
	// BlockMessage is a type alias for the block message.
	BlockMessage = asynctypes.Message[*BeaconBlock]

	// BlockBundleMessage is a type alias for the block bundle message.
	BlockBundleMessage = asynctypes.Message[*BeaconBlockBundle]

	// GenesisMessage is a type alias for the genesis message.
	GenesisMessage = asynctypes.Message[*Genesis]

	// SidecarMessage is a type alias for the sidecar message.
	SidecarMessage = asynctypes.Message[*BlobSidecars]

	// SlotMessage is a type alias for the slot message.
	SlotMessage = asynctypes.Message[*SlotData]

	// StatusMessage is a type alias for the status message.
	StatusMessage = asynctypes.Message[*service.StatusEvent]

	// ValidatorUpdateMessage is a type alias for the validator update message.
	ValidatorUpdateMessage = asynctypes.Message[*transition.ValidatorUpdates]
)

/* -------------------------------------------------------------------------- */
/*                                   Publishers                               */
/* -------------------------------------------------------------------------- */

type (
	BeaconBlockFinalizedPublisher = messaging.Publisher[*FinalizedBlockEvent]
)

/* -------------------------------------------------------------------------- */
/*                                  Routes                                    */
/* -------------------------------------------------------------------------- */

type (
	// BuildBlockAndSidecarsRoute is a type alias for the build block and sidecars route.
	BuildBlockAndSidecarsRoute = messaging.Route[*SlotMessage, *BlockBundleMessage]

	// VerifyBlockRoute is a type alias for the verify block route.
	VerifyBlockRoute = messaging.Route[*BlockMessage, *BlockMessage]

	// FinalizeBlockRoute is a type alias for the finalize block route.
	FinalizeBlockRoute = messaging.Route[*BlockMessage, *ValidatorUpdateMessage]

	// ProcessGenesisDataRoute is a type alias for the process genesis data route.
	ProcessGenesisDataRoute = messaging.Route[*GenesisMessage, *ValidatorUpdateMessage]

	// ProcessSidecarsRoute is a type alias for the process blob sidecars route.
	ProcessSidecarsRoute = messaging.Route[*SidecarMessage, *SidecarMessage]

	// VerifySidecarsRoute is a type alias for the verify sidecars route.
	VerifySidecarsRoute = messaging.Route[*SidecarMessage, *SidecarMessage]
)

/* -------------------------------------------------------------------------- */
/*                                   Dispatcher                               */
/* -------------------------------------------------------------------------- */

type (
	// Dispatcher is a type alias for the dispatcher.
	Dispatcher = dispatcher.Dispatcher

	// EventServer is a type alias for the event server.
	EventServer = asyncserver.EventServer

	// MessageServer is a type alias for the messages server.
	MessageServer = asyncserver.MessageServer
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

	// ProofAPIHandler is a type alias for the proof handler.
	ProofAPIHandler = proofapi.Handler[
		NodeAPIContext, *BeaconBlockHeader, *BeaconState,
		*BeaconStateMarshallable, *ExecutionPayloadHeader, *Validator,
	]
)
