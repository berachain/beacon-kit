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

package blob

import (
	"time"

	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/primitives/math"
)

// ProcessorMetrics is a struct that contains metrics for the blob processor.
type ProcessorMetrics struct {
	VerifyBlobsDuration metrics.Summary
	ProcessBlobDuration metrics.Summary
}

// NewProcessorMetrics returns a new ProcessorMetrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewProcessorMetrics(factory metrics.Factory) *ProcessorMetrics {
	return &ProcessorMetrics{
		VerifyBlobsDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_processor_verify_blobs_duration",
				Help:       "Time taken to verify blob sidecars in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars"},
		),
		ProcessBlobDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_processor_process_blob_duration",
				Help:       "Time taken to process blob sidecars in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars"},
		),
	}
}

// measureVerifySidecarsDuration measures the duration of the blob verification.
func (m *ProcessorMetrics) measureVerifySidecarsDuration(
	startTime time.Time,
	numSidecars math.U64,
) {
	m.VerifyBlobsDuration.With("num_sidecars", numSidecars.Base10()).Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureProcessSidecarsDuration measures the duration of the blob processing.
func (m *ProcessorMetrics) measureProcessSidecarsDuration(
	startTime time.Time,
	numSidecars math.U64,
) {
	m.ProcessBlobDuration.With("num_sidecars", numSidecars.Base10()).Observe(float64(time.Since(startTime).Milliseconds()))
}
