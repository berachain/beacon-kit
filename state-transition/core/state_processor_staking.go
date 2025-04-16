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
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/ethereum/go-ethereum/params"
)

// processOperations processes the operations and ensures they match the local state.
func (sp *StateProcessor) processOperations(
	ctx ReadOnlyContext,
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	// Verify that outstanding deposits are processed up to the maximum number of deposits.
	//
	// Unlike Eth 2.0 specs we don't check that
	// `len(body.deposits) ==  min(MAX_DEPOSITS, state.eth1_data.deposit_count - state.eth1_deposit_index)`
	deposits := blk.GetBody().GetDeposits()
	if uint64(len(deposits)) > sp.cs.MaxDepositsPerBlock() {
		return errors.Wrapf(
			ErrExceedsBlockDepositLimit, "expected: %d, got: %d",
			sp.cs.MaxDepositsPerBlock(), len(deposits),
		)
	}

	// Instead we directly compare block deposits with our local store ones.
	if err := ValidateNonGenesisDeposits(
		ctx.ConsensusCtx(),
		st,
		sp.ds,
		sp.cs.MaxDepositsPerBlock(),
		deposits,
		blk.GetBody().GetEth1Data().DepositRoot,
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
func (sp *StateProcessor) processDeposit(st *state.StateDB, dep *ctypes.Deposit) error {
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
	if err = sp.applyDeposit(st, dep); err != nil {
		return fmt.Errorf("failed to apply deposit: %w", err)
	}
	return nil
}

// applyDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor) applyDeposit(st *state.StateDB, dep *ctypes.Deposit) error {
	idx, err := st.ValidatorIndexByPubkey(dep.GetPubkey())
	if err != nil {
		sp.logger.Info("Validator does not exist so creating",
			"pubkey", dep.GetPubkey(), "index", dep.GetIndex(), "deposit_amount", dep.GetAmount())
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
		"deposit_amount", float64(dep.GetAmount().Unwrap())/params.GWei,
		"validator_index", idx,
	)
	return nil
}

// createValidator creates a validator if the deposit is valid.
func (sp *StateProcessor) createValidator(st *state.StateDB, dep *ctypes.Deposit) error {
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

	// Check that the deposit has the ETH1 withdrawal credentials.
	if !dep.HasEth1WithdrawalCredentials() {
		sp.logger.Warn(
			"adding validator with non-ETH1 withdrawal credentials -- NOT withdrawable",
			"pubkey", dep.GetPubkey().String(),
			"deposit_index", dep.GetIndex(),
			"amount_gwei", dep.GetAmount().Unwrap(),
		)
		sp.metrics.incrementValidatorNotWithdrawable()
	}

	// Verify that the message was signed correctly.
	err = dep.VerifySignature(
		ctypes.NewForkData(
			// Deposits must be signed with GENESIS_FORK_VERSION.
			sp.cs.GenesisForkVersion(),
			genesisValidatorsRoot,
		),
		sp.cs.DomainTypeDeposit(),
		sp.signer.VerifySignature,
	)
	if err != nil {
		// Ignore deposits that fail the signature check.
		sp.logger.Warn(
			"failed deposit signature verification",
			"pubkey", dep.GetPubkey().String(),
			"deposit_index", dep.GetIndex(),
			"amount_gwei", dep.GetAmount().Unwrap(),
			"error", err,
		)
		sp.metrics.incrementDepositStakeLost()
		return nil
	}

	// Add the validator to the registry.
	return sp.addValidatorToRegistry(st, dep)
}

// addValidatorToRegistry adds a validator to the registry.
func (sp *StateProcessor) addValidatorToRegistry(st *state.StateDB, dep *ctypes.Deposit) error {
	val := ctypes.NewValidatorFromDeposit(
		dep.GetPubkey(),
		dep.GetWithdrawalCredentials(),
		dep.GetAmount(),
		math.Gwei(sp.cs.EffectiveBalanceIncrement()),
		math.Gwei(sp.cs.MaxEffectiveBalance()),
	)

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
		"deposit_amount", float64(dep.GetAmount().Unwrap())/params.GWei,
		"validator_index", idx, "withdrawal_epoch", val.GetWithdrawableEpoch(),
	)
	return nil
}
