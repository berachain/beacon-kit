// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package notify

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/event"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// Service represents the BeaconKit notification service. It is used to register
// event feeds and handlers, and dispatch events to the handlers.
// It leverages GrandCentralDispatch to dispatch events to handlers.
type Service struct {
	service.BaseService

	// running is a boolean that indicates whether the service is running or
	// not.
	running bool

	// feeds is a map that holds the event feeds.
	feeds map[string]*event.Feed

	// handlers is a map that holds the event handler queue pairs.
	handlers map[string][]eventHandlerQueuePair

	// dispatch is an instance of GrandCentralDispatch.
	gcd GrandCentralDispatch
}

// Start spawns any goroutines required by the service.
func (s *Service) Start(ctx context.Context) {
	for name, handlers := range s.handlers {
		feed, ok := s.feeds[name]
		if !ok {
			continue
		}

		for _, pair := range handlers {
			// Create a channel for the handler
			ch := make(chan *Event)
			subscription := feed.Subscribe(ch)

			// Start a goroutine to listen for events and call the handler
			go func(
				pair eventHandlerQueuePair,
				ch <-chan *Event,
				subscription ethereum.Subscription,
			) {
				for {
					select {
					case event := <-ch:
						// Use the dispatch queue to call the handler's Handle
						// method
						// asynchronously
						if err := s.gcd.GetQueue(pair.queueID).Async(func() {
							pair.handler.HandleNotification(event)
						}); err != nil {
							// Choosing to panic here because it doesn't make
							// sense for the service we're controlling to have
							// stopped the queue
							panic(err)
						}
					case <-subscription.Err():
						return
					case <-ctx.Done():
						s.running = false
						// This will receive a value when the stop channel is
						// closed
						return
					}
				}
			}(pair, ch, subscription)
		}
	}
	s.running = true
}

// Status returns error if the service is not considered healthy.
func (s *Service) Status() error { return nil }

// RegisterFeed registers a new feed associated with the provided key.
func (s *Service) RegisterFeed(name string) {
	if s.running {
		panic(ErrRegisterFeedServiceStarted)
	}
	if _, ok := s.feeds[name]; !ok {
		s.feeds[name] = &event.Feed{}
	}
}

// RegisterHandler registers a new handler associated with the provided key.
// It also takes a queueID which is used to dispatch the handler on.
func (s *Service) RegisterHandler(
	name string, queueID string, handler EventHandler,
) error {
	if s.running {
		panic(ErrRegisterFeedServiceStarted)
	}

	_, found := s.feeds[name]
	if !found {
		return ErrFeedNotFound
	}

	s.handlers[name] = append(s.handlers[name], eventHandlerQueuePair{
		handler: handler,
		queueID: queueID,
	})
	return nil
}

// Dispatch dispatches an event to all handlers associated
// with the provided key.
func (s *Service) Dispatch(feedName string, event any) {
	feed, ok := s.feeds[feedName]
	if ok {
		feed.Send(event)
	}
}
