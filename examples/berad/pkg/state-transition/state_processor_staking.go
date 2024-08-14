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

package transition

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/davecgh/go-spew/spew"
)

// processDeposits processes the deposits and ensures  they match the
// local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, _, _, _, _,
]) processDeposits(
	st BeaconStateT,
	deposits []DepositT,
) error {
	// Ensure the deposits match the local state.
	for _, dep := range deposits {
		if err := sp.processDeposit(st, dep); err != nil {
			return err
		}
	}
	return nil
}

// processDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, _, _, _, _,
]) processDeposit(
	st BeaconStateT,
	dep DepositT,
) error {
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	if err = st.SetEth1DepositIndex(
		depositIndex + 1,
	); err != nil {
		return err
	}

	return sp.applyDeposit(st, dep)
}

// applyDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, ValidatorT, _, _, _, _,
]) applyDeposit(
	st BeaconStateT,
	dep DepositT,
) error {
	idx, err := st.ValidatorIndexByPubkey(dep.GetPubkey())
	// If the validator already exists, we update the balance.
	if err == nil {
		var val ValidatorT
		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return err
		}

		// TODO: Modify balance here and then effective balance once per epoch.
		val.SetEffectiveBalance(min(val.GetEffectiveBalance()+dep.GetAmount(),
			math.Gwei(sp.cs.MaxEffectiveBalance())))
		return st.UpdateValidatorAtIndex(idx, val)
	}

	// If the validator does not exist, we add the validator.
	// Add the validator to the registry.
	return sp.createValidator(st, dep)
}

// createValidator creates a validator if the deposit is valid.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, ForkDataT, _, _, _, _, _, _,
]) createValidator(
	st BeaconStateT,
	dep DepositT,
) error {
	var (
		genesisValidatorsRoot common.Root
		epoch                 uint64
		err                   error
	)

	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// At genesis, the validators sign over an empty root.
	if slot == 0 {
		genesisValidatorsRoot = common.Root{}
	} else {
		// Get the genesis validators root to be used to find fork data later.
		genesisValidatorsRoot, err = st.GetGenesisValidatorsRoot()
		if err != nil {
			return err
		}
	}

	// Get the current epoch.
	epoch = sp.cs.SlotToEpoch(slot.Unwrap())

	// Verify that the message was signed correctly.
	var d ForkDataT
	if err = dep.VerifySignature(
		d.New(
			version.FromUint32[common.Version](
				sp.cs.ActiveForkVersionForEpoch(epoch),
			), genesisValidatorsRoot,
		),
		sp.cs.DomainTypeDeposit(),
		sp.signer.VerifySignature,
	); err != nil {
		return err
	}

	// Add the validator to the registry.
	return sp.addValidatorToRegistry(st, dep)
}

// addValidatorToRegistry adds a validator to the registry.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, ValidatorT, _, _, _, _,
]) addValidatorToRegistry(
	st BeaconStateT,
	dep DepositT,
) error {
	var val ValidatorT
	val = val.New(
		dep.GetPubkey(),
		dep.GetWithdrawalCredentials(),
		dep.GetAmount(),
		math.Gwei(sp.cs.EffectiveBalanceIncrement()),
		math.Gwei(sp.cs.MaxEffectiveBalance()),
	)

	// TODO: This is a bug that lives on bArtio. Delete this eventually.
	const bArtioChainID = 80084
	if sp.cs.DepositEth1ChainID() == bArtioChainID {
		if err := st.AddValidatorBartio(val); err != nil {
			return err
		}
	} else if err := st.AddValidator(val); err != nil {
		return err
	}

	idx, err := st.ValidatorIndexByPubkey(val.GetPubkey())
	if err != nil {
		return err
	}

	return st.IncreaseBalance(idx, dep.GetAmount())
}

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
//nolint:lll
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) processWithdrawals(
	st BeaconStateT,
	body BeaconBlockBodyT,
) error {
	// Dequeue and verify the logs.
	var (
		nextValidatorIndex math.ValidatorIndex
		payload            = body.GetExecutionPayload()
		payloadWithdrawals = payload.GetWithdrawals()
	)

	// Get the expected withdrawals.
	expectedWithdrawals, err := sp.expectedWithdrawals(st)
	if err != nil {
		return err
	}
	numWithdrawals := len(expectedWithdrawals)

	// Ensure the withdrawals have the same length
	if numWithdrawals != len(payloadWithdrawals) {
		return errors.Newf(
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Compare and process each withdrawal.
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return errors.Newf(
				"withdrawals do not match expected %s, got %s",
				spew.Sdump(wd), spew.Sdump(payloadWithdrawals[i]),
			)
		}

		// Then we process the withdrawal.
		if err = st.DecreaseBalance(
			wd.GetValidatorIndex(), wd.GetAmount(),
		); err != nil {
			return err
		}
	}

	// Update the next withdrawal index if this block contained withdrawals
	if numWithdrawals != 0 {
		// Next sweep starts after the latest withdrawal's validator index
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

	// Update the next validator index to start the next withdrawal sweep
	//#nosec:G701 // won't overflow in practice.
	if numWithdrawals == int(sp.cs.MaxWithdrawalsPerPayload()) {
		// Next sweep starts after the latest withdrawal's validator index
		nextValidatorIndex =
			(expectedWithdrawals[len(expectedWithdrawals)-1].GetIndex() + 1) %
				math.U64(totalValidators)
	} else {
		// Advance sweep by the max length of the sweep if there was not
		// a full set of withdrawals
		nextValidatorIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValidatorIndex += math.ValidatorIndex(
			sp.cs.MaxValidatorsPerWithdrawalsSweep())
		nextValidatorIndex %= math.ValidatorIndex(totalValidators)
	}

	return st.SetNextWithdrawalValidatorIndex(nextValidatorIndex)
}

// processForcedWithdrawals is a helper function to process forced withdrawals.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) processForcedWithdrawals(
	st BeaconStateT,
	valUpdates transition.ValidatorUpdates,
) error {
	// Process the forced withdrawals.
	for _, valUpdate := range valUpdates {
		idx, err := st.ValidatorIndexByPubkey(valUpdate.Pubkey)
		if err != nil {
			return err
		}

		// Set the validator's effective balance to 0.
		if err = st.SetBalance(idx, 0); err != nil {
			return err
		}
	}
	return nil
}

