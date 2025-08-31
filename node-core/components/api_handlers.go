// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

type NodeAPIHandlersInput struct {
	depinject.In
	BeaconAPIHandler  *beaconapi.Handler
	BuilderAPIHandler *builderapi.Handler
	ConfigAPIHandler  *configapi.Handler
	DebugAPIHandler   *debugapi.Handler
	EventsAPIHandler  *eventsapi.Handler
	NodeAPIHandler    *nodeapi.Handler
	ProofAPIHandler   *proofapi.Handler
}

func ProvideNodeAPIHandlers(in NodeAPIHandlersInput) []handlers.Handlers {
	return []handlers.Handlers{
		in.BeaconAPIHandler,
		in.BuilderAPIHandler,
		in.ConfigAPIHandler,
		in.DebugAPIHandler,
		in.EventsAPIHandler,
		in.NodeAPIHandler,
		in.ProofAPIHandler,
	}
}

func ProvideNodeAPIBeaconHandler(b NodeAPIBackend) *beaconapi.Handler {
	return beaconapi.NewHandler(b)
}

func ProvideNodeAPIBuilderHandler() *builderapi.Handler {
	return builderapi.NewHandler()
}

func ProvideNodeAPIConfigHandler(b NodeAPIBackend) *configapi.Handler {
	return configapi.NewHandler(b)
}

func ProvideNodeAPIDebugHandler(b NodeAPIBackend) *debugapi.Handler {
	return debugapi.NewHandler(b)
}

func ProvideNodeAPIEventsHandler() *eventsapi.Handler {
	return eventsapi.NewHandler()
}

func ProvideNodeAPINodeHandler(b NodeAPINodeBackend) *nodeapi.Handler {
	return nodeapi.NewHandler(b)
}

func ProvideNodeAPIProofHandler(b NodeAPIBackend) *proofapi.Handler {
	return proofapi.NewHandler(b)
}
