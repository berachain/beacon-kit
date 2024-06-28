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
	broker "github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/beacon"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/service"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/middleware"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
)

type (
	// ABCIMiddleware is a type alias for the ABCIMiddleware.
	ABCIMiddleware = middleware.ABCIMiddleware[
		*AvailabilityStore,
		*BeaconBlock,
		BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*Genesis,
	]

	// AttributesFactory is a type alias for the attributes factory.
	AttributesFactory = attributes.Factory[
		BeaconState,
		*engineprimitives.PayloadAttributes[*Withdrawal],
	]

	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = dastore.Store[*BeaconBlockBody]

	// BeaconBlock type aliases.
	BeaconBlock       = types.BeaconBlock
	BeaconBlockBody   = types.BeaconBlockBody
	BeaconBlockHeader = types.BeaconBlockHeader

	// BeaconState is a type alias for the BeaconState.
	BeaconState = core.BeaconState[
		*BeaconBlockHeader, *types.Eth1Data,
		*ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *Withdrawal,
	]

	// BeaconStateMarshallable is a type alias for the BeaconStateMarshallable.
	BeaconStateMarshallable = state.BeaconStateMarshallable[
		*BeaconBlockHeader, *types.Eth1Data, *ExecutionPayloadHeader,
		*types.Fork, *types.Validator,
	]

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = datypes.BlobSidecars

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = dablob.Processor[
		*AvailabilityStore,
		*BeaconBlockBody,
	]

	// ChainService is a type alias for the chain service.
	ChainService = blockchain.Service[
		*AvailabilityStore,
		*BeaconBlock,
		*BeaconBlockBody,
		*BeaconBlockHeader,
		BeaconState,
		*BlobSidecars,
		*Deposit,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*Genesis,
		*engineprimitives.PayloadAttributes[*Withdrawal],
		*Withdrawal,
	]

	// DAService is a type alias for the DA service.
	DAService = da.Service[
		*dastore.Store[*BeaconBlockBody],
		*BeaconBlockBody,
		*BlobSidecars,
		*broker.Broker[*SidecarEvent],
		*ExecutionPayload,
	]

	// DBManager is a type alias for the database manager.
	DBManager = manager.DBManager[
		*BeaconBlock,
		*BlockEvent,
	]

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

	// DepositService is a type alias for the deposit service.
	DepositService = deposit.Service[
		*BeaconBlock,
		*BeaconBlockBody,
		*BlockEvent,
		*Deposit,
		*ExecutionPayload,
		types.WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient[
		*ExecutionPayload, *engineprimitives.PayloadAttributes[*Withdrawal]]

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine[
		*ExecutionPayload, *engineprimitives.PayloadAttributes[*Withdrawal],
		engineprimitives.PayloadID, *Withdrawal,
	]

	// ExecutionPayload type aliases.
	ExecutionPayload       = types.ExecutionPayload
	ExecutionPayloadHeader = types.ExecutionPayloadHeader

	// Genesis is a type alias for the genesis.
	Genesis = genesis.Genesis[*Deposit, *ExecutionPayloadHeader]

	// KVStore is a type alias for the KV store.
	KVStore = beacondb.KVStore[
		*BeaconBlockHeader, *types.Eth1Data, *ExecutionPayloadHeader,
		*types.Fork, *types.Validator,
	]

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = payloadbuilder.PayloadBuilder[
		BeaconState,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*engineprimitives.PayloadAttributes[*Withdrawal],
		engineprimitives.PayloadID,
	]

	// StateProcessor is the type alias for the state processor interface.
	StateProcessor = blockchain.StateProcessor[
		*BeaconBlock,
		BeaconState,
		*BlobSidecars,
		*transition.Context,
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = beacon.StorageBackend[
		*AvailabilityStore,
		*BeaconBlockBody,
		BeaconState,
		*BlobSidecars,
		*Deposit,
		*DepositStore,
	]

	// ValidatorService is a type alias for the validator service.
	ValidatorService = validator.Service[
		*BeaconBlock,
		*BeaconBlockBody,
		BeaconState,
		*BlobSidecars,
		*Deposit,
		*DepositStore,
		*types.Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*types.ForkData,
	]

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal
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
	SlotEvent = asynctypes.Event[math.Slot]

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
/*                                  Services                                  */
/* -------------------------------------------------------------------------- */

// PayloadAttributes is the interface for the payload attributes.
type PayloadAttributes[SelfT any] interface {
	engineprimitives.PayloadAttributer
	// New creates a new payload attributes instance.
	New(
		uint32,
		uint64,
		common.Bytes32,
		common.ExecutionAddress,
		[]*engineprimitives.Withdrawal,
		common.Root,
	) (SelfT, error)
}
