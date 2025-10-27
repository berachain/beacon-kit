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

package validator

import (
	"time"

	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Metrics is a struct that contains metrics for the validator service.
type Metrics struct {
	// RequestBlockForProposalDuration tracks time to request block for proposal
	// Using Summary for backward compatibility with cosmos-sdk/telemetry.
	RequestBlockForProposalDuration metrics.Summary

	// StateRootComputationDuration tracks time to compute state root
	// Using Summary for backward compatibility with cosmos-sdk/telemetry.
	StateRootComputationDuration metrics.Summary

	// FailedToRetrievePayload tracks failed payload retrievals
	FailedToRetrievePayload metrics.Counter
}

// NewMetrics returns a new Metrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		RequestBlockForProposalDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_validator_request_block_for_proposal_duration",
				Help:       "Time taken to request block for proposal in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		StateRootComputationDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_validator_state_root_computation_duration",
				Help:       "Time taken to compute state root in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		FailedToRetrievePayload: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_validator_failed_to_retrieve_payload",
				Help: "Number of times validator failed to retrieve payload",
			},
			[]string{"slot", "error"},
		),
	}
}

// measureRequestBlockForProposalTime measures the time taken to request block for proposal.
func (m *Metrics) measureRequestBlockForProposalTime(start time.Time) {
	m.RequestBlockForProposalDuration.Observe(float64(time.Since(start).Milliseconds()))
}

// measureStateRootComputationTime measures the time taken to compute the state root of a block.
func (m *Metrics) measureStateRootComputationTime(start time.Time) {
	m.StateRootComputationDuration.Observe(float64(time.Since(start).Milliseconds()))
}

// failedToRetrievePayload increments the counter for the number of times the validator
// failed to retrieve payloads.
func (m *Metrics) failedToRetrievePayload(slot math.Slot, err error) {
	m.FailedToRetrievePayload.With("slot", slot.Base10(), "error", err.Error()).Add(1)
}
