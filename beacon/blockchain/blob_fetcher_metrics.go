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
	"github.com/berachain/beacon-kit/observability/metrics"
)

// Metric reason constants for blob fetcher.
const (
	expiredReasonOutsideDA  = "outside_da_period"
	expiredReasonMaxRetries = "max_retries"
)

// BlobFetcherMetrics contains metrics for the blob fetcher queue and retry operations.
type BlobFetcherMetrics struct {
	RetriesTotal           metrics.Counter
	RequestsExpiredTotal   metrics.Counter
	RequestsCompletedTotal metrics.Counter
	RequestsQueuedTotal    metrics.Counter
	QueueDepth             metrics.Gauge
}

// NewBlobFetcherMetrics returns a new BlobFetcherMetrics instance with metrics from the provided factory.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewBlobFetcherMetrics(factory metrics.Factory) *BlobFetcherMetrics {
	return &BlobFetcherMetrics{
		RetriesTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blob_fetcher",
				Name:      "retries_total",
				Help:      "Number of times a blob request was retried after failure",
			},
			nil,
		),
		RequestsExpiredTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blob_fetcher",
				Name:      "requests_expired_total",
				Help:      "Number of blob fetch requests that expired before completion",
			},
			[]string{"reason"},
		),
		RequestsCompletedTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blob_fetcher",
				Name:      "requests_completed_total",
				Help:      "Number of blob fetch requests that completed successfully",
			},
			nil,
		),
		RequestsQueuedTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blob_fetcher",
				Name:      "requests_queued_total",
				Help:      "Number of new blob fetch requests added to the queue",
			},
			nil,
		),
		QueueDepth: factory.NewGauge(
			metrics.GaugeOpts{
				Subsystem: "blob_fetcher",
				Name:      "queue_depth",
				Help:      "Current depth of the blob fetcher queue",
			},
			nil,
		),
	}
}

// recordRetry increments counter when a blob request is retried after failure.
func (m *BlobFetcherMetrics) recordRetry() {
	m.RetriesTotal.Add(1)
}

// recordRequestExpired increments counter when request expires before completion.
// Reason: "outside_da_period", "max_retries"
func (m *BlobFetcherMetrics) recordRequestExpired(reason string) {
	m.RequestsExpiredTotal.With("reason", reason).Add(1)
}

// recordRequestComplete increments counter when request completes successfully.
func (m *BlobFetcherMetrics) recordRequestComplete() {
	m.RequestsCompletedTotal.Add(1)
}

// recordRequestQueued increments counter when a new request is added to queue.
func (m *BlobFetcherMetrics) recordRequestQueued() {
	m.RequestsQueuedTotal.Add(1)
}

// setQueueDepth sets the current depth of the blob fetcher queue.
func (m *BlobFetcherMetrics) setQueueDepth(depth int) {
	m.QueueDepth.Set(float64(depth))
}
