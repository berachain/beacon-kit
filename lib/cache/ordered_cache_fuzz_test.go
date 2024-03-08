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
	"math/rand"
	"sync"
	"testing"

	"github.com/berachain/beacon-kit/lib/cache"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

var _ cache.Comparable[int] = IntComparable{}

type IntComparable struct{}

func (IntComparable) Compare(lhs, rhs int) int {
	if lhs < rhs {
		return -1
	} else if lhs > rhs {
		return 1
	}
	return 0
}

func FuzzOrderedCacheSimple(f *testing.F) {
	// Create a new ordered cache.
	cache := cache.NewOrderedCache[int](IntComparable{})

	f.Add(10)
	f.Fuzz(func(t *testing.T, n int) {
		if n < 0 {
			t.Skip()
		}

		for _, elem := range rand.Perm(n) {
			cache.Insert(elem)
		}

		for i := range n {
			if i%2 == 0 {
				// i: 0 2 4 6 8
				// e: 0 1 2 3 4
				e, err := cache.RemoveFront()
				require.NoError(t, err)
				require.Equal(t, i/2, e)
			} else {
				// i: 1 3 5 7 9
				// e: 9 8 7 6 5
				e, err := cache.RemoveBack()
				require.NoError(t, err)
				require.Equal(t, n-(i+1)/2, e)
			}
			require.Equal(t, n-i-1, cache.Len())
		}
	})
}

func FuzzOrderedCacheConcurrencySafety(f *testing.F) {
	f.Fuzz(func(t *testing.T, n int) {
		if n <= 0 {
			t.Skip()
		}

		cache := cache.NewOrderedCache(IntComparable{})
		numGoroutines := 10
		numOperations := 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					switch rand.Intn(3) {
					case 0: // Insert
						cache.Insert(
							rand.Intn(n),
						) // Use n to generate random inputs
					case 1: // RemoveFront
						_, _ = cache.RemoveFront() // Ignore errors for simplicity
					case 2: // RemoveBack
						_, _ = cache.RemoveBack() // Ignore errors for simplicity
					}
				}
			}()
		}

		wg.Wait()

		// After all operations, the cache should be in a consistent state.
		// We can't assert the exact contents of the cache, but we can check
		// if the cache is sorted correctly.
		var prev int
		first := true
		_, err := cache.Front()
		for err == nil {
			// Remove the element to get the next in the next iteration
			value, _ := cache.RemoveFront()
			if !first && prev > value {
				t.Errorf(
					"cache is not in ascending order: found %d after %d",
					value,
					prev,
				)
			}
			prev = value
			first = false
			_, err = cache.Front()
		}
	})
}
