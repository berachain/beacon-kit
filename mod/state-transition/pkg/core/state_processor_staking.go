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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/davecgh/go-spew/spew"
)

// processOperations processes the operations and ensures they match the
// local state.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processOperations(
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	deposits := blk.GetBody().GetDeposits()
	if err := sp.validateNonGenesisDeposits(st, deposits); err != nil {
		return err
	}
	for _, dep := range deposits {
		if err := sp.processDeposit(st, dep); err != nil {
			return err
		}
	}
	return nil
}

// processDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, _, _, _, _, _,
]) processDeposit(
	st BeaconStateT,
	dep DepositT,
) error {
	if err := st.SetEth1DepositIndex(dep.GetIndex().Unwrap()); err != nil {
		return err
	}

	return sp.applyDeposit(st, dep)
}

// applyDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) applyDeposit(
	st BeaconStateT,
	dep DepositT,
) error {
	idx, err := st.ValidatorIndexByPubkey(dep.GetPubkey())
	if err != nil {
		// If the validator does not exist, we add the validator.
		// Add the validator to the registry.
		return sp.createValidator(st, dep)
	}

	// If the validator already exists, we update the balance.
	var val ValidatorT
	val, err = st.ValidatorByIndex(idx)
	if err != nil {
		return err
	}

	// TODO: Modify balance here and then effective balance once per epoch.
	updatedBalance := types.ComputeEffectiveBalance(
		val.GetEffectiveBalance()+dep.GetAmount(),
		math.Gwei(sp.cs.EffectiveBalanceIncrement()),
		math.Gwei(sp.cs.MaxEffectiveBalance()),
	)
	val.SetEffectiveBalance(updatedBalance)
	if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
		return err
	}
	return st.IncreaseBalance(idx, dep.GetAmount())
}

// createValidator creates a validator if the deposit is valid.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, ForkDataT, _, _, _, _, _, _,
]) createValidator(
	st BeaconStateT,
	dep DepositT,
) error {
	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// At genesis, the validators sign over an empty root.
	genesisValidatorsRoot := common.Root{}
	if slot != 0 {
		// Get the genesis validators root to be used to find fork data later.
		genesisValidatorsRoot, err = st.GetGenesisValidatorsRoot()
		if err != nil {
			return err
		}
	}

	// Get the current epoch.
	epoch := sp.cs.SlotToEpoch(slot)

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
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, ValidatorT, _, _, _, _,
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
	if sp.cs.DepositEth1ChainID() == spec.BartioChainID {
		// Note in AddValidatorBartio we implicitly increase
		// the balance from state st. This is unlike AddValidator.
		return st.AddValidatorBartio(val)
	}

	if err := st.AddValidator(val); err != nil {
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
// NOTE: Modified from the Ethereum 2.0 specification to support EVM inflation:
// 1. The first withdrawal MUST be a fixed EVM inflation withdrawal
// 2. Subsequent withdrawals (if any) are processed as validator withdrawals
// 3. This modification reduces the maximum validator withdrawals per block by
// one
//
//nolint:lll,funlen,gocognit // TODO: Simplify when dropping special cases.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processWithdrawals(
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	// Dequeue and verify the logs.
	var (
		body               = blk.GetBody()
		slot               = blk.GetSlot().Unwrap()
		nextValidatorIndex math.ValidatorIndex
		payload            = body.GetExecutionPayload()
		payloadWithdrawals = payload.GetWithdrawals()
	)

	// Get the expected withdrawals.
	expectedWithdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		return err
	}
	numWithdrawals := len(expectedWithdrawals)
	if sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot >= spec.BoonetFork2Height {
		// Enforce that there is at least one withdrawal for EVM inflation.
		if numWithdrawals == 0 {
			return ErrZeroWithdrawals
		}
	}

	// Ensure the expected and payload withdrawals have the same length.
	if numWithdrawals != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			numWithdrawals, len(payloadWithdrawals),
		)
	}

	if sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot == spec.BoonetFork1Height {
		// Slot used to emergency mint EVM tokens on Boonet.
		if !expectedWithdrawals[0].Equals(payloadWithdrawals[0]) {
			return fmt.Errorf(
				"minting withdrawal does not match expected %s, got %s",
				spew.Sdump(expectedWithdrawals[0]),
				spew.Sdump(payloadWithdrawals[0]),
			)
		}

		// No processing needed.
		return nil
	}

	// Compare and process each withdrawal.
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrWithdrawalMismatch,
				"withdrawal at index %d does not match expected %s, got %s",
				i, spew.Sdump(wd), spew.Sdump(payloadWithdrawals[i]),
			)
		}

		if sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
			slot >= spec.BoonetFork2Height {
			// The first withdrawal is the EVM inflation withdrawal. Aside from
			// simple validation, no processing to the state is needed.
			if i == 0 {
				if !wd.Equals(st.EVMInflationWithdrawal()) {
					return ErrFirstWithdrawalNotEVMInflation
				}
				continue
			}
		}

		// Process the validator withdrawal.
		if err = st.DecreaseBalance(
			wd.GetValidatorIndex(), wd.GetAmount(),
		); err != nil {
			return err
		}
	}

	if sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot >= spec.BoonetFork2Height {
		// If there is only the EVM inflation withdrawal, no state update is
		// needed.
		if numWithdrawals == 1 {
			return nil
		}
	}

	// Next sweep starts after the latest withdrawal's validator index.
	if numWithdrawals != 0 {
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
	//#nosec:G701 // won't overflow in practice.
	if numWithdrawals == int(sp.cs.MaxWithdrawalsPerPayload()) {
		if (sp.cs.DepositEth1ChainID() == spec.BartioChainID) ||
			(sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
				slot < spec.BoonetFork2Height) {
			nextValidatorIndex =
				(expectedWithdrawals[numWithdrawals-1].GetIndex() + 1) %
					math.ValidatorIndex(totalValidators)
		} else {
			// Next sweep starts after the latest withdrawal's validator index.
			nextValidatorIndex = (expectedWithdrawals[numWithdrawals-1].
				GetValidatorIndex() + 1) % math.ValidatorIndex(totalValidators)
		}
	} else {
		// Advance sweep by the max length of the sweep if there was not a full
		// set of withdrawals.
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
