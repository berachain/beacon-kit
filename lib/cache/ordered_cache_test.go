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

package cache_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/lib/cache"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

var _ cache.Comparable[uint64] = Uint64Comparable{}

type Uint64Comparable struct{}

func (Uint64Comparable) Compare(lhs, rhs uint64) int {
	if lhs < rhs {
		return -1
	} else if lhs > rhs {
		return 1
	}
	return 0
}

func TestOrderedCache(t *testing.T) {
	// Create a new ordered cache.
	cache := cache.NewOrderedCache[uint64](Uint64Comparable{})

	// Insert elements.
	cache.Insert(2)
	cache.Insert(5)
	cache.Insert(4)
	cache.Insert(1)
	cache.Insert(3)

	// Remove elements.
	// 1 2 3 4 5
	i, err := cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(1), i)
	require.Equal(t, 4, cache.Len())

	// 2 3 4 5
	i, err = cache.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, uint64(5), i)
	require.Equal(t, 3, cache.Len())

	// 2 3 4
	cache.Insert(2)
	require.Equal(t, 3, cache.Len())

	// 2 3 4 5
	cache.Insert(5)
	require.Equal(t, 4, cache.Len())

	// 2 3 4 5
	i, err = cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(2), i)
	require.Equal(t, 3, cache.Len())

	// 3 4 5
	i, err = cache.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, uint64(5), i)
	require.Equal(t, 2, cache.Len())

	// 3 4
	i, err = cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(3), i)
	require.Equal(t, 1, cache.Len())
}

type Log struct {
	blockNumber uint64
	logIndex    uint64
}

func NewLog(blockNumber, logIndex uint64) Log {
	return Log{blockNumber, logIndex}
}

var _ cache.Comparable[Log] = LogComparable{}

type LogComparable struct{}

// Compare returns the lexicographic comparison of the two logs.
func (LogComparable) Compare(lhs, rhs Log) int {
	c := Uint64Comparable{}
	compareBlockNumber := c.Compare(lhs.blockNumber, rhs.blockNumber)
	if compareBlockNumber != 0 {
		return compareBlockNumber
	}
	return c.Compare(lhs.logIndex, rhs.logIndex)
}

func TestLogCache(t *testing.T) {
	// Create a new ordered cache.
	cache := cache.NewOrderedCache[Log](LogComparable{})

	// Insert elements.
	cache.Insert(NewLog(2, 2))
	cache.Insert(NewLog(1, 1))
	cache.Insert(NewLog(2, 1))
	cache.Insert(NewLog(1, 2))

	// Remove elements.
	// (1, 1) (1, 2) (2, 1) (2, 2)
	i, err := cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(1, 1), i)
	require.Equal(t, 3, cache.Len())

	// (1, 2) (2, 1) (2, 2)
	i, err = cache.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 2), i)
	require.Equal(t, 2, cache.Len())

	// (1, 2) (2, 1)
	cache.Insert(NewLog(1, 2))
	require.Equal(t, 2, cache.Len())

	// (1, 2) (2, 1) (2, 2)
	cache.Insert(NewLog(2, 2))
	require.Equal(t, 3, cache.Len())

	// (1, 2) (2, 1) (2, 2)
	i, err = cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(1, 2), i)
	require.Equal(t, 2, cache.Len())

	// (2, 1) (2, 2)
	i, err = cache.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 2), i)
	require.Equal(t, 1, cache.Len())

	// (2, 1)
	i, err = cache.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 1), i)
	require.Equal(t, 0, cache.Len())
}
