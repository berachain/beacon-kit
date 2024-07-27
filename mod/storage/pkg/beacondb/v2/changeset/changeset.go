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

package changeset

import (
	"sync"

	"cosmossdk.io/core/store"
	db "github.com/cosmos/cosmos-db"
)

// Changeset is a wrapper around store.Changeset that holds a map of changes
// for more efficient querying
// INVARIANT: changes map and Changeset are always in sync
type Changeset struct {
	*store.Changeset
	*db.MemDB

	mu *sync.RWMutex
}

// New initializes a new Changeset with an empty store.Changeset and
// changes map.
func New() *Changeset {
	return &Changeset{
		Changeset: store.NewChangeset(),
		MemDB:     db.NewMemDB(),
		mu:        &sync.RWMutex{},
	}
}

func (cs *Changeset) GetChanges() *store.Changeset {
	return cs.Changeset
}

// NewWithPairs creates a new changeset with the given pairs.
func NewWithPairs(pairs map[string]store.KVPairs) *Changeset {
	cs := &Changeset{
		Changeset: store.NewChangesetWithPairs(pairs),
		MemDB:     db.NewMemDB(),
		mu:        &sync.RWMutex{},
	}
	for _, kvPairs := range pairs {
		for _, pair := range kvPairs {
			if err := cs.Set(
				pair.Key,
				pair.Value,
			); err != nil {
				panic(err)
			}
		}
	}
	return cs
}

// Add adds a change to the changeset and changes map
func (cs *Changeset) Add(storeKey, key, value []byte, remove bool) error {
	defer cs.mu.Unlock()
	cs.mu.Lock()
	// add/remove the change to the map of changes
	if remove {
		if err := cs.MemDB.Delete(key); err != nil {
			return err
		}
	} else {
		if err := cs.MemDB.Set(key, value); err != nil {
			return err
		}
	}
	cs.Changeset.Add(storeKey, key, value, remove)
	return nil
}

// AddKVPair adds a KVPair to the Changeset and changes map
func (cs *Changeset) AddKVPair(storeKey []byte, pair store.KVPair) {
	cs.Add(storeKey, pair.Key, pair.Value, pair.Remove)
}

// Query queries the changeset with the given store key and key
func (cs *Changeset) Query(storeKey, key []byte) ([]byte, bool) {
	// Note: MemDB returns no error but value is nil if key is not found,
	// so we need to check if value is nil
	if value, err := cs.MemDB.Get(key); err == nil {
		return value, value != nil
	}
	return nil, false
}

// Flush resets the changeset and changes map.
func (cs *Changeset) Flush() {
	cs.Changeset = store.NewChangeset()
	cs.MemDB.Close()
	cs.MemDB = db.NewMemDB()
}

func (cs *Changeset) Copy() *Changeset {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	csCopy := New()
	for _, change := range cs.Changeset.Changes {
		for _, kvpair := range change.StateChanges {
			csCopy.Add(change.Actor, kvpair.Key, kvpair.Value, kvpair.Remove)
		}
	}

	return csCopy
}
