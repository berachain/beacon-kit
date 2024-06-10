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
	"testing"

	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/stretchr/testify/require"
)

func TestPayloadIDCache(t *testing.T) {
	cacheUnderTest := cache.NewPayloadIDCache[[8]byte, [32]byte, uint64]()

	t.Run("Get from empty cache", func(t *testing.T) {
		var r [32]byte
		p, ok := cacheUnderTest.Get(0, r)
		require.False(t, ok)
		require.Equal(t, [8]byte{}, p)
	})

	t.Run("Set and Get", func(t *testing.T) {
		slot := uint64(1234)
		r := [32]byte{1, 2, 3}
		pid := [8]byte{1, 2, 3, 3, 7, 8, 7, 8}
		cacheUnderTest.Set(slot, r, pid)

		p, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(t, pid, p)
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		slot := uint64(1234)
		r := [32]byte{1, 2, 3}
		newPid := [8]byte{9, 9, 9, 9, 9, 9, 9, 9}
		cacheUnderTest.Set(slot, r, newPid)

		p, ok := cacheUnderTest.Get(slot, r)
		require.True(t, ok)
		require.Equal(t, newPid, p)
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		slot := uint64(9456456)
		r := [32]byte{4, 5, 6}
		pid := [8]byte{4, 5, 6, 6, 9, 0, 9, 0}
		cacheUnderTest.Set(slot, r, pid)

		// Prune and attempt to retrieve pruned entry
		cacheUnderTest.UnsafePrunePrior(slot + 1)
		p, ok := cacheUnderTest.Get(slot, r)
		require.False(t, ok)
		require.Equal(t, [8]byte{}, p)
	})

	t.Run("Multiple entries and prune", func(t *testing.T) {
		// Set multiple entries
		for i := range uint8(5) {
			slot := uint64(i)
			r := [32]byte{i, i + 1, i + 2}
			pid := [8]byte{
				i, i, i, i, i, i, i, i,
			}
			cacheUnderTest.Set(slot, r, pid)
		}

		// Prune and check if only the last two entries exist
		cacheUnderTest.UnsafePrunePrior(3)
		for i := range uint8(3) {
			slot := uint64(i)
			r := [32]byte{i, i + 1, i + 2}
			_, ok := cacheUnderTest.Get(slot, r)
			require.False(t, ok, "Expected entry to be pruned for slot", slot)
		}

		for i := uint8(3); i < 5; i++ {
			slot := uint64(i)
			r := [32]byte{i, i + 1, i + 2}
			_, ok := cacheUnderTest.Get(slot, r)
			require.True(t, ok, "Expected entry to exist for slot", slot)
		}
	})
}
