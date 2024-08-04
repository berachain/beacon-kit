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
	"github.com/berachain/beacon-kit/mod/depinject"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	builderapi "github.com/berachain/beacon-kit/mod/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/mod/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/mod/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/mod/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/mod/node-api/handlers/node"
)

type NodeAPIHandlersInput struct {
	depinject.In

	BeaconAPIHandler  *BeaconAPIHandler
	BuilderAPIHandler *BuilderAPIHandler
	ConfigAPIHandler  *ConfigAPIHandler
	DebugAPIHandler   *DebugAPIHandler
	EventsAPIHandler  *EventsAPIHandler
	NodeAPIHandler    *NodeAPIHandler
}

func ProvideNodeAPIHandlers(
	in NodeAPIHandlersInput,
) []handlers.Handlers[NodeAPIContext] {
	return []handlers.Handlers[NodeAPIContext]{
		in.BeaconAPIHandler,
		in.BuilderAPIHandler,
		in.ConfigAPIHandler,
		in.DebugAPIHandler,
		in.EventsAPIHandler,
		in.NodeAPIHandler,
	}
}

// func ProvideNodeAPIBeaconHandler(b *NodeAPIBackend) *BeaconAPIHandler {
// 	return beaconapi.NewHandler[
// 		*BeaconBlockHeader,
// 		NodeAPIContext,
// 		*Fork,
// 		*Validator,
// 	](b)
// }

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

func ProvideNodeAPIHandler() *NodeAPIHandler {
	return nodeapi.NewHandler[NodeAPIContext]()
}

func ProvideNodeAPIProofHandler(b *NodeAPIBackend) *ProofAPIHandler {
	return proofapi.NewHandler[NodeAPIContext](b)
}

func DefaultNodeAPIHandlers() []any {
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
