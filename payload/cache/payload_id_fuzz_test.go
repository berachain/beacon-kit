// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/stretchr/testify/require"
)

func FuzzPayloadIDCacheBasic(f *testing.F) {
	f.Add(uint64(1), []byte{1, 2, 3}, []byte{1, 2, 3, 4, 5, 6, 7, 8})
	f.Add(uint64(2), []byte{4, 5, 6}, []byte{9, 10, 11, 12, 13, 14, 15, 16})
	f.Add(uint64(3), []byte{7, 8, 9}, []byte{17, 18, 19, 20, 21, 22, 23, 24})
	f.Fuzz(func(t *testing.T, s uint64, _r, _p []byte) {
		var r [32]byte
		copy(r[:], _r)
		slot := s
		pid := [8]byte(_p[:8])
		cacheUnderTest := cache.NewPayloadIDCache[[8]byte, [32]byte, uint64]()
		cacheUnderTest.Set(slot, r, pid)

		p, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(t, pid, p)

		// Test overwriting the same slot and root with a different PayloadID
		newPid := [8]byte{}
		for i := range pid {
			newPid[i] = pid[i] + 1 // Simple mutation for a new PayloadID
		}
		cacheUnderTest.Set((slot), r, newPid)

		p, ok = cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(
			t, newPid, p, "PayloadID should be overwritten with the new value")

		// Prune and verify deletion
		cacheUnderTest.UnsafePrunePrior((slot) + 1)
		_, ok = cacheUnderTest.Get(slot, r)
		require.False(t, ok, "Entry should be pruned and not found")
	})
}

func FuzzPayloadIDInvalidInput(f *testing.F) {
	// Intentionally invalid inputs
	f.Add(uint64(1), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, []byte{1, 2, 3})

	f.Fuzz(func(t *testing.T, s uint64, _r, _p []byte) {
		var r [32]byte
		slot := s
		if len(_r) > 32 {
			// Expecting an error or specific handling of oversized input
			t.Skip(
				"Skipping test due to intentionally invalid input size for root.",
			)
		}
		copy(r[:], _r)
		var paddedPayload [8]byte
		copy(paddedPayload[:], _p[:min(len(_p), 8)])
		pid := [8]byte(paddedPayload[:])
		cacheUnderTest := cache.NewPayloadIDCache[[8]byte, [32]byte, uint64]()
		cacheUnderTest.Set(slot, r, pid)

		_, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
	})
}

func FuzzPayloadIDCacheConcurrency(f *testing.F) {
	f.Add(uint64(1), []byte{1, 2, 3}, []byte{1, 2, 3, 4})

	f.Fuzz(func(t *testing.T, s uint64, _r, _p []byte) {
		cacheUnderTest := cache.NewPayloadIDCache[[8]byte, [32]byte, uint64]()
		slot := s
		var wg sync.WaitGroup
		wg.Add(2)

		// Set operation in one goroutine
		go func() {
			defer wg.Done()
			var r [32]byte
			copy(r[:], _r)
			var paddedPayload [8]byte
			copy(paddedPayload[:], _p[:min(len(_p), 8)])
			pid := [8]byte(paddedPayload[:])
			cacheUnderTest.Set((slot), r, pid)
		}()

		// Get operation in another goroutine
		var ok bool
		go func() {
			defer wg.Done()
			time.Sleep(
				10 * time.Millisecond,
			) // Small delay to let the Set operation proceed
			var r [32]byte
			copy(r[:], _r)
			_, ok = cacheUnderTest.Get(slot, r)
		}()

		wg.Wait()
		require.True(t, ok)
	})
}
