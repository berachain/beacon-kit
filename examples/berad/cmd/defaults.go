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
	"github.com/berachain/beacon-kit/examples/berad/cmd/components"
	beacondcomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
)

//nolint:funlen // happens
func DefaultComponents() []any {
	c := []any{
		beacondcomponents.ProvideABCIMiddleware[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BlobSidecar, *BlobSidecars, *Deposit, *ExecutionPayloadHeader,
			*Genesis, *Logger,
		],
		beacondcomponents.ProvideAttributesFactory[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, *Logger,
		],
		beacondcomponents.ProvideAvailibilityStore[*BeaconBlockBody, *Logger],
		beacondcomponents.ProvideAvailabilityPruner[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlobSidecars, *Logger,
		],
		beacondcomponents.ProvideBeaconDepositContract[
			*Deposit, *ExecutionPayload, *ExecutionPayloadHeader,
		],
		beacondcomponents.ProvideBlockStore[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader, *Logger,
		],
		beacondcomponents.ProvideBlockStoreService[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BlockStore, *Logger,
		],
		beacondcomponents.ProvideBlsSigner,
		beacondcomponents.ProvideBlobProcessor[
			*AvailabilityStore, *BeaconBlockBody, *BeaconBlockHeader,
			*BlobSidecar, *BlobSidecars, *Logger,
		],
		beacondcomponents.ProvideBlobProofVerifier,
		beacondcomponents.ProvideBlobVerifier[
			*BeaconBlockHeader, *BlobSidecar, *BlobSidecars,
		],
		beacondcomponents.ProvideChainService[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlobSidecars, *BlockStore, *Deposit, *DepositStore,
			*ExecutionPayload, *ExecutionPayloadHeader, *Genesis,
			*KVStore, *Logger, *StorageBackend,
		],
		beacondcomponents.ProvideNode,
		beacondcomponents.ProvideChainSpec,
		beacondcomponents.ProvideConfig,
		beacondcomponents.ProvideServerConfig,
		// beacondcomponents.ProvideConsensusEngine[
		// 	*AvailabilityStore, *BeaconBlockHeader, *BeaconState,
		// 	*BeaconStateMarshallable, *BlockStore, *KVStore, *StorageBackend,
		// ],
		beacondcomponents.ProvideDAService[
			*AvailabilityStore, *BeaconBlockBody, *BlobSidecar,
			*BlobSidecars, *Logger,
		],
		beacondcomponents.ProvideDBManager[*AvailabilityStore, *DepositStore, *Logger],
		beacondcomponents.ProvideDepositPruner[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*Deposit, *DepositStore, *Logger,
		],
		beacondcomponents.ProvideDepositService[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader, *Deposit,
			*DepositContract, *DepositStore, *ExecutionPayload,
			*ExecutionPayloadHeader, *Logger,
		],
		beacondcomponents.ProvideDepositStore[*Deposit],
		beacondcomponents.ProvideDispatcher[
			*BeaconBlock, *BlobSidecars, *Genesis, *Logger,
		],
		beacondcomponents.ProvideEngineClient[
			*ExecutionPayload, *ExecutionPayloadHeader, *Logger,
		],
		beacondcomponents.ProvideExecutionEngine[
			*ExecutionPayload, *ExecutionPayloadHeader, *Logger,
		],
		beacondcomponents.ProvideJWTSecret,
		beacondcomponents.ProvideLocalBuilder[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*ExecutionPayload, *ExecutionPayloadHeader, *KVStore, *Logger,
		],
		beacondcomponents.ProvideReportingService[*Logger],
		beacondcomponents.ProvideCometBFTService[*Logger],
		beacondcomponents.ProvideServiceRegistry[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockStore, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecar, *BlobSidecars,
			*Deposit, *DepositStore, *ExecutionPayload, *ExecutionPayloadHeader,
			*Genesis, *KVStore, *Logger,
			NodeAPIContext,
		],
		beacondcomponents.ProvideSidecarFactory[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
		],
		components.ProvideStateProcessor[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BeaconState, *BeaconStateMarshallable, *Deposit, *ExecutionPayload,
			*ExecutionPayloadHeader, *KVStore,
		],
		beacondcomponents.ProvideKVStore[*BeaconBlockHeader, *ExecutionPayloadHeader],
		beacondcomponents.ProvideStorageBackend[
			*AvailabilityStore, *BlockStore, *BeaconState,
			*KVStore, *DepositStore,
		],
		beacondcomponents.ProvideTelemetrySink,
		beacondcomponents.ProvideTelemetryService,
		beacondcomponents.ProvideTrustedSetup,
		beacondcomponents.ProvideValidatorService[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlockStore, *BlobSidecars, *Deposit, *DepositStore,
			*ExecutionPayload, *ExecutionPayloadHeader, *KVStore, *Logger,
			*StorageBackend,
		],
		// TODO Hacks
		beacondcomponents.ProvideKVStoreService,
		beacondcomponents.ProvideKVStoreKey,
	}
	c = append(c,
		beacondcomponents.ProvideNodeAPIServer[*Logger, NodeAPIContext],
		beacondcomponents.ProvideNodeAPIEngine,
		beacondcomponents.ProvideNodeAPIBackend[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockStore, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecars, *Deposit, *DepositStore,
			*ExecutionPayloadHeader, *KVStore, *CometBFTService, *StorageBackend,
		],
	)

	c = append(c,
		beacondcomponents.ProvideNodeAPIHandlers[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, NodeAPIContext,
		],
		beacondcomponents.ProvideNodeAPIBeaconHandler[
			*BeaconBlockHeader, *BeaconState, *CometBFTService, NodeAPIContext,
		],
		beacondcomponents.ProvideNodeAPIBuilderHandler[NodeAPIContext],
		beacondcomponents.ProvideNodeAPIConfigHandler[NodeAPIContext],
		beacondcomponents.ProvideNodeAPIDebugHandler[NodeAPIContext],
		beacondcomponents.ProvideNodeAPIEventsHandler[NodeAPIContext],
		beacondcomponents.ProvideNodeAPINodeHandler[NodeAPIContext],
		beacondcomponents.ProvideNodeAPIProofHandler[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, *CometBFTService, NodeAPIContext,
		],
	)

	return c
}
