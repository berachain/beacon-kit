// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package blob

import (
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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

// measureBuildSidecarDuration measures the duration of the build sidecar.
func (fm *factoryMetrics) measureBuildSidecarsDuration(
	startTime time.Time, numSidecars math.U64,
) {
	fm.sink.MeasureSince(
		"beacon_kit.da.blob.factory.build_sidecar_duration",
		startTime,
		"num_sidecars",
		string(numSidecars.String()),
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
