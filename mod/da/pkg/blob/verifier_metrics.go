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

// verifierMetrics is a struct that contains metrics for the verifier.
type verifierMetrics struct {
	// TelemetrySink is the sink for the metrics.
	sink TelemetrySink
}

// newVerifierMetrics creates a new verifierMetrics.
func newVerifierMetrics(
	sink TelemetrySink,
) *verifierMetrics {
	return &verifierMetrics{
		sink: sink,
	}
}

// measureVerifyBlobsDuration measures the duration of the blob verification.
func (vm *verifierMetrics) measureVerifyBlobsDuration(
	startTime time.Time,
	numSidecars math.U64,
	kzgImplementation string,
) {
	vm.sink.MeasureSince(
		"beacon_kit.da.blob.verifier.verify_blobs_duration",
		startTime,
		"num_sidecars",
		string(numSidecars.String()),
		"kzg_implementation",
		kzgImplementation,
	)
}

// measureVerifyInclusionProofsDuration measures the duration of the inclusion
// proofs verification.
func (vm *verifierMetrics) measureVerifyInclusionProofsDuration(
	startTime time.Time,
	numSidecars math.U64,
) {
	vm.sink.MeasureSince(
		"beacon_kit.da.blob.verifier.verify_inclusion_proofs_duration",
		startTime,
		"num_sidecars",
		string(numSidecars.String()),
	)
}

// measureVerifyKZGProofsDuration measures the duration of the KZG proofs
// verification.
func (vm *verifierMetrics) measureVerifyKZGProofsDuration(
	startTime time.Time,
	numSidecars math.U64,
	kzgImplementation string,
) {
	vm.sink.MeasureSince(
		"beacon_kit.da.blob.verifier.verify_kzg_proofs_duration",
		startTime,
		"num_sidecars",
		string(numSidecars.String()),
		"kzg_implementation",
		kzgImplementation,
	)
}
