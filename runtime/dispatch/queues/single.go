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

package queues

import (
	"sync"
	"time"
)

// SingleDispatchQueue is a dispatch queue that dispatches a single item at a time.
// It respects the order of items added to the queue and will always
// process the freshest item that was MOST recently added to the queue.
type SingleDispatchQueue struct {
	mu       sync.Mutex
	item     chan WorkItem
	wg       sync.WaitGroup // WaitGroup for tracking in-flight work items.
	stopChan chan struct{}  // Channel for signaling stop.
}

// NewSingleDispatchQueue creates a new SingleDispatchQueue.
func NewSingleDispatchQueue() *SingleDispatchQueue {
	q := &SingleDispatchQueue{
		item:     make(chan WorkItem, 1),
		stopChan: make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-q.stopChan:
				return
			// No race condition exists between case item, ok := <-q.item in NewSingleDispatchQueue
			// and the q.item <- item in Async. The select statement in NewSingleDispatchQueue
			// listens for incoming items on the q.item channel in a separate goroutine and thus
			// the operation is not affected by the mutex lock in Async.
			case item, ok := <-q.item:
				if ok {
					item()
					q.wg.Done()
				}
			}
		}
	}()
	return q
}

// Async adds a work item to the queue to be executed asynchronously.
func (q *SingleDispatchQueue) Async(item WorkItem) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Remove the currently pending item before
	// adding the new one to the channel.
	select {
	case <-q.item:
		// Decrement the WaitGroup as the corresponding wg.Add(1) from the item
		// that is being removed from the channel is never called.
		q.wg.Done()
	default:
	}

	// Push the new item.
	q.wg.Add(1)
	q.item <- item
}

// AsyncAfter adds a work item to the queue to be executed after a specified duration.
func (q *SingleDispatchQueue) AsyncAfter(deadline time.Duration, execute WorkItem) {
	q.wg.Add(1)
	go func() {
		time.Sleep(deadline)
		q.Async(execute)
	}()
}

// Sync adds a work item to the queue and waits for its execution to complete.
func (q *SingleDispatchQueue) Sync(execute WorkItem) {
	done := make(chan struct{})
	q.Async(func() {
		execute()
		close(done)
	})
	<-done
}

// AsyncAndWait adds a work item to the queue to be executed asynchronously
// and waits for its execution to complete.
func (q *SingleDispatchQueue) AsyncAndWait(execute WorkItem) {
	q.Async(execute)
	q.wg.Wait()
}
