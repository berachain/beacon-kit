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

// Metric reason constants for blob fetcher.
const (
	expiredReasonOutsideDA  = "outside_da_period"
	expiredReasonMaxRetries = "max_retries"
)

// blobFetcherMetrics contains metrics for the blob fetcher queue and retry operations.
type blobFetcherMetrics struct {
	sink TelemetrySink
}

// newBlobFetcherMetrics creates a new blobFetcherMetrics instance.
func newBlobFetcherMetrics(sink TelemetrySink) *blobFetcherMetrics {
	return &blobFetcherMetrics{sink: sink}
}

// recordRetry increments counter when a blob request is retried after failure.
func (m *blobFetcherMetrics) recordRetry() {
	m.sink.IncrementCounter("beacon_kit.blob_fetcher.retries_total")
}

// recordRequestExpired increments counter when request expires before completion.
// Reason: "outside_da_period", "max_retries"
func (m *blobFetcherMetrics) recordRequestExpired(reason string) {
	m.sink.IncrementCounter("beacon_kit.blob_fetcher.requests_expired_total", "reason", reason)
}

// recordRequestComplete increments counter when request completes successfully.
func (m *blobFetcherMetrics) recordRequestComplete() {
	m.sink.IncrementCounter("beacon_kit.blob_fetcher.requests_completed_total")
}

// recordRequestQueued increments counter when a new request is added to queue.
func (m *blobFetcherMetrics) recordRequestQueued() {
	m.sink.IncrementCounter("beacon_kit.blob_fetcher.requests_queued_total")
}

// setQueueDepth sets the current depth of the blob fetcher queue.
func (m *blobFetcherMetrics) setQueueDepth(depth int) {
	m.sink.SetGauge("beacon_kit.blob_fetcher.queue_depth", int64(depth))
}
