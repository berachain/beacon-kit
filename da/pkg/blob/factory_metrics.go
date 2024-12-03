// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blob

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/pkg/math"
)

// factoryMetrics is a struct that contains metrics for the factory.
type factoryMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
}

// newFactoryMetrics creates a new factoryMetrics.
func newFactoryMetrics(
	sink TelemetrySink,
) *factoryMetrics {
	return &factoryMetrics{
		sink: sink,
	}
}

// measureBuildSidecarsDuration measures the duration of the build sidecars.
func (fm *factoryMetrics) measureBuildSidecarsDuration(
	startTime time.Time, numSidecars math.U64,
) {
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_sidecar_duration",
		startTime,
		"num_sidecars",
		numSidecars.Base10(),
	)
}

// measureBuildKZGInclusionProofDuration measures the duration of the build KZG
// inclusion proof.
func (fm *factoryMetrics) measureBuildKZGInclusionProofDuration(
	startTime time.Time,
) {
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_kzg_inclusion_proof_duration",
		startTime,
	)
}

// measureBuildBlockBodyProofDuration measures the duration of the build block
// body proof.
func (fm *factoryMetrics) measureBuildBlockBodyProofDuration(
	startTime time.Time,
) {
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_block_body_proof_duration",
		startTime,
	)
}

// measureBuildCommitmentProofDuration measures the duration of the build
// commitment proof.
func (fm *factoryMetrics) measureBuildCommitmentProofDuration(
	startTime time.Time,
) {
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_commitment_proof_duration",
		startTime,
	)
}
