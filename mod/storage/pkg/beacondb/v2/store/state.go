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
	"context"

	sdkcollections "cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/store"
	storev2 "cosmossdk.io/store/v2"
	"github.com/berachain/beacon-kit/mod/errors"
	storectx "github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/context"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/v2/iterator"
	"github.com/berachain/beacon-kit/mod/storage/pkg/collections"
)

// TODO: pull this from some cfg or something idk remind me to do it after
// im done getting shafted by server v2
const (
	ModuleName = "beacon"
)

// StateStore is a store for the state of the beacon
type StateStore struct {
	ctx             *storectx.Context
	emphemeralState *EphemeralStore
	state           storev2.RootStore
}

// NewStore creates a new state store
func NewStore(
	ephemeralStore *EphemeralStore,
	state storev2.RootStore,
) *StateStore {
	return &StateStore{
		emphemeralState: ephemeralStore,
		state:           state,
	}
}

// AddChange adds a change to the changeset
func (s *StateStore) AddChange(storeKey []byte, key []byte, value []byte) {
	s.ctx.Changeset.Add(storeKey, key, value, false)
}

// TODO: CHANGE NAME TO QUERY
// QueryState queries the state store for the given key
// It first queries the changeset, then the block store, then the chain store
func (s *StateStore) QueryState(storeKey, key []byte) ([]byte, error) {
	// QUERY CHANGESET
	value, found := s.ctx.Changeset.Query(storeKey, key)
	if found {
		return value, nil
	}

	// NOT FOUND IN CHANGESET -> QUERY BLOCK STORE
	value, found = s.emphemeralState.Query(storeKey, key)
	if found {
		return value, nil
	}

	// NOT FOUND IN BLOCK STORE -> QUERY COMMIT STORE
	version, err := s.state.GetLatestVersion()
	if err != nil {
		return nil, err
	}
	// if the version is 0, we're in genesis
	if version != 0 {
		var resp storev2.QueryResult
		resp, err = s.state.Query(storeKey, version, key, false)
		if err != nil {
			// TODO: clean this up
			if errors.Is(err, sdkcollections.ErrNotFound) {
				return nil, collections.ErrNotFound
			}
			return nil, err
		}
		if resp.Value != nil {
			return resp.Value, nil
		}
	}
	return nil, collections.ErrNotFound
}

// Iterator returns a combined iterator over the the changeset, the block store,
// and the chain store
func (s *StateStore) Iterator(start, end []byte) (store.Iterator, error) {
	// get chain store iterator
	_, readerMap, err := s.state.StateLatest()
	if err != nil {
		return nil, err
	}
	// get reader with the storeKey from reader map
	reader, err := readerMap.GetReader([]byte(ModuleName))
	if err != nil {
		return nil, err
	}
	// get iterator from reader with the prefixed start and end
	stateIter, err := reader.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	// get chainset iterator
	changeSetIter, err := s.ctx.Changeset.Iterator(start, end)
	if err != nil {
		return nil, err
	}

	// get block store iterator
	blockStoreIter, err := s.emphemeralState.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return iterator.New(start, end, stateIter, changeSetIter, blockStoreIter), nil
}

func (s *StateStore) WorkingHash() ([]byte, error) {
	if s.ctx.Changeset.Size() > 0 {
		return s.state.WorkingHash(s.ctx.Changeset.GetChanges())
	}
	return s.state.WorkingHash(s.emphemeralState.GetChanges())
}

func (s *StateStore) Commit() ([]byte, error) {
	hash, err := s.state.Commit(s.ctx.Changeset.GetChanges())
	return hash, err
}

// Save commits the changeset to the BlockStore and resets the changeset
func (s *StateStore) Save() {
	// reset the changeset following the commit
	defer func() {
		s.ctx.Changeset.Flush()
	}()
	size := s.ctx.Changeset.Size()
	if size == 0 {
		return
	}
	// commit changes to block store
	s.emphemeralState.Commit(s.ctx.Changeset.GetChanges())
}

// Context returns the context of the StateStore
func (s *StateStore) Context() *storectx.Context {
	return s.ctx
}

// WithContext returns a new StateStore with the given context
// Invariant: the blockStore in the retuned copy must always point to the same
// blockStore as the current state
func (s *StateStore) WithContext(ctx context.Context) *StateStore {
	storeCtx, ok := ctx.(*storectx.Context)
	if !ok {
		storeCtx = storectx.New(ctx)
	}
	// initialize a new StateStore with the current block store and chain store
	storeCopy := NewStore(s.emphemeralState, s.state)
	// set the context to the given context
	storeCopy.ctx = storeCtx
	return storeCopy
}

// TODO V2: i'll remove this once i confirm whatever sdfjklhsdfjkl
// TODO: hack to get around IAVL versioning issues (a call to save must be made
// before any queries to SC)
func (s *StateStore) Init() {
	s.ctx = storectx.New(context.Background())
	// if err := s.state.SetInitialVersion(uint64(1)); err != nil {
	// 	panic(err)
	// }
	// s.AddChange([]byte("beacon"), []byte("4206969666"), []byte("69"))
	// s.Save()
}

// TODO: remove this when no needed
// StateLatest returns the latest state from the chain store
func (s *StateStore) StateLatest() (uint64, corestore.ReaderMap, error) {
	return s.state.StateLatest()
}
