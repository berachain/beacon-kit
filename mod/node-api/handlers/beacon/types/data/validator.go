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

package data

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
)

type ValidatorData[
	ValidatorT types.Validator[WithdrawalCredentialsT],
	WithdrawalCredentialsT types.WithdrawalCredentials,
] struct {
	ValidatorBalanceData
	Status    string     `json:"status"`
	Validator ValidatorT `json:"validator"`
}

type validatorJSON struct {
	PublicKey                  string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

type responseJSON struct {
	Index     string        `json:"index"`
	Balance   string        `json:"balance"`
	Status    string        `json:"status"`
	Validator validatorJSON `json:"validator"`
}

func (vd ValidatorData[
	ValidatorT, WithdrawalCredentialsT,
]) MarshalJSON() ([]byte, error) {
	withdrawalCredentials := vd.Validator.GetWithdrawalCredentials()
	withdrawalCredentialsBytes := withdrawalCredentials.Bytes()

	return json.Marshal(responseJSON{
		Index:   strconv.FormatUint(vd.Index, 10),
		Balance: strconv.FormatUint(vd.Balance, 10),
		Status:  vd.Status,
		Validator: validatorJSON{
			PublicKey: vd.Validator.GetPubkey().String(),
			WithdrawalCredentials: "0x" + hex.EncodeToString(
				withdrawalCredentialsBytes,
			),
			EffectiveBalance: strconv.FormatUint(
				vd.Validator.GetEffectiveBalance().Unwrap(), 10,
			),
			Slashed: vd.Validator.IsSlashed(),
			ActivationEligibilityEpoch: strconv.FormatUint(
				vd.Validator.GetActivationEligibilityEpoch().Unwrap(), 10,
			),
			ActivationEpoch: strconv.FormatUint(
				vd.Validator.GetActivationEpoch().Unwrap(), 10,
			),
			ExitEpoch: strconv.FormatUint(
				vd.Validator.GetExitEpoch().Unwrap(), 10,
			),
			WithdrawableEpoch: strconv.FormatUint(
				vd.Validator.GetWithdrawableEpoch().Unwrap(), 10,
			),
		},
	})
}
