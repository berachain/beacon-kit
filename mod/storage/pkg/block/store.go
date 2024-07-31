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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

// var _ pruner.Prunable = (*KVStore[BeaconBlock])(nil)

const (
	KeyBlockPrefix      = "block"
	KeyRootsPrefix      = "roots"
	KeyTimestampsPrefix = "timestamps"
)

type KVStoreProvider struct {
	store.KVStoreWithBatch
}

// OpenKVStore opens a new KV store.
func (p *KVStoreProvider) OpenKVStore(context.Context) store.KVStore {
	return p.KVStoreWithBatch
}

// KVStore is a simple KV store based implementation that stores
// beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	blocks     sdkcollections.Map[math.Slot, BeaconBlockT]
	roots      sdkcollections.Map[[]byte, math.Slot]
	timestamps sdkcollections.Map[math.U64, math.Slot]

	mu           sync.RWMutex
	cdc          *encoding.SSZInterfaceCodec[BeaconBlockT]
	earliestSlot uint64
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock[BeaconBlockT]](
	kvsp store.KVStoreService,
) *KVStore[BeaconBlockT] {
	schemaBuilder := sdkcollections.NewSchemaBuilder(kvsp)
	cdc := &encoding.SSZInterfaceCodec[BeaconBlockT]{}
	return &KVStore[BeaconBlockT]{
		blocks: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyBlockPrefix)),
			KeyBlockPrefix,
			encoding.U64Key,
			cdc,
		),
		roots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyRootsPrefix)),
			KeyRootsPrefix,
			sdkcollections.BytesKey,
			encoding.U64Value,
		),
		timestamps: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte(KeyTimestampsPrefix)),
			KeyTimestampsPrefix,
			encoding.U64Key,
			encoding.U64Value,
		),
		cdc: cdc,
	}
}

// Get retrieves the block by a given index from the store.
func (kv *KVStore[BeaconBlockT]) Get(slot math.Slot) (BeaconBlockT, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return kv.blocks.Get(context.TODO(), slot)
}

// Set sets the block by a given index in the store and also stores the
// block root.
func (kv *KVStore[BeaconBlockT]) Set(slot math.Slot, blk BeaconBlockT) error {
	root, err := blk.HashTreeRoot()
	if err != nil {
		return err
	}
	ctx := context.TODO()

	kv.mu.Lock()
	defer kv.mu.Unlock()

	if err = kv.roots.Set(ctx, root[:], slot); err != nil {
		return err
	}

	if err = kv.timestamps.Set(ctx, blk.GetTimestamp(), slot); err != nil {
		return err
	}

	kv.cdc.SetActiveForkVersion(blk.Version())
	return kv.blocks.Set(ctx, slot, blk)
}

// GetSlotByRoot retrieves the slot by a given root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByRoot(
	root [32]byte,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	slot, err := kv.roots.Get(context.TODO(), root[:])
	if err != nil {
		return 0, err
	}
	return math.Slot(slot), nil
}

// Prune removes the [start, end) blocks from the store.
func (kv *KVStore[BeaconBlockT]) Prune(start, end uint64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for i := max(start, kv.earliestSlot); i < end; i++ {
		nextBlock, err := kv.blocks.Get(context.TODO(), math.Slot(i+1))
		if !errors.Is(err, sdkcollections.ErrNotFound) {
			if err != nil {
				return err
			}

			root := nextBlock.GetParentBlockRoot()
			if err = kv.roots.Remove(context.TODO(), root[:]); err != nil {
				return err
			}
		}

		// This only errors if the key passed in cannot be encoded.
		if err = kv.blocks.Remove(context.TODO(), math.Slot(i)); err != nil {
			return err
		}
	}
	kv.earliestSlot = end
	return nil
}
