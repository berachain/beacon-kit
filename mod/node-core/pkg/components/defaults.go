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
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/geth-primitives/pkg/bind"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	serviceprimitives "github.com/berachain/beacon-kit/mod/primitives/pkg/service"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
)

func DefaultComponents[
	AttestationDataT AttestationData[AttestationDataT],
	AttributesFactoryT AttributesFactory[BeaconStateT, PayloadAttributesT],
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[
		BeaconBlockT, AttestationDataT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, AttestationDataT, DepositT,
		Eth1DataT, ExecutionPayloadT, SlashingInfoT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobProcessorT BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT, BlobSidecarsT,
	],
	BlobSidecarT BlobSidecar[BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecars[BlobSidecarT, BlobSidecarsT],
	BlobFactoryT BlobFactory[BeaconBlockT, BlobSidecarsT],
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT Context[ContextT],
	DepositT Deposit[DepositT, ForkDataT, WithdrawalCredentialsT],
	DepositStoreT DepositStore[DepositT],
	Eth1DataT Eth1Data[Eth1DataT],
	EthClientT bind.ContractFilterer,
	EngineClientT EngineClient[
		ExecutionPayloadT,
		PayloadAttributesT,
		PayloadIDT,
	],
	ExecutionEngineT ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, PayloadAttributesT,
		PayloadIDT, WithdrawalT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT Fork[ForkT],
	ForkDataT ForkData[ForkDataT],
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	IndexDBT IndexDB,
	KVStoreT BeaconStore[
		KVStoreT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT, WithdrawalT,
	],
	KZGBlobProofVerifierT kzg.BlobProofVerifier,
	LegacyKeyT ~[32]byte,
	LocalBuilderT LocalBuilder[BeaconStateT, ExecutionPayloadT],
	LoggerT log.AdvancedLogger[any, LoggerT],
	MiddlewareT Middleware[SlotDataT],
	NodeT nodetypes.Node,
	PayloadAttributesT PayloadAttributes[PayloadAttributesT, WithdrawalT],
	PayloadIDT ~[8]byte,
	SlashingInfoT SlashingInfo[SlashingInfoT],
	SlotDataT SlotData[SlotDataT, AttestationDataT, SlashingInfoT],
	StateProcessorT StateProcessor[
		BeaconBlockT, BeaconStateT, ContextT,
		DepositT, ExecutionPayloadHeaderT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT Validators[ValidatorT],
	ValidatorUpdateT any,
	ValidatorUpdatesT any,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalsT Withdrawals[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
	// Services
	BlockBrokerT service.Basic,
	BlockStoreServiceT service.Basic,
	ChainServiceT service.Basic,
	DAServiceT service.Basic,
	DBManagerT service.Basic,
	DepositServiceT service.Basic,
	GenesisBrokerT service.Basic,
	NodeAPIServerT service.Basic,
	ReportingServiceT service.Basic,
	SidecarsBrokerT service.Basic,
	SlotBrokerT service.Basic,
	ValidatorServiceT service.Basic,
	ValidatorUpdateBrokerT service.Basic,
	// Events
	BeaconBlockEventT Event[BeaconBlockT],
	StatusEventT Event[*serviceprimitives.StatusEvent],
	// Pruners
	AvailabilityPrunerT pruner.Pruner[IndexDBT],
	BlockPrunerT pruner.Pruner[BlockStoreT],
	DepositPrunerT pruner.Pruner[DepositStoreT],
	// Node API Types
	BeaconAPIHandlerT handlers.Handlers[NodeAPIContextT],
	BuilderAPIHandlerT handlers.Handlers[NodeAPIContextT],
	ConfigAPIHandlerT handlers.Handlers[NodeAPIContextT],
	DebugAPIHandlerT handlers.Handlers[NodeAPIContextT],
	EventsAPIHandlerT handlers.Handlers[NodeAPIContextT],
	NodeAPIHandlerT handlers.Handlers[NodeAPIContextT],
	ProofAPIHandlerT handlers.Handlers[NodeAPIContextT],
	NodeAPIContextT NodeAPIContext,
	NodeAPIEngineT NodeAPIEngine[NodeAPIContextT],
	BeaconBackendT NodeAPIBeaconBackend[
		BeaconStateT, BeaconBlockHeaderT, ForkT, ValidatorT,
	],
	NodeAPIProofBackendT NodeAPIProofBackend[
		BeaconBlockHeaderT, BeaconStateT, ForkT, ValidatorT,
	],
]() []any {
	components := []any{
		ProvideABCIMiddleware[
			AttestationDataT, AvailabilityStoreT, BeaconBlockT,
			BeaconBlockBodyT, BeaconBlockHeaderT, BlobSidecarT,
			BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
			ExecutionPayloadHeaderT, GenesisT, LoggerT, SlashingInfoT,
			SlotDataT,
		],
		ProvideAttributesFactory[
			BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
			Eth1DataT, ExecutionPayloadHeaderT, ForkT, KVStoreT, LoggerT,
			PayloadAttributesT, ValidatorT, ValidatorsT, WithdrawalT,
		],
		ProvideAvailabilityPruner[
			AttestationDataT, AvailabilityStoreT, BeaconBlockT,
			BeaconBlockBodyT, BeaconBlockHeaderT, DepositT, Eth1DataT,
			ExecutionPayloadT, ExecutionPayloadHeaderT, IndexDBT, LoggerT,
			SlashingInfoT, WithdrawalsT,
		],
		ProvideAvailibilityStore[
			AttestationDataT, BeaconBlockBodyT, DepositT, Eth1DataT,
			ExecutionPayloadT, LoggerT, SlashingInfoT,
		],
		ProvideBeaconDepositContract[
			DepositT, EthClientT, ExecutionPayloadT, ForkDataT,
			WithdrawalCredentialsT,
		],
		ProvideBlockPruner[
			AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BlockStoreT, DepositT, Eth1DataT,
			ExecutionPayloadT, LoggerT, SlashingInfoT,
		],
		ProvideBlockStore[
			AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BlockStoreT, DepositT, Eth1DataT,
			ExecutionPayloadT, LoggerT, SlashingInfoT,
		],
		ProvideBlockStoreService[BeaconBlockT, BlockStoreT, LoggerT],
		ProvideBlsSigner[LegacyKeyT],
		ProvideBlobProcessor[
			AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
			BlobSidecarT, BlobSidecarsT, KZGBlobProofVerifierT, LoggerT,
		],
		ProvideBlobProofVerifier,
		ProvideBlobVerifier[
			BeaconBlockHeaderT, BlobSidecarT,
			BlobSidecarsT, KZGBlobProofVerifierT,
		],
		ProvideChainService[
			AttestationDataT, AvailabilityStoreT, BeaconBlockT,
			BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
			BeaconStateMarshallableT, BlobSidecarsT, BlobFactoryT,
			BlockStoreT, ContextT, DepositT, DepositStoreT, Eth1DataT,
			ExecutionEngineT, ExecutionPayloadT, ExecutionPayloadHeaderT,
			ForkT, GenesisT, KVStoreT, LocalBuilderT, LoggerT,
			PayloadAttributesT, PayloadIDT, SlashingInfoT, SlotDataT,
			StateProcessorT, StorageBackendT, ValidatorT, ValidatorsT,
			WithdrawalT, WithdrawalsT,
		],
		ProvideChainSpec,
		ProvideConfig,
		ProvideConsensusEngine[
			AttestationDataT, BeaconStateT, MiddlewareT, SlashingInfoT,
			SlotDataT, StorageBackendT, ValidatorUpdateT,
		],
		ProvideDAService[
			AvailabilityStoreT, BeaconBlockBodyT, BlobProcessorT,
			BlobSidecarT, BlobSidecarsT, ExecutionPayloadT, LoggerT,
		],
		ProvideDBManager[
			IndexDBT, AvailabilityPrunerT, BlockStoreT, BlockPrunerT,
			DepositStoreT, DepositPrunerT, LoggerT,
		],
		ProvideDepositPruner[
			BeaconBlockT, BeaconBlockBodyT, BeaconBlockEventT, DepositT,
			DepositStoreT, ExecutionPayloadT, ExecutionPayloadHeaderT,
			ForkDataT, WithdrawalsT, WithdrawalCredentialsT,
		],
		ProvideDepositService[
			AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BeaconBlockEventT, DepositT, DepositStoreT,
			Eth1DataT, ExecutionPayloadT, ExecutionPayloadHeaderT, ForkDataT,
			LoggerT, SlashingInfoT, WithdrawalsT, WithdrawalCredentialsT,
		],
		ProvideDepositStore[DepositT, ForkDataT, WithdrawalCredentialsT],
		ProvideEngineClient[
			ExecutionPayloadT, PayloadAttributesT, LoggerT, WithdrawalT,
		],
		ProvideExecutionEngine[
			EngineClientT, ExecutionPayloadT, ExecutionPayloadHeaderT,
			LoggerT, PayloadAttributesT, PayloadIDT, WithdrawalT,
			WithdrawalsT,
		],
		ProvideJWTSecret,
		ProvideLocalBuilder[
			AttributesFactoryT, BeaconBlockHeaderT, BeaconStateT,
			BeaconStateMarshallableT, Eth1DataT, ExecutionEngineT,
			ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, KVStoreT,
			LoggerT, PayloadAttributesT, PayloadIDT, ValidatorT,
			ValidatorsT, WithdrawalT, WithdrawalsT,
		],
		ProvideReportingService[LoggerT],
		ProvideServiceRegistry[
			MiddlewareT, BlockBrokerT, BlockStoreServiceT, ChainServiceT,
			DAServiceT, DBManagerT, DepositServiceT, EngineClientT,
			GenesisBrokerT, LoggerT, NodeAPIServerT, ReportingServiceT,
			SidecarsBrokerT, SlotBrokerT, ValidatorServiceT,
			ValidatorUpdateBrokerT,
		],
		ProvideSidecarFactory[
			AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, DepositT, Eth1DataT, ExecutionPayloadT,
			ExecutionPayloadHeaderT, SlashingInfoT, WithdrawalsT,
		],
		ProvideStateProcessor[
			AttestationDataT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
			ContextT, DepositT, Eth1DataT, ExecutionEngineT,
			ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, ForkDataT,
			KVStoreT, PayloadAttributesT, PayloadIDT, SlashingInfoT,
			ValidatorT, ValidatorsT, WithdrawalT, WithdrawalsT,
			WithdrawalCredentialsT,
		],
		ProvideKVStore[
			BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
			ValidatorT, ValidatorsT, WithdrawalCredentialsT,
		],
		ProvideStorageBackend[
			AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
			BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
			BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT, Eth1DataT,
			ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
			ValidatorsT, WithdrawalT, WithdrawalCredentialsT,
		],
		ProvideTelemetrySink,
		ProvideTrustedSetup,
		ProvideValidatorService[
			AttestationDataT, AvailabilityStoreT, BeaconBlockT,
			BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
			BeaconStateMarshallableT, BlobSidecarsT, BlobFactoryT,
			BlockStoreT, ContextT, DepositT, DepositStoreT, Eth1DataT,
			ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, ForkDataT,
			KVStoreT, LoggerT, LocalBuilderT, SlashingInfoT, SlotDataT,
			StateProcessorT, StorageBackendT, ValidatorT, ValidatorsT,
			WithdrawalT,
		],
		// TODO Hacks
		ProvideKVStoreService,
		ProvideKVStoreKey,
	}

	// components = append(components, DefaultNodeAPIComponents[
	// 	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	// 	BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
	// 	NodeAPIContextT, DepositT, DepositStoreT, NodeAPIEngineT, Eth1DataT,
	// 	ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, KVStoreT, LoggerT,
	// 	NodeT, StateProcessorT, StorageBackendT, ContextT, ValidatorT,
	// 	ValidatorsT, WithdrawalT, WithdrawalCredentialsT,
	// ]()...)

	// components = append(components, DefaultNodeAPIHandlers[
	// 	BeaconAPIHandlerT, BuilderAPIHandlerT, ConfigAPIHandlerT,
	// 	DebugAPIHandlerT, EventsAPIHandlerT, NodeAPIHandlerT, ProofAPIHandlerT,
	// 	NodeAPIContextT, BeaconBackendT, BeaconBlockHeaderT, BeaconStateT,
	// 	BeaconStateMarshallableT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
	// 	KVStoreT, NodeAPIProofBackendT, ValidatorT, ValidatorsT, WithdrawalT,
	// 	WithdrawalCredentialsT,
	// ]()...)

	components = append(components, DefaultBrokerProviders[
		BlobSidecarT, BeaconBlockT, GenesisT,
		SlotDataT, StatusEventT, ValidatorUpdatesT,
	]()...)
	return components
}
