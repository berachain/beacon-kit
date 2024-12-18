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
	"github.com/berachain/beacon-kit/node-api/handlers"
	beaconapi "github.com/berachain/beacon-kit/node-api/handlers/beacon"
	builderapi "github.com/berachain/beacon-kit/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/node-api/handlers/node"
	proofapi "github.com/berachain/beacon-kit/node-api/handlers/proof"
)

type NodeAPIHandlersInput[
	KVStoreT any,
	NodeAPIContextT NodeAPIContext,
] struct {
	depinject.In
	BeaconAPIHandler  *beaconapi.Handler[NodeAPIContextT]
	BuilderAPIHandler *builderapi.Handler[NodeAPIContextT]
	ConfigAPIHandler  *configapi.Handler[NodeAPIContextT]
	DebugAPIHandler   *debugapi.Handler[NodeAPIContextT]
	EventsAPIHandler  *eventsapi.Handler[NodeAPIContextT]
	NodeAPIHandler    *nodeapi.Handler[NodeAPIContextT]
	ProofAPIHandler   *proofapi.Handler[NodeAPIContextT]
}

func ProvideNodeAPIHandlers[
	KVStoreT any,
	NodeAPIContextT NodeAPIContext,
](
	in NodeAPIHandlersInput[
		KVStoreT,
		NodeAPIContextT,
	],
) []handlers.Handlers[NodeAPIContextT] {
	return []handlers.Handlers[NodeAPIContextT]{
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
	NodeT any,
	NodeAPIContextT NodeAPIContext,
](b NodeAPIBackend[NodeT]) *beaconapi.Handler[NodeAPIContextT] {
	return beaconapi.NewHandler[NodeAPIContextT](b)
}

func ProvideNodeAPIBuilderHandler[
	NodeAPIContextT NodeAPIContext,
]() *builderapi.Handler[NodeAPIContextT] {
	return builderapi.NewHandler[NodeAPIContextT]()
}

func ProvideNodeAPIConfigHandler[
	NodeAPIContextT NodeAPIContext,
]() *configapi.Handler[NodeAPIContextT] {
	return configapi.NewHandler[NodeAPIContextT]()
}

func ProvideNodeAPIDebugHandler[
	NodeAPIContextT NodeAPIContext,
]() *debugapi.Handler[NodeAPIContextT] {
	return debugapi.NewHandler[NodeAPIContextT]()
}

func ProvideNodeAPIEventsHandler[
	NodeAPIContextT NodeAPIContext,
]() *eventsapi.Handler[NodeAPIContextT] {
	return eventsapi.NewHandler[NodeAPIContextT]()
}

func ProvideNodeAPINodeHandler[
	NodeAPIContextT NodeAPIContext,
]() *nodeapi.Handler[NodeAPIContextT] {
	return nodeapi.NewHandler[NodeAPIContextT]()
}

func ProvideNodeAPIProofHandler[
	KVStoreT any,
	NodeT any,
	NodeAPIContextT NodeAPIContext,
](b NodeAPIBackend[NodeT]) *proofapi.Handler[NodeAPIContextT] {
	return proofapi.NewHandler[NodeAPIContextT](b)
}
