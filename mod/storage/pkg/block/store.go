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

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	storage "github.com/berachain/beacon-kit/mod/storage/pkg"
	lru "github.com/hashicorp/golang-lru/v2"
)

// KVStore is a simple memory store based implementation that stores metadata of
// beacon blocks.
type KVStore[BeaconBlockT BeaconBlock] struct {
	blockRoots       *lru.Cache[common.Root, math.Slot]
	executionNumbers *lru.Cache[math.U64, math.Slot]
	stateRoots       *lru.Cache[common.Root, math.Slot]

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
	executionNumbers, err := lru.New[math.U64, math.Slot](availabilityWindow)
	if err != nil {
		panic(err)
	}
	stateRoots, err := lru.New[common.Root, math.Slot](availabilityWindow)
	if err != nil {
		panic(err)
	}
	return &KVStore[BeaconBlockT]{
		blockRoots:       blockRoots,
		executionNumbers: executionNumbers,
		stateRoots:       stateRoots,
		logger:           logger,
	}
}

// Set sets the block by a given index in the store, storing the block root,
// execution number, and state root. Only this function may potentially evict
// entries from the store if the availability window is reached.
func (kv *KVStore[BeaconBlockT]) Set(blk BeaconBlockT) error {
	slot := blk.GetSlot()
	kv.blockRoots.Add(blk.HashTreeRoot(), slot)
	kv.executionNumbers.Add(blk.GetExecutionNumber(), slot)
	kv.stateRoots.Add(blk.GetStateRoot(), slot)
	return nil
}

// GetSlotByRoot retrieves the slot by a given block root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByBlockRoot(
	blockRoot common.Root,
) (math.Slot, error) {
	slot, ok := kv.blockRoots.Peek(blockRoot)
	if !ok {
		return 0, fmt.Errorf(
			"%w, block root: %s",
			storage.ErrNotFound,
			blockRoot,
		)
	}
	return slot, nil
}

// GetSlotByExecutionNumber retrieves the slot by a given execution number from
// the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByExecutionNumber(
	executionNumber math.U64,
) (math.Slot, error) {
	slot, ok := kv.executionNumbers.Peek(executionNumber)
	if !ok {
		return 0, fmt.Errorf(
			"%w, execution number: %d",
			storage.ErrNotFound,
			executionNumber,
		)
	}
	return slot, nil
}

// GetSlotByStateRoot retrieves the slot by a given state root from the store.
func (kv *KVStore[BeaconBlockT]) GetSlotByStateRoot(
	stateRoot common.Root,
) (math.Slot, error) {
	slot, ok := kv.stateRoots.Peek(stateRoot)
	if !ok {
		return 0, fmt.Errorf(
			"%w, state root: %s",
			storage.ErrNotFound,
			stateRoot,
		)
	}
	return slot, nil
}
