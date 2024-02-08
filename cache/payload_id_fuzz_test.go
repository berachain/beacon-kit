// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

	"github.com/itsdevbear/bolaris/cache"
	"github.com/itsdevbear/bolaris/types/primitives"
	"github.com/stretchr/testify/require"
)

func FuzzPayloadIDCache(f *testing.F) {
	f.Add(uint64(1), []byte{1, 2, 3}, []byte{1, 2, 3, 4, 5, 6, 7, 8})
	f.Fuzz(func(t *testing.T, slot uint64, _r, _p []byte) {
		var r [32]byte
		copy(r[:], _r)
		pid := primitives.PayloadID(_p[:8])
		cacheUnderTest := cache.NewPayloadIDCache()
		cacheUnderTest.Set(primitives.Slot(slot), r, pid)

		p, ok := cacheUnderTest.Get(primitives.Slot(slot), r)
		require.True(t, ok)
		require.Equal(t, pid, p)

		// Test overwriting the same slot and root with a different PayloadID
		newPid := primitives.PayloadID{}
		for i := range pid {
			newPid[i] = pid[i] + 1 // Simple mutation for a new PayloadID
		}
		cacheUnderTest.Set(primitives.Slot(slot), r, newPid)

		p, ok = cacheUnderTest.Get(primitives.Slot(slot), r)
		require.True(t, ok)
		require.Equal(t, newPid, p, "PayloadID should be overwritten with the new value")

		// Prune and verify deletion
		cacheUnderTest.UnsafePrune(primitives.Slot(slot) + 1)
		_, ok = cacheUnderTest.Get(primitives.Slot(slot), r)
		require.False(t, ok, "Entry should be pruned and not found")
	})
}
