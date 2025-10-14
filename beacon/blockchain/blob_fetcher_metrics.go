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
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metric reason constants for blob fetcher.
const (
	expiredReasonOutsideDA  = "outside_da_period"
	expiredReasonMaxRetries = "max_retries"
)

// blobFetcherMetrics contains native Prometheus metrics for the blob fetcher.
type blobFetcherMetrics struct {
	retriesTotal           prometheus.Counter
	requestsExpiredVec     *prometheus.CounterVec
	requestsCompletedTotal prometheus.Counter
	requestsQueuedTotal    prometheus.Counter
	queueDepth             prometheus.Gauge
}

// newBlobFetcherMetrics creates a new native Prometheus metrics instance.
func newBlobFetcherMetrics() *blobFetcherMetrics {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	// Use promauto to automatically register metrics with prometheus.DefaultRegisterer
	return &blobFetcherMetrics{
		retriesTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "beacond_node_beacon_kit_blob_fetcher_retries_total",
			Help: "Total number of blob fetch retry attempts",
			ConstLabels: prometheus.Labels{
				"host": hostname,
			},
		}),
		requestsExpiredVec: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "beacond_node_beacon_kit_blob_fetcher_requests_expired_total",
				Help: "Total number of expired blob fetch requests",
				ConstLabels: prometheus.Labels{
					"host": hostname,
				},
			},
			[]string{"reason"},
		),
		requestsCompletedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "beacond_node_beacon_kit_blob_fetcher_requests_completed_total",
			Help: "Total number of successfully completed blob fetch requests",
			ConstLabels: prometheus.Labels{
				"host": hostname,
			},
		}),
		requestsQueuedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "beacond_node_beacon_kit_blob_fetcher_requests_queued_total",
			Help: "Total number of blob fetch requests added to queue",
			ConstLabels: prometheus.Labels{
				"host": hostname,
			},
		}),
		queueDepth: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "beacond_node_beacon_kit_blob_fetcher_queue_depth",
			Help: "Current depth of the blob fetcher queue",
			ConstLabels: prometheus.Labels{
				"host": hostname,
			},
		}),
	}
}

// recordRetry increments counter when a blob request is retried after failure.
func (m *blobFetcherMetrics) recordRetry() {
	m.retriesTotal.Inc()
}

// recordRequestExpired increments counter when request expires before completion.
// Reason: "outside_da_period", "max_retries"
func (m *blobFetcherMetrics) recordRequestExpired(reason string) {
	m.requestsExpiredVec.WithLabelValues(reason).Inc()
}

// recordRequestComplete increments counter when request completes successfully.
func (m *blobFetcherMetrics) recordRequestComplete() {
	m.requestsCompletedTotal.Inc()
}

// recordRequestQueued increments counter when a new request is added to queue.
func (m *blobFetcherMetrics) recordRequestQueued() {
	m.requestsQueuedTotal.Inc()
}

// setQueueDepth sets the current depth of the blob fetcher queue.
func (m *blobFetcherMetrics) setQueueDepth(depth int) {
	m.queueDepth.Set(float64(depth))
}
