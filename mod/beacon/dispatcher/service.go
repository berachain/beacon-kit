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

package dispatcher

import (
	"context"
	"sync"
)

const Name = "dispatcher"

type Service struct {
	mu sync.RWMutex

	registry map[EventType][]chan (Event)
}

func NewService() (*Service, error) {
	dispatcher := &Service{
		mu:       sync.RWMutex{},
		registry: make(map[EventType][]chan (Event)),
	}

	return dispatcher, nil
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

// RegisterHandler registers a handler for a given event type
func (s *Service) RegisterHandler(eventType EventType) chan (Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	handler := make(chan (Event), 1)
	s.registry[eventType] = append(s.registry[eventType], handler)

	return handler
}

// Notify sends an event to all handlers registered for the event type
func (s *Service) Notify(event Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	handlers, ok := s.registry[event.Type()]
	if !ok {
		return ErrHandlerNotFound
	}

	for _, handler := range handlers {
		handler <- event
	}

	return nil
}

func (s *Service) Name() string {
	return Name
}

func (s *Service) Status() error {
	return nil
}

func (s *Service) WaitForHealthy(ctx context.Context) {}
