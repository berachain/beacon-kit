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

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
)

// Dispatcher faciliates asynchronous communication between components, typically
// services. It acts as an API facade to the underlying event and message
// servers.
type Dispatcher struct {
	eventServer EventServer
	msgServer   MessageServer
	logger      log.Logger[any]
}

// NewDispatcher creates a new dispatcher.
func NewDispatcher(
	eventServer EventServer,
	msgServer MessageServer,
	logger log.Logger[any],
) *Dispatcher {
	eventServer.SetLogger(logger)
	msgServer.SetLogger(logger)
	return &Dispatcher{
		eventServer: eventServer,
		msgServer:   msgServer,
		logger:      logger,
	}
}

// Start starts the dispatcher.
func (d *Dispatcher) Start(ctx context.Context) error {
	d.eventServer.Start(ctx)
	return nil
}

// DispatchEvent dispatches the given event to the event server.
// It will error if the <event> type is inconsistent with the publisher
// registered for the given eventID.
func (d *Dispatcher) DispatchEvent(event types.MessageI) error {
	return d.eventServer.Publish(event)
}

// DispatchRequest dispatches the given request to the message server.
// It will error if the <req> and <resp> types are inconsistent with the
// route registered for the given messageID.
func (d *Dispatcher) DispatchRequest(req types.MessageI, resp any) error {
	return d.msgServer.Request(req, resp)
}

// DispatchResponse dispatches the given response to the message server.
// It will error if the <resp> type is inconsistent with the route registered
// for the given messageID.
func (d *Dispatcher) DispatchResponse(resp types.MessageI) error {
	return d.msgServer.Respond(resp)
}

// ============================== Events ===================================

// RegisterPublisher registers the given publisher with the given eventID.
// Any subsequent events with <eventID> dispatched to this Dispatcher must be
// consistent with the type expected by <publisher>.
func (d *Dispatcher) RegisterPublisher(eventID types.MessageID, publisher types.Publisher) {
	d.logger.Info("Publisher registered", "eventID", eventID)
	d.eventServer.RegisterPublisher(eventID, publisher)
}

// Subscribe subscribes the given channel to the event with the given <eventID>.
// It will error if the channel type does not match the event type corresponding
// to the <eventID>.
func (d *Dispatcher) Subscribe(eventID types.MessageID, ch any) error {
	return d.eventServer.Subscribe(eventID, ch)
}

// ================================ Messages ================================

// RegisterMsgRecipient registers the given channel to the message with the
// given <messageID>.
func (d *Dispatcher) RegisterMsgReceiver(messageID types.MessageID, ch any) error {
	d.logger.Info("Message receiver registered", "messageID", messageID)
	return d.msgServer.RegisterReceiver(messageID, ch)
}

// RegisterRoute registers the given route with the given messageID.
// Any subsequent messages with <messageID> sent to this Dispatcher must be
// consistent with the type expected by <route>.
func (d *Dispatcher) RegisterRoute(messageID types.MessageID, route types.MessageRoute) {
	d.logger.Info("Route registered", "messageID", messageID)
	d.msgServer.RegisterRoute(messageID, route)
}
