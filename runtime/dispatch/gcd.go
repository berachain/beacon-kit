// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package dispatch

import (
	"fmt"
	"sync"

	"cosmossdk.io/log"
)

const Global = "global"

const (
	QueueTypeSingle = "single"
	QueueTypeSerial = "serial"
	QueueTypeConcur = "concurrent"
)

type GrandCentralDispatch struct {
	mu     sync.Mutex
	logger log.Logger
	queues map[string]Queue
}

// NewGrandCentralDispatch creates a new instance of GrandCentralDispatch
// and applies the provided options. The system and it's queue are inspired by
// Apple's Grand Central Dispatch, which is a system for managing concurrent
// code execution on darwin systems.
// https://developer.apple.com/documentation/dispatch
func NewGrandCentralDispatch(opts ...Option) (*GrandCentralDispatch, error) {
	gcd := &GrandCentralDispatch{
		queues: make(map[string]Queue),
	}

	// We create a global queue
	gcd.queues[Global] = NewSerialQueue()

	for _, opt := range opts {
		if err := opt(gcd); err != nil {
			return nil, err
		}
	}

	return gcd, nil
}

// Dispatch sends a value to the feed associated with the provided key.
func (gcd *GrandCentralDispatch) CreateQueue(style, id string) Queue {
	gcd.mu.Lock()
	defer gcd.mu.Unlock()

	// Check to make sure the queue doesn't already exist.
	_, ok := gcd.queues[id]
	if !ok {
		panic(fmt.Sprintf("queue already exists: %s", id))
	}

	var queue Queue
	switch style {
	case QueueTypeSingle:
		queue = NewSingleDispatchQueue()
	case QueueTypeSerial:
		queue = NewSerialQueue()
	case QueueTypeConcur:
		panic("not implemented")
		// queue = NewConcurrentDispatchQueue()
	default:
		panic("unknown queue type")
	}

	gcd.queues[id] = queue
	return queue
}

// GetQueue returns the feed associated with the provided key.
func (gcd *GrandCentralDispatch) GetQueue(id string) Queue {
	// Get the feed from the map.
	queue, ok := gcd.queues[id]
	if !ok {
		gcd.queues[id] = NewSerialQueue()
	}
	return queue
}
