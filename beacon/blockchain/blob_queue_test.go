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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//nolint:testpackage // Testing internal components
package blockchain

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func createTestBlobRequest(slot math.Slot, blobCount int) BlobFetchRequest {
	header := &ctypes.BeaconBlockHeader{Slot: slot}
	commitments := make(eip4844.KZGCommitments[common.ExecutionHash], blobCount)
	return BlobFetchRequest{Header: header, Commitments: commitments}
}

// Test that successful write produces valid JSON and cleans up temp files
func TestBlobQueue_SuccessfulWrite(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	slot := math.Slot(100)
	err = queue.Add(slot, createTestBlobRequest(slot, 3))
	require.NoError(t, err)

	// Verify no temp file exists after successful write
	tmpFile := filepath.Join(tmpDir, "0000000100.json.tmp")
	_, err = os.Stat(tmpFile)
	require.True(t, os.IsNotExist(err), "temp file should be cleaned up")

	// Verify final file is valid JSON
	data, err := os.ReadFile(filepath.Join(tmpDir, "0000000100.json"))
	require.NoError(t, err)
	var request BlobFetchRequest
	err = json.Unmarshal(data, &request)
	require.NoError(t, err, "file should contain valid JSON")
}

// Test that recent blob requests are skipped until retry interval passes
func TestBlobQueue_RetryLogic(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	withinDA := func(_, _ math.Slot) bool { return true }
	maxRetries := 72

	// Request with recent retry should be skipped
	request := createTestBlobRequest(math.Slot(100), 1)
	request.LastRetryTime = time.Now()
	data, marshalErr := json.Marshal(request)
	require.NoError(t, marshalErr)
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "0000000100.json"), data, 0600))

	_, _, err = queue.GetNext(math.Slot(200), 5*time.Minute, maxRetries, withinDA)
	require.Error(t, err)
	require.Equal(t, errNoMoreRequests, err, "should skip request not ready for retry")

	// Request with old retry should be returned
	request.LastRetryTime = time.Now().Add(-6 * time.Minute)
	data, marshalErr = json.Marshal(request)
	require.NoError(t, marshalErr)
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "0000000100.json"), data, 0600))

	got, _, err := queue.GetNext(math.Slot(200), 5*time.Minute, maxRetries, withinDA)
	require.NoError(t, err, "should return request ready for retry")
	require.Equal(t, math.Slot(100), got.Header.Slot)
}

// Test that blob requests outside availability window are deleted
func TestBlobQueue_AvailabilityWindow(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	// Add old request
	slot := math.Slot(50)
	err = queue.Add(slot, createTestBlobRequest(slot, 1))
	require.NoError(t, err)

	filename := filepath.Join(tmpDir, "0000000050.json")
	_, err = os.Stat(filename)
	require.NoError(t, err, "file should exist before cleanup")

	// GetNext with request outside DA window should delete it
	withinDAPeriod := func(_, _ math.Slot) bool { return false }
	maxRetries := 72
	_, _, err = queue.GetNext(math.Slot(1000), 1*time.Minute, maxRetries, withinDAPeriod)
	require.Error(t, err)
	require.Equal(t, errNoMoreRequests, err)

	// Verify file was deleted
	_, err = os.Stat(filename)
	require.True(t, os.IsNotExist(err), "old request should be deleted")
}

// Test that failure count increments correctly for retry logic
func TestBlobQueue_UpdateRetry(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	request := createTestBlobRequest(math.Slot(100), 2)
	request.FailureCount = 3

	err = queue.Add(math.Slot(100), request)
	require.NoError(t, err)

	filename := filepath.Join(tmpDir, "0000000100.json")
	err = queue.UpdateRetry(filename, request)
	require.NoError(t, err)

	// Verify failure count incremented
	data, readErr := os.ReadFile(filename)
	require.NoError(t, readErr)
	var updated BlobFetchRequest
	require.NoError(t, json.Unmarshal(data, &updated))
	require.Equal(t, 4, updated.FailureCount, "failure count should increment")
	require.False(t, updated.LastRetryTime.IsZero(), "retry time should be set")
}

