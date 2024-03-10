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

package skiplist_test

import (
	"testing"

	"github.com/berachain/beacon-kit/lib/skiplist"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

// var _ skiplist.Comparable[Uint64Comparable] = Uint64Comparable{}

type Uint64Comparable uint64

func (i Uint64Comparable) Compare(rhs Uint64Comparable) int {
	if uint64(i) < uint64(rhs) {
		return -1
	} else if uint64(i) > uint64(rhs) {
		return 1
	}
	return 0
}

func TestSkiplist(t *testing.T) {
	// Create a new ordered skiplist.
	skiplist := skiplist.New[Uint64Comparable]()

	// Insert elements.
	skiplist.Insert(2)
	skiplist.Insert(5)
	skiplist.Insert(4)
	skiplist.Insert(1)
	skiplist.Insert(3)

	// Remove elements.
	// 1 2 3 4 5
	i, err := skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(1), uint64(i))
	require.Equal(t, 4, skiplist.Len())

	// 2 3 4 5
	i, err = skiplist.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, uint64(5), uint64(i))
	require.Equal(t, 3, skiplist.Len())

	// 2 3 4
	skiplist.Insert(2)
	require.Equal(t, 3, skiplist.Len())

	// 2 3 4 5
	skiplist.Insert(5)
	require.Equal(t, 4, skiplist.Len())

	// 2 3 4 5
	i, err = skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(2), uint64(i))
	require.Equal(t, 3, skiplist.Len())

	// 3 4 5
	i, err = skiplist.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, uint64(5), uint64(i))
	require.Equal(t, 2, skiplist.Len())

	// 3 4
	i, err = skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, uint64(3), uint64(i))
	require.Equal(t, 1, skiplist.Len())
}

type Log struct {
	blockNumber uint64
	logIndex    uint64
}

func NewLog(blockNumber, logIndex uint64) Log {
	return Log{blockNumber, logIndex}
}

func (l Log) Compare(rhs Log) int {
	if c := Uint64Comparable(l.blockNumber).
		Compare(Uint64Comparable(rhs.blockNumber)); c != 0 {
		return c
	}
	return Uint64Comparable(l.logIndex).
		Compare(Uint64Comparable(rhs.logIndex))
}

func TestLogskiplist(t *testing.T) {
	// Create a new ordered skiplist.
	skiplist := skiplist.New[Log]()

	// Insert elements.
	skiplist.Insert(NewLog(2, 2))
	skiplist.Insert(NewLog(1, 1))
	skiplist.Insert(NewLog(2, 1))
	skiplist.Insert(NewLog(1, 2))

	// Remove elements.
	// (1, 1) (1, 2) (2, 1) (2, 2)
	i, err := skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(1, 1), i)
	require.Equal(t, 3, skiplist.Len())

	// (1, 2) (2, 1) (2, 2)
	i, err = skiplist.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 2), i)
	require.Equal(t, 2, skiplist.Len())

	// (1, 2) (2, 1)
	skiplist.Insert(NewLog(1, 2))
	require.Equal(t, 2, skiplist.Len())

	// (1, 2) (2, 1) (2, 2)
	skiplist.Insert(NewLog(2, 2))
	require.Equal(t, 3, skiplist.Len())

	// (1, 2) (2, 1) (2, 2)
	i, err = skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(1, 2), i)
	require.Equal(t, 2, skiplist.Len())

	// (2, 1) (2, 2)
	i, err = skiplist.RemoveBack()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 2), i)
	require.Equal(t, 1, skiplist.Len())

	// (2, 1)
	i, err = skiplist.RemoveFront()
	require.NoError(t, err)
	require.Equal(t, NewLog(2, 1), i)
	require.Equal(t, 0, skiplist.Len())
}
