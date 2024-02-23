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

package queue

// SingleDispatchQueue dispatches a single item at a time, maintaining order.
type SingleDispatchQueue struct {
	*DispatchQueue
}

// NewSingleDispatchQueue creates a new instance.
func NewSingleDispatchQueue() *SingleDispatchQueue {
	return &SingleDispatchQueue{DispatchQueue: NewDispatchQueue(1, 1)}
}

// Async executes a work item asynchronously, replacing any pending item.
func (q *SingleDispatchQueue) Async(item WorkItem) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Replace the pending item with the new one.
	select {
	case <-q.queue:
		// Adjust WaitGroup for the removed item.
		q.wg.Done()
	default:
		// No action for an empty queue.
	}

	// Queue the new item.
	q.wg.Add(1)
	q.queue <- item
	return nil
}
