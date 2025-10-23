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

package blobreactor

import (
	"time"

	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Metric status constants for blob reactor requests.
const (
	statusSuccess         = "success"
	statusTimeout         = "timeout"
	statusPeerNotFound    = "peer_not_found"
	statusSendFailed      = "send_failed"
	statusAllPeersFailed  = "all_peers_failed"
	statusMarshalFailed   = "marshal_failed"
	statusInvalidResponse = "invalid_response"
	statusVerifyFailed    = "verification_failed"
	messageTypeRequest    = "request"
	messageTypeResponse   = "response"
)

// Metrics contains metrics for the blob reactor P2P operations.
type Metrics struct {
	RequestTotal        metrics.Counter
	RequestDuration     metrics.Histogram
	PeerAttemptsTotal   metrics.Counter
	WorkerPoolFullTotal metrics.Counter
	ActiveRequests      metrics.Gauge
	PeersAvailable      metrics.Gauge
	PeersTotal          metrics.Gauge
}

// NewMetrics returns a new Metrics instance with Prometheus metrics.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		RequestTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blobreactor",
				Name:      "request_total",
				Help:      "Total number of blob requests completed",
			},
			[]string{"status"},
		),
		RequestDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "blobreactor",
				Name:      "request_duration",
				Help:      "Time taken to complete blob requests in seconds",
				Buckets:   prometheus.ExponentialBucketsRange(0.001, 10, 10),
			},
			[]string{"status"},
		),
		PeerAttemptsTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blobreactor",
				Name:      "peer_attempts_total",
				Help:      "Total number of peer attempts for blob requests",
			},
			[]string{"status"},
		),
		WorkerPoolFullTotal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blobreactor",
				Name:      "worker_pool_full_total",
				Help:      "Number of times worker pool was full and messages were dropped",
			},
			[]string{"message_type"},
		),
		ActiveRequests: factory.NewGauge(
			metrics.GaugeOpts{
				Subsystem: "blobreactor",
				Name:      "active_requests",
				Help:      "Number of currently active blob requests",
			},
			nil,
		),
		PeersAvailable: factory.NewGauge(
			metrics.GaugeOpts{
				Subsystem: "blobreactor",
				Name:      "peers_available",
				Help:      "Number of available peers for blob requests",
			},
			nil,
		),
		PeersTotal: factory.NewGauge(
			metrics.GaugeOpts{
				Subsystem: "blobreactor",
				Name:      "peers_total",
				Help:      "Total number of connected peers",
			},
			nil,
		),
	}
}

// recordOverallRequestComplete records completion of entire blob request (may try multiple peers).
func (m *Metrics) recordOverallRequestComplete(status string, start time.Time) {
	m.RequestTotal.With("status", status).Add(1)
	m.RequestDuration.With("status", status).Observe(time.Since(start).Seconds())
}

// recordPeerAttempt records a single peer attempt with status (no duration to avoid high cardinality).
func (m *Metrics) recordPeerAttempt(status string) {
	m.PeerAttemptsTotal.With("status", status).Add(1)
}

// observeWorkerPoolFull increments counter when worker pool is full and messages are dropped.
func (m *Metrics) observeWorkerPoolFull(messageType string) {
	m.WorkerPoolFullTotal.With("message_type", messageType).Add(1)
}

// setActiveRequests sets gauge for currently active blob requests.
func (m *Metrics) setActiveRequests(count int) {
	m.ActiveRequests.Set(float64(count))
}

// setPeerPoolSize sets gauges for peer pool statistics.
func (m *Metrics) setPeerPoolSize(available, total int) {
	m.PeersAvailable.Set(float64(available))
	m.PeersTotal.Set(float64(total))
}
