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

import (
	"context"
	"time"
)

type IndexDB interface {
	Get(index uint64, key []byte) ([]byte, error)
	Has(index uint64, key []byte) (bool, error)
	Set(index uint64, key []byte, value []byte) error
	Delete(index uint64, key []byte) error
	DeleteRange(start, end uint64) error
}

// DB is a wrapper around filedb.RangeDB that keeps track of the latest index.
type DB struct {
	IndexDB
	ticker          *time.Ticker
	windowSize      uint64
	highestSetIndex uint64

	lastDeletedIndex uint64
}

// New creates a new DB.
func New(
	db IndexDB,
	pruneInterval time.Duration,
	windowSize uint64,
) *DB {
	return &DB{
		windowSize:       windowSize,
		IndexDB:          db,
		ticker:           time.NewTicker(pruneInterval),
		lastDeletedIndex: 0,
	}
}

func (db *DB) Start(ctx context.Context) {
	go func() {
		defer db.ticker.Stop()

		for {
			select {
			case <-db.ticker.C:
				// Do the pruning
				if err := db.prune(); err != nil {
					// db.Logger().Error("Error pruning: ", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Set sets the key and value at the given index and updates the latest index.
func (p *DB) Set(index uint64, key []byte, value []byte) error {
	if err := p.IndexDB.Set(index, key, value); err != nil {
		return err
	}

	// Update the highest seen index.
	p.highestSetIndex = max(p.highestSetIndex, index)
	return nil
}

func (db *DB) prune() error {
	// If we haven't used windowSize number of indexes, we can skip
	// the pruning.
	if db.highestSetIndex < db.windowSize {
		return nil
	}

	// TODO: Optimize the underlying DeleteRange to snap to lowest
	// index in O(1).
	if err := db.DeleteRange(db.lastDeletedIndex, db.highestSetIndex-db.windowSize); err != nil {
		db.lastDeletedIndex = 0
		return err
	}
	db.lastDeletedIndex = db.highestSetIndex - db.windowSize - 1

	return nil
}
