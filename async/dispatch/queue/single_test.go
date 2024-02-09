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

package queue_test

import (
	"sync"
	"testing"

	"github.com/itsdevbear/bolaris/async/dispatch/queue"
)

func TestSingleDispatchQueueReplace(t *testing.T) {
	q := queue.NewSingleDispatchQueue()

	var (
		firstWorkStarted sync.WaitGroup
		allWorkDone      sync.WaitGroup
		mu               sync.Mutex
		cond             = sync.NewCond(&mu)
		output           = []int{}
	)
	// We only want to track if the first async function has started.
	firstWorkStarted.Add(1)
	// We want to wait for all the work to be done before we stop the queue.
	// We will only execute two tasks in total.
	allWorkDone.Add(2)

	if err := q.Async(func() {
		defer allWorkDone.Done()
		mu.Lock()
		defer mu.Unlock()

		// Mark the first async function as started.
		firstWorkStarted.Done()

		// Block on the condition variable to simulate a long-running task.
		cond.Wait()
		output = append(output, 1)
	}); err != nil {
		t.Errorf("Unexpected err %v", err)
	}

	// Wait for the first async function to start
	// before enqueueing the next two.
	firstWorkStarted.Wait()

	// These tasks should get replaced over and over by each other.
	for i := 2; i < 69; i++ {
		if err := q.Async(func() {
			defer allWorkDone.Done()

			mu.Lock()
			defer mu.Unlock()
			output = append(output, i)
		}); err != nil {
			t.Errorf("Unexpected err %v", err)
		}
	}

	// Since the first Async called hasn't exited yet (it's waiting on the condition variable),
	// the last Async should be enqueued and all others should've been replaced.
	if err := q.Async(func() {
		defer allWorkDone.Done()

		mu.Lock()
		defer mu.Unlock()
		output = append(output, 69)
	}); err != nil {
		t.Errorf("Unexpected err %v", err)
	}

	// Signal the condition variable to wake up the first async function.
	cond.Signal()

	// Wait for all running tasks to finish.
	allWorkDone.Wait()

	// The length of the output should be 2.
	if len(output) != 2 {
		t.Errorf("Expected output array length of 2, got %d", len(output))
	}

	// The first element should be 1 and the second should be 69.
	if output[0] != 1 || output[1] != 69 {
		t.Errorf("Expected output to be [1,69], got %d", output)
	}

	q.Stop()
}
