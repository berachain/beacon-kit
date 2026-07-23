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

package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BlobFetchRequest is one persisted background blob-fetch request. It carries everything needed to verify fetched sidecars without the
// block being around: the trusted header, the commitment list, and the block signature the canonical sidecars must embed.
type BlobFetchRequest struct {
	Header        *ctypes.BeaconBlockHeader                    `json:"header"`
	Commitments   eip4844.KZGCommitments[common.ExecutionHash] `json:"commitments"`
	Signature     crypto.BLSSignature                          `json:"signature"`
	LastRetryTime time.Time                                    `json:"last_retry_time"`
	FailureCount  int                                          `json:"failure_count"`
}

// blobQueue is a crash-safe, file-backed queue of pending blob-fetch requests, one JSON file per slot. Entries persist across restarts
// and are only removed on success or when their slot leaves the DA window.
type blobQueue struct {
	queueDir string
	logger   log.Logger
	sink     TelemetrySink
	// pendingCount tracks the live queue depth; it gates the node's synced status.
	pendingCount atomic.Int64
}

func newBlobQueue(queueDir string, logger log.Logger, sink TelemetrySink) (*blobQueue, error) {
	if err := os.MkdirAll(queueDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create blob download queue directory: %w", err)
	}

	q := &blobQueue{
		queueDir: queueDir,
		logger:   logger,
		sink:     sink,
	}

	// Clean up any leftover temp files from a crash mid-write, and seed the pending count from the queue
	// contents. os.ReadDir (not filepath.Glob) so a home path containing glob metacharacters like [ or * does
	// not silently match nothing.
	entries, err := os.ReadDir(queueDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob download queue directory: %w", err)
	}
	jsonFiles := 0
	for _, e := range entries {
		switch {
		case strings.HasSuffix(e.Name(), ".tmp"):
			_ = os.Remove(filepath.Join(queueDir, e.Name()))
		case strings.HasSuffix(e.Name(), ".json"):
			jsonFiles++
		}
	}
	q.pendingCount.Store(int64(jsonFiles))
	sink.SetGauge(metricsBlobQueueDepth, int64(jsonFiles))

	return q, nil
}

// listRequestFiles returns the absolute paths of the queued request files (*.json), in name order (which is
// slot order, since filenames are zero-padded slots). It uses os.ReadDir so paths with glob metacharacters
// are handled correctly.
func (q *blobQueue) listRequestFiles() ([]string, error) {
	entries, err := os.ReadDir(q.queueDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read queue directory: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			files = append(files, filepath.Join(q.queueDir, e.Name()))
		}
	}
	return files, nil
}

// PendingCount returns the number of queued fetch requests.
func (q *blobQueue) PendingCount() int {
	return int(q.pendingCount.Load())
}

func (q *blobQueue) filename(slot math.Slot) string {
	return filepath.Join(q.queueDir, fmt.Sprintf("%020d.json", slot.Unwrap()))
}

// Add persists a fetch request for the slot; adding an already-queued slot is a no-op.
func (q *blobQueue) Add(slot math.Slot, request BlobFetchRequest) error {
	filename := q.filename(slot)
	if _, err := os.Stat(filename); err == nil {
		q.logger.Debug("Blob fetch request already queued for slot, skipping", "slot", slot.Unwrap())
		return nil
	}
	if err := q.write(filename, request); err != nil {
		return err
	}
	q.pendingCount.Add(1)
	q.sink.SetGauge(metricsBlobQueueDepth, q.pendingCount.Load())
	return nil
}

// write serializes the request to a tmp file and renames it into place, so a crash mid-write can never leave a truncated request behind.
func (q *blobQueue) write(filename string, request BlobFetchRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal blob fetch request: %w", err)
	}

	// Durable write: fsync the temp file before the atomic rename, and fsync the directory after, so a crash
	// or power loss cannot leave a half-written file that later fails to parse and abandons the request.
	tempFile := filename + ".tmp"
	if writeErr := writeFileSync(tempFile, data); writeErr != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to write temp blob fetch request: %w", writeErr)
	}
	if renameErr := os.Rename(tempFile, filename); renameErr != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to rename blob fetch request: %w", renameErr)
	}
	return syncDir(q.queueDir)
}

// writeFileSync writes data to path and fsyncs it before returning.
func writeFileSync(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600) // #nosec G304 -- path constructed from the queue directory
	if err != nil {
		return err
	}
	if _, err = f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// syncDir fsyncs a directory so a rename into it is durable.
func syncDir(dir string) error {
	d, err := os.Open(dir) // #nosec G304 -- queue directory owned by this process
	if err != nil {
		return err
	}
	if err = d.Sync(); err != nil {
		_ = d.Close()
		return err
	}
	return d.Close()
}

// queuedRequest pairs a loaded request with its backing file.
type queuedRequest struct {
	request  BlobFetchRequest
	filename string
}

