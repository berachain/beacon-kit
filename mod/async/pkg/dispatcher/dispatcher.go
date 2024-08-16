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

// Ensure Dispatcher implements DispatcherI.
var _ types.Dispatcher = (*Dispatcher)(nil)

// Dispatcher faciliates asynchronous communication between components,
// typically services. It acts as an API facade to the underlying eventserver.
type Dispatcher struct {
	eventServer EventServer
	logger      log.Logger[any]
}

// NewDispatcher creates a new dispatcher.
func NewDispatcher(
	eventServer EventServer,
	logger log.Logger[any],
) *Dispatcher {
	eventServer.SetLogger(logger)
	return &Dispatcher{
		eventServer: eventServer,
		logger:      logger,
	}
}

// Start starts the underlying event server.
func (d *Dispatcher) Start(ctx context.Context) error {
	d.eventServer.Start(ctx)
	return nil
}

// Subscribe subscribes the given channel to the event with the given <eventID>.
// It will error if the channel type does not match the event type corresponding
// to the <eventID>.
// Contract: the channel must be a Subscription[T], where T is the expected
// type of the event data.
func (d *Dispatcher) Subscribe(eventID types.EventID, ch any) error {
	return d.eventServer.Subscribe(eventID, ch)
}

// PublishEvent dispatches the given event to the event server.
// It will error if the <event> type is inconsistent with the publisher
// registered for the given eventID.
func (d *Dispatcher) PublishEvent(event types.BaseEvent) error {
	return d.eventServer.Publish(event)
}

// RegisterPublishers registers the given publisher with the given eventID.
// Any subsequent events with <eventID> dispatched to this Dispatcher must be
// consistent with the type expected by <publisher>.
func (d *Dispatcher) RegisterPublishers(
	publishers ...types.Publisher,
) error {
	return d.eventServer.RegisterPublishers(publishers...)
}

// Name returns the name of this service.
func (d *Dispatcher) Name() string {
	return "dispatcher"
}
