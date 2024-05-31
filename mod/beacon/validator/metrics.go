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

package validator

import (
	"time"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// validatorMetrics is a struct that contains metrics for the chain.
type validatorMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newValidatorMetrics creates a new validatorMetrics.
func newValidatorMetrics(
	sink TelemetrySink,
) *validatorMetrics {
	return &validatorMetrics{
		sink: sink,
	}
}

// measureRequestBestBlockTime measures the time taken to run the request best
// block function.
func (cm *validatorMetrics) measureRequestBestBlockTime(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.validator.request_best_block_duration", start,
	)
}

// measureStateRootVerificationTime measures the time taken to verify the state
// root of a block.
// It records the duration from the provided start time to the current time.
func (cm *validatorMetrics) measureStateRootVerificationTime(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.validator.state_root_verification_duration", start,
	)
}

// measureStateRootComputationTime measures the time taken to compute the state
// root of a block.
// It records the duration from the provided start time to the current time.
func (cm *validatorMetrics) measureStateRootComputationTime(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.validator.state_root_computation_duration", start,
	)
}

// failedToRetrieveOptimisticPayload increments the counter for the number of
// times the
// validator failed to retrieve payloads.
func (cm *validatorMetrics) failedToRetrieveOptimisticPayload(
	slot math.Slot, blkRoot primitives.Root, err error,
) {
	cm.sink.IncrementCounter(
		"beacon_kit.validator.failed_to_retrieve_optimistic_payload",
		"slot",
		string(slot.String()),
		"block_root",
		blkRoot.String(),
		"error",
		err.Error(),
	)
}
