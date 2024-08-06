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

package dispatcher

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// MsgHandler is a function that handles a message.
type MsgHandler func(req types.Message[any]) (types.Message[any], error)

// Message server implementation with handlers
// MsgServer2 is a server for handling messages.
// don't use this, just an idea
type MsgServer2 struct {
	// routes is a map of message types to the corresponding receiving channels.
	routes map[types.MessageID]MsgHandler
}

// NewMsgServer2 creates a new message server.
func NewMsgServer2() *MsgServer2 {
	return &MsgServer2{
		routes: make(map[types.MessageID]MsgHandler),
	}
}

// RegisterHandler registers a handler for a message type.
func (s *MsgServer2) RegisterHandler(messageType types.MessageID, handler MsgHandler) {
	s.routes[messageType] = handler
}

// DispatchAndAwait dispatches a message to the server and returns the response.
func (s *MsgServer2) DispatchAndAwait(ctx context.Context, msg types.Message[any]) (types.Message[any], error) {
	handler, ok := s.routes[msg.ID()]
	if !ok {
		return types.Message[any]{}, fmt.Errorf("no handler found for message type %s", msg.ID())
	}
	return handler(msg)
}

// Dispatch dispatches a message to the server and handles the response.
func (s *MsgServer2) Dispatch(
	ctx context.Context,
	msg types.Message[any],
	handleResp types.ResponseHandler[any],

) error {
	handler, ok := s.routes[msg.ID()]
	if !ok {
		return fmt.Errorf("no handler found for message type %s", msg.ID())
	}
	resp, err := handler(msg)
	if err != nil {
		return err
	}
	handleResp(resp)
	return nil
}
