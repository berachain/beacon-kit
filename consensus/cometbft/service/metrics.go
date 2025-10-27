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
)

// Metrics holds metrics for the CometBFT service.
type Metrics struct {
	// QueryCount tracks the number of ABCI queries received, labeled by query path.
	QueryCount metrics.Counter

	// QueryDuration tracks the time taken to process ABCI queries, labeled by query path.
	QueryDuration metrics.Summary

	// PrepareProposalDuration tracks the time taken to prepare a proposal.
	PrepareProposalDuration metrics.Summary

	// ProcessProposalDuration tracks the time taken to process a proposal.
	ProcessProposalDuration metrics.Summary
}

// NewMetrics creates a new Metrics instance using the provided factory. The factory determines
// whether real Prometheus metrics or no-op metrics are created.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		QueryCount: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_comet_query_count",
				Help: "Total number of ABCI queries received",
			},
			[]string{"path"},
		),
		QueryDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_comet_query_duration",
				Help:       "Time taken to process ABCI queries in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"path"},
		),
		PrepareProposalDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_runtime_prepare_proposal_duration",
				Help:       "Time taken to prepare a proposal in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		ProcessProposalDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_runtime_process_proposal_duration",
				Help:       "Time taken to process a proposal in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
	}
}

// measureQueryDuration is a helper to measure query duration.
func (m *Metrics) measureQueryDuration(start time.Time, path string) {
	m.QueryDuration.With("path", path).Observe(float64(time.Since(start).Milliseconds()))
}

// measurePrepareProposalDuration is a helper to measure prepare proposal duration.
func (m *Metrics) measurePrepareProposalDuration(start time.Time) {
	m.PrepareProposalDuration.Observe(float64(time.Since(start).Milliseconds()))
}

// measureProcessProposalDuration is a helper to measure process proposal duration.
func (m *Metrics) measureProcessProposalDuration(start time.Time) {
	m.ProcessProposalDuration.Observe(float64(time.Since(start).Milliseconds()))
}
