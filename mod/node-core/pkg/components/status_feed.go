package components

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

// ProvideStatusFeed provides a status feed.
func ProvideStatusFeed() *event.FeedOf[types.EventID, *service.StatusEvent] {
	return &event.FeedOf[types.EventID, types.Event]{}
}
