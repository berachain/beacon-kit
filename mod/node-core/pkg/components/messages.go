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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

type MessageServerInput struct {
	depinject.In
	BuildBlockRoute          *messaging.Route[*SlotMessage, *BlockMessage]
	BuildSidecarsRoute       *messaging.Route[*SlotMessage, *SidecarMessage]
	VerifyBlockRoute         *messaging.Route[*BlockMessage, *BlockMessage]
	FinalizeBlockRoute       *messaging.Route[*BlockMessage, *ValidatorUpdateMessage]
	ProcessGenesisDataRoute  *messaging.Route[*GenesisMessage, *ValidatorUpdateMessage]
	ProcessBlobSidecarsRoute *messaging.Route[*SidecarMessage, *SidecarMessage]
}

func ProvideMessageServer(in MessageServerInput) *server.MessageServer {
	ms := server.NewMessageServer()
	ms.RegisterRoute(in.BuildBlockRoute.MessageID(), in.BuildBlockRoute)
	ms.RegisterRoute(in.BuildSidecarsRoute.MessageID(), in.BuildSidecarsRoute)
	ms.RegisterRoute(in.VerifyBlockRoute.MessageID(), in.VerifyBlockRoute)
	ms.RegisterRoute(in.FinalizeBlockRoute.MessageID(), in.FinalizeBlockRoute)
	ms.RegisterRoute(in.ProcessGenesisDataRoute.MessageID(), in.ProcessGenesisDataRoute)
	ms.RegisterRoute(in.ProcessBlobSidecarsRoute.MessageID(), in.ProcessBlobSidecarsRoute)
	return ms
}

// ProvideBuildBlockRoute provides a route for building a beacon block.
func ProvideBuildBlockRoute() *messaging.Route[*SlotMessage, *BlockMessage] {
	return messaging.NewRoute[*SlotMessage, *BlockMessage](events.BuildBeaconBlock)
}

// ProvideBuildSidecarsRoute provides a route for building sidecars.
func ProvideBuildSidecarsRoute() *messaging.Route[*SlotMessage, *SidecarMessage] {
	return messaging.NewRoute[*SlotMessage, *SidecarMessage](events.BuildBlobSidecars)
}

// ProvideVerifyBeaconBlockRoute provides a route for verifying a beacon block.
func ProvideVerifyBeaconBlockRoute() *messaging.Route[*BlockMessage, *BlockMessage] {
	return messaging.NewRoute[*BlockMessage, *BlockMessage](events.VerifyBeaconBlock)
}

// ProvideFinalizeBeaconBlockRoute provides a route for finalizing a beacon block.
func ProvideFinalizeBeaconBlockRoute() *messaging.Route[*BlockMessage, *ValidatorUpdateMessage] {
	return messaging.NewRoute[*BlockMessage, *ValidatorUpdateMessage](events.FinalizeBeaconBlock)
}

// ProvideProcessGenesisDataRoute provides a route for processing genesis data.
func ProvideProcessGenesisDataRoute() *messaging.Route[*GenesisMessage, *ValidatorUpdateMessage] {
	return messaging.NewRoute[*GenesisMessage, *ValidatorUpdateMessage](events.ProcessGenesisData)
}

// ProvideProcessBlobSidecarsRoute provides a route for processing blob sidecars.
func ProvideProcessBlobSidecarsRoute() *messaging.Route[*SidecarMessage, *SidecarMessage] {
	return messaging.NewRoute[*SidecarMessage, *SidecarMessage](events.VerifyBlobSidecars)
}
