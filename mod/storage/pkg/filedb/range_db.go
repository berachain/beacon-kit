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

package filedb

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	db "github.com/berachain/beacon-kit/mod/storage/pkg/interfaces"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
)

// two is a constant for the number 2.
const two = 2

// Compile-time assertion of prunable interface.
var _ pruner.Prunable = (*RangeDB)(nil)

// RangeDB is a database that stores versioned data.
// It prefixes keys with an index.
// Invariant: No index below firstNonNilIndex should be populated.
type RangeDB struct {
	db.DB
	firstNonNilIndex uint64
}

// NewRangeDB creates a new RangeDB.
func NewRangeDB(db db.DB) *RangeDB {
	return &RangeDB{
		DB:               db,
		firstNonNilIndex: 0,
	}
}

// Get retrieves the value associated with the given index and key.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB) Get(index uint64, key []byte) ([]byte, error) {
	return db.DB.Get(db.prefix(index, key))
}

// Has checks if the given index and key exist in the database.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB) Has(index uint64, key []byte) (bool, error) {
	return db.DB.Has(db.prefix(index, key))
}

// Set stores the value with the given index and key in the database.
// It prefixes the key with the index and a slash before storing it in the
// underlying database.
func (db *RangeDB) Set(index uint64, key []byte, value []byte) error {
	// enforce invariant
	if index < db.firstNonNilIndex {
		db.firstNonNilIndex = index
	}
	return db.DB.Set(db.prefix(index, key), value)
}

// Delete removes the value associated with the given index and key from the
// database. It prefixes the key with the index and a slash before deleting it
// from the underlying database.
func (db *RangeDB) Delete(index uint64, key []byte) error {
	return db.DB.Delete(db.prefix(index, key))
}

// DeleteRange removes all values associated with the given index from the
// filesystem. It is INCLUSIVE of the `from` index and EXCLUSIVE of
// the `to“ index.
func (db *RangeDB) DeleteRange(from, to uint64) error {
	f, ok := db.DB.(*DB)
	if !ok {
		return errors.New("rangedb: delete range not supported for this db")
	}
	for ; from < to; from++ {
		if err := f.fs.RemoveAll(fmt.Sprintf("%d/", from)); err != nil {
			return err
		}
	}
	return nil
}

// Prune removes all values in the given range [start, end) from the db.
func (db *RangeDB) Prune(start, end uint64) error {
	start = max(start, db.firstNonNilIndex)
	if err := db.DeleteRange(start, end); err != nil {
		// Resets last pruned index in case Delete somehow populates indices on
		// err. This will cause the next prune operation is O(n), but next
		// successful prune will set it to the correct value, so runtime is
		// ammortized
		db.firstNonNilIndex = 0
		return err
	}
	db.firstNonNilIndex = end
	return nil
}

// prefix prefixes the given key with the index and a slash.
func (db *RangeDB) prefix(index uint64, key []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s", index, hex.FromBytes(key).Unwrap()))
}

// ExtractIndex extracts the index from a prefixed key.
func ExtractIndex(prefixedKey []byte) (uint64, error) {
	parts := bytes.SplitN(prefixedKey, []byte("/"), two)
	if len(parts) < two {
		return 0, errors.New("invalid key format")
	}

	indexStr := string(parts[0])
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		return 0, errors.Newf("invalid index: %w", err)
	}

	//#nosec:g
	return index, nil
}
