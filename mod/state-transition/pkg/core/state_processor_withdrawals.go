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

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/davecgh/go-spew/spew"
)

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
// NOTE: Modified from the Ethereum 2.0 specification to support EVM inflation:
// 1. The first withdrawal MUST be a fixed EVM inflation withdrawal
// 2. Subsequent withdrawals (if any) are processed as validator withdrawals
// 3. This modification reduces the maximum validator withdrawals per block by
// one
//
//nolint:lll // TODO: Simplify when dropping special cases.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processWithdrawals(
	st BeaconStateT,
	blk BeaconBlockT,
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

	return sp.processWithdrawalsByFork(
		st, expectedWithdrawals, payloadWithdrawals)
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _,
	_, _, _, _, WithdrawalT, WithdrawalsT, _,
]) processWithdrawalsByFork(
	st BeaconStateT,
	expectedWithdrawals []WithdrawalT,
	payloadWithdrawals []WithdrawalT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return errors.Wrap(
			err, "failed loading slot while processing withdrawals",
		)
	}

	// Common validations
	if len(expectedWithdrawals) != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Chain/Fork specific processing
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		return sp.processWithdrawalsBartio(
			st,
			expectedWithdrawals,
			payloadWithdrawals,
			slot,
		)

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot == math.U64(spec.BoonetFork1Height):
		// Slot used to emergency mint EVM tokens on Boonet.
		if !expectedWithdrawals[0].Equals(payloadWithdrawals[0]) {
			return fmt.Errorf(
				"minting withdrawal does not match expected %s, got %s",
				spew.Sdump(expectedWithdrawals[0]),
				spew.Sdump(payloadWithdrawals[0]),
			)
		}

		return nil // No processing needed.

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height):
		// Boonet inherited the Bartio behaviour pre BoonetFork2Height
		// nothing specific to do
		return sp.processWithdrawalsBartio(
			st,
			expectedWithdrawals,
			payloadWithdrawals,
			slot,
		)

	default:
		return sp.processWithdrawalsDefault(
			st,
			expectedWithdrawals,
			payloadWithdrawals,
			slot,
		)
	}
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _,
	_, _, _, _, WithdrawalT, WithdrawalsT, _,
]) processWithdrawalsBartio(
	st BeaconStateT,
	expectedWithdrawals []WithdrawalT,
	payloadWithdrawals []WithdrawalT,
	slot math.Slot,
) error {
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrWithdrawalMismatch,
				"withdrawal at index %d does not match expected %s, got %s",
				i, spew.Sdump(wd), spew.Sdump(payloadWithdrawals[i]),
			)
		}

		// Process the validator withdrawal.
		if err := st.DecreaseBalance(
			wd.GetValidatorIndex(), wd.GetAmount(),
		); err != nil {
			return err
		}
	}

	if len(expectedWithdrawals) != 0 {
		if err := st.SetNextWithdrawalIndex(
			(expectedWithdrawals[len(expectedWithdrawals)-1].
				GetIndex() + 1).Unwrap(),
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

	//#nosec:G701 // won't overflow in practice.
	if len(expectedWithdrawals) == int(sp.cs.MaxWithdrawalsPerPayload()) {
		nextValidatorIndex =
			(expectedWithdrawals[len(expectedWithdrawals)-1].GetIndex() + 1) %
				math.ValidatorIndex(totalValidators)
		// Note: this is a bug, we should have used ValidatorIndex instead of
		// GetIndex. processWithdrawalsDefault fixes it
	} else {
		// Advance sweep by the max length of the sweep if there was not a full
		// set of withdrawals.
		nextValidatorIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValidatorIndex += math.ValidatorIndex(
			sp.cs.MaxValidatorsPerWithdrawalsSweep(
				state.IsPostUpgrade, spec.BartioChainID, slot,
			))
		nextValidatorIndex %= math.ValidatorIndex(totalValidators)
	}

	return st.SetNextWithdrawalValidatorIndex(nextValidatorIndex)
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _,
	_, _, _, _, WithdrawalT, WithdrawalsT, _,
]) processWithdrawalsDefault(
	st BeaconStateT,
	expectedWithdrawals []WithdrawalT,
	payloadWithdrawals []WithdrawalT,
	slot math.Slot,
) error {
	// Enforce that first withdrawal is EVM inflation
	if len(payloadWithdrawals) == 0 {
		return ErrZeroWithdrawals
	}
	if !payloadWithdrawals[0].Equals(st.EVMInflationWithdrawal()) {
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

		if err := st.DecreaseBalance(
			expectedWithdrawals[i].GetValidatorIndex(),
			expectedWithdrawals[i].GetAmount(),
		); err != nil {
			return err
		}
	}

	if numWithdrawals > 1 {
		if err := st.SetNextWithdrawalIndex(
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

	//#nosec:G701 // won't overflow in practice.
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
		nextValidatorIndex += math.ValidatorIndex(
			sp.cs.MaxValidatorsPerWithdrawalsSweep(
				state.IsPostUpgrade, sp.cs.DepositEth1ChainID(), slot,
			))
		nextValidatorIndex %= math.ValidatorIndex(totalValidators)
	}

	return st.SetNextWithdrawalValidatorIndex(nextValidatorIndex)
}
