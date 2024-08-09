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
	"errors"
	"sync"

	sdkcollections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

const StoreName = "blocks"

// KVStore is a simple KV store based implementation that stores beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	blocks *sdkcollections.IndexedMap[
		uint64, BeaconBlockT, indexes[BeaconBlockT],
	]
	prevBlockSlot uint64

	mu           sync.RWMutex
	cs           common.ChainSpec
	blockCodec   *encoding.SSZInterfaceCodec[BeaconBlockT]
	earliestSlot uint64
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock[BeaconBlockT]](
	kvsp store.KVStoreService,
	cs common.ChainSpec,
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	blockCodec := &encoding.SSZInterfaceCodec[BeaconBlockT]{}
	return &KVStore[BeaconBlockT]{
		blocks: sdkcollections.NewIndexedMap(
			schemaBuilder,
			sdkcollections.NewPrefix(StoreName),
			StoreName,
			sdkcollections.Uint64Key,
			blockCodec,
			newIndexes[BeaconBlockT](schemaBuilder),
		),
		blockCodec:    blockCodec,
		cs:            cs,
		earliestSlot:  1,
		prevBlockSlot: 0,
	}
}

// Get retrieves the block by a given index from the store.
func (kv *KVStore[BeaconBlockT]) Get(slot math.Slot) (BeaconBlockT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	kv.blockCodec.SetActiveForkVersion(kv.cs.ActiveForkVersionForSlot(slot))
	return kv.blocks.Get(context.TODO(), slot.Unwrap())
}

// Set sets the block by a given index in the store and also stores the
// block root.
func (kv *KVStore[BeaconBlockT]) Set(
	slot math.Slot,
	prevStateRoot common.Root,
	blk BeaconBlockT,
) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if err := kv.catchupBlockData(prevStateRoot); err != nil {
		return err
	}
	kv.blockCodec.SetActiveForkVersion(kv.cs.ActiveForkVersionForSlot(slot))
	if err := kv.blocks.Set(context.TODO(), slot.Unwrap(), blk); err != nil {
		return err
	}
	kv.prevBlockSlot = slot.Unwrap()
	return nil
}

// Prune removes the [start, end) blocks from the store.
func (kv *KVStore[BeaconBlockT]) Prune(start, end uint64) error {
	var ctx = context.TODO()
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// We only return early from this loop with an error if the key
	// passed in cannot be encoded.
	s := max(start, kv.earliestSlot)
	for i := s; i < end; i++ {
		kv.blockCodec.SetActiveForkVersion(
			kv.cs.ActiveForkVersionForSlot(math.Slot(i)),
		)
		if err := kv.blocks.Remove(ctx, i); err != nil {
			if errors.Is(err, sdkcollections.ErrNotFound) {
				// Either the slot was missed or we never stored
				// the block to begin with, either way it's ok.
				continue
			}
			return err
		}
		// Update earliest slot as we go, an error will still
		// have committed any removes up to that point.
		kv.earliestSlot = i
	}

	// If we successfully pruned, update the earliest slot.
	kv.earliestSlot = end
	return nil
}

// catchupBlockData updates the block data for the previous block slot.
//
// This is used to catch up the block data for the previous block slot after
// the block is complete and 'correct'.
func (kv *KVStore[BeaconBlockT]) catchupBlockData(
	prevStateRoot common.Root,
) error {
	if kv.prevBlockSlot == 0 {
		return nil
	}
	prevBlock, err := kv.blocks.Get(context.TODO(), kv.prevBlockSlot)
	if err != nil {
		return err
	}
	prevBlock.SetStateRoot(prevStateRoot)
	return kv.blocks.Set(context.TODO(), kv.prevBlockSlot, prevBlock)
}
