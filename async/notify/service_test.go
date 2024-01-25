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
	defer service.Stop()

	// Dispatch an event
	var event = &feed.Event{
		Type: 1,
		Data: "test",
	}
	service.Dispatch(feedName, event)

	time.Sleep(500 * time.Millisecond)

	// Check if the handler received the event
	if len(handler.receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(handler.receivedEvents))
	}
	if handler.receivedEvents[0] != event {
		t.Fatalf("Expected event %v, got %v", event, handler.receivedEvents[0])
	}
}
