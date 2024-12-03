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

package validator

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/pkg/math"
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

// measureRequestBlockForProposalTime measures the time taken to run the request
// best
// block function.
func (cm *validatorMetrics) measureRequestBlockForProposalTime(
	start time.Time,
) {
	cm.sink.MeasureSince(
		"beacon_kit.validator.request_block_for_proposal_duration", start,
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

// failedToRetrievePayload increments the counter for the number of
// times the validator failed to retrieve payloads.
func (cm *validatorMetrics) failedToRetrievePayload(
	slot math.Slot, err error,
) {
	cm.sink.IncrementCounter(
		"beacon_kit.validator.failed_to_retrieve_payload",
		"slot",
		slot.Base10(),
		"error",
		err.Error(),
	)
}
