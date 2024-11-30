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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
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
			math.Gwei(sp.cs.GetMaxEffectiveBalance())))
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
	epoch = sp.cs.GetSlotToEpoch(slot.Unwrap())

	// Verify that the message was signed correctly.
	var d ForkDataT
	if err = dep.VerifySignature(
		d.New(
			version.FromUint32[common.Version](
				sp.cs.GetActiveForkVersionForEpoch(epoch),
			), genesisValidatorsRoot,
		),
		sp.cs.GetDomainTypeDeposit(),
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
		math.Gwei(sp.cs.GetEffectiveBalanceIncrement()),
		math.Gwei(sp.cs.GetMaxEffectiveBalance()),
	)

	return st.AddValidator(val)
}

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
//nolint:lll
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) processWithdrawals(
	_ BeaconStateT,
	_ BeaconBlockBodyT,
) error {
	// TODO: implement
	return nil
}

// processForcedWithdrawals is a helper function to process forced withdrawals.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _,
]) processForcedWithdrawals(
	_ BeaconStateT,
	_ transition.ValidatorUpdates,
) error {
	// TODO: Implement this function.
	return nil
}

// ExpectedWithdrawals retrieves the expected withdrawals.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _,
	_, _, _, _, _, _, _, WithdrawalsT, _,
]) ExpectedWithdrawals(st BeaconStateT) (WithdrawalsT, error) {
	return st.GetWithdrawals()
}
