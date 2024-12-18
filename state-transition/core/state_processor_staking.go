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
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core/state"
)

// processOperations processes deposits with basic validation. Other features found in the
// Ethereum 2.0 specification are not implemented currently.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#operations
func (sp *StateProcessor[
	BeaconBlockT, BeaconStateT, _, _,
]) processOperations(st BeaconStateT, blk BeaconBlockT) error {
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	// Verify that the provided deposit root in eth1data is consistent with our local view of
	// the deposit tree.
	localDeposits, localDepositsRoot, err := sp.ds.GetDepositsByIndex(
		depositIndex, sp.cs.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}
	eth1Data := blk.GetBody().GetEth1Data()
	if eth1Data.DepositRoot != localDepositsRoot {
		return errors.New("local deposit tree root does not match the block deposit tree root")
	}

	// Verify that the provided deposit count is consistent with our local view of the
	// deposit tree.
	if uint64(len(localDeposits)) != min(
		sp.cs.MaxDepositsPerBlock(),
		eth1Data.DepositCount.Unwrap()-depositIndex,
	) {
		return errors.Wrapf(
			ErrDepositCountMismatch, "expected: %d, got: %d",
			min(sp.cs.MaxDepositsPerBlock(), eth1Data.DepositCount.Unwrap()-depositIndex),
			len(localDeposits),
		)
	}

	// The provided eth1data is valid, accept it and set locally.
	if err = st.SetEth1Data(eth1Data); err != nil {
		return err
	}

	// Process each deposit in the block.
	for _, dep := range blk.GetBody().GetDeposits() {
		if err = sp.processDeposit(st, dep); err != nil {
			return err
		}
	}
	return nil
}

// processDeposit processes the deposit similarly to the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#deposits
func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) processDeposit(st BeaconStateT, dep *ctypes.Deposit) error {
	// Verify proof of deposit inclusion.
	eth1DepositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}
	eth1Data, err := st.GetEth1Data()
	if err != nil {
		return err
	}
	if !merkle.VerifyProof(
		eth1Data.DepositRoot, dep.Data.HashTreeRoot(), eth1DepositIndex, dep.GetProof(),
	) {
		return errors.Wrapf(ErrInvalidDepositProof, "deposit: %+v", dep.Data)
	}

	// Update the deposit index.
	newDepositIndex := eth1DepositIndex + 1
	if err = st.SetEth1DepositIndex(newDepositIndex); err != nil {
		return err
	}
	sp.logger.Info(
		"Processed deposit to update Eth 1 deposit index",
		"previous", eth1DepositIndex, "new", newDepositIndex,
	)

	// Apply the deposit.
	return sp.applyDeposit(st, dep.Data)
}

// applyDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) applyDeposit(st BeaconStateT, dep *ctypes.DepositData) error {
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
	_, BeaconStateT, _, _,
]) createValidator(st BeaconStateT, dep *ctypes.DepositData) error {
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
			"failed deposit signature verification", "deposit_index", dep.GetIndex(), "error", err,
		)
		return nil
	}

	// Add the validator to the registry.
	return sp.addValidatorToRegistry(st, dep, slot)
}

// addValidatorToRegistry adds a validator to the registry.
func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) addValidatorToRegistry(
	st BeaconStateT,
	dep *ctypes.DepositData,
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
