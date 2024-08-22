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

import "github.com/berachain/beacon-kit/mod/node-core/pkg/components"

//nolint:funlen // happens
func DefaultComponents() []any {
	c := []any{
		components.ProvideABCIMiddleware[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BlobSidecar, *BlobSidecars, *Deposit, *Genesis, *Logger,
		],
		components.ProvideAttributesFactory[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*KVStore, *Logger,
		],
		components.ProvideAvailibilityStore[*BeaconBlockBody, *Logger],
		components.ProvideAvailabilityPruner[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlobSidecars, *Logger,
		],
		components.ProvideBeaconDepositContract[*Deposit],
		components.ProvideBlockStore[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader, *Logger,
		],
		components.ProvideBlockStoreService[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BlockStore, *Logger,
		],
		components.ProvideBlsSigner,
		components.ProvideBlobProcessor[
			*AvailabilityStore, *BeaconBlockBody, *BeaconBlockHeader,
			*BlobSidecar, *BlobSidecars, *Logger,
		],
		components.ProvideBlobProofVerifier,
		components.ProvideBlobVerifier[
			*BeaconBlockHeader, *BlobSidecar, *BlobSidecars,
		],
		components.ProvideChainService[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlobSidecars, *BlockStore, *Deposit, *DepositStore, *Genesis,
			*KVStore, *Logger, *StorageBackend,
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
			*AvailabilityStore, *BeaconBlockBody, *BlobSidecar,
			*BlobSidecars, *Logger,
		],
		components.ProvideDBManager[*AvailabilityStore, *DepositStore, *Logger],
		components.ProvideDepositPruner[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*Deposit, *DepositStore, *Logger,
		],
		components.ProvideDepositService[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader, *Deposit,
			*DepositContract, *DepositStore, *Logger,
		],
		components.ProvideDepositStore[*Deposit],
		components.ProvideDispatcher[
			*BeaconBlock, *BlobSidecars, *Genesis, *Logger,
		],
		components.ProvideEngineClient[*Logger],
		components.ProvideExecutionEngine[*Logger],
		components.ProvideJWTSecret,
		components.ProvideLocalBuilder[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*KVStore, *Logger,
		],
		components.ProvideReportingService[*Logger],
		components.ProvideServiceRegistry[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockStore, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecar, *BlobSidecars,
			*Deposit, *DepositStore, *Genesis, *KVStore, *Logger,
			NodeAPIContext,
		],
		components.ProvideSidecarFactory[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
		],
		components.ProvideStateProcessor[
			*BeaconBlock, *BeaconBlockBody, *BeaconBlockHeader,
			*BeaconState, *BeaconStateMarshallable, *Deposit, *KVStore,
		],
		components.ProvideKVStore[*BeaconBlockHeader],
		components.ProvideStorageBackend[
			*AvailabilityStore, *BlockStore, *BeaconState,
			*KVStore, *DepositStore,
		],
		components.ProvideTelemetrySink,
		components.ProvideTelemetryService,
		components.ProvideTrustedSetup,
		components.ProvideValidatorService[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlockStore, *BlobSidecars, *Deposit, *DepositStore,
			*KVStore, *Logger, *StorageBackend,
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
			*BeaconBlockHeader, *BlockStore, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecars, *Deposit, *DepositStore,
			*KVStore, Node, *StorageBackend,
		],
	)

	c = append(c,
		components.ProvideNodeAPIHandlers[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*KVStore, NodeAPIContext,
		],
		components.ProvideNodeAPIBeaconHandler[
			*BeaconBlockHeader, *BeaconState, Node, NodeAPIContext,
		],
		components.ProvideNodeAPIBuilderHandler[NodeAPIContext],
		components.ProvideNodeAPIConfigHandler[NodeAPIContext],
		components.ProvideNodeAPIDebugHandler[NodeAPIContext],
		components.ProvideNodeAPIEventsHandler[NodeAPIContext],
		components.ProvideNodeAPINodeHandler[NodeAPIContext],
		components.ProvideNodeAPIProofHandler[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*KVStore, Node, NodeAPIContext,
		],
	)

	return c
}
