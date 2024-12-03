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

// measureVerifySidecarsDuration measures the duration of the blob verification.
func (vm *verifierMetrics) measureVerifySidecarsDuration(
	startTime time.Time,
	numSidecars math.U64,
	kzgImplementation string,
) {
	vm.sink.MeasureSince(
		"beacon_kit.da.blob.verifier.verify_blobs_duration",
		startTime,
		"num_sidecars",
		numSidecars.Base10(),
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
		numSidecars.Base10(),
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
		numSidecars.Base10(),
		"kzg_implementation",
		kzgImplementation,
	)
}
