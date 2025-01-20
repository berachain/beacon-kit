// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package deposit

import (
	"cosmossdk.io/core/store"
	dbm "github.com/cosmos/cosmos-db"
)

var _ store.KVStoreWithBatch = &SyncedDB{}

// We have verified experimentally that deposits are often *not* flushed
// as soon as they are enqueue when pebbleDB is chosed as backend. This may
// cause an issue with ungraceful restarts, which may lead to loss of deposits,
// resulting in the node being unable to verify any incoming deposit.
// SyncedDB solves the issues since it maps the Set call to a SetSync call
// which ensure that every single deposit is flushed when enqueued.
type SyncedDB struct {
	db dbm.DB
}

func NewSynced(db dbm.DB) *SyncedDB {
	return &SyncedDB{db: db}
}

func (s *SyncedDB) Get(key []byte) ([]byte, error) {
	return s.db.Get(key)
}

func (s *SyncedDB) Has(key []byte) (bool, error) {
	return s.db.Has(key)
}

func (s *SyncedDB) Set(key, value []byte) error {
	return s.db.SetSync(key, value)
}

func (s *SyncedDB) Delete(key []byte) error {
	return s.db.Delete(key)
}

func (s *SyncedDB) Iterator(start, end []byte) (store.Iterator, error) {
	return s.db.Iterator(start, end)
}

func (s *SyncedDB) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return s.db.ReverseIterator(start, end)
}

func (s *SyncedDB) NewBatch() store.Batch {
	return s.db.NewBatch()
}

func (s *SyncedDB) NewBatchWithSize(i int) store.Batch {
	return s.db.NewBatchWithSize(i)
}

func (s *SyncedDB) Close() error {
	return s.db.Close()
}
