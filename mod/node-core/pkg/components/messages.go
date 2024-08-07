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
	"github.com/berachain/beacon-kit/mod/async/pkg/messaging"
	"github.com/berachain/beacon-kit/mod/async/pkg/server"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
)

type MessageServerInput struct {
	depinject.In
	BuildBlockAndSidecarsRoute *BuildBlockAndSidecarsRoute
	VerifyBlockRoute           *VerifyBlockRoute
	FinalizeBlockRoute         *FinalizeBlockRoute
	ProcessGenesisDataRoute    *ProcessGenesisDataRoute
	ProcessBlobSidecarsRoute   *ProcessBlobSidecarsRoute
}

func ProvideMessageServer(in MessageServerInput) *MessageServer {
	ms := server.NewMessageServer()
	ms.RegisterRoute(in.BuildBlockAndSidecarsRoute.MessageID(), in.BuildBlockAndSidecarsRoute)
	ms.RegisterRoute(in.VerifyBlockRoute.MessageID(), in.VerifyBlockRoute)
	ms.RegisterRoute(in.FinalizeBlockRoute.MessageID(), in.FinalizeBlockRoute)
	ms.RegisterRoute(in.ProcessGenesisDataRoute.MessageID(), in.ProcessGenesisDataRoute)
	ms.RegisterRoute(in.ProcessBlobSidecarsRoute.MessageID(), in.ProcessBlobSidecarsRoute)
	return ms
}

// ProvideBuildBlockAndSidecarsRoute provides a route for building a beacon block.
func ProvideBuildBlockAndSidecarsRoute() *BuildBlockAndSidecarsRoute {
	return messaging.NewRoute[*SlotMessage, *BlockBundleMessage](messages.BuildBeaconBlockAndSidecars)
}

// ProvideVerifyBeaconBlockRoute provides a route for verifying a beacon block.
func ProvideVerifyBeaconBlockRoute() *VerifyBlockRoute {
	return messaging.NewRoute[*BlockMessage, *BlockMessage](messages.VerifyBeaconBlock)
}

// ProvideFinalizeBeaconBlockRoute provides a route for finalizing a beacon block.
func ProvideFinalizeBeaconBlockRoute() *FinalizeBlockRoute {
	return messaging.NewRoute[*BlockMessage, *ValidatorUpdateMessage](messages.FinalizeBeaconBlock)
}

// ProvideProcessGenesisDataRoute provides a route for processing genesis data.
func ProvideProcessGenesisDataRoute() *ProcessGenesisDataRoute {
	return messaging.NewRoute[*GenesisMessage, *ValidatorUpdateMessage](messages.ProcessGenesisData)
}

// ProvideProcessBlobSidecarsRoute provides a route for processing blob sidecars.
func ProvideProcessBlobSidecarsRoute() *ProcessBlobSidecarsRoute {
	return messaging.NewRoute[*SidecarMessage, *SidecarMessage](messages.VerifyBlobSidecars)
}

// MessageServerComponents returns all the depinject providers for the message
// server.
func MessageServerComponents() []any {
	return []any{
		ProvideMessageServer,
		ProvideBuildBlockAndSidecarsRoute,
		ProvideVerifyBeaconBlockRoute,
		ProvideFinalizeBeaconBlockRoute,
		ProvideProcessGenesisDataRoute,
		ProvideProcessBlobSidecarsRoute,
	}
}
