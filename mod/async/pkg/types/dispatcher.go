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

package types

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

// Dispatcher is the full API for a dispatcher that facilitates the publishing
// of async events and the sending and receiving of async messages.
type Dispatcher interface {
	EventDispatcher
	// Start starts the dispatcher.
	Start(ctx context.Context) error
	// RegisterBrokers registers brokers to the dispatcher.
	RegisterBrokers(brokers ...Broker) error
	// Name returns the name of the dispatcher.
	Name() string
}

// EventDispatcher is the API for a dispatcher that facilitates the publishing
// of async events.
type EventDispatcher interface {
	// Publish publishes an event to the dispatcher.
	Publish(event events.BaseEvent) error
	// Subscribe subscribes the given channel to all events with the given event
	// ID.
	// Contract: the channel must be a Subscription[T], where T is the expected
	// type of the event data.
	Subscribe(eventID events.EventID, ch any) error
	// TODO: add unsubscribe
}

// publisher is the interface that supports basic event publisher operations.
type Broker interface {
	// Start starts the event publisher.
	Start(ctx context.Context)
	// Publish publishes the given event to the event publisher.
	Publish(event events.BaseEvent) error
	// Subscribe subscribes the given channel to the event publisher.
	Subscribe(ch any) error
	// Unsubscribe unsubscribes the given channel from the event publisher.
	Unsubscribe(ch any) error
	// EventID returns the event ID that the publisher is responsible for.
	EventID() events.EventID
}
