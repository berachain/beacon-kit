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
	"time"

	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics is a struct that contains metrics for the blockchain service.
type Metrics struct {
	StateTransitionDuration               metrics.Histogram
	RebuildPayloadForRejectedBlockSuccess metrics.Counter
	RebuildPayloadForRejectedBlockFailure metrics.Counter
	OptimisticPayloadBuildSuccess         metrics.Counter
	OptimisticPayloadBuildFailure         metrics.Counter
	StateRootVerificationDuration         metrics.Histogram
	FailedToGetBlockLogs                  metrics.Counter
	FailedToEnqueueDeposits               metrics.Counter
}

// NewPrometheusMetrics returns a new Metrics instance with Prometheus metrics.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		StateTransitionDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "beacon_blockchain",
				Name:      "state_transition_duration",
				Help:      "Time taken to process state transition in seconds",
				Buckets:   prometheus.ExponentialBucketsRange(0.001, 10, 10),
			},
			nil,
		),
		RebuildPayloadForRejectedBlockSuccess: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blockchain",
				Name:      "rebuild_payload_for_rejected_block_success",
				Help:      "Number of successful payload rebuilds for rejected blocks",
			},
			nil,
		),
		RebuildPayloadForRejectedBlockFailure: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blockchain",
				Name:      "rebuild_payload_for_rejected_block_failure",
				Help:      "Number of failed payload rebuilds for rejected blocks",
			},
			nil,
		),
		OptimisticPayloadBuildSuccess: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blockchain",
				Name:      "optimistic_payload_build_success",
				Help:      "Number of successful optimistic payload builds",
			},
			nil,
		),
		OptimisticPayloadBuildFailure: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "blockchain",
				Name:      "optimistic_payload_build_failure",
				Help:      "Number of failed optimistic payload builds",
			},
			nil,
		),
		StateRootVerificationDuration: factory.NewHistogram(
			metrics.HistogramOpts{
				Subsystem: "blockchain",
				Name:      "state_root_verification_duration",
				Help:      "Time taken to verify state root in seconds",
				Buckets:   prometheus.ExponentialBucketsRange(0.001, 10, 10),
			},
			nil,
		),
		FailedToGetBlockLogs: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "execution_deposit",
				Name:      "failed_to_get_block_logs",
				Help:      "Number of times failed to read deposits from execution layer block logs",
			},
			nil,
		),
		FailedToEnqueueDeposits: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "execution_deposit",
				Name:      "failed_to_enqueue_deposits",
				Help:      "Number of times failed to enqueue deposits to storage",
			},
			nil,
		),
	}
}

// measureStateTransitionDuration measures the time to process the state transition for a block.
func (m *Metrics) measureStateTransitionDuration(start time.Time) {
	m.StateTransitionDuration.Observe(time.Since(start).Seconds())
}

// markRebuildPayloadForRejectedBlockSuccess increments the counter for the number of times
// the validator successfully rebuilt the payload for a rejected block.
func (m *Metrics) markRebuildPayloadForRejectedBlockSuccess() {
	m.RebuildPayloadForRejectedBlockSuccess.Add(1)
}

// markRebuildPayloadForRejectedBlockFailure increments the counter for the number of times
// the validator failed to build an optimistic payload due to a failure.
func (m *Metrics) markRebuildPayloadForRejectedBlockFailure() {
	m.RebuildPayloadForRejectedBlockFailure.Add(1)
}

// markOptimisticPayloadBuildSuccess increments the counter for the number of times
// the validator successfully built an optimistic payload.
func (m *Metrics) markOptimisticPayloadBuildSuccess() {
	m.OptimisticPayloadBuildSuccess.Add(1)
}

// markOptimisticPayloadBuildFailure increments the counter for the number of times
// the validator failed to build an optimistic payload.
func (m *Metrics) markOptimisticPayloadBuildFailure() {
	m.OptimisticPayloadBuildFailure.Add(1)
}

// TODO: remove once state caching is activated
// measureStateRootVerificationTime measures the time taken to verify the state root of a block.
// It records the duration from the provided start time to the current time.
func (m *Metrics) measureStateRootVerificationTime(start time.Time) {
	m.StateRootVerificationDuration.Observe(time.Since(start).Seconds())
}
