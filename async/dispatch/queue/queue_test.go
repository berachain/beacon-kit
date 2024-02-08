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
	"sync/atomic"
	"testing"
	"time"

	"github.com/itsdevbear/bolaris/async/dispatch/queue"
	"github.com/stretchr/testify/assert"
)

func TestDispatchQueueConcurrent_Async(t *testing.T) {
	q := queue.NewDispatchQueue(4, 4)

	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)

	wg.Add(10)

	for i := 0; i < 10; i++ {
		if err := q.Async(func() {
			defer wg.Done()
			counter.Add(1)
		}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	wg.Wait()

	if counter.Load() != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter.Load())
	}

	q.Stop()
}

func TestDispatchQueueConcurrent_AsyncAfter(t *testing.T) {
	q := queue.NewDispatchQueue(4, 4)

	var asyncAfterExecuted bool

	wg := &sync.WaitGroup{}
	wg.Add(1)

	startTime := time.Now()
	waitTime := time.Millisecond * 100
	err := q.AsyncAfter(waitTime, func() {
		asyncAfterExecuted = true
		wg.Done()
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wg.Wait()

	if !asyncAfterExecuted {
		t.Errorf("AsyncAfter function did not execute")
	}

	if time.Since(startTime) < waitTime {
		t.Errorf("AsyncAfter function executed earlier than expected")
	}

	q.Stop()
}

func TestDispatchQueueConcurrent_Sync(t *testing.T) {
	q := queue.NewDispatchQueue(4, 4)

	var syncExecuted bool

	err := q.Sync(func() {
		syncExecuted = true
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !syncExecuted {
		t.Errorf("Sync function did not execute")
	}

	q.Stop()
}

func TestDispatchQueueConcurrent_AsyncAndWait(t *testing.T) {
	q := queue.NewDispatchQueue(4, 4)

	var asyncAndWaitExecuted bool

	err := q.AsyncAndWait(func() {
		asyncAndWaitExecuted = true
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !asyncAndWaitExecuted {
		t.Errorf("AsyncAndWait function did not execute")
	}

	q.Stop()
}

func TestDispatchQueueConcurrent_Stop(t *testing.T) {
	q := queue.NewDispatchQueue(4, 4)

	// Add some items to the queue
	for i := 0; i < 10; i++ {
		if err := q.Async(func() {
			time.Sleep(time.Millisecond * 100)
		}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	// Stop the queue
	q.Stop()

	// Try to add another item to the queue, it should panic
	defer func() {
	}()

	if err := q.Async(func() {
		t.Errorf("Async function executed after Stop")
	}); err == nil {
		assert.Equal(t, err, queue.ErrAddToStoppedQueue)
	}
}
