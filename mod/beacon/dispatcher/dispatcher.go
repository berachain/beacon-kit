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

type Dispatcher struct {
	mu sync.RWMutex

	registry map[EventType]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		registry: make(map[EventType]Handler),
	}
}

func (d *Dispatcher) Start(ctx context.Context, event EventType) error {
	return d.registry[event].Start(ctx)
}

func (d *Dispatcher) StartAll(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, handler := range d.registry {
		if err := handler.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) RegisterHandler(eventType EventType, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.registry[eventType] = handler
}

func (d *Dispatcher) Notify(event Event) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handler, ok := d.registry[event.Type()]
	if !ok {
		return ErrHandlerNotFound
	}

	handler.Notify(event)
	return nil
}