// Test that blob queue processes requests in order alphabetically by filename
func TestBlobQueue_ProcessingOrder(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	withinDA := func(_, _ math.Slot) bool { return true }
	maxRetries := 72

	// Add requests out of order
	slots := []math.Slot{200, 100, 150}
	for _, slot := range slots {
		require.NoError(t, queue.Add(slot, createTestBlobRequest(slot, 1)))
	}

	// Should process in ascending slot order (lexicographic filename order)
	expected := []math.Slot{100, 150, 200}
	for _, expectedSlot := range expected {
		got, filename, getErr := queue.GetNext(math.Slot(300), 1*time.Minute, maxRetries, withinDA)
		require.NoError(t, getErr)
		require.Equal(t, expectedSlot, got.Header.Slot)
		require.NoError(t, queue.Remove(filename))
	}

	// Queue should be empty
	_, _, err = queue.GetNext(math.Slot(300), 1*time.Minute, maxRetries, withinDA)
	require.Error(t, err)
	require.Equal(t, errNoMoreRequests, err)
}

// Test that when blob requests exceed the max retry limit they are deleted
func TestBlobQueue_MaxRetryLimit(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	withinDA := func(_, _ math.Slot) bool { return true }
	maxRetries := 72

	// Create request that has exceeded retry limit
	request := createTestBlobRequest(math.Slot(100), 2)
	request.FailureCount = maxRetries // At the limit
	data, marshalErr := json.Marshal(request)
	require.NoError(t, marshalErr)

	filename := filepath.Join(tmpDir, "0000000100.json")
	require.NoError(t, os.WriteFile(filename, data, 0600))

	// GetNext should delete the request and return errNoMoreRequests
	_, _, err = queue.GetNext(math.Slot(200), 1*time.Minute, maxRetries, withinDA)
	require.Error(t, err)
	require.Equal(t, errNoMoreRequests, err)

	// Verify file was deleted
	_, statErr := os.Stat(filename)
	require.True(t, os.IsNotExist(statErr), "request should be deleted after exceeding retry limit")
}

// Test that requests under retry limit are still processed
func TestBlobQueue_UnderRetryLimit(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	withinDA := func(_, _ math.Slot) bool { return true }
	maxRetries := 72

	// Create request with failures but under the limit
	request := createTestBlobRequest(math.Slot(100), 2)
	request.FailureCount = maxRetries - 1                     // One below limit
	request.LastRetryTime = time.Now().Add(-10 * time.Minute) // Ready to retry
	data, marshalErr := json.Marshal(request)
	require.NoError(t, marshalErr)

	filename := filepath.Join(tmpDir, "0000000100.json")
	require.NoError(t, os.WriteFile(filename, data, 0600))

	// GetNext should return the request (not delete it)
	got, gotFilename, err := queue.GetNext(math.Slot(200), 1*time.Minute, maxRetries, withinDA)
	require.NoError(t, err)
	require.Equal(t, math.Slot(100), got.Header.Slot)
	require.Equal(t, filename, gotFilename)
	require.Equal(t, maxRetries-1, got.FailureCount)
}

// Test that corrupted JSON files are renamed to .corrupted and not processed
func TestBlobQueue_CorruptedFileHandling(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	queue, err := newBlobQueue(tmpDir, log.NewNopLogger())
	require.NoError(t, err)

	withinDA := func(_, _ math.Slot) bool { return true }
	maxRetries := 72

	// Create a corrupted JSON file (invalid JSON syntax)
	corruptedFilename := filepath.Join(tmpDir, "0000000100.json")
	corruptedData := []byte(`{"header":{"slot":100},"commitments":[INVALID JSON}`)
	require.NoError(t, os.WriteFile(corruptedFilename, corruptedData, 0600))

	// Create a valid request to ensure queue continues processing
	validSlot := math.Slot(101)
	require.NoError(t, queue.Add(validSlot, createTestBlobRequest(validSlot, 1)))

	// Verify corrupted file exists before processing
	_, err = os.Stat(corruptedFilename)
	require.NoError(t, err, "corrupted file should exist")

	// GetNext should skip corrupted file and process valid one
	got, _, err := queue.GetNext(math.Slot(200), 1*time.Minute, maxRetries, withinDA)
	require.NoError(t, err, "should process valid request despite corrupted file")
	require.Equal(t, validSlot, got.Header.Slot)

	// Verify corrupted file was renamed to .corrupted
	_, err = os.Stat(corruptedFilename)
	require.True(t, os.IsNotExist(err), "original corrupted file should not exist")

	renamedFile := corruptedFilename + ".corrupted"
	_, err = os.Stat(renamedFile)
	require.NoError(t, err, "corrupted file should be renamed to .corrupted")
}
