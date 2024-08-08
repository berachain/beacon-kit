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

package runtime

import (
	"context"
	"io"

	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewKVStoreService(storeKey *storetypes.KVStoreKey) store.KVStoreService {
	return &kvStoreService{key: storeKey}
}

type kvStoreService struct {
	key *storetypes.KVStoreKey
}

func (k kvStoreService) OpenKVStore(ctx context.Context) store.KVStore {
	return newKVStore(sdk.UnwrapSDKContext(ctx).KVStore(k.key))
}

type memStoreService struct {
	key *storetypes.MemoryStoreKey
}

func NewMemStoreService(
	storeKey *storetypes.MemoryStoreKey,
) store.MemoryStoreService {
	return &memStoreService{key: storeKey}
}

func (m memStoreService) OpenMemoryStore(ctx context.Context) store.KVStore {
	return newKVStore(sdk.UnwrapSDKContext(ctx).KVStore(m.key))
}

func NewTransientStoreService(
	storeKey *storetypes.TransientStoreKey,
) store.TransientStoreService {
	return &transientStoreService{key: storeKey}
}

type transientStoreService struct {
	key *storetypes.TransientStoreKey
}

func (t transientStoreService) OpenTransientStore(
	ctx context.Context,
) store.KVStore {
	return newKVStore(sdk.UnwrapSDKContext(ctx).KVStore(t.key))
}

type failingStoreService struct{}

func (failingStoreService) OpenKVStore(ctx context.Context) store.KVStore {
	panic(
		"kv store service not available for this module: verify runtime `skip_store_keys` app config if not expected",
	)
}

func (failingStoreService) OpenMemoryStore(ctx context.Context) store.KVStore {
	panic(
		"memory kv store service not available for this module: verify runtime `skip_store_keys` app config if not expected",
	)
}

func (failingStoreService) OpenTransientStore(
	ctx context.Context,
) store.KVStore {
	panic(
		"transient kv store service not available for this module: verify runtime `skip_store_keys` app config if not expected",
	)
}

// CoreKVStore is a wrapper of Core/Store kvstore interface
// Remove after https://github.com/cosmos/cosmos-sdk/issues/14714 is closed.
type coreKVStore struct {
	kvStore storetypes.KVStore
}

// NewKVStore returns a wrapper of Core/Store kvstore interface
// Remove once store migrates to core/store kvstore interface.
func newKVStore(store storetypes.KVStore) store.KVStore {
	return coreKVStore{kvStore: store}
}

// Get returns nil iff key doesn't exist. Errors on nil key.
func (store coreKVStore) Get(key []byte) ([]byte, error) {
	return store.kvStore.Get(key), nil
}

// Has checks if a key exists. Errors on nil key.
func (store coreKVStore) Has(key []byte) (bool, error) {
	return store.kvStore.Has(key), nil
}

// Set sets the key. Errors on nil key or value.
func (store coreKVStore) Set(key, value []byte) error {
	store.kvStore.Set(key, value)
	return nil
}

// Delete deletes the key. Errors on nil key.
func (store coreKVStore) Delete(key []byte) error {
	store.kvStore.Delete(key)
	return nil
}

// Iterator iterates over a domain of keys in ascending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
// To iterate over entire domain, use store.Iterator(nil, nil)
// CONTRACT: No writes may happen within a domain while an iterator exists over
// it.
// Exceptionally allowed for cachekv.Store, safe to write in the modules.
func (store coreKVStore) Iterator(start, end []byte) (store.Iterator, error) {
	return store.kvStore.Iterator(start, end), nil
}

// ReverseIterator iterates over a domain of keys in descending order. End is
// exclusive.
// Start must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
// CONTRACT: No writes may happen within a domain while an iterator exists over
// it.
// Exceptionally allowed for cachekv.Store, safe to write in the modules.
func (store coreKVStore) ReverseIterator(
	start, end []byte,
) (store.Iterator, error) {
	return store.kvStore.ReverseIterator(start, end), nil
}

// Adapter.
var _ storetypes.KVStore = kvStoreAdapter{}

type kvStoreAdapter struct {
	store store.KVStore
}

func (kvStoreAdapter) CacheWrap() storetypes.CacheWrap {
	panic("unimplemented")
}

func (kvStoreAdapter) CacheWrapWithTrace(
	w io.Writer,
	tc storetypes.TraceContext,
) storetypes.CacheWrap {
	panic("unimplemented")
}

func (kvStoreAdapter) GetStoreType() storetypes.StoreType {
	panic("unimplemented")
}

func (s kvStoreAdapter) Delete(key []byte) {
	err := s.store.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (s kvStoreAdapter) Get(key []byte) []byte {
	bz, err := s.store.Get(key)
	if err != nil {
		panic(err)
	}
	return bz
}

func (s kvStoreAdapter) Has(key []byte) bool {
	has, err := s.store.Has(key)
	if err != nil {
		panic(err)
	}
	return has
}

func (s kvStoreAdapter) Set(key, value []byte) {
	err := s.store.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (s kvStoreAdapter) Iterator(start, end []byte) dbm.Iterator {
	it, err := s.store.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	return it
}

func (s kvStoreAdapter) ReverseIterator(start, end []byte) dbm.Iterator {
	it, err := s.store.ReverseIterator(start, end)
	if err != nil {
		panic(err)
	}
	return it
}

func KVStoreAdapter(store store.KVStore) storetypes.KVStore {
	return &kvStoreAdapter{store}
}
