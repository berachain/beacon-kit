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
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/service"
)

func DefaultComponents() []any {
	components := []any{
		nodecomponents.ProvideABCIMiddleware[
			*AttestationData, *AvailabilityStore, *BeaconBlock,
			*BeaconBlockBody, *BeaconBlockHeader, *BlobSidecar,
			*BlobSidecars, *Deposit, *Eth1Data, *ExecutionPayload,
			*ExecutionPayloadHeader, *Genesis, *Logger, *SlashingInfo,
			*SlotData,
		],
		nodecomponents.ProvideAttributesFactory[
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*Eth1Data, *ExecutionPayloadHeader, *Fork, *KVStore, *Logger,
			*PayloadAttributes, *Validator, Validators, *Withdrawal,
		],
		nodecomponents.ProvideAvailabilityPruner[
			*AttestationData, *AvailabilityStore, *BeaconBlock,
			*BeaconBlockBody, *BeaconBlockHeader, *Deposit, *Eth1Data,
			*ExecutionPayload, *ExecutionPayloadHeader, *IndexDB, *Logger,
			*SlashingInfo, Withdrawals,
		],
		nodecomponents.ProvideAvailibilityStore[
			*AttestationData, *BeaconBlockBody, *Deposit, *Eth1Data,
			*ExecutionPayload, *Logger, *SlashingInfo,
		],
		nodecomponents.ProvideBeaconDepositContract[
			*Deposit, *EngineClient, *ExecutionPayload, *ForkData,
			*PayloadAttributes, PayloadID, WithdrawalCredentials,
		],
		nodecomponents.ProvideBlockPruner[
			*AttestationData, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockStore, *Deposit, *Eth1Data,
			*ExecutionPayload, *Logger, *SlashingInfo,
		],
		nodecomponents.ProvideBlockStore[
			*AttestationData, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockStore, *Deposit, *Eth1Data,
			*ExecutionPayload, *Logger, *SlashingInfo,
		],
		nodecomponents.ProvideBlockStoreService[*BeaconBlock, *BlockStore, *Logger],
		nodecomponents.ProvideBlsSigner[LegacyKey],
		nodecomponents.ProvideBlobProcessor[
			*AvailabilityStore, *BeaconBlockBody, *BeaconBlockHeader,
			*BlobSidecar, *BlobSidecars, KZGBlobProofVerifier, *Logger,
		],
		nodecomponents.ProvideBlobProofVerifier,
		nodecomponents.ProvideBlobVerifier[
			*BeaconBlockHeader, *BlobSidecar,
			*BlobSidecars, KZGBlobProofVerifier,
		],
		nodecomponents.ProvideChainService[
			*AttestationData, *AvailabilityStore, *BeaconBlock,
			*BeaconBlockBody, *BeaconBlockHeader, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecars, *BlobFactory,
			*BlockStore, *Context, *Deposit, *DepositStore, *Eth1Data,
			*ExecutionEngine, *ExecutionPayload, *ExecutionPayloadHeader,
			*Fork, *Genesis, *KVStore, *LocalBuilder, *Logger,
			*PayloadAttributes, PayloadID, *SlashingInfo, *SlotData,
			*StateProcessor, *StorageBackend, *Validator, Validators,
			*Withdrawal, Withdrawals,
		],
		nodecomponents.ProvideChainSpec,
		nodecomponents.ProvideConfig,
		nodecomponents.ProvideConsensusEngine[
			*AttestationData, *BeaconState, *ABCIMiddleware, *SlashingInfo,
			*SlotData, *StorageBackend, *ValidatorUpdate,
		],
		nodecomponents.ProvideDAService[
			*AvailabilityStore, *BeaconBlockBody, *BlobProcessor,
			*BlobSidecar, *BlobSidecars, *ExecutionPayload, *Logger,
		],
		nodecomponents.ProvideDBManager[
			*IndexDB, DAPruner, *BlockStore, BlockPruner,
			*DepositStore, DepositPruner, *Logger,
		],
		nodecomponents.ProvideDepositPruner[
			*BeaconBlock, *BeaconBlockBody, *BlockEvent, *Deposit,
			*DepositStore, *ExecutionPayload, *ExecutionPayloadHeader,
			*ForkData, *Logger, Withdrawals, WithdrawalCredentials,
		],
		nodecomponents.ProvideDepositService[
			*AttestationData, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BlockEvent, *Deposit, *DepositStore,
			*Eth1Data, *ExecutionPayload, *ExecutionPayloadHeader, *ForkData,
			*Logger, *SlashingInfo, Withdrawals, WithdrawalCredentials,
		],
		nodecomponents.ProvideDepositStore[*Deposit, *ForkData, WithdrawalCredentials],
		nodecomponents.ProvideEngineClient[
			*ExecutionPayload, *PayloadAttributes, *Logger, *Withdrawal,
		],
		nodecomponents.ProvideExecutionEngine[
			*EngineClient, *ExecutionPayload, *ExecutionPayloadHeader,
			*Logger, *PayloadAttributes, PayloadID, *Withdrawal,
			Withdrawals,
		],
		nodecomponents.ProvideJWTSecret,
		nodecomponents.ProvideLocalBuilder[
			*AttributesFactory, *BeaconBlockHeader, *BeaconState,
			*BeaconStateMarshallable, *Eth1Data, *ExecutionEngine,
			*ExecutionPayload, *ExecutionPayloadHeader, *Fork, *KVStore,
			*Logger, *PayloadAttributes, PayloadID, *Validator,
			Validators, *Withdrawal, Withdrawals,
		],
		nodecomponents.ProvideReportingService[*Logger],
		nodecomponents.ProvideServiceRegistry[
			*ABCIMiddleware, *BlockBroker, *BlockStoreService, *ChainService,
			*DAService, *DBManager, *DepositService, *EngineClient,
			*GenesisBroker, *Logger, *NodeAPIServer, *ReportingService,
			*SidecarsBroker, *SlotBroker, *ValidatorService,
			*ValidatorUpdateBroker,
		],
		nodecomponents.ProvideSidecarFactory[
			*AttestationData, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *Deposit, *Eth1Data, *ExecutionPayload,
			*ExecutionPayloadHeader, *SlashingInfo, Withdrawals,
		],
		nodecomponents.ProvideStateProcessor[
			*AttestationData, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*Context, *Deposit, *Eth1Data, *ExecutionEngine,
			*ExecutionPayload, *ExecutionPayloadHeader, *Fork, *ForkData,
			*KVStore, *PayloadAttributes, PayloadID, *SlashingInfo,
			*Validator, Validators, *Withdrawal, Withdrawals,
			WithdrawalCredentials,
		],
		nodecomponents.ProvideStore[
			*BeaconBlockHeader, *Eth1Data, *ExecutionPayloadHeader, *Fork,
			*Validator, Validators, WithdrawalCredentials,
		],
		nodecomponents.ProvideStorageBackend[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlobSidecars, *BlockStore, *Deposit, *DepositStore, *Eth1Data,
			*ExecutionPayloadHeader, *Fork, *ForkData, *KVStore, *Validator,
			Validators, *Withdrawal, WithdrawalCredentials,
		],
		nodecomponents.ProvideTelemetrySink,
		nodecomponents.ProvideTrustedSetup,
		nodecomponents.ProvideValidatorService[
			*AttestationData, *AvailabilityStore, *BeaconBlock,
			*BeaconBlockBody, *BeaconBlockHeader, *BeaconState,
			*BeaconStateMarshallable, *BlobSidecars, *BlobFactory,
			*BlockStore, *Context, *Deposit, *DepositStore, *Eth1Data,
			*ExecutionPayload, *ExecutionPayloadHeader, *Fork, *ForkData,
			*KVStore, *Logger, *LocalBuilder, *SlashingInfo, *SlotData,
			*StateProcessor, *StorageBackend, *Validator, Validators,
			*Withdrawal,
		],
		// TODO Hacks
		nodecomponents.ProvideKVStoreService,
		nodecomponents.ProvideKVStoreKey,
	}

	components = append(components,
		nodecomponents.ProvideNodeAPIHandlers[
			*BeaconAPIHandler, *BuilderAPIHandler, *ConfigAPIHandler,
			NodeAPIContext, *DebugAPIHandler, *EventsAPIHandler,
			*NodeAPIHandler, *ProofAPIHandler,
		],
		nodecomponents.ProvideNodeAPIBeaconHandler[
			*NodeAPIBackend, *BeaconState, *BeaconBlockHeader,
			NodeAPIContext, *Fork, *Validator,
		],
		nodecomponents.ProvideNodeAPIBuilderHandler[NodeAPIContext],
		nodecomponents.ProvideNodeAPIConfigHandler[NodeAPIContext],
		nodecomponents.ProvideNodeAPIDebugHandler[NodeAPIContext],
		nodecomponents.ProvideNodeAPIEventsHandler[NodeAPIContext],
		nodecomponents.ProvideNodeAPINodeHandler[NodeAPIContext],
		nodecomponents.ProvideNodeAPIProofHandler[
			NodeAPIContext, *BeaconBlockHeader, *BeaconState,
			*BeaconStateMarshallable, *Eth1Data,
			*ExecutionPayloadHeader, *Fork, *KVStore,
			*NodeAPIBackend, *Validator, Validators,
			*Withdrawal, WithdrawalCredentials,
		],
	)

	components = append(components,
		nodecomponents.ProvideNodeAPIServer[
			NodeAPIContext, *NodeAPIEngine, *Logger,
		],
		nodecomponents.ProvideNodeAPIEngine,
		nodecomponents.ProvideNodeAPIBackend[
			*AvailabilityStore, *BeaconBlock, *BeaconBlockBody,
			*BeaconBlockHeader, *BeaconState, *BeaconStateMarshallable,
			*BlobSidecars, *BlockStore, *Deposit, *DepositStore,
			*Eth1Data, *ExecutionPayload, *ExecutionPayloadHeader, *Fork,
			*KVStore, Node, *StateProcessor, *StorageBackend,
			*Context, *Validator, Validators, *Withdrawal,
			WithdrawalCredentials,
		],
	)

	components = append(components,
		nodecomponents.ProvideBlobBroker[*BlobSidecars],
		nodecomponents.ProvideBlockBroker[*BeaconBlock],
		nodecomponents.ProvideGenesisBroker[*Genesis],
		nodecomponents.ProvideSlotBroker[*SlotData],
		nodecomponents.ProvideStatusBroker[*service.StatusEvent],
		nodecomponents.ProvideValidatorUpdateBroker[ValidatorUpdates],
	)
	return components
}
