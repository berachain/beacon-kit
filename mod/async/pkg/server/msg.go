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

package server

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
)

// MessageServer is a server for sending and receiving messages.
type MessageServer struct {
	routes map[types.MessageID]types.MessageRoute
	logger log.Logger[any]
}

// NewMessageServer creates a new message server.
func NewMessageServer() *MessageServer {
	return &MessageServer{
		routes: make(map[types.MessageID]types.MessageRoute),
	}
}

// Request sends a message to the server and awaits for a response.
// The response is written to the provided response pointer.
func (ms *MessageServer) Request(req types.MessageI, respPtr any) error {
	// send request and await response
	route, ok := ms.routes[req.ID()]
	if !ok {
		return ErrRouteNotFound
	}
	if err := route.SendRequest(req); err != nil {
		return err
	}
	return route.AwaitResponse(respPtr)
}

// Respond sends a response to the route that corresponds to the response's
// messageID.
func (ms *MessageServer) Respond(resp types.MessageI) error {
	route, ok := ms.routes[resp.ID()]
	if !ok {
		return ErrRouteNotFound
	}
	return route.SendResponse(resp)
}

// RegisterRoute registers the route with the given messageID.
// Any subsequent messages with <messageID> sent to this MessageServer must be
// consistent with the type expected by <route>.
func (ms *MessageServer) RegisterRoute(
	messageID types.MessageID, route types.MessageRoute,
) error {
	if ms.routes[messageID] != nil {
		return ErrRouteAlreadyRegistered(messageID)
	}
	ms.routes[messageID] = route
	return nil
}

// SetRecipient sets the recipient for the route with the given messageID.
// Errors if the route with the given messageID is not found or the route
// already has a registered recipient.
func (ms *MessageServer) RegisterReceiver(
	messageID types.MessageID, ch any,
) error {
	route, ok := ms.routes[messageID]
	if !ok {
		return ErrRouteNotFound
	}
	return route.RegisterReceiver(ch)
}

// SetLogger sets the logger for the message server.
func (ms *MessageServer) SetLogger(logger log.Logger[any]) {
	ms.logger = logger
}
