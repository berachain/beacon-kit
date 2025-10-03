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

package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
)

var (
	errNoMoreRequests = errors.New("no more requests in queue")
)

// BlobFetchRequest contains the minimal data needed to fetch and validate blobs.
type BlobFetchRequest struct {
	Header        *ctypes.BeaconBlockHeader                    `json:"header"`
	Commitments   eip4844.KZGCommitments[common.ExecutionHash] `json:"commitments"`
	LastRetryTime time.Time                                    `json:"last_retry_time"`
	FailureCount  int                                          `json:"failure_count"`
}

// blobQueue handles persistent queue operations using the filesystem.
// This struct has no external dependencies and can be tested without mocks.
type blobQueue struct {
	queueDir string
	logger   log.Logger
}

// newBlobQueue creates a new blob queue with the given directory.
// It creates the directory if it doesn't exist and cleans up orphaned temp files.
func newBlobQueue(queueDir string, logger log.Logger) (*blobQueue, error) {
	// Create queue directory
	if err := os.MkdirAll(queueDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create blob download queue directory: %w", err)
	}

	// Clean up any leftover temp files in the unlikely event we crashed while writing a request
	tmpFiles, _ := filepath.Glob(filepath.Join(queueDir, "*.tmp"))
	for _, tmpFile := range tmpFiles {
		_ = os.Remove(tmpFile)
	}

	return &blobQueue{
		queueDir: queueDir,
		logger:   logger,
	}, nil
}

// Add queues a new blob fetch request.
func (q *blobQueue) Add(slot math.Slot, request BlobFetchRequest) error {
	// Serialize to JSON file with slot as filename
	filename := filepath.Join(q.queueDir, fmt.Sprintf("%010d.json", slot.Unwrap()))

	// Check if request already exists for this slot
	if _, err := os.Stat(filename); err == nil {
		q.logger.Info("Blob fetch request already queued for slot, skipping", "slot", slot.Unwrap())
		return nil
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal blob fetch request: %w", err)
	}

	// Write the request to a tmp file first, then rename atomically. This prevents
	// reading a partially written file causing JSON unmarshal errors.
	tempFile := filename + ".tmp"
	if writeErr := os.WriteFile(tempFile, data, 0600); writeErr != nil {
		return fmt.Errorf("failed to write temp blob fetch request: %w", writeErr)
	}
	if renameErr := os.Rename(tempFile, filename); renameErr != nil {
		_ = os.Remove(tempFile)
		return fmt.Errorf("failed to rename blob fetch request: %w", renameErr)
	}

	return nil
}

// GetNext returns the next request that is ready to be processed.
// It skips requests outside the availability window and requests not ready for retry.
func (q *blobQueue) GetNext(
	headSlot math.Slot,
	retryInterval time.Duration,
	maxRetries int,
	withinDAPeriod func(block, current math.Slot) bool,
) (BlobFetchRequest, string, error) {
	files, err := os.ReadDir(q.queueDir)
	if err != nil {
		return BlobFetchRequest{}, "", fmt.Errorf("failed to read queue directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := filepath.Join(q.queueDir, file.Name())
		fileData, readErr := os.ReadFile(filename) // #nosec G304 // filename is constructed from queueDir
		if readErr != nil {
			return BlobFetchRequest{}, filename, fmt.Errorf("failed to read request file: %w", readErr)
		}

		var request BlobFetchRequest
		if err = json.Unmarshal(fileData, &request); err != nil {
			return BlobFetchRequest{}, filename, fmt.Errorf("failed to unmarshal request: %w", err)
		}

		// Check if request is outside availability window
		if headSlot > 0 && !withinDAPeriod(request.Header.Slot, headSlot) {
			q.logger.Warn("Request is outside availability window, deleting",
				"slot", request.Header.Slot.Unwrap(),
				"head_slot", headSlot.Unwrap(),
				"failure_count", request.FailureCount)
			_ = q.Remove(filename)
			continue
		}

		// Check if request has exceeded max retry limit
		if request.FailureCount >= maxRetries {
			q.logger.Warn("Request exceeded max retry limit, deleting",
				"slot", request.Header.Slot.Unwrap(),
				"failure_count", request.FailureCount,
				"max_retries", maxRetries)
			_ = q.Remove(filename)
			continue
		}

		// Check if this request needs to wait before retry
		if !request.LastRetryTime.IsZero() && time.Since(request.LastRetryTime) < retryInterval {
			continue // Skip, not ready to retry yet
		}

		return request, filename, nil
	}

	return BlobFetchRequest{}, "", errNoMoreRequests
}

// Remove deletes a request file from the queue.
func (q *blobQueue) Remove(filename string) error {
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		q.logger.Error("Failed to delete request file", "file", filename, "error", err)
		return err
	}
	return nil
}

// UpdateRetry updates the retry metadata for a failed request.
func (q *blobQueue) UpdateRetry(filename string, request BlobFetchRequest) error {
	request.FailureCount++
	request.LastRetryTime = time.Now()

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if err = os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	return nil
}
