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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/storage/pkg/encoding"
)

// KVStore is a simple KV store based implementation that stores beacon blocks.
type KVStore[BeaconBlockT BeaconBlock] struct {
	blockRoots             sdkcollections.Map[[]byte, math.Slot]
	blockRootsLookup       sdkcollections.Map[math.Slot, []byte]
	executionNumbers       sdkcollections.Map[math.U64, math.Slot]
	executionNumbersLookup sdkcollections.Map[math.Slot, math.U64]
	stateRoots             sdkcollections.Map[[]byte, math.Slot]
	stateRootsLookup       sdkcollections.Map[math.Slot, []byte]

	// w is the availabilityWindow, the number of slots to keep in the store.
	w math.U64
	// l is the latest slot in the store.
	latest math.Slot

	mu     sync.RWMutex
	logger log.Logger[any]
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock](
	kvsp store.KVStoreService,
	logger log.Logger[any],
	availabilityWindow uint64,
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
		blockRootsLookup: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(blockRootsLookupPrefix),
			blockRootsLookupName,
			encoding.U64Key,
			sdkcollections.BytesValue,
		),
		executionNumbers: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(executionNumbersPrefix),
			executionNumbersName,
			encoding.U64Key,
			encoding.U64Value,
		),
		executionNumbersLookup: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(executionNumbersLookupPrefix),
			executionNumbersLookupName,
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
		stateRootsLookup: sdkcollections.NewMap(
			schemaBuilder,
			sdkcollections.NewPrefix(stateRootsLookupPrefix),
			stateRootsLookupName,
			encoding.U64Key,
			sdkcollections.BytesValue,
		),
		w:      math.U64(availabilityWindow),
		logger: logger,
	}
}

// Set sets the block by a given index in the store, storing the block root,
// execution number, and state root.
func (kv *KVStore[BeaconBlockT]) Set(blk BeaconBlockT) error {
	ctx := context.TODO()

	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Set the latest slot.
	kv.latest = blk.GetSlot()

	// Get the index of the block in the availability window.
	idx := kv.latest % kv.w

	// Set the block root in the store.
	blockRoot := blk.HashTreeRoot()
	if err := kv.blockRoots.Set(ctx, blockRoot[:], idx); err != nil {
		return err
	}
	if err := kv.blockRootsLookup.Set(ctx, idx, blockRoot[:]); err != nil {
		return err
	}

	// Set the execution number in the store.
	executionNumber := blk.GetExecutionNumber()
	if err := kv.executionNumbers.Set(ctx, executionNumber, idx); err != nil {
		return err
	}
	if err := kv.executionNumbersLookup.Set(ctx, idx, executionNumber); err != nil {
		return err
	}

	// Set the state root in the store.
	stateRoot := blk.GetStateRoot()
	if err := kv.stateRoots.Set(ctx, stateRoot[:], idx); err != nil {
		return err
	}
	if err := kv.stateRootsLookup.Set(ctx, idx, stateRoot[:]); err != nil {
		return err
	}

	return nil
}

// GetSlotByRoot retrieves the slot by a given block root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByBlockRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	// Get the index of the slot in the availability window.
	idx, err := kv.blockRoots.Get(context.TODO(), root[:])
	if err != nil {
		return 0, err
	}

	return kv.GetSlotFromWindow(idx)
}

// GetSlotByStateRoot retrieves the slot by a given state root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByStateRoot(
	root common.Root,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	// Get the index of the slot in the availability window.
	idx, err := kv.stateRoots.Get(context.TODO(), root[:])
	if err != nil {
		return 0, err
	}

	return kv.GetSlotFromWindow(idx)
}

// GetSlotByExecutionNumber retrieves the slot by a given execution number from
// the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByExecutionNumber(
	executionNumber math.U64,
) (math.Slot, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	idx, err := kv.executionNumbers.Get(context.TODO(), executionNumber)
	if err != nil {
		return 0, err
	}

	return kv.GetSlotFromWindow(idx)
}

// GetSlotFromWindow calculates the correct slot based on the index in the
// availability window and the latest slot seen.
//
// NOTE: This function is not thread safe and should be called with the mutex
// locked.
func (kv *KVStore[BeaconBlockT]) GetSlotFromWindow(
	idx math.U64,
) (math.Slot, error) {
	if idx >= kv.w {
		return 0, errors.New("index out of bounds")
	}

	slot := idx + ((kv.latest / kv.w) * kv.w)

	// If the index is past the index of the latest slot, we must
	// subtract the window size to get the correct slot.
	if idx > kv.latest%kv.w {
		slot -= kv.w
	}

	return slot, nil
}
