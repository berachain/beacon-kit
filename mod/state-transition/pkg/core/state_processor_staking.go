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
	"slices"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
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

	// Unlike Eth 2.0 specs we don't check that
	// len(body.deposits) ==  min(MAX_DEPOSITS,
	// state.eth1_data.deposit_count - state.eth1_deposit_index)
	// Instead we directly compare block deposits with store ones.
	deposits := blk.GetBody().GetDeposits()
	if uint64(len(deposits)) > sp.cs.MaxDepositsPerBlock() {
		return errors.Wrapf(
			ErrExceedsBlockDepositLimit, "expected: %d, got: %d",
			sp.cs.MaxDepositsPerBlock(), len(deposits),
		)
	}
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
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	depositIndex := dep.GetIndex().Unwrap()

	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// On bartio we set the deposit index to the last processed deposit
		// index + 1.
		depositIndex++
	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height):
		// On boonet pre fork 2, we set the deposit index to the last processed
		// deposit index + 1.
		depositIndex++
	default:
		// Nothing to do. We correctly set the deposit index to the last
		// processed deposit index.
	}

	// Set the deposit index in beacon state.
	if err = st.SetEth1DepositIndex(depositIndex); err != nil {
		return err
	}

	sp.logger.Info(
		"Processed deposit to set Eth 1 deposit index",
		"deposit_index", depositIndex,
	)

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
		// ErrNotFound from other kind of errors
		return sp.createValidator(st, dep)
	}

	// if validator exist, just update its balance
	if err = st.IncreaseBalance(idx, dep.GetAmount()); err != nil {
		return errors.Wrap(err, "failed to increase balance")
	}

	sp.logger.Info(
		"Processed deposit to increase balance",
		"deposit_amount", float64(dep.GetAmount().Unwrap())/math.GweiPerWei,
		"validator_index", idx,
	)

	return nil
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
	var candidateVal ValidatorT
	candidateVal = candidateVal.New(
		dep.GetPubkey(),
		dep.GetWithdrawalCredentials(),
		dep.GetAmount(),
		math.Gwei(sp.cs.EffectiveBalanceIncrement()),
		math.Gwei(sp.cs.MaxEffectiveBalance()),
	)

	// BeaconKit enforces a cap on the validator set size. If the deposit
	// breaches the cap, we find the validator with the smallest stake and
	// mark it as withdrawable so that it will be evicted next epoch and
	// its deposits returned.

	nextEpochVals, err := sp.nextEpochValidatorSet(st)
	if err != nil {
		return err
	}
	//#nosec:G701 // no overflow risk here
	if uint64(len(nextEpochVals)) < sp.cs.ValidatorSetCap() {
		// cap not hit, just add the validator
		return sp.addValidatorInternal(st, candidateVal, dep.GetAmount())
	}

	// Adding the validator would breach the cap. Find the validator
	// with the smallest stake among current and candidate validators
	// and kick it out.
	lowestStakeVal, err := sp.lowestStakeVal(nextEpochVals)
	if err != nil {
		return err
	}

	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	nextEpoch := sp.cs.SlotToEpoch(slot) + 1

	if candidateVal.GetEffectiveBalance() <= lowestStakeVal.GetEffectiveBalance() {
		// in case of tie-break among candidate validator we prefer
		// existing one so we mark candidate as withdrawable
		// We wait next epoch to return funds, as a way to curb spamming
		candidateVal.SetWithdrawableEpoch(nextEpoch)
		return sp.addValidatorInternal(st, candidateVal, dep.GetAmount())
	}

	// mark existing validator for eviction and add candidate
	lowestStakeVal.SetWithdrawableEpoch(nextEpoch)
	idx, err := st.ValidatorIndexByPubkey(lowestStakeVal.GetPubkey())
	if err != nil {
		return err
	}
	if err = st.UpdateValidatorAtIndex(idx, lowestStakeVal); err != nil {
		return err
	}
	return sp.addValidatorInternal(st, candidateVal, dep.GetAmount())
}

// nextEpochValidatorSet returns the current estimation of what next epoch
// validator set would be.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) nextEpochValidatorSet(st BeaconStateT) ([]ValidatorT, error) {
	slot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}
	nextEpoch := sp.cs.SlotToEpoch(slot) + 1

	vals, err := st.GetValidators()
	if err != nil {
		return nil, err
	}
	activeVals := make([]ValidatorT, 0, len(vals))
	for _, val := range vals {
		if val.GetEffectiveBalance() <= math.U64(sp.cs.EjectionBalance()) {
			continue
		}
		if val.GetWithdrawableEpoch() == nextEpoch {
			continue
		}
		activeVals = append(activeVals, val)
	}

	return activeVals, nil
}

// TODO: consider moving this to BeaconState directly
func (*StateProcessor[
	_, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) lowestStakeVal(currentVals []ValidatorT) (
	ValidatorT,
	error,
) {
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

	if err = st.IncreaseBalance(idx, depositAmount); err != nil {
		return errors.Wrap(err, "failed to increase balance")
	}

	sp.logger.Info(
		"Processed deposit to create new validator",
		"deposit_amount", float64(depositAmount.Unwrap())/math.GweiPerWei,
		"validator_index", idx,
	)

	return nil
}
