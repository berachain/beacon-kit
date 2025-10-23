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

package state

import (
	"github.com/berachain/beacon-kit/observability/metrics"
)

// Metrics is a struct that contains metrics for the StateDB.
type Metrics struct {
	// PartialWithdrawalRequestInvalid tracks invalid partial withdrawal requests
	PartialWithdrawalRequestInvalid metrics.Counter

	// ExcessStakePartialWithdrawal tracks withdrawals due to excess stake
	ExcessStakePartialWithdrawal metrics.Counter
}

// NewMetrics returns a new Metrics instance.
// Metric names are kept identical to cosmos-sdk/telemetry output for Grafana compatibility.
func NewMetrics(factory metrics.Factory) *Metrics {
	return &Metrics{
		PartialWithdrawalRequestInvalid: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "statedb",
				Name:      "partial_withdrawal_request_invalid",
				Help:      "Number of invalid partial withdrawal requests",
			},
			nil,
		),
		ExcessStakePartialWithdrawal: factory.NewCounter(
			metrics.CounterOpts{
				Subsystem: "statedb",
				Name:      "excess_stake_partial_withdrawal",
				Help:      "Number of withdrawals created due to validator stake exceeding MaxEffectiveBalance",
			},
			nil,
		),
	}
}

// incrementPartialWithdrawalRequestInvalid increments the counter for invalid
// partial withdrawal requests.
func (s *StateDB) incrementPartialWithdrawalRequestInvalid() {
	s.metrics.PartialWithdrawalRequestInvalid.Add(1)
}

// incrementExcessStakePartialWithdrawal increments the counter when a withdrawal
// is created because a validator's stake went over the MaxEffectiveBalance.
func (s *StateDB) incrementExcessStakePartialWithdrawal() {
	s.metrics.ExcessStakePartialWithdrawal.Add(1)
}
