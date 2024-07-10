package event

import "context"

// Feed is a generic interface for sending events.
type Feed[EventT any] interface {
	// Send sends an event and returns the number of
	// subscribers that received it.
	Publish(ctx context.Context, event EventT) error
	// Subscribe returns a channel that will receive events.
	Subscribe() (chan EventT, error)
}
