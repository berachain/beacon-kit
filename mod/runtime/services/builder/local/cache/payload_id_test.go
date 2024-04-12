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

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/engine"
	"github.com/berachain/beacon-kit/mod/runtime/services/builder/local/cache"
	"github.com/stretchr/testify/require"
)

func TestPayloadIDCache(t *testing.T) {
	cacheUnderTest := cache.NewPayloadIDCache()

	t.Run("Get from empty cache", func(t *testing.T) {
		var r [32]byte
		p, ok := cacheUnderTest.Get(0, r)
		require.False(t, ok)
		require.Equal(t, engine.PayloadID{}, p)
	})

	t.Run("Set and Get", func(t *testing.T) {
		slot := primitives.Slot(1234)
		r := [32]byte{1, 2, 3}
		pid := engine.PayloadID{1, 2, 3, 3, 7, 8, 7, 8}
		cacheUnderTest.Set(slot, r, pid)

		p, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(t, pid, p)
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		slot := primitives.Slot(1234)
		r := [32]byte{1, 2, 3}
		newPid := engine.PayloadID{9, 9, 9, 9, 9, 9, 9, 9}
		cacheUnderTest.Set(slot, r, newPid)

		p, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(t, newPid, p)
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		slot := primitives.Slot(9456456)
		r := [32]byte{4, 5, 6}
		pid := engine.PayloadID{4, 5, 6, 6, 9, 0, 9, 0}
		cacheUnderTest.Set(slot, r, pid)

		// Prune and attempt to retrieve pruned entry
		cacheUnderTest.UnsafePrunePrior(slot + 1)
		p, ok := cacheUnderTest.Get(slot, r)
		require.False(t, ok)
		require.Equal(t, engine.PayloadID{}, p)
	})

	t.Run("Multiple entries and prune", func(t *testing.T) {
		// Set multiple entries
		for i := range uint8(5) {
			slot := primitives.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			pid := engine.PayloadID{
				i, i, i, i, i, i, i, i,
			}
			cacheUnderTest.Set(slot, r, pid)
		}

		// Prune and check if only the last two entries exist
		cacheUnderTest.UnsafePrunePrior(3)
		for i := range uint8(3) {
			slot := primitives.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			_, ok := cacheUnderTest.Get(slot, r)
			require.False(t, ok, "Expected entry to be pruned for slot", slot)
		}

		for i := uint8(3); i < 5; i++ {
			slot := primitives.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			_, ok := cacheUnderTest.Get(slot, r)
			require.True(t, ok, "Expected entry to exist for slot", slot)
		}
	})
}