// LoadReady returns the queued requests that are ready for a fetch attempt: still within the DA window relative to headSlot and past
// their retry backoff. Requests whose slot has left the DA window are deleted (with a metric); corrupted files are renamed aside for
// inspection. There is no retry cap: a request lives until it succeeds or its slot leaves the window.
//
// It does NOT recompute the pending count from the directory scan: doing so would clobber a concurrent Add from
// the consensus goroutine. The count is delta-driven (Add on enqueue, Remove/drop on dequeue), so every path
// here that removes a file also decrements it.
func (q *blobQueue) LoadReady(
	headSlot math.Slot,
	retryInterval time.Duration,
	withinDAPeriod func(block, current math.Slot) bool,
) ([]queuedRequest, error) {
	files, err := q.listRequestFiles()
	if err != nil {
		return nil, err
	}

	ready := make([]queuedRequest, 0, len(files))
	for _, filename := range files {
		fileData, readErr := os.ReadFile(filename) // #nosec G304 -- filename comes from queueDir
		if readErr != nil {
			q.logger.Error("Failed to read blob fetch request", "file", filename, "error", readErr)
			continue
		}

		var request BlobFetchRequest
		if err = json.Unmarshal(fileData, &request); err != nil || request.Header == nil {
			q.handleCorrupted(filename, headSlot, withinDAPeriod, err)
			continue
		}

		// Drop requests whose slot left the DA window; nobody is required to serve them anymore and the node no longer needs them to be synced.
		if headSlot > 0 && !withinDAPeriod(request.Header.GetSlot(), headSlot) {
			q.sink.IncrementCounter(metricsBlobRequestsExpired)
			q.logger.Warn("Blob fetch request left the DA window, dropping",
				"slot", request.Header.GetSlot().Unwrap(),
				"head_slot", headSlot.Unwrap(),
				"failure_count", request.FailureCount)
			q.Remove(filename)
			continue
		}

		if !request.LastRetryTime.IsZero() && time.Since(request.LastRetryTime) < retryInterval {
			continue // Not ready to retry yet.
		}
		ready = append(ready, queuedRequest{request: request, filename: filename})
	}

	return ready, nil
}

// handleCorrupted quarantines an unparseable queue file. fsynced writes make this essentially unreachable, so
// an in-window corruption is an alert-worthy data-availability event: the request (including the block
// signature) cannot be recovered from this layer, and the block's blobs may be missing until an operator
// intervenes. Out-of-window corruption is harmless.
func (q *blobQueue) handleCorrupted(
	filename string,
	headSlot math.Slot,
	withinDAPeriod func(block, current math.Slot) bool,
	cause error,
) {
	inWindow := headSlot == 0 || slotFromFilename(filename, withinDAPeriod, headSlot)
	corruptedFile := filename + ".corrupted"
	if renameErr := os.Rename(filename, corruptedFile); renameErr != nil {
		_ = os.Remove(filename)
	}
	q.decrement()
	if inWindow {
		q.logger.Error("Corrupted in-window blob fetch request; blobs for this slot may be missing",
			"file", filename, "corrupted_file", corruptedFile, "error", cause)
	} else {
		q.logger.Warn("Corrupted out-of-window blob fetch request, quarantining",
			"file", filename, "error", cause)
	}
}

// slotFromFilename parses the zero-padded slot out of a queue filename and reports whether it is still within
// the DA window. On a parse failure it conservatively reports in-window (treat as the alert-worthy case).
func slotFromFilename(filename string, withinDAPeriod func(block, current math.Slot) bool, headSlot math.Slot) bool {
	base := strings.TrimSuffix(filepath.Base(filename), ".json")
	n, err := strconv.ParseUint(base, 10, 64)
	if err != nil {
		return true
	}
	return withinDAPeriod(math.Slot(n), headSlot)
}

// Remove deletes a fulfilled request and decrements the pending count.
func (q *blobQueue) Remove(filename string) {
	if err := os.Remove(filename); err != nil && !errors.Is(err, os.ErrNotExist) {
		q.logger.Error("Failed to delete blob fetch request", "file", filename, "error", err)
		return
	}
	q.decrement()
}

// decrement lowers the pending count by one, clamped at zero.
func (q *blobQueue) decrement() {
	if q.pendingCount.Add(-1) < 0 {
		q.pendingCount.Store(0)
	}
	q.sink.SetGauge(metricsBlobQueueDepth, q.pendingCount.Load())
}

// UpdateRetry bumps the failure count and retry timestamp after a failed attempt.
func (q *blobQueue) UpdateRetry(entry queuedRequest) {
	entry.request.FailureCount++
	entry.request.LastRetryTime = time.Now()
	if err := q.write(entry.filename, entry.request); err != nil {
		q.logger.Error("Failed to update blob fetch retry metadata",
			"file", entry.filename, "error", err)
	}
}
