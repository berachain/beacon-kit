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

package blockchain

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/math"
)

// chainMetrics is a struct that contains metrics for the chain.
type chainMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newChainMetrics creates a new chainMetrics.
func newChainMetrics(
	sink TelemetrySink,
) *chainMetrics {
	return &chainMetrics{
		sink: sink,
	}
}

// measureStateTransitionDuration measures the time to process
// the state transition for a block.
func (cm *chainMetrics) measureStateTransitionDuration(
	start time.Time,
) {
	cm.sink.MeasureSince(
		"beacon_kit.beacon.blockchain.state_transition_duration",
		start,
	)
}

// markRebuildPayloadForRejectedBlockSuccess increments the counter for the
// number of times
// the validator successfully rebuilt the payload for a rejected block.
func (cm *chainMetrics) markRebuildPayloadForRejectedBlockSuccess(
	slot math.Slot,
) {
	cm.sink.IncrementCounter(
		"beacon_kit.blockchain.rebuild_payload_for_rejected_block_success",
		"slot",
		slot.Base10(),
	)
}

// markRebuildPayloadForRejectedBlockFailure increments the counter for the
// number of times
// the validator failed to build an optimistic payload due to a failure.
func (cm *chainMetrics) markRebuildPayloadForRejectedBlockFailure(
	slot math.Slot,
	err error,
) {
	cm.sink.IncrementCounter(
		"beacon_kit.blockchain.rebuild_payload_for_rejected_block_failure",
		"slot",
		slot.Base10(),
		"error",
		err.Error(),
	)
}

// markOptimisticPayloadBuildSuccess increments the counter for the number of
// times
// the validator successfully built an optimistic payload.
func (cm *chainMetrics) markOptimisticPayloadBuildSuccess(slot math.Slot) {
	cm.sink.IncrementCounter(
		"beacon_kit.blockchain.optimistic_payload_build_success",
		"slot",
		slot.Base10(),
	)
}

// markOptimisticPayloadBuildFailure increments the counter for the number of
// times
// the validator failed to build an optimistic payload.
func (cm *chainMetrics) markOptimisticPayloadBuildFailure(
	slot math.Slot,
	err error,
) {
	cm.sink.IncrementCounter(
		"beacon_kit.blockchain.optimistic_payload_build_failure",
		"slot",
		slot.Base10(),
		"error",
		err.Error(),
	)
}

// measureStateRootVerificationTime measures the time taken to verify the state
// root of a block.
// It records the duration from the provided start time to the current time.
func (cm *chainMetrics) measureStateRootVerificationTime(start time.Time) {
	cm.sink.MeasureSince(
		"beacon_kit.blockchain.state_root_verification_duration", start,
	)
}
