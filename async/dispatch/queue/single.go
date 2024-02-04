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

package queue

// SingleDispatchQueue is a dispatch queue that dispatches a single item at a time.
// It respects the order of items added to the queue and will always
// process the freshest item that was MOST recently added to the queue.
type SingleDispatchQueue struct {
	*DispatchQueue
}

// NewSingleDispatchQueue creates a new SingleDispatchQueue.
func NewSingleDispatchQueue() *SingleDispatchQueue {
	q := &SingleDispatchQueue{
		DispatchQueue: NewDispatchQueue(1, 1),
	}
	return q
}

// Async adds a work item to the queue to be executed asynchronously.
func (q *SingleDispatchQueue) Async(item WorkItem) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Remove the currently pending item before
	// adding the new one to the channel.
	select {
	case <-q.queue:
		// Decrement the WaitGroup as the corresponding wg.Add(1) from the item
		// that is being removed from the channel is never called.
		q.wg.Done()
	default:
		// If there is no item in the channel, do nothing.
	}

	// Push the new item.
	q.wg.Add(1)
	q.queue <- item
}
