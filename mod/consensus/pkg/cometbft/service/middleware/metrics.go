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

package middleware

import (
	"time"
)

// ABCIMiddlewareMetrics is a struct that contains metrics for the chain.
type ABCIMiddlewareMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newABCIMiddlewareMetrics creates a new ABCIMiddlewareMetrics.
func newABCIMiddlewareMetrics(
	sink TelemetrySink,
) *ABCIMiddlewareMetrics {
	return &ABCIMiddlewareMetrics{
		sink: sink,
	}
}

// measurePrepareProposalDuration measures the time to prepare.
func (cm *ABCIMiddlewareMetrics) measurePrepareProposalDuration(
	start time.Time,
) {
	cm.sink.MeasureSince(
		"beacon_kit.runtime.prepare_proposal_duration", start,
	)
}

// measureProcessProposalDuration measures the time to process.
func (cm *ABCIMiddlewareMetrics) measureProcessProposalDuration(
	start time.Time,
) {
	cm.sink.MeasureSince(
		"beacon_kit.runtime.process_proposal_duration", start,
	)
}
