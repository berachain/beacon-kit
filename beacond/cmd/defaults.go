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
	"github.com/berachain/beacon-kit/beacond/cmd/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
)

//nolint:funlen // happens
func DefaultComponents() []any {
	c := []any{
		components.ProvideABCIMiddleware[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.BlobSidecar, *types.BlobSidecars,
			*types.Deposit, *types.ExecutionPayloadHeader, *types.Genesis,
			*types.Logger,
		],
		components.ProvideAttributesFactory[
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.BeaconStateMarshallable, *types.ExecutionPayloadHeader,
			*types.KVStore, *types.Logger,
		],
		components.ProvideAvailibilityStore[
			*types.BeaconBlockBody, *types.Logger,
		],
		components.ProvideAvailabilityPruner[
			*types.AvailabilityStore, *types.BeaconBlock,
			*types.BeaconBlockBody, *types.BeaconBlockHeader,
			*types.BlobSidecars, *types.Logger,
		],
		components.ProvideBeaconDepositContract[
			*types.Deposit, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader,
		],
		components.ProvideBlockStore[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.Logger,
		],
		components.ProvideBlockStoreService[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.BlockStore, *types.Logger,
		],
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor[
			*types.AvailabilityStore, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.BlobSidecar, *types.BlobSidecars,
			*types.Logger,
		],
		components.ProvideBlobProofVerifier,
		components.ProvideBlobVerifier[
			*types.BeaconBlockHeader, *types.BlobSidecar, *types.BlobSidecars,
		],
		components.ProvideChainService[
			*types.AvailabilityStore, *types.BeaconBlock,
			*types.BeaconBlockBody, *types.BeaconBlockHeader,
			*types.BeaconState, *types.BeaconStateMarshallable,
			*types.BlobSidecars, *types.BlockStore, *types.Deposit,
			*types.DepositStore, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader, *types.Genesis, *types.KVStore,
			*types.Logger, *types.StorageBackend,
		],
		components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideServerConfig,
		// components.ProvideConsensusEngine[
		// 	*AvailabilityStore, *BeaconBlockHeader, *BeaconState,
		// 	*BeaconStateMarshallable, *BlockStore, *KVStore, *StorageBackend,
		// ],
		components.ProvideDAService[
			*types.AvailabilityStore, *types.BeaconBlockBody,
			*types.BlobSidecar, *types.BlobSidecars, *types.Logger,
		],
		components.ProvideDBManager[
			*types.AvailabilityStore, *types.DepositStore, *types.Logger,
		],
		components.ProvideDepositPruner[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.Deposit, *types.DepositStore,
			*types.Logger,
		],
		components.ProvideDepositService[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.Deposit, *types.DepositContract,
			*types.DepositStore, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader, *types.Logger,
		],
		components.ProvideDepositStore[*types.Deposit],
		components.ProvideDispatcher[
			*types.BeaconBlock, *types.BlobSidecars, *types.Genesis,
			*types.Logger,
		],
		components.ProvideEngineClient[
			*types.ExecutionPayload, *types.ExecutionPayloadHeader,
			*types.Logger,
		],
		components.ProvideExecutionEngine[
			*types.ExecutionPayload, *types.ExecutionPayloadHeader,
			*types.Logger,
		],
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder[
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.BeaconStateMarshallable, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader, *types.KVStore, *types.Logger,
		],
		components.ProvideReportingService[*types.Logger],
		components.ProvideCometBFTService[*types.Logger],
		components.ProvideServiceRegistry[
			*types.AvailabilityStore, *types.BeaconBlock,
			*types.BeaconBlockBody, *types.BeaconBlockHeader, *types.BlockStore,
			*types.BeaconState, *types.BeaconStateMarshallable,
			*types.BlobSidecar, *types.BlobSidecars, *types.Deposit,
			*types.DepositStore, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader, *types.Genesis, *types.KVStore,
			*types.Logger, types.NodeAPIContext,
		],
		components.ProvideSidecarFactory[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader,
		],
		components.ProvideStateProcessor[
			*types.BeaconBlock, *types.BeaconBlockBody,
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.BeaconStateMarshallable, *types.Deposit,
			*types.ExecutionPayload, *types.ExecutionPayloadHeader,
			*types.KVStore,
		],
		components.ProvideKVStore[
			*types.BeaconBlockHeader, *types.ExecutionPayloadHeader,
		],
		components.ProvideStorageBackend[
			*types.AvailabilityStore, *types.BlockStore, *types.BeaconState,
			*types.KVStore, *types.DepositStore,
		],
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService[
			*types.AvailabilityStore, *types.BeaconBlock,
			*types.BeaconBlockBody, *types.BeaconBlockHeader,
			*types.BeaconState, *types.BeaconStateMarshallable,
			*types.BlockStore, *types.BlobSidecars, *types.Deposit,
			*types.DepositStore, *types.ExecutionPayload,
			*types.ExecutionPayloadHeader, *types.KVStore, *types.Logger,
			*types.StorageBackend,
		],
		// TODO Hacks
		components.ProvideKVStoreService,
		components.ProvideKVStoreKey,
	}
	c = append(c,
		components.ProvideNodeAPIServer[*types.Logger, types.NodeAPIContext],
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend[
			*types.AvailabilityStore, *types.BeaconBlock,
			*types.BeaconBlockBody, *types.BeaconBlockHeader, *types.BlockStore,
			*types.BeaconState, *types.BeaconStateMarshallable,
			*types.BlobSidecars, *types.Deposit, *types.DepositStore,
			*types.ExecutionPayloadHeader, *types.KVStore,
			*types.CometBFTService, *types.StorageBackend,
		],
	)

	c = append(c,
		components.ProvideNodeAPIHandlers[
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.BeaconStateMarshallable, *types.ExecutionPayloadHeader,
			*types.KVStore, types.NodeAPIContext,
		],
		components.ProvideNodeAPIBeaconHandler[
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.CometBFTService, types.NodeAPIContext,
		],
		components.ProvideNodeAPIBuilderHandler[types.NodeAPIContext],
		components.ProvideNodeAPIConfigHandler[types.NodeAPIContext],
		components.ProvideNodeAPIDebugHandler[types.NodeAPIContext],
		components.ProvideNodeAPIEventsHandler[types.NodeAPIContext],
		components.ProvideNodeAPINodeHandler[types.NodeAPIContext],
		components.ProvideNodeAPIProofHandler[
			*types.BeaconBlockHeader, *types.BeaconState,
			*types.BeaconStateMarshallable, *types.ExecutionPayloadHeader,
			*types.KVStore, *types.CometBFTService, types.NodeAPIContext,
		],
	)

	return c
}
