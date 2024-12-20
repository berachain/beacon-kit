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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core/state"
)

// processOperations processes the operations and ensures they match the
// local state.
func (sp *StateProcessor[
	_, _,
]) processOperations(
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
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
	if err := sp.validateNonGenesisDeposits(
		st, deposits, blk.GetBody().GetEth1Data().DepositRoot,
	); err != nil {
		return err
	}
	for _, dep := range deposits {
		if err := sp.processDeposit(st, dep); err != nil {
			return err
		}
	}
	return st.SetEth1Data(blk.GetBody().Eth1Data)
}

// processDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _,
]) processDeposit(
	st *state.StateDB,
	dep *ctypes.Deposit,
) error {
	eth1DepositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	if err = st.SetEth1DepositIndex(eth1DepositIndex + 1); err != nil {
		return err
	}

	sp.logger.Info(
		"Processed deposit to set Eth 1 deposit index",
		"previous", eth1DepositIndex, "new", eth1DepositIndex+1,
	)

	return sp.applyDeposit(st, dep)
}

// applyDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, _,
]) applyDeposit(
	st *state.StateDB,
	dep *ctypes.Deposit,
) error {
	idx, err := st.ValidatorIndexByPubkey(dep.GetPubkey())
	if err != nil {
		// If the validator does not exist, we add the validator.
		// TODO: improve error handling by distinguishing
		// ErrNotFound from other kind of errors
		return sp.createValidator(st, dep)
	}

	// The validator already exist and we need to update its balance.
	// EffectiveBalance must be updated in processEffectiveBalanceUpdates
	// However before BoonetFork2Height we mistakenly update EffectiveBalance
	// every slot. We must preserve backward compatibility so we special case
	// Boonet to allow proper bootstrapping.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	if sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height) {
		var val *ctypes.Validator
		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return err
		}

		updatedBalance := ctypes.ComputeEffectiveBalance(
			val.GetEffectiveBalance()+dep.GetAmount(),
			math.Gwei(sp.cs.EffectiveBalanceIncrement()),
			math.Gwei(sp.cs.MaxEffectiveBalance(false)),
		)
		val.SetEffectiveBalance(updatedBalance)
		if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
			return err
		}
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
	_, _,
]) createValidator(
	st *state.StateDB,
	dep *ctypes.Deposit,
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

	// Verify that the deposit has the ETH1 withdrawal credentials.
	if !dep.HasEth1WithdrawalCredentials() {
		// Ignore deposits with non-ETH1 withdrawal credentials.
		sp.logger.Info(
			"ignoring deposit with non-ETH1 withdrawal credentials",
			"deposit_index", dep.GetIndex(),
		)
		return nil
	}

	// Verify that the message was signed correctly.
	err = dep.VerifySignature(
		ctypes.NewForkData(
			version.FromUint32[common.Version](
				sp.cs.ActiveForkVersionForEpoch(epoch),
			), genesisValidatorsRoot,
		),
		sp.cs.DomainTypeDeposit(),
		sp.signer.VerifySignature,
	)
	if err != nil {
		// Ignore deposits that fail the signature check.
		sp.logger.Info(
			"failed deposit signature verification",
			"deposit_index", dep.GetIndex(),
			"error", err,
		)
		return nil
	}

	// Add the validator to the registry.
	return sp.addValidatorToRegistry(st, dep, slot)
}

// addValidatorToRegistry adds a validator to the registry.
func (sp *StateProcessor[
	_, _,
]) addValidatorToRegistry(
	st *state.StateDB,
	dep *ctypes.Deposit,
	slot math.Slot,
) error {
	var val *ctypes.Validator
	val = val.New(
		dep.GetPubkey(),
		dep.GetWithdrawalCredentials(),
		dep.GetAmount(),
		math.Gwei(sp.cs.EffectiveBalanceIncrement()),
		math.Gwei(sp.cs.MaxEffectiveBalance(
			state.IsPostFork3(sp.cs.DepositEth1ChainID(), slot),
		)),
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
		"validator_index", idx, "withdrawal_epoch", val.GetWithdrawableEpoch(),
	)
	return nil
}
