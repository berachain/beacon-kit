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
	"github.com/berachain/beacon-kit/node-core/components"
)

//nolint:funlen // happens
func DefaultComponents() []any {
	c := []any{
		components.ProvideAttributesFactory[
			*BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, *Logger,
		],
		components.ProvideAvailibilityStore[*BeaconBlockBody, *Logger],
		components.ProvideDepositContract[
			*Deposit, *ExecutionPayload, *ExecutionPayloadHeader,
		],
		components.ProvideBlockStore[
			*BeaconBlock, *BeaconBlockBody, *Logger,
		],
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor[
			*AvailabilityStore, *BeaconBlockBody,
			*ConsensusSidecars, *BlobSidecar, *BlobSidecars, *Logger,
		],
		components.ProvideBlobProofVerifier,
		components.ProvideChainService[
			*AvailabilityStore,
			*ConsensusBlock, *BeaconBlock, *BeaconBlockBody,
			*BeaconState, *BeaconStateMarshallable,
			*BlobSidecar, *BlobSidecars, *ConsensusSidecars, *BlockStore,
			*Deposit,
			*DepositStore, *DepositContract,
			*ExecutionPayload, *ExecutionPayloadHeader, *Genesis,
			*KVStore, *Logger, *StorageBackend, *BlockStore,
		],
		components.ProvideNode,
		components.ProvideChainSpec,
		components.ProvideConfig,
		components.ProvideServerConfig,
		// components.ProvideConsensusEngine[
		// 	*AvailabilityStore, *BeaconBlockHeader, *BeaconState,
		// 	*BeaconStateMarshallable, *BlockStore, *KVStore, *StorageBackend,
		// ],
		components.ProvideDepositStore[*Deposit, *Logger],
		components.ProvideEngineClient[
			*ExecutionPayload, *ExecutionPayloadHeader, *Logger,
		],
		components.ProvideExecutionEngine[
			*ExecutionPayload, *ExecutionPayloadHeader, *Logger,
		],
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder[
			*BeaconState, *BeaconStateMarshallable,
			*ExecutionPayload, *ExecutionPayloadHeader, *KVStore, *Logger,
		],
		components.ProvideReportingService[
			*ExecutionPayload, *PayloadAttributes, *Logger,
		],
		components.ProvideCometBFTService[*Logger],
		components.ProvideServiceRegistry[
			*AvailabilityStore,
			*ConsensusBlock, *BeaconBlock, *BeaconBlockBody,
			*BlockStore, *BeaconState,
			*BeaconStateMarshallable,
			*ConsensusSidecars, *BlobSidecar, *BlobSidecars,
			*Deposit, *DepositStore, *ExecutionPayload, *ExecutionPayloadHeader,
			*Genesis, *KVStore, *Logger,
			NodeAPIContext,
		],
		components.ProvideSidecarFactory[
			*BeaconBlock, *BeaconBlockBody,
		],
		components.ProvideStateProcessor[
			*Logger, *BeaconBlock, *BeaconBlockBody,
			*BeaconState, *BeaconStateMarshallable, *Deposit, *DepositStore,
			*ExecutionPayload, *ExecutionPayloadHeader, *KVStore,
		],
		components.ProvideKVStore[*ExecutionPayloadHeader],
		components.ProvideStorageBackend[
			*AvailabilityStore, *BlockStore, *BeaconState,
			*KVStore, *DepositStore,
		],
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconState, *BeaconStateMarshallable,
			*BlockStore, *BlobSidecar, *BlobSidecars, *Deposit, *DepositStore,
			*ExecutionPayload, *ExecutionPayloadHeader, *KVStore, *Logger,
			*StorageBackend,
		],
		// TODO Hacks
		components.ProvideKVStoreService,
		components.ProvideKVStoreKey,
	}
	c = append(c,
		components.ProvideNodeAPIServer[*Logger, NodeAPIContext],
		components.ProvideNodeAPIEngine,
		components.ProvideNodeAPIBackend[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BlockStore, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecars, *Deposit, *DepositStore,
			*ExecutionPayloadHeader, *KVStore, *CometBFTService, *StorageBackend,
		],
	)

	c = append(c,
		components.ProvideNodeAPIHandlers[
			*BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, NodeAPIContext,
		],
		components.ProvideNodeAPIBeaconHandler[
			*BeaconState, *CometBFTService, NodeAPIContext,
		],
		components.ProvideNodeAPIBuilderHandler[NodeAPIContext],
		components.ProvideNodeAPIConfigHandler[NodeAPIContext],
		components.ProvideNodeAPIDebugHandler[NodeAPIContext],
		components.ProvideNodeAPIEventsHandler[NodeAPIContext],
		components.ProvideNodeAPINodeHandler[NodeAPIContext],
		components.ProvideNodeAPIProofHandler[
			*BeaconState, *BeaconStateMarshallable,
			*ExecutionPayloadHeader, *KVStore, *CometBFTService, NodeAPIContext,
		],
	)

	return c
}
