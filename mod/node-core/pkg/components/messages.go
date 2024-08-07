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
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
)

type MessageServerInput struct {
	depinject.In
	Routes []asynctypes.MessageRoute
}

// ProvideMessageServer provides a message server.
func ProvideMessageServer(in MessageServerInput) *server.MessageServer {
	ms := server.NewMessageServer()
	for _, route := range in.Routes {
		ms.RegisterRoute(route.MessageID(), route)
	}
	return ms
}

// RouteFactory creates a new route for the given message ID.
func RouteFactory(mID string) asynctypes.MessageRoute {
	switch mID {
	case messages.BuildBeaconBlockAndSidecars:
		return messaging.NewRoute[
			*SlotMessage, *BlockBundleMessage,
		](messages.BuildBeaconBlockAndSidecars)
	case messages.VerifyBeaconBlock:
		return messaging.NewRoute[
			*BlockMessage, *BlockMessage,
		](messages.VerifyBeaconBlock)
	case messages.FinalizeBeaconBlock:
		return messaging.NewRoute[
			*BlockMessage, *ValidatorUpdateMessage,
		](messages.FinalizeBeaconBlock)
	case messages.ProcessGenesisData:
		return messaging.NewRoute[
			*GenesisMessage, *ValidatorUpdateMessage,
		](messages.ProcessGenesisData)
	case messages.VerifySidecars:
		return messaging.NewRoute[
			*SidecarMessage, *SidecarMessage,
		](messages.VerifySidecars)
	case messages.ProcessSidecars:
		return messaging.NewRoute[
			*SidecarMessage, *SidecarMessage,
		](messages.ProcessSidecars)
	default:
		return nil
	}
}

// ProvideMessageRoutes provides all the message routes.
func ProvideMessageRoutes() []asynctypes.MessageRoute {
	return []asynctypes.MessageRoute{
		RouteFactory(messages.BuildBeaconBlockAndSidecars),
		RouteFactory(messages.VerifyBeaconBlock),
		RouteFactory(messages.FinalizeBeaconBlock),
		RouteFactory(messages.ProcessGenesisData),
		RouteFactory(messages.VerifySidecars),
		RouteFactory(messages.ProcessSidecars),
	}
}

// MessageServerComponents returns all the depinject providers for the message
// server.
func MessageServerComponents() []any {
	return []any{
		ProvideMessageServer,
		ProvideMessageRoutes,
	}
}
