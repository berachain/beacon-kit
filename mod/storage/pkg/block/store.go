// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package block

import (
	"context"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

// KVStore is a simple KV store based implementation that stores beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	blocks *sdkcollections.IndexedMap[
		math.Slot, BeaconBlockT, indexes[BeaconBlockT],
	]

	blockRoots       sdkcollections.Map[[]byte, math.Slot]
	executionNumbers sdkcollections.Map[math.U64, math.Slot]
	stateRoots       sdkcollections.Map[[]byte, math.Slot]

	mu     sync.RWMutex
	cs     common.ChainSpec
	logger log.Logger[any]
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock[BeaconBlockT]](
	kvsp store.KVStoreService,
	cs common.ChainSpec,
	logger log.Logger[any],
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	return &KVStore[BeaconBlockT]{
		blockRoots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(blockRootsPrefix),
			blockRootsName,
			sdkcollections.BytesKey,
			encoding.U64Value,
		),
		executionNumbers: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(executionNumbersPrefix),
			executionNumbersName,
			encoding.U64Key,
			encoding.U64Value,
		),
		stateRoots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(stateRootsPrefix),
			stateRootsName,
			sdkcollections.BytesKey,
			encoding.U64Value,
		),
		cs:     cs,
		logger: logger,
	}
}

// Set sets the block by a given index in the store and also stores the
// block root.
func (kv *KVStore[BeaconBlockT]) Set(blk BeaconBlockT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	slot := blk.GetSlot()
	return kv.blocks.Set(context.TODO(), slot, blk)
}

// GetSlotByRoot retrieves the slot by a given root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByBlockRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.BlockRoots.MatchExact(
		context.TODO(), root[:],
	)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func (kv *KVStore[BeaconBlockT]) GetSlotByStateRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.StateRoots.MatchExact(
		context.TODO(), root[:],
	)
	if err != nil {
		return 0, err
	}
	return slot, nil
}

// GetSlotByExecutionNumber retrieves the slot by a given execution number from
// the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByExecutionNumber(
	executionNumber math.U64,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.blocks.Indexes.ExecutionNumbers.MatchExact(
		context.TODO(), executionNumber,
	)
	if err != nil {
		return 0, err
	}
	return slot, nil
}
