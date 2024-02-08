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

	cache "github.com/itsdevbear/bolaris/cache"
	"github.com/itsdevbear/bolaris/types/primitives"
	"github.com/stretchr/testify/require"
)

func TestValidatorPayloadIDsCache_GetAndSaveValidatorPayloadIDs(t *testing.T) {
	cacheUnderTest := cache.NewPayloadIDCache()
	var r [32]byte
	p, ok := cacheUnderTest.Get(0, r)
	require.False(t, ok)
	require.Equal(t, primitives.PayloadID{}, p)

	slot := primitives.Slot(1234)
	pid := primitives.PayloadID{1, 2, 3, 3, 7, 8, 7, 8}
	r = [32]byte{1, 2, 3}
	cacheUnderTest.Set(slot, r, pid)
	p, ok = cacheUnderTest.Get(slot, r)
	require.True(t, ok)
	require.Equal(t, pid, p)

	slot = primitives.Slot(9456456)
	r = [32]byte{4, 5, 6}
	cacheUnderTest.Set(slot, r, primitives.PayloadID{})
	p, ok = cacheUnderTest.Get(slot, r)
	require.True(t, ok)
	require.Equal(t, primitives.PayloadID{}, p)

	// reset cache without pid
	slot = primitives.Slot(9456456)
	r = [32]byte{7, 8, 9}
	pid = [8]byte{3, 2, 3, 33, 72, 8, 7, 8}
	cacheUnderTest.Set(slot, r, pid)
	p, ok = cacheUnderTest.Get(slot, r)
	require.True(t, ok)
	require.Equal(t, pid, p)

	// Forked chain
	r = [32]byte{1, 2, 3}
	p, ok = cacheUnderTest.Get(slot, r)
	require.False(t, ok)
	require.Equal(t, primitives.PayloadID{}, p)

	// existing pid - change the cache
	slot = primitives.Slot(9456456)
	r = [32]byte{7, 8, 9}
	newPid := primitives.PayloadID{1, 2, 3, 33, 72, 8, 7, 1}
	cacheUnderTest.Set(slot, r, newPid)
	p, ok = cacheUnderTest.Get(slot, r)
	require.True(t, ok)
	require.Equal(t, newPid, p)

	// remove cache entry
	cacheUnderTest.Prune(slot + 1)
	p, ok = cacheUnderTest.Get(slot, r)
	require.Equal(t, primitives.PayloadID{}, p)
	require.False(t, ok)
}
