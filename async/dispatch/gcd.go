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
	dqueue "github.com/itsdevbear/bolaris/async/dispatch/queue"
)

const (
	// GlobalQueueID is the identifier for the global queue.
	GlobalQueueID = "global"

	// DefaultQueueSize is the default size of a queue.
	DefaultQueueSize = uint16(64) // todo: make this configurable.

	// DefaultConcurrentQueueWorkerCount is the default size of a concurrent queue.
	DefaultConcurrentQueueWorkerCount = uint16(64) // todo: make this configurable.
)

// QueueType represents the type of a queue.
type QueueType string

// Constants for the different types of queues.
const (
	// QueueTypeSingle represents a single queue.
	QueueTypeSingle QueueType = "single"

	// QueueTypeSerial represents a serial queue.
	QueueTypeSerial QueueType = "serial"

	// QueueTypeConcur represents a concurrent queue.
	QueueTypeConcur QueueType = "concurrent"
)

// GrandCentralDispatch is a structure that holds the mutex, logger and queues.
type GrandCentralDispatch struct {
	mu     sync.RWMutex
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
	gcd.queues[GlobalQueueID] = dqueue.NewDispatchQueue(
		1,
		DefaultQueueSize,
	)

	for _, opt := range opts {
		if err := opt(gcd); err != nil {
			return nil, err
		}
	}

	return gcd, nil
}

// Dispatch sends a value to the feed associated with the provided key.
func (gcd *GrandCentralDispatch) CreateQueue(id string, queueType QueueType) Queue {
	gcd.mu.Lock()
	defer gcd.mu.Unlock()

	// Check to make sure the queue doesn't already exist.
	_, found := gcd.queues[id]
	if found {
		panic(fmt.Sprintf("queue already exists: %s", id))
	}

	var queue Queue
	switch queueType {
	case QueueTypeSingle:
		queue = dqueue.NewSingleDispatchQueue()
	case QueueTypeSerial:
		queue = dqueue.NewDispatchQueue(1, DefaultQueueSize)
	case QueueTypeConcur:

		queue = dqueue.NewDispatchQueue(DefaultConcurrentQueueWorkerCount, DefaultQueueSize)
	default:
		panic("unknown queue type")
	}

	gcd.logger.Info("intialized new dispatch queue", "id", id, "type", queueType)
	gcd.queues[id] = queue
	return queue
}

// GetQueue returns the feed associated with the provided key.
func (gcd *GrandCentralDispatch) GetQueue(id string) Queue {
	gcd.mu.RLock()
	defer gcd.mu.RUnlock()
	queue, ok := gcd.queues[id]
	if !ok {
		return nil
	}
	return queue
}
