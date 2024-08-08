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

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// Publisher is responsible for broadcasting all events corresponding to the
// <eventID> to all registered client channels.
type Publisher[T any] struct {
	eventID types.EventID
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
func NewPublisher[T any](eventID string) *Publisher[T] {
	return &Publisher[T]{
		eventID: types.EventID(eventID),
		clients: make(map[chan T]struct{}),
		msgs:    make(chan T, defaultBufferSize),
		timeout: defaultPublisherTimeout,
		mu:      sync.Mutex{},
	}
}

func (p *Publisher[T]) EventID() types.EventID {
	return p.eventID
}

// Start starts the publisher loop.
func (p *Publisher[T]) Start(ctx context.Context) {
	go p.start(ctx)
}

// start starts the publisher loop.
func (p *Publisher[T]) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// close all leftover clients and break the publisher loop
			for client := range p.clients {
				p.Unsubscribe(client)
			}
			return
		case msg := <-p.msgs:
			// broadcast published msg to all clients
			for client := range p.clients {
				// send msg to client (or discard msg after timeout)
				select {
				case client <- msg:
				case <-time.After(p.timeout):
				}
			}
		}
	}
}

// Publish publishes a msg to the b.
// Returns ErrTimeout on timeout.
func (p *Publisher[T]) Publish(msg types.MessageI) error {
	typedMsg, err := ensureType[T](msg)
	if err != nil {
		return err
	}
	ctx := msg.Context()
	select {
	case p.msgs <- typedMsg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe registers the provided channel to the publisher,
// Returns ErrTimeout on timeout.
// TODO: see if its possible to accept a channel instead of any
func (p *Publisher[T]) Subscribe(ch any) error {
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[client] = struct{}{}
	return nil
}

// Unsubscribe removes a client from the publisher.
// Returns an error if the provided channel is not of type chan T.
func (p *Publisher[T]) Unsubscribe(ch any) error {
	client, err := ensureType[chan T](ch)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, client)
	close(client)
	return nil
}
