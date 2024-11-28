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

package core

type stateProcessorMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newStateProcessorMetrics creates a new stateProcessorMetrics.
func newStateProcessorMetrics(
	sink TelemetrySink,
) *stateProcessorMetrics {
	return &stateProcessorMetrics{
		sink: sink,
	}
}

func (s *stateProcessorMetrics) gaugeTimestamps(
	payloadTimestamp uint64,
	consensusTimestamp uint64,
) {
	// the diff can be positive or negative depending on whether the payload
	// timestamp is ahead or behind the consensus timestamp
	diff := int64(payloadTimestamp) - int64(consensusTimestamp) //#nosec:G701
	s.sink.SetGauge("beacon_kit.state.payload_consensus_timestamp_diff", diff)
}
