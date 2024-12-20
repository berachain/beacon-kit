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

package main

import (
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	"github.com/berachain/beacon-kit/consensus-types/types"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	consruntimetypes "github.com/berachain/beacon-kit/consensus/types"
	dablob "github.com/berachain/beacon-kit/da/blob"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	engineclient "github.com/berachain/beacon-kit/execution/client"
	"github.com/berachain/beacon-kit/execution/deposit"
	execution "github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/services/version"
	"github.com/berachain/beacon-kit/payload/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/block"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
	"github.com/berachain/beacon-kit/storage/filedb"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* -------------------------------------------------------------------------- */
/*                                  Services                                  */
/* -------------------------------------------------------------------------- */

type (
	// AttributesFactory is a type alias for the attributes factory.
	AttributesFactory = attributes.Factory

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = dablob.Processor[
		*AvailabilityStore,
		*ConsensusSidecars,
	]

	// ChainService is a type alias for the chain service.
	ChainService = blockchain.Service[
		*AvailabilityStore,
		*DepositStore,
		*ConsensusBlock,
		*BlockStore,
		*Genesis,
		*ConsensusSidecars,
	]

	// CometBFTService is a type alias for the CometBFT service.
	CometBFTService = cometbft.Service[*Logger]

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine

	// IndexDB is a type alias for the range DB.
	IndexDB = filedb.RangeDB

	// KVStore is a type alias for the KV store.
	KVStore = beacondb.KVStore

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = payloadbuilder.PayloadBuilder

	// NodeAPIEngine is a type alias for the node API engine.
	NodeAPIEngine = echo.Engine

	// NodeAPIServer is a type alias for the node API server.
	NodeAPIServer = server.Server[NodeAPIContext]

	// ReportingService is a type alias for the reporting service.
	ReportingService = version.ReportingService

	// SidecarFactory is a type alias for the sidecar factory.
	SidecarFactory = dablob.SidecarFactory

	// StateProcessor is the type alias for the state processor interface.
	StateProcessor = core.StateProcessor[
		*Context,
		*KVStore,
	]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = storage.Backend[
		*AvailabilityStore,
		*BlockStore,
		*DepositStore,
		*KVStore,
	]

	// ValidatorService is a type alias for the validator service.
	ValidatorService = validator.Service[*DepositStore]
)

/* -------------------------------------------------------------------------- */
/*                                    Types                                   */
/* -------------------------------------------------------------------------- */

type (
	// AvailabilityStore is a type alias for the availability store.
	AvailabilityStore = dastore.Store

	// BeaconBlock type aliases.
	ConsensusBlock = consruntimetypes.ConsensusBlock

	// BlobSidecars type aliases.
	ConsensusSidecars = consruntimetypes.ConsensusSidecars
	BlobSidecar       = datypes.BlobSidecar
	BlobSidecars      = datypes.BlobSidecars

	// BlockStore is a type alias for the block store.
	BlockStore = block.KVStore[*types.BeaconBlock]

	// Context is a type alias for the transition context.
	Context = transition.Context

	// DepositContract is a type alias for the deposit contract.
	DepositContract = deposit.WrappedDepositContract

	// DepositStore is a type alias for the deposit store.
	DepositStore = depositdb.KVStore

	// Eth1Data is a type alias for the eth1 data.
	Eth1Data = types.Eth1Data

	// Fork is a type alias for the fork.
	Fork = types.Fork

	// ForkData is a type alias for the fork data.
	ForkData = types.ForkData

	// Genesis is a type alias for the Genesis type.
	Genesis = types.Genesis

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// LoggerConfig is a type alias for the logger config.
	LoggerConfig = phuslu.Config

	// SlotData is a type alias for the incoming slot.
	SlotData = consruntimetypes.SlotData

	// LegacyKey type alias to LegacyKey used for LegacySinger construction.
	LegacyKey = signer.LegacyKey

	// NodeAPIBackend is a type alias for the node API backend.
	NodeAPIBackend = backend.Backend[
		*AvailabilityStore,
		*BlockStore,
		sdk.Context,
		*DepositStore,
		*CometBFTService,
		*KVStore,
		*StorageBackend,
	]

	// NodeAPIContext is a type alias for the node API context.
	NodeAPIContext = echo.Context
)
