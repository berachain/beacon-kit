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
	"runtime/debug"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

const StoreName = "blocks"

// KVStore is a simple KV store based implementation that stores beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	blocks *sdkcollections.IndexedMap[
		math.Slot, BeaconBlockT, indexes[BeaconBlockT],
	]
	nextToPrune math.Slot

	mu         sync.RWMutex
	cs         common.ChainSpec
	blockCodec *encoding.SSZInterfaceCodec[BeaconBlockT]
	logger     log.Logger[any]
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock[BeaconBlockT]](
	kvsp store.KVStoreService,
	cs common.ChainSpec,
	logger log.Logger[any],
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	blockCodec := &encoding.SSZInterfaceCodec[BeaconBlockT]{}
	return &KVStore[BeaconBlockT]{
		blocks: sdkcollections.NewIndexedMap(
			schemaBuilder,
			sdkcollections.NewPrefix(StoreName),
			StoreName,
			encoding.U64Key,
			blockCodec,
			newIndexes[BeaconBlockT](schemaBuilder),
		),
		blockCodec: blockCodec,
		cs:         cs,
		logger:     logger,
	}
}

// Get retrieves the block by a given index from the store.
func (kv *KVStore[BeaconBlockT]) Get(slot math.Slot) (BeaconBlockT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	kv.blockCodec.SetActiveForkVersion(kv.cs.ActiveForkVersionForSlot(slot))
	return kv.blocks.Get(context.TODO(), slot)
}

// Set sets the block by a given index in the store and also stores the
// block root.
func (kv *KVStore[BeaconBlockT]) Set(blk BeaconBlockT) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	slot := blk.GetSlot()
	kv.blockCodec.SetActiveForkVersion(kv.cs.ActiveForkVersionForSlot(slot))
	return kv.blocks.Set(context.TODO(), slot, blk)
}

// Prune removes the [start, end) blocks from the store.
func (kv *KVStore[BeaconBlockT]) Prune(start, end uint64) error {
	ctx := context.TODO()

	kv.mu.Lock()
	defer kv.mu.Unlock()

	s := max(math.Slot(start), kv.nextToPrune)
	for kv.nextToPrune = s; kv.nextToPrune < math.Slot(end); kv.nextToPrune++ {
		kv.blockCodec.SetActiveForkVersion(
			kv.cs.ActiveForkVersionForSlot(kv.nextToPrune),
		)
		kv.prune(ctx, kv.nextToPrune)
	}

	return nil
}

// prune removes the block at the given slot from the store. It handles panics
// and errors to avoid fatal crashes.
//
// NOTE: assumes the kvstore lock is held.
func (kv *KVStore[BeaconBlockT]) prune(ctx context.Context, slot math.Slot) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: add metrics here.
			// TODO: should also handle deleting the value manually from the db?

			kv.logger.Error(
				"‼️ panic occurred while pruning block",
				"slot", slot,
				"panic", r,
				"stack", debug.Stack(),
			)
		}
	}()

	if err := kv.blocks.Remove(ctx, slot); err != nil {
		// This can error for 2 reasons:
		// 1. The slot was not found -- either the slot was missed or we
		//    never stored the block to begin with, either way it's ok.
		if !errors.Is(err, sdkcollections.ErrNotFound) {
			// 2. The slot was found but (en/de)coding failed. In this
			//    case, we choose not to retry removal and instead
			//    continue. This means this slot may never be pruned, but
			//    ensures that we always get to pruning subsequent slots.
			kv.logger.Error(
				"‼️ failed to prune block", "slot", slot, "err", err,
			)
		}
	}
}
