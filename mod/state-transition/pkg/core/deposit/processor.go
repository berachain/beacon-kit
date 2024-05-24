// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package deposit

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// ProcessDeposits processes the deposits and ensures they match the
// local state.
func (h *Handler[BeaconStateT]) processDeposits(
	st BeaconStateT,
	deposits []*types.Deposit,
) error {
	// Ensure the deposits match the local state.
	for _, dep := range deposits {
		if err := h.processDeposit(st, dep); err != nil {
			return err
		}
		// TODO: unhood this in better spot later
		if err := st.SetEth1DepositIndex(dep.Index); err != nil {
			return err
		}
	}
	return nil
}

// processDeposit processes the deposit and ensures it matches the local state.
func (h *Handler[BeaconStateT]) processDeposit(
	st BeaconStateT,
	dep *types.Deposit,
) error {
	// TODO: fill this in properly
	// if !h.isValidMerkleBranch(
	// 	leaf,
	// 	dep.Credentials,
	// 	32 + 1,
	// 	dep.Index,
	// 	st.root,
	// ) {
	// 	return errors.New("invalid merkle branch")
	// }
	idx, err := st.ValidatorIndexByPubkey(dep.Pubkey)
	// If the validator already exists, we update the balance.
	if err == nil {
		var val *types.Validator
		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return err
		}

		// TODO: Modify balance here and then effective balance once per epoch.
		val.EffectiveBalance = min(val.EffectiveBalance+dep.Amount,
			math.Gwei(h.cs.MaxEffectiveBalance()))
		return st.UpdateValidatorAtIndex(idx, val)
	}
	// If the validator does not exist, we add the validator.
	// Add the validator to the registry.
	return h.createValidator(st, dep)
}

// createValidator creates a validator if the deposit is valid.
func (h *Handler[BeaconStateT]) createValidator(
	st BeaconStateT,
	dep *types.Deposit,
) error {
	var (
		genesisValidatorsRoot primitives.Root
		epoch                 math.Epoch
		err                   error
	)

	// Get the genesis validators root to be used to find fork data later.
	genesisValidatorsRoot, err = st.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}

	// Get the current epoch.
	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	epoch = h.cs.SlotToEpoch(slot)

	// Get the fork data for the current epoch.
	fd := types.NewForkData(
		version.FromUint32[primitives.Version](
			h.cs.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	depositMessage := types.DepositMessage{
		Pubkey:      dep.Pubkey,
		Credentials: dep.Credentials,
		Amount:      dep.Amount,
	}
	if err = depositMessage.VerifyCreateValidator(
		fd, dep.Signature, h.signer.VerifySignature, h.cs.DomainTypeDeposit(),
	); err != nil {
		return err
	}

	// Add the validator to the registry.
	return h.addValidatorToRegistry(st, dep)
}

// addValidatorToRegistry adds a validator to the registry.
func (h *Handler[BeaconStateT]) addValidatorToRegistry(
	st BeaconStateT,
	dep *types.Deposit,
) error {
	val := types.NewValidatorFromDeposit(
		dep.Pubkey,
		dep.Credentials,
		dep.Amount,
		math.Gwei(h.cs.EffectiveBalanceIncrement()),
		math.Gwei(h.cs.MaxEffectiveBalance()),
	)
	if err := st.AddValidator(val); err != nil {
		return err
	}

	idx, err := st.ValidatorIndexByPubkey(val.Pubkey)
	if err != nil {
		return err
	}
	return st.IncreaseBalance(idx, dep.Amount)
}
