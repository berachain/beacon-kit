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

package broker

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/primitives/async"
)

// Broker is the unique publisher for broadcasting all events corresponding
// to the <eventID> to all registered client channels.
// There should be exactly one broker responsible for every eventID.
type Broker[T async.BaseEvent] struct {
	// eventID is a unique identifier for the event that this broker is
	// responsible for.
	eventID async.EventID
	// subscriptions is a map of subscribed subscriptions.
	subscriptions map[chan T]struct{}
	// msgs is the channel for publishing new messages.
	msgs chan T
	// timeout is the timeout for sending a msg to a client.
	timeout time.Duration
}

// New creates a new broker publishing events of type T for the
// provided eventID.
func New[T async.BaseEvent](eventID string) *Broker[T] {
	return &Broker[T]{
		eventID:       async.EventID(eventID),
		subscriptions: make(map[chan T]struct{}),
		msgs:          make(chan T, defaultBufferSize),
		timeout:       defaultBrokerTimeout,
	}
}

// EventID returns the event ID that the broker is responsible for.
func (b *Broker[T]) EventID() async.EventID {
	return b.eventID
}

// Start starts the broker loop in a goroutine to listen and broadcast events
// to all subscribers.
func (b *Broker[T]) Start(ctx context.Context) {
	go b.start(ctx)
}

// start is a helper function to listen and broadcast events.
func (b *Broker[T]) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// close all leftover clients and break the broker loop
			b.shutdown()
			return
		case msg := <-b.msgs:
			// broadcast published msg to all clients
			b.broadcast(msg)
		}
	}
}

// Publish publishes a msg to all subscribers.
// Errors if the message is not of type T, or if the context is canceled.
func (b *Broker[T]) Publish(msg async.BaseEvent) error {
	// assert that the message is of type T
	typedMsg, err := ensureType[T](msg)
	if err != nil {
		return err
	}
	ctx := msg.Context()
	select {
	case b.msgs <- typedMsg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe registers the provided channel to the broker.
// Errors if the channel is not of type chan T.
// Contract: the channel must be a Chan[T], where T is the expected
// type of the event data.
func (b *Broker[T]) Subscribe(ch any) error {
	// assert that the channel is of type chan T
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	b.subscriptions[client] = struct{}{}
	return nil
}

// Unsubscribe removes a client from the broker.
// Returns an error if the provided channel is not of type chan T.
func (b *Broker[T]) Unsubscribe(ch any) error {
	// assert that the channel is of type chan T
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	delete(b.subscriptions, client)
	close(client)
	return nil
}

func (b *Broker[T]) broadcast(msg T) {
	for client := range b.subscriptions {
		// send msg to client (or discard msg after timeout)
		// we could consider using a go routine for each client to allow
		// for concurrent notification attempts, while respecting the timeout
		// for each client individually
		select {
		case client <- msg:
		case <-time.After(b.timeout):
		}
	}
}

// shutdown closes all leftover clients.
func (b *Broker[T]) shutdown() {
	for client := range b.subscriptions {
		if err := b.Unsubscribe(client); err != nil {
			panic(err)
		}
	}
}
