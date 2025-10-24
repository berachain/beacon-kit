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

// VerifierMetrics is a struct that contains metrics for the blob verifier.
type VerifierMetrics struct {
	VerifyBlobsDuration           metrics.Summary
	VerifyInclusionProofsDuration metrics.Summary
	VerifyKZGProofsDuration       metrics.Summary
}

// NewVerifierMetrics returns a new VerifierMetrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewVerifierMetrics(factory metrics.Factory) *VerifierMetrics {
	return &VerifierMetrics{
		VerifyBlobsDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_verifier_verify_blobs_duration",
				Help:       "Time taken to verify blobs in seconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars", "kzg_implementation"},
		),
		VerifyInclusionProofsDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_verifier_verify_inclusion_proofs_duration",
				Help:       "Time taken to verify inclusion proofs in seconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars"},
		),
		VerifyKZGProofsDuration: factory.NewSummary(
			metrics.SummaryOpts{
				Name:       "beacon_kit_da_blob_verifier_verify_kzg_proofs_duration",
				Help:       "Time taken to verify KZG proofs in seconds",
				Objectives: metrics.QuantilesP50P90P99,
			},
			[]string{"num_sidecars", "kzg_implementation"},
		),
	}
}

// measureVerifySidecarsDuration measures the duration of the blob verification.
func (m *VerifierMetrics) measureVerifySidecarsDuration(
	startTime time.Time,
	numSidecars math.U64,
	kzgImplementation string,
) {
	m.VerifyBlobsDuration.With(
		"num_sidecars", numSidecars.Base10(),
		"kzg_implementation", kzgImplementation,
	).Observe(time.Since(startTime).Seconds())
}

// measureVerifyInclusionProofsDuration measures the duration of the inclusion proofs verification.
func (m *VerifierMetrics) measureVerifyInclusionProofsDuration(
	startTime time.Time,
	numSidecars math.U64,
) {
	m.VerifyInclusionProofsDuration.With("num_sidecars", numSidecars.Base10()).Observe(time.Since(startTime).Seconds())
}

// measureVerifyKZGProofsDuration measures the duration of the KZG proofs verification.
func (m *VerifierMetrics) measureVerifyKZGProofsDuration(
	startTime time.Time,
	numSidecars math.U64,
	kzgImplementation string,
) {
	m.VerifyKZGProofsDuration.With(
		"num_sidecars", numSidecars.Base10(),
		"kzg_implementation", kzgImplementation,
	).Observe(time.Since(startTime).Seconds())
}
