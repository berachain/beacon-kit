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
	"bytes"
	"fmt"
	"slices"

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
	// Verify that outstanding deposits are processed
	// up to the maximum number of deposits

	// TODO we should assert  here that
	// len(body.deposits) ==  min(MAX_DEPOSITS,
	// state.eth1_data.deposit_count - state.eth1_deposit_index)
	// Until we fix eth1Data we do a partial check
	deposits := blk.GetBody().GetDeposits()
	if uint64(len(deposits)) > sp.cs.MaxDepositsPerBlock() {
		return errors.Wrapf(
			ErrExceedsBlockDepositLimit, "expected: %d, got: %d",
			sp.cs.MaxDepositsPerBlock(), len(deposits),
		)
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
		// TODO: improve error handling by distinguishing
		// validator not found from other kind of errors
		return sp.createValidator(st, dep)
	}

	// if validator exist, just update its balance
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

	// BeaconKit enforces a cap on the validator set size. If the deposit
	// breaches the cap, we find the validator with the smallest stake and
	// mark it as withdrawable so that it will be eventually evicted and
	// its deposits returned.
	capHit, err := sp.isValidatorCapHit(st)
	if err != nil {
		return err
	}
	if !capHit {
		return sp.addValidatorInternal(st, val, dep.GetAmount())
	}

	// Adding the validator would breach the cap. Find the validator
	// with the smallest stake among current and candidate validators
	// and kit it out.
	currSmallestVal, err := sp.findSmallestValidator(st)
	if err != nil {
		return err
	}

	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	epoch := sp.cs.SlotToEpoch(slot)

	if val.GetEffectiveBalance() <= currSmallestVal.GetEffectiveBalance() {
		// in case of tie-break among candidate validator we prefer
		// existing one so we mark candidate as withdrawable
		val.SetWithdrawableEpoch(epoch)
		return sp.addValidatorInternal(st, val, dep.GetAmount())
	}

	// mark exiting validator for eviction and add candidate
	currSmallestVal.SetWithdrawableEpoch(epoch)
	idx, err := st.ValidatorIndexByPubkey(currSmallestVal.GetPubkey())
	if err != nil {
		return err
	}
	if err = st.UpdateValidatorAtIndex(idx, currSmallestVal); err != nil {
		return err
	}
	return sp.addValidatorInternal(st, val, dep.GetAmount())
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) isValidatorCapHit(st BeaconStateT) (bool, error) {
	if sp.processingGenesis {
		// minor optimization: we do check in Genesis that
		// validators do no exceed cap. So hitting cap here
		// should never happen
		return false, nil
	}
	validators, err := st.GetValidators()
	if err != nil {
		return false, err
	}

	//#nosec:G701 // no overflow risk here
	return uint32(len(validators)) >= sp.cs.GetValidatorSetCapSize(), nil
}

// TODO: consider moving this to BeaconState directly
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, DepositT, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) findSmallestValidator(st BeaconStateT) (
	ValidatorT,
	error,
) {
	var smallestVal ValidatorT

	// candidateVal would breach the cap. Find the validator with the
	// smallest stake among current and candidate validators.
	currentVals, err := st.GetValidatorsByEffectiveBalance()
	if err != nil {
		return smallestVal, err
	}

	// TODO: consider heapifying slice instead. We only care about the smallest
	slices.SortFunc(currentVals, func(lhs, rhs ValidatorT) int {
		var (
			val1Stake = lhs.GetEffectiveBalance()
			val2Stake = rhs.GetEffectiveBalance()
		)
		switch {
		case val1Stake < val2Stake:
			return -1
		case val1Stake > val2Stake:
			return 1
		default:
			// validators pks are guaranteed to be different
			var (
				val1Pk = lhs.GetPubkey()
				val2Pk = rhs.GetPubkey()
			)
			return bytes.Compare(val1Pk[:], val2Pk[:])
		}
	})
	return currentVals[0], nil
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) addValidatorInternal(
	st BeaconStateT,
	val ValidatorT,
	depositAmount math.Gwei,
) error {
	// TODO: This is a bug that lives on bArtio. Delete this eventually.
	if sp.cs.DepositEth1ChainID() == bArtioChainID {
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
	return st.IncreaseBalance(idx, depositAmount)
}

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
//nolint:lll
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
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
	expectedWithdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		return err
	}
	numWithdrawals := len(expectedWithdrawals)

	// Ensure the withdrawals have the same length
	if numWithdrawals != len(payloadWithdrawals) {
		return errors.Wrapf(
			ErrNumWithdrawalsMismatch,
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Compare and process each withdrawal.
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return errors.Wrapf(
				ErrNumWithdrawalsMismatch,
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
