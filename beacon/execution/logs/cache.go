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

package logs

import (
	"github.com/itsdevbear/bolaris/lib/cache"
)

type Cache struct {
	// This is an in-memory cache of logs.
	cache.Cache[LogContainer]
}

// NewCache returns a new cache for LogContainer.
func NewCache() *Cache {
	return &Cache{
		cache.NewOrderedCache[LogContainer](LogComparable{}),
	}
}

// LogComparable is a comparable for LogContainer.
type LogComparable struct{}

// Compare is a lexicographic comparison of logs.
func (LogComparable) Compare(lhs, rhs LogContainer) int {
	blockCmp := compareUint64(lhs.BlockNumber(), rhs.BlockNumber())
	if blockCmp != 0 {
		return blockCmp
	}
	return compareUint64(lhs.LogIndex(), rhs.LogIndex())
}

// compareUint64 compares two uint64 values.
func compareUint64(a, b uint64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
