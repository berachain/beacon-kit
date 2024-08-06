package dispatcher

import (
	"context"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

type MessageServer interface {
	RegisterRecipient(mID types.MessageID, ch chan any) error
	Request(req types.Message[any], resp any) error
	Respond(resp types.Message[any]) error
	RegisterRoute(mID types.MessageID, route messageRoute) error
}

type EventServer interface {
	Start(ctx context.Context)
	RegisterFeed(mID types.MessageID, feed publisher)
	Subscribe(mID types.MessageID, ch chan any) error
	Dispatch(ctx context.Context, event *types.Message[any])
}

// publisher is the interface that supports basic event feed operations.
type publisher interface {
	// Start starts the event feed.
	Start(ctx context.Context)
	// Publish publishes the given event to the event feed.
	Publish(ctx context.Context, event any)
	// Subscribe subscribes the given channel to the event feed.
	Subscribe(ch any) error
	// Unsubscribe unsubscribes the given channel from the event feed.
	Unsubscribe(ch any) error
}

// messageRoute is the interface that supports basic message route operations.
type messageRoute interface {
	// RegisterRecipient sets the recipient for the route.
	RegisterRecipient(ch chan any) error
	// SendRequest sends a request to the recipient.
	SendRequest(msg any) error
	// SendResponse sends a response to the recipient.
	SendResponse(msg any) error
	// AwaitResponse awaits a response from the route.
	AwaitResponse(emptyResp any) error
}
