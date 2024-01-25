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

package notify_test

import (
	"sync"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/async/dispatch"
	"github.com/itsdevbear/bolaris/async/notify"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/feed"
)

type TestEvent struct {
	Msg string
}

type TestHandler struct {
	receivedEvents []interface{}
	mu             sync.Mutex
}

func (h *TestHandler) HandleNotification(event interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.receivedEvents = append(h.receivedEvents, event)
}

func TestDispatch(t *testing.T) {
	// Register a handler
	handler := &TestHandler{}
	queueID := "testQueue"

	// Setup GCD and Service.
	gcd, _ := dispatch.NewGrandCentralDispatch(
		dispatch.WithLogger(log.NewNopLogger()),
		dispatch.WithDispatchQueue(queueID, dispatch.QueueTypeSerial),
	)
	service := notify.NewService(
		notify.WithLogger(log.NewNopLogger()),
		notify.WithGCD(gcd),
	)

	// Register a feed
	feedName := "testFeed"
	service.RegisterFeed(feedName)

	err := service.RegisterHandler(feedName, queueID, handler)
	if err != nil {
		t.Fatalf("Failed to register handler: %v", err)
	}

	// Start the service
	service.Start()
	defer func() {
		if err = service.Stop(); err != nil {
			t.Fatalf("Failed to stop service: %v", err)
		}
	}()

	// Dispatch an event
	var event = &feed.Event{
		Type: 1,
		Data: "test",
	}
	service.Dispatch(feedName, event)

	time.Sleep(100 * time.Millisecond)

	// Check if the handler received the event
	if len(handler.receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(handler.receivedEvents))
	}
	if handler.receivedEvents[0] != event {
		t.Fatalf("Expected event %v, got %v", event, handler.receivedEvents[0])
	}
}
