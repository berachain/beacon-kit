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

package prunedb

import "github.com/berachain/beacon-kit/mod/storage/pkg/filedb"

// DB is a wrapper around filedb.RangeDB that keeps track of the latest index.
type DB struct {
	*filedb.RangeDB
	latestIndex uint64
}

// New creates a new DB.
func New(db *filedb.RangeDB) *DB {
	return &DB{
		RangeDB: db,
	}
}

// GetLatestIndex returns the latest index.
func (p *DB) GetLatestIndex() uint64 {
	return p.latestIndex
}

// Set sets the key and value at the given index and updates the latest index.
func (p *DB) Set(index uint64, key []byte, value []byte) error {
	if err := p.RangeDB.Set(index, key, value); err != nil {
		return err
	}

	p.latestIndex = index
	return nil
}
