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

package beacondb

import (
	"fmt"

	"cosmossdk.io/core/store"
)

// Changeset is a wrapper around store.Changeset that holds a map of changes
// for more efficient querying
// INVARIANT: changes map and Changeset are always in sync
type Changeset struct {
	*store.Changeset
	changes map[string][]byte
}

// NewChangeset initializes a new Changeset with an empty store.Changeset and
// changes map.
func NewChangeset() *Changeset {
	return &Changeset{
		Changeset: store.NewChangeset(),
		changes:   make(map[string][]byte),
	}
}

// NewChangesetWithPairs creates a new changeset with the given pairs.
func NewChangesetWithPairs(pairs map[string]store.KVPairs) *Changeset {
	cs := &Changeset{
		Changeset: store.NewChangesetWithPairs(pairs),
		changes:   make(map[string][]byte),
	}
	for storeKey, kvPairs := range pairs {
		for _, pair := range kvPairs {
			cs.changes[buildPath([]byte(storeKey), pair.Key)] = pair.Value
		}
	}
	return cs
}

// Add adds a change to the changeset and changes map
func (cs *Changeset) Add(storeKey, key, value []byte, remove bool) {
	keyPath := buildPath(storeKey, key)
	// add/remove the change to the map of changes
	if remove {
		cs.changes[keyPath] = nil
	} else {
		cs.changes[keyPath] = value
	}
	fmt.Println("ADDING CHANGE TO CHANGESET", keyPath, value, remove)
	cs.Changeset.Add(storeKey, key, value, remove)
}

// AddKVPair adds a KVPair to the Changeset and changes map
func (cs *Changeset) AddKVPair(storeKey []byte, pair store.KVPair) {
	cs.Add(storeKey, pair.Key, pair.Value, pair.Remove)
}

// Query queries the changeset with the given store key and key
func (cs *Changeset) Query(storeKey []byte, key []byte) ([]byte, bool) {
	keyPath := buildPath(storeKey, key)
	fmt.Println("QUERYING CHANGESET WITH ", keyPath)
	fmt.Println("KEY", key)
	if value, found := cs.changes[keyPath]; found {
		fmt.Println("FUCKING FOUND")
		return value, true
	}
	return nil, false
}

func buildPath(storeKey []byte, key []byte) string {
	return string(append(storeKey, key...))
}
