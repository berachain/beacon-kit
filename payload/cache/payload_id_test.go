// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestPayloadIDCache(t *testing.T) {
	t.Parallel()
	cacheUnderTest := cache.NewPayloadIDCache()

	t.Run("Get from empty cache", func(t *testing.T) {
		var r [32]byte
		p, ok := cacheUnderTest.GetAndEvict(0, r)
		require.False(t, ok)
		require.Equal(t, engineprimitives.PayloadID{}, p.PayloadID)
	})

	t.Run("Set and Get", func(t *testing.T) {
		slot := math.Slot(1234)
		r := [32]byte{1, 2, 3}
		pid := engineprimitives.PayloadID{1, 2, 3, 3, 7, 8, 7, 8}
		cacheUnderTest.Set(slot, r, pid, version.Deneb())

		p, ok := cacheUnderTest.GetAndEvict(slot, r)
		require.True(t, ok)
		require.Equal(t, pid, p.PayloadID)
	})

	t.Run("Overwrite existing", func(t *testing.T) {
		slot := math.Slot(1234)
		r := [32]byte{1, 2, 3}
		newPid := engineprimitives.PayloadID{9, 9, 9, 9, 9, 9, 9, 9}
		cacheUnderTest.Set(slot, r, newPid, version.Deneb())

		p, ok := cacheUnderTest.GetAndEvict(slot, r)
		require.True(t, ok)
		require.Equal(t, newPid, p.PayloadID)
	})

	t.Run("Prune and verify deletion", func(t *testing.T) {
		slot := math.Slot(9456456)
		r := [32]byte{4, 5, 6}
		pid := engineprimitives.PayloadID{4, 5, 6, 6, 9, 0, 9, 0}
		// Set pid for slot.
		cacheUnderTest.Set(slot, r, pid, version.Deneb())

		// Set historicalPayloadIDCacheSize+1 number of pids. This should
		// prune the first slot from the cache.
		cacheUnderTest.Set(slot+1, r, pid, version.Deneb())
		cacheUnderTest.Set(slot+2, r, pid, version.Deneb())
		cacheUnderTest.Set(slot+3, r, pid, version.Deneb())

		// Attempt to retrieve pruned slot.
		ok := cacheUnderTest.Has(slot, r)
		require.False(t, ok)
	})

	t.Run("Multiple entries and prune", func(t *testing.T) {
		numEntries := 10
		historicalPayloadIDCacheSize := 2
		// Set multiple entries
		for i := range uint8(numEntries) {
			slot := math.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			pid := [8]byte{
				i, i, i, i, i, i, i, i,
			}
			cacheUnderTest.Set(slot, r, pid, version.Deneb())
		}

		// Only the last historicalPayloadIDCacheSize+1 number of entries
		// should still exist.
		for i := range uint8(numEntries - (historicalPayloadIDCacheSize + 1)) {
			slot := math.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			ok := cacheUnderTest.Has(slot, r)
			require.False(t, ok, "Expected entry to be pruned for slot", slot)
		}

		for i := uint8(numEntries - (historicalPayloadIDCacheSize + 1)); i < uint8(numEntries); i++ {
			slot := math.Slot(i)
			r := [32]byte{i, i + 1, i + 2}
			ok := cacheUnderTest.Has(slot, r)
			require.True(t, ok, "Expected entry to exist for slot", slot)
		}
	})
}
