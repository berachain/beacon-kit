package components

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

// ProvideStatusFeed provides a status feed.
func ProvideStatusFeed() *event.FeedOf[string, *service.StatusEvent] {
	return &event.FeedOf[string, *service.StatusEvent]{}
}
