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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	beaconapi "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon"
	builderapi "github.com/berachain/beacon-kit/mod/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/mod/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/mod/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/mod/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/mod/node-api/handlers/node"
	proofapi "github.com/berachain/beacon-kit/mod/node-api/handlers/proof"
)

type NodeAPIHandlersInput[
	BeaconAPIHandlerT handlers.Handlers[ContextT],
	BuilderAPIHandlerT handlers.Handlers[ContextT],
	ConfigAPIHandlerT handlers.Handlers[ContextT],
	ContextT NodeAPIContext,
	DebugAPIHandlerT handlers.Handlers[ContextT],
	EventsAPIHandlerT handlers.Handlers[ContextT],
	NodeAPIHandlerT handlers.Handlers[ContextT],
	ProofAPIHandlerT handlers.Handlers[ContextT],
] struct {
	depinject.In

	BeaconAPIHandler  BeaconAPIHandlerT
	BuilderAPIHandler BuilderAPIHandlerT
	ConfigAPIHandler  ConfigAPIHandlerT
	DebugAPIHandler   DebugAPIHandlerT
	EventsAPIHandler  EventsAPIHandlerT
	NodeAPIHandler    NodeAPIHandlerT
	ProofAPIHandler   ProofAPIHandlerT
}

func ProvideNodeAPIHandlers[
	BeaconAPIHandlerT handlers.Handlers[ContextT],
	BuilderAPIHandlerT handlers.Handlers[ContextT],
	ConfigAPIHandlerT handlers.Handlers[ContextT],
	ContextT NodeAPIContext,
	DebugAPIHandlerT handlers.Handlers[ContextT],
	EventsAPIHandlerT handlers.Handlers[ContextT],
	NodeAPIHandlerT handlers.Handlers[ContextT],
	ProofAPIHandlerT handlers.Handlers[ContextT],
](
	in NodeAPIHandlersInput[
		BeaconAPIHandlerT, BuilderAPIHandlerT, ConfigAPIHandlerT,
		ContextT, DebugAPIHandlerT, EventsAPIHandlerT, NodeAPIHandlerT,
		ProofAPIHandlerT,
	],
) []handlers.Handlers[ContextT] {
	return []handlers.Handlers[ContextT]{
		in.BeaconAPIHandler,
		in.BuilderAPIHandler,
		in.ConfigAPIHandler,
		in.DebugAPIHandler,
		in.EventsAPIHandler,
		in.NodeAPIHandler,
		in.ProofAPIHandler,
	}
}

func ProvideNodeAPIBeaconHandler[
	BeaconBackendT NodeAPIBeaconBackend[
		BeaconStateT, BeaconBlockHeaderT, ForkT, ValidatorT,
	],
	BeaconStateT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ContextT NodeAPIContext,
	ForkT any,
	ValidatorT any,
](b BeaconBackendT) *beaconapi.Handler[
	BeaconBlockHeaderT,
	ContextT,
	ForkT,
	ValidatorT,
] {
	return beaconapi.NewHandler[
		BeaconBackendT,
		BeaconBlockHeaderT,
		ContextT,
		ForkT,
		ValidatorT,
	](b)
}

func ProvideNodeAPIBuilderHandler[
	ContextT NodeAPIContext,
]() *builderapi.Handler[ContextT] {
	return builderapi.NewHandler[ContextT]()
}

func ProvideNodeAPIConfigHandler[
	ContextT NodeAPIContext,
]() *configapi.Handler[ContextT] {
	return configapi.NewHandler[ContextT]()
}

func ProvideNodeAPIDebugHandler[
	ContextT NodeAPIContext,
]() *debugapi.Handler[ContextT] {
	return debugapi.NewHandler[ContextT]()
}

func ProvideNodeAPIEventsHandler[
	ContextT NodeAPIContext,
]() *eventsapi.Handler[ContextT] {
	return eventsapi.NewHandler[ContextT]()
}

func ProvideNodeAPINodeHandler[
	ContextT NodeAPIContext,
]() *nodeapi.Handler[ContextT] {
	return nodeapi.NewHandler[ContextT]()
}

func ProvideNodeAPIProofHandler[
	ContextT NodeAPIContext,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	Eth1DataT any,
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	ForkT any,
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT Validators[ValidatorT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
]() *proofapi.Handler[
	ContextT,
	BeaconBlockHeaderT,
	BeaconStateT,
	BeaconStateMarshallableT,
	ExecutionPayloadHeaderT,
	ValidatorT,
] {
	return proofapi.NewHandler[
		ContextT,
		BeaconBlockHeaderT,
		BeaconStateT,
		BeaconStateMarshallableT,
		ExecutionPayloadHeaderT,
		ValidatorT,
	]()
}

func DefaultNodeAPIHandlers[ContextT NodeAPIContext]() []any {
	return []any{
		ProvideNodeAPIHandlers,
		ProvideNodeAPIBeaconHandler,
		ProvideNodeAPIBuilderHandler,
		ProvideNodeAPIConfigHandler,
		ProvideNodeAPIDebugHandler,
		ProvideNodeAPIEventsHandler,
		ProvideNodeAPINodeHandler,
		ProvideNodeAPIProofHandler,
	}
}
