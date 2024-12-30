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

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	db "github.com/berachain/beacon-kit/storage/interfaces"
	"github.com/berachain/beacon-kit/storage/pruner"
)

const (
	keyFormat  = "%d/%s"
	pathFormat = "%d/"
	keyParts   = 2
)

// Compile-time assertion of prunable interface.
var (
	_ pruner.Prunable = (*RangeDB)(nil)

	ErrRangeNotSupported = errors.New("RangeDB DeleteRange: delete range not supported for this db")
)

// RangeDB is a database that stores versioned data.
// It prefixes keys with an index.
// Invariant: No index below firstNonNilIndex should be populated.
type RangeDB struct {
	coreDB *DB

	// lowerBoundIndex is used as a loose check for stored indexes
	// monotonicity. The goal is to make sure we do not overwrite
	// indexes which have been or will be deleted eventually via pruning.
	lowerBoundIndex uint64
}

// NewRangeDB creates a new RangeDB.
func NewRangeDB(coreDB db.DB) *RangeDB {
	cDB, ok := coreDB.(*DB)
	if !ok {
		panic(ErrRangeNotSupported)
	}
	return &RangeDB{
		coreDB:          cDB,
		lowerBoundIndex: 0,
	}
}

// Get retrieves the value associated with the given index and key.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB) Get(index uint64, key []byte) ([]byte, error) {
	return db.coreDB.Get(prefix(index, key))
}

// Has checks if the given index and key exist in the database.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB) Has(index uint64, key []byte) (bool, error) {
	return db.coreDB.Has(prefix(index, key))
}

// Set stores the value with the given index and key in the database.
// It prefixes the key with the index and a slash before storing it in the
// underlying database.
func (db *RangeDB) Set(index uint64, key []byte, value []byte) error {
	index = max(index, db.lowerBoundIndex) // enforce invariant
	return db.coreDB.Set(prefix(index, key), value)
}

// Delete removes the value associated with the given index and key from the
// database. It prefixes the key with the index and a slash before deleting it
// from the underlying database.
func (db *RangeDB) Delete(index uint64, key []byte) error {
	return db.coreDB.Delete(prefix(index, key))
}

// DeleteRange removes all values associated with the given index from the
// filesystem. It is INCLUSIVE of the `from` index and EXCLUSIVE of
// the `to“ index.
func (db *RangeDB) DeleteRange(from, to uint64) error {
	if from > to {
		return fmt.Errorf(
			"RangeDB DeleteRange start: %d, end: %d: %w",
			from, to, pruner.ErrInvalidRange,
		)
	}
	for ; from < to; from++ {
		path := fmt.Sprintf(pathFormat, from)
		if err := db.coreDB.fs.RemoveAll(path); err != nil {
			return fmt.Errorf(
				"RangeDB DeleteRange start: %d, end: %d, failed RemoveAll: %w",
				from, to, err,
			)
		}
	}
	return nil
}

// Prune removes all values in the given range [start, end) from the db.
func (db *RangeDB) Prune(start, end uint64) error {
	start = max(start, db.lowerBoundIndex)
	if start > end {
		return fmt.Errorf(
			"RangeDB Prune start: %d, end: %d: %w",
			start, end, pruner.ErrInvalidRange,
		)
	}

	// DeleteRange may fail and so some files to be pruned may have not
	// been removed. We *do not* retry to prune those files to avoid getting
	// stuck with them. Instead we update lowerBoundIndex as if deletion
	// was successful and we return an error.
	err := db.DeleteRange(start, end)
	db.lowerBoundIndex = end
	return err
}

// prefix prefixes the given key with the index and a slash.
func prefix(index uint64, key []byte) []byte {
	return []byte(fmt.Sprintf(keyFormat, index, hex.EncodeBytes(key)))
}

// ExtractIndex extracts the index from a prefixed key.
func ExtractIndex(prefixedKey []byte) (uint64, error) {
	parts := bytes.SplitN(prefixedKey, []byte("/"), keyParts)
	if len(parts) < keyParts {
		return 0, errors.New("invalid key format")
	}

	indexStr := string(parts[0])
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return index, nil
}
