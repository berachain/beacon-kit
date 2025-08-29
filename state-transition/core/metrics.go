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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core

import (
	"github.com/berachain/beacon-kit/primitives/math"
)

type stateProcessorMetrics struct {
	// sink is the sink for the metrics.
	sink TelemetrySink
}

// newStateProcessorMetrics creates a new stateProcessorMetrics.
func newStateProcessorMetrics(sink TelemetrySink) *stateProcessorMetrics {
	return &stateProcessorMetrics{
		sink: sink,
	}
}

func (s *stateProcessorMetrics) gaugeBlockGasUsed(blockNumber, txGasUsed, blobGasUsed math.U64) {
	blockNumberStr := blockNumber.Base10()
	s.sink.SetGauge(
		"beacon_kit.state.block_tx_gas_used",
		int64(txGasUsed.Unwrap()), // #nosec G115
		"block_number",
		blockNumberStr,
	)
	s.sink.SetGauge(
		"beacon_kit.state.block_blob_gas_used",
		int64(blobGasUsed.Unwrap()), // #nosec G115
		"block_number",
		blockNumberStr,
	)
}

func (s *stateProcessorMetrics) gaugePartialWithdrawalsEnqueued(count int) {
	s.sink.SetGauge("beacon_kit.state.partial_withdrawals_enqueued", int64(count))
}

func (s *stateProcessorMetrics) gaugeTimestamps(payloadTimestamp, consensusTimestamp uint64) {
	// the diff can be positive or negative depending on whether the payload
	// timestamp is ahead or behind the consensus timestamp
	diff := int64(payloadTimestamp) - int64(consensusTimestamp) // #nosec G115
	s.sink.SetGauge("beacon_kit.state.payload_consensus_timestamp_diff", diff)
}

func (s *stateProcessorMetrics) incrementDepositStakeLost() {
	s.sink.IncrementCounter("beacon_kit.state.deposit_stake_lost")
}

func (s *stateProcessorMetrics) incrementPartialWithdrawalRequestDropped() {
	s.sink.IncrementCounter("beacon_kit.state.partial_withdrawal_request_dropped")
}

func (s *stateProcessorMetrics) incrementPartialWithdrawalRequestInvalid() {
	s.sink.IncrementCounter("beacon_kit.state.partial_withdrawal_request_invalid")
}

func (s *stateProcessorMetrics) incrementValidatorNotWithdrawable() {
	s.sink.IncrementCounter("beacon_kit.state.validator_not_withdrawable")
}
