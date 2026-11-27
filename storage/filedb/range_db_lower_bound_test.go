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

package filedb_test

import (
	"testing"
	"time"

	file "github.com/berachain/beacon-kit/storage/filedb"
	"github.com/stretchr/testify/require"
)

// TestRangeDB_LowerBoundIndexSurvivesRestart pins the invariant that the prune
// floor is durable. Before it was persisted, a new RangeDB over an existing
// store started at zero, so the next Prune re-walked every index back to
// genesis, issuing one RemoveAll per index inside FinalizeBlock.
func TestRangeDB_LowerBoundIndexSurvivesRestart(t *testing.T) {
	t.Parallel()
	fdb := newTestFDB(t.TempDir())

	rdb := file.NewRangeDB(fdb)
	require.NoError(t, rdb.Prune(0, 50_000))
	require.Equal(t, uint64(50_000), getFirstNonNilIndex(rdb))

	// Same on-disk store, fresh process.
	restarted := file.NewRangeDB(fdb)
	require.Equal(t, uint64(50_000), getFirstNonNilIndex(restarted),
		"prune floor must be restored from disk on restart")
}

// TestRangeDB_PruneAfterRestartDoesNotRewalk is the behavioural half of the
// invariant above: the first Prune after a restart must not re-walk the range
// from genesis.
func TestRangeDB_PruneAfterRestartDoesNotRewalk(t *testing.T) {
	t.Parallel()
	fdb := newTestFDB(t.TempDir())

	rdb := file.NewRangeDB(fdb)
	require.NoError(t, rdb.Prune(0, 200_000))

	restarted := file.NewRangeDB(fdb)
	start := time.Now()
	require.NoError(t, restarted.Prune(0, 200_010))
	elapsed := time.Since(start)

	// The call should walk 10 indexes, not 200_010. Allow a very generous
	// margin so this stays stable on slow CI, while still failing loudly if
	// the full-range walk ever comes back.
	require.Less(t, elapsed, 2*time.Second,
		"Prune after restart re-walked the range from genesis")
	require.Equal(t, uint64(200_010), getFirstNonNilIndex(restarted))
}

// TestRangeDB_FreshStoreStartsAtZero guards the other direction: a store that
// has never been pruned has no persisted floor and must start at zero rather
// than erroring.
func TestRangeDB_FreshStoreStartsAtZero(t *testing.T) {
	t.Parallel()
	rdb := file.NewRangeDB(newTestFDB(t.TempDir()))
	require.Equal(t, uint64(0), getFirstNonNilIndex(rdb))
}
