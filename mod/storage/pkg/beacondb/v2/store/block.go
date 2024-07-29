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

package store

import (
	"sync"

	"cosmossdk.io/core/store"
	db "github.com/cosmos/cosmos-db"
)

// BlockChanges is an extension of the changeset
type BlockChanges struct {
	*store.Changeset
}

// NewBlockChanges creates an empty blockChanges struct
func NewBlockChanges() *BlockChanges {
	return &BlockChanges{
		Changeset: store.NewChangeset(),
	}
}

// Extend extends the block changes with the given changeset
func (bc *BlockChanges) Extend(changes *store.Changeset) {
	bc.Changes = append(bc.Changes, changes.Changes...)
}

// BlockStore is a mem store for the state changes in the current block
// It should persist over the entire lifecycle of a block, and reset once
// it has been delivered
type BlockStore struct {
	blockChanges *BlockChanges
	db           *db.MemDB
	mu           sync.Mutex
}

// NewBlockStore creates a new block store.
// BlockStore is a singleton, so New should only be called once while building.
func NewBlockStore() *BlockStore {
	return &BlockStore{
		blockChanges: NewBlockChanges(),
		db:           db.NewMemDB(),
	}
}

// Add adds a change to the changeset and changes map
func (bs *BlockStore) Add(storeKey, key, value []byte, remove bool) error {
	defer bs.mu.Unlock()
	bs.mu.Lock()
	// add/remove the change to the map of changes
	if remove {
		if err := bs.db.Delete(key); err != nil {
			return err
		}
	} else {
		if err := bs.db.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Query queries the BlockStore for the given key
// return: value, found
func (bs *BlockStore) Query(storeKey, key []byte) ([]byte, bool) {
	// if not found, memdb returns value as nil
	if value, err := bs.db.Get(key); err == nil {
		return value, value != nil
	}
	return nil, false
}

// Commit adds the changes to the block changes and db
func (bs *BlockStore) Commit(changes *store.Changeset) {
	// add the changes to the mem store
	for _, change := range changes.Changes {
		for _, kvpair := range change.StateChanges {
			bs.Add(change.Actor, kvpair.Key, kvpair.Value, kvpair.Remove)
		}
	}
	// extend the slice of block changes
	bs.blockChanges.Extend(changes)
}

func (bs *BlockStore) GetChanges() *store.Changeset {
	return bs.blockChanges.Changeset
}

// Flush resets the block changes and db
func (bs *BlockStore) Flush() {
	bs.blockChanges = NewBlockChanges()
	bs.db.Close()
	bs.db = db.NewMemDB()
}

// Iterator returns an iterator over the block store memdb
func (bs *BlockStore) Iterator(start, end []byte) (store.Iterator, error) {
	return bs.db.Iterator(start, end)
}
