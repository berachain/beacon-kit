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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package core

import (
	"github.com/berachain/beacon-kit/observability/metrics"
	"github.com/berachain/beacon-kit/primitives/math"
)

// Metrics contains metrics for the state processor.
type Metrics struct {
	BlockTxGasUsed                  metrics.Gauge
	BlockBlobGasUsed                metrics.Gauge
	PartialWithdrawalsEnqueued      metrics.Gauge
	PayloadConsensusTimestampDiff   metrics.Gauge
	DepositStakeLost                metrics.Counter
	PartialWithdrawalRequestDropped metrics.Counter
	PartialWithdrawalRequestInvalid metrics.Counter
	ValidatorNotWithdrawable        metrics.Counter
}

// NewMetrics returns a new Metrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		BlockTxGasUsed: factory.NewGauge(
			metrics.GaugeOpts{
				Name: "beacon_kit_state_block_tx_gas_used",
				Help: "Transaction gas used in the block",
			},
			nil,
		),
		BlockBlobGasUsed: factory.NewGauge(
			metrics.GaugeOpts{
				Name: "beacon_kit_state_block_blob_gas_used",
				Help: "Blob gas used in the block",
			},
			nil,
		),
		PartialWithdrawalsEnqueued: factory.NewGauge(
			metrics.GaugeOpts{
				Name: "beacon_kit_state_partial_withdrawals_enqueued",
				Help: "Number of partial withdrawals enqueued",
			},
			nil,
		),
		PayloadConsensusTimestampDiff: factory.NewGauge(
			metrics.GaugeOpts{
				Name: "beacon_kit_state_payload_consensus_timestamp_diff",
				Help: "Difference between payload timestamp and consensus timestamp",
			},
			nil,
		),
		DepositStakeLost: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_state_deposit_stake_lost",
				Help: "Number of deposits with stake lost",
			},
			nil,
		),
		PartialWithdrawalRequestDropped: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_state_partial_withdrawal_request_dropped",
				Help: "Number of partial withdrawal requests dropped",
			},
			nil,
		),
		PartialWithdrawalRequestInvalid: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_state_partial_withdrawal_request_invalid",
				Help: "Number of invalid partial withdrawal requests",
			},
			nil,
		),
		ValidatorNotWithdrawable: factory.NewCounter(
			metrics.CounterOpts{
				Name: "beacon_kit_state_validator_not_withdrawable",
				Help: "Number of validators not withdrawable",
			},
			nil,
		),
	}
}

func (m *Metrics) gaugeBlockGasUsed(txGasUsed, blobGasUsed math.U64) {
	m.BlockTxGasUsed.Set(float64(txGasUsed.Unwrap()))
	m.BlockBlobGasUsed.Set(float64(blobGasUsed.Unwrap()))
}

func (m *Metrics) gaugePartialWithdrawalsEnqueued(count int) {
	m.PartialWithdrawalsEnqueued.Set(float64(count))
}

func (m *Metrics) gaugeTimestamps(payloadTimestamp, consensusTimestamp uint64) {
	// the diff can be positive or negative depending on whether the payload
	// timestamp is ahead or behind the consensus timestamp
	diff := int64(payloadTimestamp) - int64(consensusTimestamp) // #nosec G115
	m.PayloadConsensusTimestampDiff.Set(float64(diff))
}

func (m *Metrics) incrementDepositStakeLost() {
	m.DepositStakeLost.Add(1)
}

func (m *Metrics) incrementPartialWithdrawalRequestDropped() {
	m.PartialWithdrawalRequestDropped.Add(1)
}

func (m *Metrics) incrementPartialWithdrawalRequestInvalid() {
	m.PartialWithdrawalRequestInvalid.Add(1)
}

func (m *Metrics) incrementValidatorNotWithdrawable() {
	m.ValidatorNotWithdrawable.Add(1)
}
