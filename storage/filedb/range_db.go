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

package filedb

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
	db "github.com/berachain/beacon-kit/storage/interfaces"
	"github.com/berachain/beacon-kit/storage/pruner"
	"github.com/spf13/afero"
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

	rwMu sync.RWMutex

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
	db.rwMu.RLock()
	defer db.rwMu.RUnlock()
	return db.coreDB.Get(prefix(index, key))
}

// Has checks if the given index and key exist in the database.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB) Has(index uint64, key []byte) (bool, error) {
	db.rwMu.RLock()
	defer db.rwMu.RUnlock()
	return db.coreDB.Has(prefix(index, key))
}

// Set stores the value with the given index and key in the database.
// It prefixes the key with the index and a slash before storing it in the
// underlying database.
func (db *RangeDB) Set(index uint64, key []byte, value []byte) error {
	db.rwMu.Lock()
	defer db.rwMu.Unlock()

	index = max(index, db.lowerBoundIndex) // enforce invariant
	return db.coreDB.Set(prefix(index, key), value)
}

// Delete removes the value associated with the given index and key from the
// database. It prefixes the key with the index and a slash before deleting it
// from the underlying database.
func (db *RangeDB) Delete(index uint64, key []byte) error {
	db.rwMu.Lock()
	defer db.rwMu.Unlock()
	return db.coreDB.Delete(prefix(index, key))
}

// deleteRange removes all values associated with the given index from the
// filesystem. It is INCLUSIVE of the `from` index and EXCLUSIVE of
// the `to“ index.
func (db *RangeDB) deleteRange(from, to uint64) error {
	if from > to {
		return fmt.Errorf(
			"RangeDB deleteRange start: %d, end: %d: %w",
			from, to, pruner.ErrInvalidRange,
		)
	}
	for i := from; i < to; i++ {
		path := fmt.Sprintf(pathFormat, i)
		if err := db.coreDB.fs.RemoveAll(path); err != nil {
			return fmt.Errorf(
				"RangeDB DeleteRange failed RemoveAll index %d: %w",
				i, err,
			)
		}
	}
	return nil
}

// Prune removes all values in the given range [start, end) from the db.
func (db *RangeDB) Prune(start, end uint64) error {
	db.rwMu.Lock()
	defer db.rwMu.Unlock()
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
	err := db.deleteRange(start, end)
	db.lowerBoundIndex = end
	return err
}

// GetByIndex takes the database index and returns all associated entries,
// expecting database keys to follow the prefix() format. If index does not
// exist in the DB for any reason (pruned, invalid index), an empty list is
// returned with no error.
func (db *RangeDB) GetByIndex(index uint64) ([][]byte, error) {
	db.rwMu.RLock()
	defer db.rwMu.RUnlock()
	indexDir := fmt.Sprintf(pathFormat, index)
	entries, err := afero.ReadDir(db.coreDB.fs, indexDir)
	if err != nil {
		if os.IsNotExist(err) {
			return [][]byte{}, nil
		}
		return nil, err
	}
	keys := make([][]byte, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if !strings.HasSuffix(filename, db.coreDB.extension) {
			continue
		}
		var sidecarBz []byte
		sidecarBz, err = afero.ReadFile(db.coreDB.fs, filepath.Join(indexDir, filename))
		if err != nil {
			return keys, err
		}
		keys = append(keys, sidecarBz)
	}
	return keys, nil
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
	index, err := math.U64FromString(indexStr)
	if err != nil {
		return 0, err
	}

	return index.Unwrap(), nil
}
