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
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	lru "github.com/hashicorp/golang-lru/v2"
)

// KVStore is a simple memory store based implementation that stores metadata of
// beacon blocks.
type KVStore[BeaconBlockT BeaconBlock] struct {
	// Beacon block root to slot mapping is injective for finalized blocks.
	blockRoots *lru.Cache[common.Root, math.Slot]

	// Timestamp to slot mapping is injective for finalized blocks. This is
	// guaranteed by CometBFT consensus. So each slot will be associated with a
	// different timestamp (no overwriting) as we store only finalized blocks.
	timestamps *lru.Cache[math.U64, math.Slot]

	// Beacon state root to slot mapping is injective for finalized blocks.
	stateRoots *lru.Cache[common.Root, math.Slot]

	// Logger for the store.
	logger log.Logger
}

// NewStore creates a new block store.
func NewStore[BeaconBlockT BeaconBlock](
	logger log.Logger,
	availabilityWindow int,
) *KVStore[BeaconBlockT] {
	blockRoots, err := lru.New[common.Root, math.Slot](availabilityWindow)
	if err != nil {
		panic(err)
	}
	timestamps, err := lru.New[math.U64, math.Slot](availabilityWindow)
	if err != nil {
		panic(err)
	}
	stateRoots, err := lru.New[common.Root, math.Slot](availabilityWindow)
	if err != nil {
		panic(err)
	}
	return &KVStore[BeaconBlockT]{
		blockRoots: blockRoots,
		timestamps: timestamps,
		stateRoots: stateRoots,
		logger:     logger,
	}
}

// Set sets the block by a given index in the store, storing the block root,
// timestamp, and state root. Only this function may potentially evict
// entries from the store if the availability window is reached.
func (kv *KVStore[BeaconBlockT]) Set(blk BeaconBlockT) error {
	slot := blk.GetSlot()
	kv.blockRoots.Add(blk.HashTreeRoot(), slot)
	kv.timestamps.Add(blk.GetTimestamp(), slot)
	kv.stateRoots.Add(blk.GetStateRoot(), slot)
	return nil
}

// GetSlotByRoot retrieves the slot by a given block root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByBlockRoot(
	blockRoot common.Root,
) (math.Slot, error) {
	slot, ok := kv.blockRoots.Peek(blockRoot)
	if !ok {
		return 0, fmt.Errorf("slot not found at block root: %s", blockRoot)
	}
	return slot, nil
}

// GetParentSlotByTimestamp retrieves the parent slot by a given timestamp from
// the store.
func (kv *KVStore[BeaconBlockT]) GetParentSlotByTimestamp(
	timestamp math.U64,
) (math.Slot, error) {
	slot, ok := kv.timestamps.Peek(timestamp)
	if !ok {
		return slot, fmt.Errorf("slot not found at timestamp: %d", timestamp)
	}
	if slot == 0 {
		return slot, errors.New("parent slot not supported for genesis slot 0")
	}

	return slot - 1, nil
}

// GetSlotByStateRoot retrieves the slot by a given state root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByStateRoot(
	stateRoot common.Root,
) (math.Slot, error) {
	slot, ok := kv.stateRoots.Peek(stateRoot)
	if !ok {
		return 0, fmt.Errorf("slot not found at state root: %s", stateRoot)
	}
	return slot, nil
}
