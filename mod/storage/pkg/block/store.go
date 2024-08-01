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

// KVStore is a simple KV store based implementation that stores beacon blocks.
type KVStore[BeaconBlockT BeaconBlock[BeaconBlockT]] struct {
	blocks           sdkcollections.Map[math.Slot, BeaconBlockT]
	roots            sdkcollections.Map[[]byte, math.Slot]
	executionNumbers sdkcollections.Map[math.U64, math.Slot]

	mu           sync.RWMutex
	cdc          *encoding.SSZInterfaceCodec[BeaconBlockT]
	earliestSlot math.Slot
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
			sdkcollections.NewPrefix([]byte{BlockKeyPrefix}),
			BlocksMapName,
			encoding.U64Key,
			cdc,
		),
		roots: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{RootsKeyPrefix}),
			RootsMapName,
			sdkcollections.BytesKey,
			encoding.U64Value,
		),
		executionNumbers: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix([]byte{ExecutionNumbersKeyPrefix}),
			ExecutionNumbersMapName,
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
	var (
		ctx  = context.TODO()
		root = blk.HashTreeRoot()
		err  error
	)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Set the block root in the roots map.
	if err = kv.roots.Set(ctx, root[:], slot); err != nil {
		return err
	}

	// Set the block execution number in the execution numbers map.
	if err = kv.executionNumbers.Set(
		ctx, blk.GetExecutionNumber(), slot,
	); err != nil {
		return err
	}

	// Set the block in the blocks map.
	kv.cdc.SetActiveForkVersion(blk.Version())
	return kv.blocks.Set(ctx, slot, blk)
}

// GetSlotByRoot retrieves the slot by a given root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.roots.Get(context.TODO(), root[:])
}

// GetSlotByExecutionNumber retrieves the slot by a given execution number from
// the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByExecutionNumber(
	executionNumber math.U64,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.executionNumbers.Get(context.TODO(), executionNumber)
}

// Prune removes the [start, end) blocks from the store.
func (kv *KVStore[BeaconBlockT]) Prune(start, end uint64) error {
	var (
		ctx  = context.TODO()
		s, e = math.Slot(start), math.Slot(end)
	)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	// We only return early from this loop with an error if the key
	// passed in cannot be encoded.
	for i := max(s, kv.earliestSlot); i < e; i++ {
		block, err := kv.blocks.Get(ctx, i)
		if !errors.Is(err, sdkcollections.ErrNotFound) {
			// If block is found and still errors, exit and return.
			if err != nil {
				return err
			}

			// Block is found so remove from roots map.
			root := block.HashTreeRoot()
			if err = kv.roots.Remove(ctx, root[:]); err != nil {
				return err
			}

			// Block is found so also remove from execution numbers map.
			if err = kv.executionNumbers.Remove(
				ctx, block.GetExecutionNumber(),
			); err != nil {
				return err
			}
		}

		// Finally remove the block from the blocks map.
		if err = kv.blocks.Remove(ctx, i); err != nil {
			return err
		}
	}

	kv.earliestSlot = e
	return nil
}
