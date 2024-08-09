package block

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// GetSlotByRoot retrieves the slot by a given root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByBlockRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.BlockRoots.MatchExact(
		context.TODO(),
		root[:],
	)
	if err != nil {
		return 0, err
	}
	return math.Slot(slot), nil
}

func (kv *KVStore[BeaconBlockT]) GetSlotByStateRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.StateRoots.MatchExact(
		context.TODO(),
		root[:],
	)
	if err != nil {
		return 0, err
	}
	return math.Slot(slot), nil
}

// GetSlotByExecutionNumber retrieves the slot by a given execution number from
// the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByExecutionNumber(
	executionNumber math.U64,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.ExecutionNumbers.MatchExact(
		context.TODO(),
		executionNumber.Unwrap(),
	)
	if err != nil {
		return 0, err
	}
	return math.Slot(slot), nil
}
