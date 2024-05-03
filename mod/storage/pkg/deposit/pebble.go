// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package deposit

import (
	"bytes"
	"path/filepath"
	"slices"

	corestore "cosmossdk.io/core/store"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/cockroachdb/pebble"
)

const DBFileSuffix = ".db"
const three = 3

// PebbleDB implements RawDB using PebbleDB as the underlying storage engine.
// It is used for only store v2 migration, since some clients use PebbleDB as
// the IAVL v0/v1 backend.
type PebbleDB struct {
	storage *pebble.DB
}

func NewPebbleDB(name, dataDir string) (*PebbleDB, error) {
	return NewPebbleDBWithOpts(name, dataDir)
}

//nolint:gomnd // yoinked from cosmos
func NewPebbleDBWithOpts(name, dataDir string) (*PebbleDB, error) {
	do := &pebble.Options{
		MaxConcurrentCompactions: func() int { return three }, // default 1
	}

	do.EnsureDefaults()

	dbPath := filepath.Join(dataDir, name+DBFileSuffix)
	db, err := pebble.Open(dbPath, do)
	if err != nil {
		return nil, errors.Newf("failed to open PebbleDB: %w", err)
	}

	return &PebbleDB{storage: db}, nil
}

func (db *PebbleDB) Close() error {
	err := db.storage.Close()
	db.storage = nil
	return err
}

func (db *PebbleDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key cannot be empty")
	}

	bz, closer, err := db.storage.Get(key)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			// in case of a fresh database
			return nil, nil
		}

		return nil, errors.Newf("failed to perform PebbleDB read: %w", err)
	}

	if len(bz) == 0 {
		return nil, closer.Close()
	}

	return bz, closer.Close()
}

func (db *PebbleDB) Has(key []byte) (bool, error) {
	bz, err := db.Get(key)
	if err != nil {
		return false, err
	}

	return bz != nil, nil
}

func (db *PebbleDB) Delete(key []byte) error {
	err := db.storage.Delete(key, &pebble.WriteOptions{Sync: false})
	if err != nil {
		return errors.Newf("failed to delete key from PebbleDB: %w", err)
	}
	return nil
}

func (db *PebbleDB) Set(key, value []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}
	if value == nil {
		return errors.New("value cannot be nil")
	}

	err := db.storage.Set(key, value, &pebble.WriteOptions{Sync: false})
	if err != nil {
		return errors.Newf("failed to set key-value pair in PebbleDB: %w", err)
	}
	return nil
}

func (db *PebbleDB) Iterator(start, end []byte) (corestore.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, errors.New("key cannot be empty")
	}

	itr, err := db.storage.NewIter(
		&pebble.IterOptions{LowerBound: start, UpperBound: end},
	)
	if err != nil {
		return nil, errors.Newf("failed to create PebbleDB iterator: %w", err)
	}

	return newPebbleDBIterator(itr, start, end, false), nil
}

func (db *PebbleDB) ReverseIterator(
	start, end []byte,
) (corestore.Iterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, errors.New("key cannot be empty")
	}

	itr, err := db.storage.NewIter(
		&pebble.IterOptions{LowerBound: start, UpperBound: end},
	)
	if err != nil {
		return nil, errors.Newf("failed to create PebbleDB iterator: %w", err)
	}

	return newPebbleDBIterator(itr, start, end, true), nil
}

func (db *PebbleDB) NewBatch() RawBatch {
	return &pebbleDBBatch{
		db:    db,
		batch: db.storage.NewBatch(),
	}
}

func (db *PebbleDB) NewBatchWithSize(size int) RawBatch {
	return &pebbleDBBatch{
		db:    db,
		batch: db.storage.NewBatchWithSize(size),
	}
}

var _ corestore.Iterator = (*pebbleDBIterator)(nil)

type pebbleDBIterator struct {
	source  *pebble.Iterator
	start   []byte
	end     []byte
	valid   bool
	reverse bool
}

func newPebbleDBIterator(
	src *pebble.Iterator,
	start, end []byte,
	reverse bool,
) *pebbleDBIterator {
	// move the underlying PebbleDB cursor to the first key
	var valid bool
	if reverse {
		if end == nil {
			valid = src.Last()
		} else {
			valid = src.SeekLT(end)
		}
	} else {
		valid = src.First()
	}

	return &pebbleDBIterator{
		source:  src,
		start:   start,
		end:     end,
		valid:   valid,
		reverse: reverse,
	}
}

func (itr *pebbleDBIterator) Domain() ([]byte, []byte) {
	return itr.start, itr.end
}

func (itr *pebbleDBIterator) Valid() bool {
	// once invalid, forever invalid
	if !itr.valid || !itr.source.Valid() {
		itr.valid = false
		return itr.valid
	}

	// if source has error, consider it invalid
	if err := itr.source.Error(); err != nil {
		itr.valid = false
		return itr.valid
	}

	// if key is at the end or past it, consider it invalid
	if end := itr.end; end != nil {
		if bytes.Compare(end, itr.Key()) <= 0 {
			itr.valid = false
			return itr.valid
		}
	}

	return true
}

func (itr *pebbleDBIterator) Key() []byte {
	itr.assertIsValid()
	return slices.Clone(itr.source.Key())
}

func (itr *pebbleDBIterator) Value() []byte {
	itr.assertIsValid()
	return slices.Clone(itr.source.Value())
}

func (itr *pebbleDBIterator) Next() {
	itr.assertIsValid()

	if itr.reverse {
		itr.valid = itr.source.Prev()
	} else {
		itr.valid = itr.source.Next()
	}
}

func (itr *pebbleDBIterator) Error() error {
	return itr.source.Error()
}

func (itr *pebbleDBIterator) Close() error {
	err := itr.source.Close()
	itr.source = nil
	itr.valid = false

	return err
}

func (itr *pebbleDBIterator) assertIsValid() {
	if !itr.valid {
		panic("pebbleDB iterator is invalid")
	}
}

var _ RawBatch = (*pebbleDBBatch)(nil)

type pebbleDBBatch struct {
	db    *PebbleDB
	batch *pebble.Batch
}

func (b *pebbleDBBatch) Set(key, value []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}
	if value == nil {
		return errors.New("key cannot be nil")
	}
	if b.batch == nil {
		return errors.New("batch closed")
	}

	return b.batch.Set(key, value, nil)
}

func (b *pebbleDBBatch) Delete(key []byte) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}
	if b.batch == nil {
		return errors.New("batch closed")
	}

	return b.batch.Delete(key, nil)
}

func (b *pebbleDBBatch) Write() error {
	err := b.batch.Commit(&pebble.WriteOptions{Sync: false})
	if err != nil {
		return errors.Newf("failed to write PebbleDB batch: %w", err)
	}

	return nil
}

func (b *pebbleDBBatch) WriteSync() error {
	err := b.batch.Commit(&pebble.WriteOptions{Sync: true})
	if err != nil {
		return errors.Newf("failed to write PebbleDB batch: %w", err)
	}

	return nil
}

func (b *pebbleDBBatch) Close() error {
	return b.batch.Close()
}

func (b *pebbleDBBatch) GetByteSize() (int, error) {
	return b.batch.Len(), nil
}