// TODO: This is exposed for the PayloadBuilder and probably should be done in a
// better way.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _,
	_, _, _, _, ValidatorT, ValidatorsT, WithdrawalT, _, _,
]) ExpectedWithdrawals(st BeaconStateT) ([]WithdrawalT, error) {
	return sp.expectedWithdrawals(st)
}

// ExpectedWithdrawals as defined in the Ethereum 2.0 Specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-get_expected_withdrawals
//
//nolint:lll
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _, _, _, _, ValidatorT, ValidatorsT, WithdrawalT, _, _,
]) expectedWithdrawals(st BeaconStateT) ([]WithdrawalT, error) {
	var (
		balance           math.Gwei
		withdrawalAddress common.ExecutionAddress
		withdrawals       = make([]WithdrawalT, 0)
	)

	slot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	epoch := math.Epoch(uint64(slot) / sp.cs.SlotsPerEpoch())

	withdrawalIndex, err := st.GetNextWithdrawalIndex()
	if err != nil {
		return nil, err
	}

	validatorIndex, err := st.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return nil, err
	}

	totalValidators, err := st.GetTotalValidators()
	if err != nil {
		return nil, err
	}

	bound := min(
		totalValidators, sp.cs.MaxValidatorsPerWithdrawalsSweep(),
	)

	var (
		validator  ValidatorT
		withdrawal WithdrawalT
		amount     math.Gwei
	)

	// Iterate through indices to find the next validators to withdraw.
	for range bound {
		validator, err = st.ValidatorByIndex(validatorIndex)
		if err != nil {
			return nil, err
		}

		balance, err = st.GetBalance(validatorIndex)
		if err != nil {
			return nil, err
		}

		withdrawalAddress, err = validator.
			GetWithdrawalCredentials().ToExecutionAddress()
		if err != nil {
			return nil, err
		}

		// Set the amount of the withdrawal depending on the balance of the
		// validator.
		if validator.IsFullyWithdrawable(balance, epoch) {
			amount = balance
		} else if validator.IsPartiallyWithdrawable(
			balance, math.Gwei(sp.cs.MaxEffectiveBalance()),
		) {
			amount = balance - math.Gwei(sp.cs.MaxEffectiveBalance())
		}
		withdrawal = withdrawal.New(
			math.U64(withdrawalIndex),
			validatorIndex,
			withdrawalAddress,
			amount,
		)

		withdrawals = append(withdrawals, withdrawal)

		// Increment the withdrawal index to process the next withdrawal.
		withdrawalIndex++

		// Cap the number of withdrawals to the maximum allowed per payload.
		//#nosec:G701 // won't overflow in practice.
		if len(withdrawals) == int(sp.cs.MaxWithdrawalsPerPayload()) {
			break
		}

		// Increment the validator index to process the next validator.
		validatorIndex = (validatorIndex + 1) % math.ValidatorIndex(
			totalValidators,
		)
	}

	return withdrawals, nil
}
