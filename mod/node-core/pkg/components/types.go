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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	consruntimetypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	statedb "github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
)

type (
	// AttestationData is a type alias for the attestation data.
	AttestationData = types.AttestationData

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

	BeaconStore = beacondb.Store[
		*BeaconBlockHeader,
		*Eth1Data,
		*ExecutionPayloadHeader,
		*Fork,
		*Validator,
		Validators,
	]

	// BlobSidecars is a type alias for the blob sidecars.
	BlobSidecars = datypes.BlobSidecars

	// BlockStore is a type alias for the block store.
	BlockStore = block.KVStore[*BeaconBlock]

	Context = transition.Context

	// Deposit is a type alias for the deposit.
	Deposit = types.Deposit

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

	// SlashingInfo is a type alias for the slashing info.
	SlashingInfo = types.SlashingInfo

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData[
		*AttestationData,
		*SlashingInfo,
	]

	// IndexDB is a type alias for the range DB.
	IndexDB = filedb.RangeDB

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// NodeAPIBackend is a type alias for the node API backend.
	// NodeAPIBackend = backend.Backend[
	// 	*AvailabilityStore,
	// 	*BeaconBlock,
	// 	*BeaconBlockBody,
	// 	*BeaconBlockHeader,
	// 	*BeaconState,
	// 	*BeaconStateMarshallable,
	// 	*BlobSidecars,
	// 	*BlockStore,
	// 	sdk.Context,
	// 	*Deposit,
	// 	*DepositStore,
	// 	*Eth1Data,
	// 	*ExecutionPayloadHeader,
	// 	*Fork,
	// 	nodetypes.Node,
	// 	*BeaconStore,
	// 	*StorageBackend,
	// 	*Validator,
	// 	*Withdrawal,
	// 	WithdrawalCredentials,
	// ]

	// // NodeAPIContext is a type alias for the node API context.
	// NodeAPIContext = echo.Context

	// // NodeAPIEngine is a type alias for the node API engine.
	// NodeAPIEngine = echo.Engine

	// // NodeAPIServer is a type alias for the node API server.
	// NodeAPIServer = server.Server[
	// 	NodeAPIContext,
	// 	*NodeAPIEngine,
	// ]

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

	// Validators is a type alias for the validators.
	Validators = types.Validators

	// Withdrawal is a type alias for the engineprimitives withdrawal.
	Withdrawal = engineprimitives.Withdrawal

	// WithdrawalCredentials is a type alias for the withdrawal credentials.
	WithdrawalCredentials = types.WithdrawalCredentials
)

/* -------------------------------------------------------------------------- */
/*                                API Handlers                                */
/* -------------------------------------------------------------------------- */

// type (
// 	// BeaconAPIHandler is a type alias for the beacon handler.
// 	BeaconAPIHandler = beaconapi.Handler[
// 		*BeaconBlockHeader, NodeAPIContext, *Fork, *Validator,
// 	]

// 	// BuilderAPIHandler is a type alias for the builder handler.
// 	BuilderAPIHandler = builderapi.Handler[NodeAPIContext]

// 	// ConfigAPIHandler is a type alias for the config handler.
// 	ConfigAPIHandler = configapi.Handler[NodeAPIContext]

// 	// DebugAPIHandler is a type alias for the debug handler.
// 	DebugAPIHandler = debugapi.Handler[NodeAPIContext]

// 	// EventsAPIHandler is a type alias for the events handler.
// 	EventsAPIHandler = eventsapi.Handler[NodeAPIContext]

// 	// NodeAPIHandler is a type alias for the node handler.
// 	NodeAPIHandler = nodeapi.Handler[NodeAPIContext]
// )
