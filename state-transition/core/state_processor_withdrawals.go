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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
)

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
// NOTE: Modified from the Ethereum 2.0 specification to support EVM inflation:
// 1. The first withdrawal MUST be a fixed EVM inflation withdrawal
// 2. Subsequent withdrawals (if any) are processed as validator withdrawals
// 3. This modification reduces the maximum validator withdrawals per block by one.
//

func (sp *StateProcessor) processWithdrawals(
	st *state.StateDB, blk *ctypes.BeaconBlock,
) error {
	// Dequeue and verify the logs.
	var (
		body               = blk.GetBody()
		payload            = body.GetExecutionPayload()
		payloadWithdrawals = payload.GetWithdrawals()
	)

	// Get the expected withdrawals.
	expectedWithdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		return err
	}

	// Common validations
	if len(expectedWithdrawals) != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Enforce that first withdrawal is EVM inflation
	if len(payloadWithdrawals) == 0 {
		return ErrZeroWithdrawals
	}
	if !payloadWithdrawals[0].Equals(st.EVMInflationWithdrawal(blk.GetSlot())) {
		return ErrFirstWithdrawalNotEVMInflation
	}
	numWithdrawals := len(expectedWithdrawals)

	// Process all subsequent validator withdrawals.
	for i := 1; i < numWithdrawals; i++ {
		// Ensure the withdrawals match the local state.
		if !expectedWithdrawals[i].Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrWithdrawalMismatch,
				"withdrawal at index %d does not match expected %s, got %s",
				i,
				spew.Sdump(expectedWithdrawals[i]),
				spew.Sdump(payloadWithdrawals[i]),
			)
		}

		if err = st.DecreaseBalance(
			expectedWithdrawals[i].GetValidatorIndex(), expectedWithdrawals[i].GetAmount(),
		); err != nil {
			return err
		}
	}

	if numWithdrawals > 1 {
		if err = st.SetNextWithdrawalIndex(
			(expectedWithdrawals[numWithdrawals-1].GetIndex() + 1).Unwrap(),
		); err != nil {
			return err
		}
	}

	totalValidators, err := st.GetTotalValidators()
	if err != nil {
		return err
	}

	// Update the next validator index to start the next withdrawal sweep.
	var nextValidatorIndex math.ValidatorIndex

	// #nosec G115 -- won't overflow in practice.
	if numWithdrawals == int(sp.cs.MaxWithdrawalsPerPayload()) {
		// Next sweep starts after the latest withdrawal's validator index.
		nextValidatorIndex = (expectedWithdrawals[numWithdrawals-1].
			GetValidatorIndex() + 1) % math.ValidatorIndex(totalValidators)
	} else {
		// Advance sweep by the max length of the sweep if there was not a full
		// set of withdrawals.
		nextValidatorIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValidatorIndex += math.ValidatorIndex(sp.cs.MaxValidatorsPerWithdrawalsSweep())
		nextValidatorIndex %= math.ValidatorIndex(totalValidators)
	}

	if err = st.SetNextWithdrawalValidatorIndex(nextValidatorIndex); err != nil {
		return err
	}

	sp.logger.Info(
		"Processed withdrawals",
		"num_withdrawals", numWithdrawals,
		"evm_inflation", float64(payloadWithdrawals[0].GetAmount().Unwrap())/params.GWei,
	)

	return nil
}
