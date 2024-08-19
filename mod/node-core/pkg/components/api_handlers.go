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
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, *Eth1Data,
		*ExecutionPayloadHeader, *Fork, *Validator,
	],
	KVStoreT any,
] struct {
	depinject.In
	BeaconAPIHandler *beaconapi.Handler[
		BeaconBlockHeaderT, NodeAPIContext, *Fork, *Validator,
	]
	BuilderAPIHandler *BuilderAPIHandler
	ConfigAPIHandler  *ConfigAPIHandler
	DebugAPIHandler   *DebugAPIHandler
	EventsAPIHandler  *EventsAPIHandler
	NodeAPIHandler    *NodeAPIHandler
	ProofAPIHandler   *proofapi.Handler[
		BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
		NodeAPIContext, *ExecutionPayloadHeader, *Validator,
	]
}

func ProvideNodeAPIHandlers[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, *Eth1Data,
		*ExecutionPayloadHeader, *Fork, *Validator,
	],
	KVStoreT any,
](
	in NodeAPIHandlersInput[
		BeaconBlockHeaderT, BeaconStateT,
		BeaconStateMarshallableT, KVStoreT,
	],
) []handlers.Handlers[NodeAPIContext] {
	return []handlers.Handlers[NodeAPIContext]{
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
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT any,
	NodeT any,
](b NodeAPIBackend[
	BeaconBlockHeaderT,
	BeaconStateT,
	*Fork,
	NodeT,
	*Validator,
]) *beaconapi.Handler[
	BeaconBlockHeaderT, NodeAPIContext, *Fork, *Validator,
] {
	return beaconapi.NewHandler[
		BeaconBlockHeaderT,
		NodeAPIContext,
		*Fork,
		*Validator,
	](b)
}

func ProvideNodeAPIBuilderHandler() *BuilderAPIHandler {
	return builderapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPIConfigHandler() *ConfigAPIHandler {
	return configapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPIDebugHandler() *DebugAPIHandler {
	return debugapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPIEventsHandler() *EventsAPIHandler {
	return eventsapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPINodeHandler() *NodeAPIHandler {
	return nodeapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPIProofHandler[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, *Eth1Data,
		*ExecutionPayloadHeader, *Fork, *Validator,
	],
	KVStoreT any,
	NodeT any,
](b NodeAPIBackend[
	BeaconBlockHeaderT,
	BeaconStateT,
	*Fork,
	NodeT,
	*Validator,
]) *proofapi.Handler[
	BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
	NodeAPIContext, *ExecutionPayloadHeader, *Validator,
] {
	return proofapi.NewHandler[
		BeaconBlockHeaderT,
		BeaconStateT,
		BeaconStateMarshallableT,
		NodeAPIContext,
		*ExecutionPayloadHeader,
		*Validator,
	](b)
}

// func DefaultNodeAPIHandlers[
// 	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
// ]() []any {
// 	return []any{
// 		ProvideNodeAPIHandlers,
// 		ProvideNodeAPIBeaconHandler[BeaconBlockHeaderT],
// 		ProvideNodeAPIBuilderHandler,
// 		ProvideNodeAPIConfigHandler,
// 		ProvideNodeAPIDebugHandler,
// 		ProvideNodeAPIEventsHandler,
// 		ProvideNodeAPINodeHandler,
// 		ProvideNodeAPIProofHandler,
// 	}
// }
