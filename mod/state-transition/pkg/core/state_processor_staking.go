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
	// Verify that outstanding deposits are processed up to the maximum number
	// of deposits.
	deposits := blk.GetBody().GetDeposits()
	index, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	eth1Data, err := st.GetEth1Data()
	if err != nil {
		return err
	}
	depositCount := min(
		sp.cs.MaxDepositsPerBlock(),
		eth1Data.GetDepositCount().Unwrap()-index,
	)
	_ = depositCount
	// TODO: Update eth1data count and check this.
	// if uint64(len(deposits)) != depositCount {
	// 	return errors.New("deposit count mismatch")
	// }
	return sp.processDeposits(st, deposits)
}

// processDeposits processes the deposits and ensures  they match the
// local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, _, _, _, _, _,
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
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, _, _, _, _, _,
]) processDeposit(
	st BeaconStateT,
	dep DepositT,
) error {
	var nextDepositIndex uint64
	switch depositIndex, err := st.GetEth1DepositIndex(); {
	case err == nil:
		// just increment the deposit index if no error
		nextDepositIndex = depositIndex + 1
	case sp.processingGenesis && err != nil:
		// If errored and still processing genesis,
		// Eth1DepositIndex may have not been set yet.
		nextDepositIndex = 0
	default:
		// Failed retrieving Eth1DepositIndex outside genesis is an error
		return fmt.Errorf(
			"failed retrieving eth1 deposit index outside of processing genesis: %w",
			err,
		)
	}

	if err := st.SetEth1DepositIndex(nextDepositIndex); err != nil {
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
// NOTE: This function is modified from the spec to allow a fixed withdrawal
// (must be the first withdrawal) used for EVM inflation.
//
//nolint:lll
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processWithdrawals(
	st BeaconStateT,
	body BeaconBlockBodyT,
) error {
	var (
		payload            = body.GetExecutionPayload()
		payloadWithdrawals = payload.GetWithdrawals()
	)

	// Get the expected withdrawals.
	expectedWithdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		return err
	}
	numWithdrawals := len(expectedWithdrawals)

	// Ensure the expected and payload withdrawals have the same length.
	if numWithdrawals != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Ensure the EVM inflation withdrawal is the first withdrawal.
	if numWithdrawals == 0 {
		return ErrNoWithdrawals
	}
	if !expectedWithdrawals[0].GetAddress().Equals(
		sp.cs.EVMInflationAddress(),
	) ||
		expectedWithdrawals[0].GetAmount() != math.U64(
			sp.cs.EVMInflationPerBlock(),
		) {
		return ErrFirstWithdrawalNotEVMInflation
	}

	// Compare and process each withdrawal.
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrWithdrawalMismatch,
				"withdrawals do not match expected %s, got %s",
				spew.Sdump(wd), spew.Sdump(payloadWithdrawals[i]),
			)
		}

		// The first withdrawal is the EVM inflation withdrawal. No processing
		// is needed.
		if i == 0 {
			continue
		}

		// Process the validator withdrawal.
		if err = st.DecreaseBalance(
			wd.GetValidatorIndex(), wd.GetAmount(),
		); err != nil {
			return err
		}
	}

	// Update the next withdrawal index if this block contained withdrawals
	// (excluding the EVM inflation withdrawal).
	if numWithdrawals > 1 {
		// Next sweep starts after the latest withdrawal's validator index.
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
	var nextValIndex math.ValidatorIndex
	if numWithdrawals > 1 &&
		//#nosec:G701 // won't overflow in practice.
		numWithdrawals == int(sp.cs.MaxWithdrawalsPerPayload()) {
		// Next sweep starts after the latest withdrawal's validator index.
		nextValIndex =
			expectedWithdrawals[numWithdrawals-1].GetValidatorIndex() + 1
	} else {
		// Advance sweep by the max length of the sweep if there was not a full
		// set of withdrawals.
		nextValIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValIndex += math.U64(sp.cs.MaxValidatorsPerWithdrawalsSweep())
	}
	nextValIndex %= math.U64(totalValidators)
	return st.SetNextWithdrawalValidatorIndex(nextValIndex)
}
