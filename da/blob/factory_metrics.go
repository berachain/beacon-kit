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

// FactoryMetrics is a struct that contains metrics for the sidecar factory.
type FactoryMetrics struct {
	BuildSidecarDuration           metrics.Summary
	BuildKZGInclusionProofDuration metrics.Summary
	BuildBlockBodyProofDuration    metrics.Summary
	BuildCommitmentProofDuration   metrics.Summary
}

// NewFactoryMetrics returns a new FactoryMetrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewFactoryMetrics(factory metrics.Factory) *FactoryMetrics {
	return &FactoryMetrics{
		BuildSidecarDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_factory_build_sidecar_duration",
				Help:       "Time taken to build blob sidecars in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars"},
		),
		BuildKZGInclusionProofDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_factory_build_kzg_inclusion_proof_duration",
				Help:       "Time taken to build KZG inclusion proof in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		BuildBlockBodyProofDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_factory_build_block_body_proof_duration",
				Help:       "Time taken to build block body proof in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
		BuildCommitmentProofDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_factory_build_commitment_proof_duration",
				Help:       "Time taken to build commitment proof in milliseconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			nil,
		),
	}
}

// measureBuildSidecarsDuration measures the duration of the build sidecars.
func (m *FactoryMetrics) measureBuildSidecarsDuration(
	startTime time.Time, numSidecars math.U64,
) {
	m.BuildSidecarDuration.With("num_sidecars", numSidecars.Base10()).Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureBuildKZGInclusionProofDuration measures the duration of the build KZG inclusion proof.
func (m *FactoryMetrics) measureBuildKZGInclusionProofDuration(
	startTime time.Time,
) {
	m.BuildKZGInclusionProofDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureBuildBlockBodyProofDuration measures the duration of the build block body proof.
func (m *FactoryMetrics) measureBuildBlockBodyProofDuration(
	startTime time.Time,
) {
	m.BuildBlockBodyProofDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}

// measureBuildCommitmentProofDuration measures the duration of the build commitment proof.
func (m *FactoryMetrics) measureBuildCommitmentProofDuration(
	startTime time.Time,
) {
	m.BuildCommitmentProofDuration.Observe(float64(time.Since(startTime).Milliseconds()))
}
