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
	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/beacon"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/payload/pkg/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
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
	AttributesFactory = attributes.Factory[BeaconState, *Withdrawal]

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

	// BlobFeed is a type alias for the blob feed.
	BlobFeed = event.FeedOf[asynctypes.EventID, *asynctypes.Event[*BlobSidecars]]

	// BlockEvent is a type alias for the block event.
	BlockEvent = asynctypes.Event[*BeaconBlock]

	// BlockFeed is a type alias for the block feed.
	BlockFeed = event.FeedOf[asynctypes.EventID, *BlockEvent]

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
	]

	// DBManager is a type alias for the database manager.
	DBManager = manager.DBManager[
		*BeaconBlock,
		*BlockEvent,
		event.Subscription,
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
		event.Subscription,
		types.WithdrawalCredentials,
	]

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore[*Deposit]

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient[*ExecutionPayload]

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine[
		*ExecutionPayload, engineprimitives.PayloadID,
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
		BeaconState, *ExecutionPayload,
		*ExecutionPayloadHeader, engineprimitives.PayloadID,
	]

	// StateProcessor is the type alias for the state processor inteface.
	StateProcessor = blockchain.StateProcessor[
		*BeaconBlock,
		BeaconState,
		*BlobSidecars,
		*transition.Context,
		*Deposit,
		*ExecutionPayloadHeader,
	]

	// SlotFeed is a type alias for the slot feed.
	SlotFeed = event.FeedOf[asynctypes.EventID, *asynctypes.Event[math.Slot]]

	// StatusFeed is a type alias for the status feed.
	StatusFeed = event.FeedOf[
		asynctypes.EventID, *asynctypes.Event[*service.StatusEvent],
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
