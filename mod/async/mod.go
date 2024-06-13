package async

import "github.com/ethereum/go-ethereum/event"

type (
	Feed          = event.Feed
	FeedOf[T any] struct {
		event.FeedOf[T]
	}
)
