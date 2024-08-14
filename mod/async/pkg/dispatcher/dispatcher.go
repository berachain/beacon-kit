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

// Dispatcher faciliates asynchronous communication between components,
// typically services. It acts as an API facade to the underlying event and
// message servers.

// Ensure Dispatcher implements DispatcherI.
var _ types.Dispatcher = (*Dispatcher)(nil)

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

// PublishEvent dispatches the given event to the event server.
// It will error if the <event> type is inconsistent with the publisher
// registered for the given eventID.
func (d *Dispatcher) PublishEvent(event types.BaseMessage) error {
	return d.eventServer.Publish(event)
}

// SendRequest dispatches the given request to the message server.
// It will error if the <req> and <resp> types are inconsistent with the
// route registered for the given messageID.
func (d *Dispatcher) SendRequest(req types.BaseMessage, future any) error {
	return d.msgServer.Request(req, future)
}

// SendResponse dispatches the given response to the message server.
// It will error if the <resp> type is inconsistent with the route registered
// for the given messageID.
func (d *Dispatcher) SendResponse(resp types.BaseMessage) error {
	return d.msgServer.Respond(resp)
}

// ============================== Events ===================================

// RegisterPublishers registers the given publisher with the given eventID.
// Any subsequent events with <eventID> dispatched to this Dispatcher must be
// consistent with the type expected by <publisher>.
func (d *Dispatcher) RegisterPublishers(
	publishers ...types.Publisher,
) error {
	var err error
	for _, publisher := range publishers {
		d.logger.Info("Publisher registered", "eventID", publisher.EventID())
		err = d.eventServer.RegisterPublisher(publisher.EventID(), publisher)
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe subscribes the given channel to the event with the given <eventID>.
// It will error if the channel type does not match the event type corresponding
// to the <eventID>.
// Contract: the channel must be a Subscription[T], where T is the expected
// type of the event data.
func (d *Dispatcher) Subscribe(eventID types.MessageID, ch any) error {
	return d.eventServer.Subscribe(eventID, ch)
}

// ================================ Messages ================================

// RegisterMsgRecipient registers the given channel to the message with the
// given <messageID>.
func (d *Dispatcher) RegisterMsgReceiver(
	messageID types.MessageID, ch any,
) error {
	d.logger.Info("Message receiver registered", "messageID", messageID)
	return d.msgServer.RegisterReceiver(messageID, ch)
}

// RegisterRoutes registers the given route with the given messageID.
// Any subsequent messages with <messageID> sent to this Dispatcher must be
// consistent with the type expected by <route>.
func (d *Dispatcher) RegisterRoutes(
	routes ...types.MessageRoute,
) error {
	var err error
	for _, route := range routes {
		d.logger.Info("Route registered", "messageID", route.MessageID())
		err = d.msgServer.RegisterRoute(route.MessageID(), route)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dispatcher) Name() string {
	return "dispatcher"
}
