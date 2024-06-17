package components

import (
	"github.com/berachain/beacon-kit/mod/async/pkg/event"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

// ProvideStatusFeed provides a status feed.
func ProvideStatusFeed() *event.FeedOf[
	asynctypes.EventID, *asynctypes.Event[*service.StatusEvent],
] {
	return &event.FeedOf[
		asynctypes.EventID, *asynctypes.Event[*service.StatusEvent],
	]{}
}
