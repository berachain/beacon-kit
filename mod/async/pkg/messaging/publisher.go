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

package messaging

import (
	"context"
	"sync"
	"time"
)

// Publisher broadcasts msgs to registered clients.
type Publisher[T any] struct {
	// clients is a map of subscribed clients.
	clients map[chan T]struct{}
	// msgs is the channel for publishing new messages.
	msgs chan T
	// timeout is the timeout for sending a msg to a client.
	timeout time.Duration
	// mu is the mutex for the clients map.
	mu sync.Mutex
}

// NewPublisher creates a new b.
func NewPublisher[T any](name string) *Publisher[T] {
	return &Publisher[T]{
		clients: make(map[chan T]struct{}),
		msgs:    make(chan T, defaultBufferSize),
		timeout: defaultPublisherTimeout,
		mu:      sync.Mutex{},
	}
}

// Start starts the broker loop.
func (b *Publisher[T]) Start(ctx context.Context) {
	go b.start(ctx)
}

// start starts the broker loop.
func (b *Publisher[T]) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// close all leftover clients and break the broker loop
			for client := range b.clients {
				b.Unsubscribe(client)
			}
			return
		case msg := <-b.msgs:
			// broadcast published msg to all clients
			for client := range b.clients {
				// send msg to client (or discard msg after timeout)
				select {
				case client <- msg:
				case <-time.After(b.timeout):
				}
			}
		}
	}
}

// Publish publishes a msg to the b.
// Returns ErrTimeout on timeout.
func (b *Publisher[T]) Publish(ctx context.Context, msg any) error {
	typedMsg, err := ensureType[T](msg)
	if err != nil {
		return err
	}
	select {
	case b.msgs <- typedMsg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe registers the provided channel to the publisher,
// Returns ErrTimeout on timeout.
// TODO: see if its possible to accept a channel instead of any
func (b *Publisher[T]) Subscribe(ch any) error {
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = struct{}{}
	return nil
}

// Unsubscribe removes a client from the publisher.
// Returns an error if the provided channel is not of type chan T.
func (f *Publisher[T]) Unsubscribe(ch any) error {
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.clients, client)
	close(client)
	return nil
}
