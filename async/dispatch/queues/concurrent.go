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

// defaultQueueSize is the default size of the queue.
const defaultQueueSize = 32

// ConcurrentQueue is a concurrent queue for dispatching work items.
type ConcurrentQueue struct {
	queue    chan WorkItem  // Channel for dispatching work items.
	wg       sync.WaitGroup // WaitGroup for tracking in-flight work items.
	stopChan chan struct{}  // Channel for signaling stop.
	stopped  bool           // Flag indicating if the queue has been stopped.
	mu       sync.Mutex     // Mutex for protecting stopped flag.
}

// NewConcurrentDispatchQueue creates a new Queue and starts its worker goroutines.
func NewConcurrentDispatchQueue(workerCount int) *ConcurrentQueue {
	q := &ConcurrentQueue{
		queue:    make(chan WorkItem, defaultQueueSize),
		stopChan: make(chan struct{}),
	}

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-q.stopChan:
					return
				case item := <-q.queue:
					if item != nil {
						item()
						q.wg.Done()
					}
				}
			}
		}()
	}

	return q
}

// Async adds a work item to the queue to be executed asynchronously.
func (q *ConcurrentQueue) Async(execute WorkItem) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.stopped {
		panic("Queue has been stopped")
	}

	q.wg.Add(1)
	q.queue <- execute
}

// AsyncAfter adds a work item to the queue to be executed after a specified duration.
func (q *ConcurrentQueue) AsyncAfter(deadline time.Duration, execute WorkItem) {
	q.wg.Add(1)
	go func() {
		time.Sleep(deadline)
		q.queue <- execute
	}()
}

// Sync adds a work item to the queue and waits for its execution to complete.
func (q *ConcurrentQueue) Sync(execute WorkItem) {
	done := make(chan struct{})
	q.Async(func() {
		execute()
		close(done)
	})
	<-done
}

// AsyncAndWait adds a work item to the queue and waits for all work items to complete.
func (q *ConcurrentQueue) AsyncAndWait(execute WorkItem) {
	q.Async(execute)
	q.wg.Wait()
}

// Stop stops the queue, preventing new work items from being added and waits for all
// in-flight work items to complete.
func (q *ConcurrentQueue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.stopped {
		return
	}

	// Wait for all tasks to complete
	q.wg.Wait()

	q.stopped = true

	// Close the queue channel to stop receiving new tasks
	close(q.queue)

	// Stop the workers
	close(q.stopChan)
}
