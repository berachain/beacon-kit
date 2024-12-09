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
	"github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
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
		// Bartio has a bug which makes DepositEth1ChainID point to
		// next deposit index, not latest processed deposit index.
		// We keep it for backward compatibility.
		depositIndex++
	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot != 0 && slot < math.U64(spec.BoonetFork2Height):
		// Boonet pre fork2 has a bug which makes DepositEth1ChainID point to
		// next deposit index, not latest processed deposit index.
		// We keep it for backward compatibility.
		depositIndex++
	default:
		// Nothing to do. We correctly set the deposit index to the last
		// processed deposit index.
	}

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
		return err
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
		// Ignore deposits that fail the signature check.
		sp.logger.Info(
			"failed deposit signature verification",
			"deposit_index", dep.GetIndex(),
			"error", err,
		)

		return nil
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
	if err = st.IncreaseBalance(idx, dep.GetAmount()); err != nil {
		return err
	}
	sp.logger.Info(
		"Processed deposit to create new validator",
		"deposit_amount", float64(dep.GetAmount().Unwrap())/math.GweiPerWei,
		"validator_index", idx,
		"withdrawal_epoch", val.GetWithdrawableEpoch(),
	)
	return nil
}
