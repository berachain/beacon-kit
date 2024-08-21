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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

var _ types.Dispatcher = (*Dispatcher)(nil)

// Dispatcher faciliates asynchronous communication between components,
// typically services.
type Dispatcher struct {
	brokers map[async.EventID]types.Broker
	logger  log.Logger[any]
}

// NewDispatcher creates a new event server.
func New(
	logger log.Logger[any],
	options ...Option,
) (*Dispatcher, error) {
	d := &Dispatcher{
		brokers: make(map[async.EventID]types.Broker),
		logger:  logger,
	}
	for _, option := range options {
		if err := option(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// Publish dispatches the given event to the broker with the given eventID.
func (d *Dispatcher) Publish(event async.BaseEvent) error {
	broker, ok := d.brokers[event.ID()]
	if !ok {
		return errBrokerNotFound(event.ID())
	}
	return broker.Publish(event)
}

// Subscribe subscribes the given channel to the broker with the given
// eventID. It will error if the channel type does not match the event type
// corresponding to the broker.
// Contract: the channel must be a Subscription[T], where T is the expected
// type of the event data.
func (d *Dispatcher) Subscribe(eventID async.EventID, ch any) error {
	broker, ok := d.brokers[eventID]
	if !ok {
		return errBrokerNotFound(eventID)
	}
	return broker.Subscribe(ch)
}

// Start will start all the brokers in the Dispatcher.
func (d *Dispatcher) Start(ctx context.Context) error {
	for _, broker := range d.brokers {
		go broker.Start(ctx)
	}
	return nil
}

func (d *Dispatcher) Name() string {
	return "dispatcher"
}

// RegisterBrokers registers the given brokers with the given eventIDs.
// Any subsequent events with <eventID> dispatched to this Dispatcher must be
// consistent with the type expected by <broker>.
func (d *Dispatcher) RegisterBrokers(
	brokers ...types.Broker,
) error {
	var ok bool
	for _, broker := range brokers {
		if _, ok = d.brokers[broker.EventID()]; ok {
			return errBrokerAlreadyExists(broker.EventID())
		}
		d.brokers[broker.EventID()] = broker
	}
	return nil
}
