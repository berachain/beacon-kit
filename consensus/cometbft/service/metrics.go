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

package cometbft

import (
	"time"

	"github.com/berachain/beacon-kit/observability/metrics"
	prominternal "github.com/prometheus/client_golang/prometheus"
)

// Metrics holds metrics for the CometBFT service.
//
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
type Metrics struct {
	// QueryCount tracks the number of ABCI queries received, labeled by query path.
	QueryCount metrics.Counter

	// QueryDuration tracks the time taken to process ABCI queries, labeled by query path.
	QueryDuration metrics.Histogram

	// PrepareProposalDuration tracks the time taken to prepare a proposal.
	PrepareProposalDuration metrics.Histogram

	// ProcessProposalDuration tracks the time taken to process a proposal.
	ProcessProposalDuration metrics.Histogram
}

// NewMetrics creates a new Metrics instance using the provided factory.
// The factory determines whether real Prometheus metrics or no-op metrics are created.
//
//nolint:mnd // magic numbers are histogram bucket ranges for timing metrics
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		QueryCount: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "comet",
				Name:      "query_count",
				Help:      "Total number of ABCI queries received",
			},
			[]string{"path"},
		),
		QueryDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "comet",
				Name:      "query_duration",
				Help:      "Time taken to process ABCI queries in seconds",
				Buckets:   prominternal.ExponentialBucketsRange(0.001, 10, 10),
			},
			[]string{"path"},
		),
		PrepareProposalDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "runtime",
				Name:      "prepare_proposal_duration",
				Help:      "Time taken to prepare a proposal in seconds",
				Buckets:   prominternal.ExponentialBucketsRange(0.001, 10, 10),
			},
			nil,
		),
		ProcessProposalDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "runtime",
				Name:      "process_proposal_duration",
				Help:      "Time taken to process a proposal in seconds",
				Buckets:   prominternal.ExponentialBucketsRange(0.001, 10, 10),
			},
			nil,
		),
	}
}

// measureQueryDuration is a helper to measure query duration.
func (m *Metrics) measureQueryDuration(start time.Time, path string) {
	m.QueryDuration.With("path", path).Observe(time.Since(start).Seconds())
}

// measurePrepareProposalDuration is a helper to measure prepare proposal duration.
func (m *Metrics) measurePrepareProposalDuration(start time.Time) {
	m.PrepareProposalDuration.Observe(time.Since(start).Seconds())
}

// measureProcessProposalDuration is a helper to measure process proposal duration.
func (m *Metrics) measureProcessProposalDuration(start time.Time) {
	m.ProcessProposalDuration.Observe(time.Since(start).Seconds())
}
