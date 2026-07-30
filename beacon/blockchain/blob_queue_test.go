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

//nolint:testpackage // we test the unexported blobQueue.
package blockchain

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func newTestQueue(t *testing.T) *blobQueue {
	t.Helper()
	q, err := newBlobQueue(filepath.Join(t.TempDir(), "download_queue"), log.NewTestLogger(t), noopSink{})
	require.NoError(t, err)
	return q
}

type noopSink struct{}

func (noopSink) IncrementCounter(string, ...string)        {}
func (noopSink) SetGauge(string, int64, ...string)         {}
func (noopSink) MeasureSince(string, time.Time, ...string) {}

func testRequest(slot math.Slot, numBlobs int) BlobFetchRequest {
	commitments := make(eip4844.KZGCommitments[common.ExecutionHash], numBlobs)
	return BlobFetchRequest{
		Header:      &ctypes.BeaconBlockHeader{Slot: slot},
		Commitments: commitments,
		Signature:   crypto.BLSSignature{0x01},
	}
}

// withinAlways treats every slot as inside the DA window.
func withinAlways(math.Slot, math.Slot) bool { return true }

// withinNever treats every slot as outside the DA window.
func withinNever(math.Slot, math.Slot) bool { return false }

func TestBlobQueue_AddLoadRemove(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)

	require.NoError(t, q.Add(5, testRequest(5, 2)))
	require.NoError(t, q.Add(7, testRequest(7, 1)))
	require.Equal(t, 2, q.PendingCount())

	// Duplicate adds are no-ops.
	require.NoError(t, q.Add(5, testRequest(5, 2)))
	require.Equal(t, 2, q.PendingCount())

	ready, err := q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 2)
	// Entries come back in slot order (zero-padded filenames).
	require.Equal(t, math.Slot(5), ready[0].request.Header.GetSlot())
	require.Equal(t, math.Slot(7), ready[1].request.Header.GetSlot())
	// The signature survives the round trip.
	require.Equal(t, crypto.BLSSignature{0x01}, ready[0].request.Signature)

	q.Remove(ready[0].filename)
	require.Equal(t, 1, q.PendingCount())

	ready, err = q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	require.Equal(t, math.Slot(7), ready[0].request.Header.GetSlot())
}

// Entries persist across queue restarts (crash safety).
func TestBlobQueue_SurvivesRestart(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "download_queue")
	logger := log.NewTestLogger(t)
	sink := noopSink{}

	q, err := newBlobQueue(dir, logger, sink)
	require.NoError(t, err)
	require.NoError(t, q.Add(9, testRequest(9, 3)))

	q2, err := newBlobQueue(dir, logger, sink)
	require.NoError(t, err)
	require.Equal(t, 1, q2.PendingCount())

	ready, err := q2.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	require.Len(t, ready[0].request.Commitments, 3)
}

// A request is retried indefinitely while in the DA window: UpdateRetry never
// causes deletion, only backoff.
func TestBlobQueue_RetryHasNoCap(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))

	for range 50 {
		ready, err := q.LoadReady(100, 0, withinAlways)
		require.NoError(t, err)
		require.Len(t, ready, 1, "request must survive arbitrarily many retries")
		q.UpdateRetry(ready[0])
	}

	ready, err := q.LoadReady(100, 0, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	require.Equal(t, 50, ready[0].request.FailureCount)
}

// Backoff: an entry that failed recently is not ready again until the retry
// interval elapses.
func TestBlobQueue_RetryBackoff(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))

	ready, err := q.LoadReady(100, time.Hour, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	q.UpdateRetry(ready[0])

	ready, err = q.LoadReady(100, time.Hour, withinAlways)
	require.NoError(t, err)
	require.Empty(t, ready, "entry must respect the retry backoff")
	// Still pending, though (not deleted).
	require.Equal(t, 1, q.PendingCount())
}

