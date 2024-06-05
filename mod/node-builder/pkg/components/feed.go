package components

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/ethereum/go-ethereum/event"
)

// ProvideBlockFeed provides a block feed for the depinject framework.
func ProvideBlockFeed() *event.FeedOf[events.Block[*types.BeaconBlock]] {
	return &event.FeedOf[events.Block[*types.BeaconBlock]]{}
}
