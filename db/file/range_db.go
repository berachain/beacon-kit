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

package file

import (
	"fmt"
)

// numeric is a type that represents a numeric type.
type numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// RangeDB is a database that stores versioned data.
// It prefixes keys with an index.
type RangeDB[T numeric] struct {
	*DB
}

// NewRangeDB creates a new RangeDB.
func NewRangeDB[T numeric](db *DB) *RangeDB[T] {
	return &RangeDB[T]{
		DB: db,
	}
}

// Get retrieves the value associated with the given index and key.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB[T]) Get(index T, key []byte) ([]byte, error) {
	return db.DB.Get(db.prefix(index, key))
}

// Has checks if the given index and key exist in the database.
// It prefixes the key with the index and a slash before querying the underlying
// database.
func (db *RangeDB[T]) Has(index T, key []byte) (bool, error) {
	return db.DB.Has(db.prefix(index, key))
}

// Set stores the value with the given index and key in the database.
// It prefixes the key with the index and a slash before storing it in the
// underlying database.
func (db *RangeDB[T]) Set(index T, key []byte, value []byte) error {
	return db.DB.Set(db.prefix(index, key), value)
}

// Delete removes the value associated with the given index and key from the
// database. It prefixes the key with the index and a slash before deleting it
// from the underlying database.
func (db *RangeDB[T]) Delete(index T, key []byte) error {
	return db.DB.Delete(db.prefix(index, key))
}

// DeleteRange removes all values associated with the given index from the
// filesystem. It is INCLUSIVE of the `from` index and EXCLUSIVE of
// the `to“ index.
func (db *RangeDB[T]) DeleteRange(from, to T) error {
	for ; from < to; from++ {
		err := db.fs.RemoveAll(string(db.prefix(from, nil)))
		if err != nil {
			return err
		}
	}
	return nil
}

// prefix prefixes the given key with the index and a slash.
func (db *RangeDB[T]) prefix(index T, key []byte) []byte {
	return append([]byte(fmt.Sprintf("%d/", index)), key...)
}
