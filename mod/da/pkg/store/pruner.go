package store

import (
	"github.com/berachain/beacon-kit/mod/primitives"
)

func GetPruneParamsFn[
	BeaconBlockT BeaconBlock,
	BlockEventT BlockEvent[BeaconBlockT],
](cs primitives.ChainSpec) func(BlockEventT) (uint64, uint64) {
	return func(event BlockEventT) (uint64, uint64) {
		blk := event.Block()
		window := cs.MinEpochsForBlobsSidecarsRequest() * cs.SlotsPerEpoch()
		startIndex := blk.GetSlot().Unwrap() - window
		return startIndex, window
	}
}