// A request whose slot left the DA window is dropped, and stops counting
// against the node's synced status.
func TestBlobQueue_ExpiresOutsideDAWindow(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))

	ready, err := q.LoadReady(100, time.Minute, withinNever)
	require.NoError(t, err)
	require.Empty(t, ready)
	require.Equal(t, 0, q.PendingCount())

	// The file is gone for good.
	ready, err = q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Empty(t, ready)
}

// A head slot of zero (startup, nothing finalized yet) must not expire
// anything.
func TestBlobQueue_ZeroHeadSlotDoesNotExpire(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))

	ready, err := q.LoadReady(0, time.Minute, withinNever)
	require.NoError(t, err)
	require.Len(t, ready, 1)
}

// Corrupted queue files are moved aside instead of wedging the queue.
func TestBlobQueue_CorruptedFileSetAside(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))
	require.NoError(t, os.WriteFile(filepath.Join(q.queueDir, "00000000000000000007.json"), []byte("{garbage"), 0600))

	ready, err := q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	require.Equal(t, math.Slot(5), ready[0].request.Header.GetSlot())

	corrupted, err := filepath.Glob(filepath.Join(q.queueDir, "*.corrupted"))
	require.NoError(t, err)
	require.Len(t, corrupted, 1)
}

// Leftover tmp files from a crash mid-write are cleaned up at startup.
func TestBlobQueue_CleansTmpFilesAtStartup(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "download_queue")
	require.NoError(t, os.MkdirAll(dir, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "x.json.tmp"), []byte("partial"), 0600))

	q, err := newBlobQueue(dir, log.NewTestLogger(t), noopSink{})
	require.NoError(t, err)
	require.Equal(t, 0, q.PendingCount())

	tmpFiles, err := filepath.Glob(filepath.Join(dir, "*.tmp"))
	require.NoError(t, err)
	require.Empty(t, tmpFiles)
}

// The queue must work when the home path contains glob metacharacters: os.ReadDir (not filepath.Glob) means a
// path like beacon[prod] does not silently match zero files and lose queued fetches.
func TestBlobQueue_MetacharPath(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "beacon[prod]", "download_queue")
	q, err := newBlobQueue(dir, log.NewTestLogger(t), noopSink{})
	require.NoError(t, err)

	require.NoError(t, q.Add(5, testRequest(5, 1)))
	require.Equal(t, 1, q.PendingCount())

	ready, err := q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1)
	require.Equal(t, math.Slot(5), ready[0].request.Header.GetSlot())
}

// A corrupted queue file is quarantined and the pending count is decremented, and LoadReady does not clobber a
// concurrent Add's increment (the count is delta-driven, not recomputed from a scan).
func TestBlobQueue_CorruptedDecrementsAndNoClobber(t *testing.T) {
	t.Parallel()
	q := newTestQueue(t)
	require.NoError(t, q.Add(5, testRequest(5, 1)))
	require.NoError(t, q.Add(6, testRequest(6, 1)))
	// Corrupt slot 5's file.
	require.NoError(t, os.WriteFile(q.filename(5), []byte("{garbage"), 0600))

	ready, err := q.LoadReady(100, time.Minute, withinAlways)
	require.NoError(t, err)
	require.Len(t, ready, 1, "only the valid entry is ready")
	require.Equal(t, math.Slot(6), ready[0].request.Header.GetSlot())
	require.Equal(t, 1, q.PendingCount(), "corrupted entry must be decremented")

	corrupted, err := filepath.Glob(filepath.Join(q.queueDir, "*.corrupted"))
	require.NoError(t, err)
	require.Len(t, corrupted, 1)

	// An Add concurrent with the scan is not lost: after LoadReady the count reflects deltas, so adding again
	// increments from the live value rather than a stale scan count.
	require.NoError(t, q.Add(7, testRequest(7, 1)))
	require.Equal(t, 2, q.PendingCount())
}
