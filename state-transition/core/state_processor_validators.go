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
	stdbytes "bytes"
	"fmt"
	"slices"

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/sourcegraph/conc/iter"
)

//nolint:lll // let it be
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) processRegistryUpdates(
	st BeaconStateT,
) error {
	vals, err := st.GetValidators()
	if err != nil {
		return fmt.Errorf("registry update, failed listing validators: %w", err)
	}

	slot, err := st.GetSlot()
	if err != nil {
		return fmt.Errorf("registry update, failed loading slot: %w", err)
	}
	currEpoch := sp.cs.SlotToEpoch(slot)
	nextEpoch := currEpoch + 1

	minEffectiveBalance := math.Gwei(sp.cs.EjectionBalance() + sp.cs.EffectiveBalanceIncrement())

	// We do not currently have a cap on validator churn,
	// so we can process validators activations in a single loop
	var idx math.ValidatorIndex
	for si, val := range vals {
		valModified := false
		if val.IsEligibleForActivationQueue(minEffectiveBalance) {
			val.SetActivationEligibilityEpoch(nextEpoch)
			valModified = true
		}
		if val.IsEligibleForActivation(currEpoch) {
			val.SetActivationEpoch(nextEpoch)
			valModified = true
		}
		// Note: without slashing and voluntary withdrawals, there is no way
		// for an activa validator to have its balance less or equal to EjectionBalance

		if valModified {
			idx, err = st.ValidatorIndexByPubkey(val.GetPubkey())
			if err != nil {
				return fmt.Errorf("registry update, failed loading validator index, state index %d: %w", si, err)
			}
			if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
				return fmt.Errorf("registry update, failed updating validator idx %d: %w", idx, err)
			}
		}
	}

	// validators registry will be possibly further modified in order to enforce
	// validators set cap. We will do that at the end of processEpoch, once all
	// Eth 2.0 like transitions has been done (notable EffectiveBalances handling).
	return nil
}

//nolint:lll // let it be
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) processValidatorSetCap(
	st BeaconStateT,
) error {
	// Enforce the validator set cap by:
	// 1- retrieving validators active next epoch
	// 2- sorting them by stake
	// 3- dropping enough validators to fulfill the cap

	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	nextEpoch := sp.cs.SlotToEpoch(slot) + 1

	nextEpochVals, err := sp.getActiveVals(st, nextEpoch)
	if err != nil {
		return fmt.Errorf("registry update, failed retrieving next epoch vals: %w", err)
	}

	if uint64(len(nextEpochVals)) <= sp.cs.ValidatorSetCap() {
		// nothing to eject
		return nil
	}

	slices.SortFunc(nextEpochVals, func(lhs, rhs ValidatorT) int {
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
			return stdbytes.Compare(val1Pk[:], val2Pk[:])
		}
	})

	// We do not currently have a cap on validators churn, so we stop
	// validators this epoch and we withdraw them next epoch
	var idx math.ValidatorIndex
	for li := range uint64(len(nextEpochVals)) - sp.cs.ValidatorSetCap() {
		valToEject := nextEpochVals[li]
		valToEject.SetExitEpoch(nextEpoch)
		valToEject.SetWithdrawableEpoch(nextEpoch + 1)
		idx, err = st.ValidatorIndexByPubkey(valToEject.GetPubkey())
		if err != nil {
			return fmt.Errorf("registry update, failed loading validator index: %w", err)
		}
		if err = st.UpdateValidatorAtIndex(idx, valToEject); err != nil {
			return fmt.Errorf("registry update, failed ejecting validator idx %d: %w", li, err)
		}
	}

	return nil
}

// Note: validatorSetsDiffs does not need to be a StateProcessor method
// but it helps simplifying generic instantiation.
func (*StateProcessor[
	_, _, _, _, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) validatorSetsDiffs(
	prevEpochValidators []ValidatorT,
	currEpochValidator []ValidatorT,
) transition.ValidatorUpdates {
	currentValSet := iter.Map(
		currEpochValidator,
		func(val *ValidatorT) *transition.ValidatorUpdate {
			v := (*val)
			return &transition.ValidatorUpdate{
				Pubkey:           v.GetPubkey(),
				EffectiveBalance: v.GetEffectiveBalance(),
			}
		},
	)

	res := make([]*transition.ValidatorUpdate, 0)
	prevValsSet := make(map[string]math.Gwei, len(prevEpochValidators))
	for _, v := range prevEpochValidators {
		pk := v.GetPubkey()
		prevValsSet[string(pk[:])] = v.GetEffectiveBalance()
	}

	for _, newVal := range currentValSet {
		key := string(newVal.Pubkey[:])
		oldBal, found := prevValsSet[key]
		if !found {
			// new validator, we add it with its weight
			res = append(res, newVal)
			continue
		}
		if oldBal != newVal.EffectiveBalance {
			// validator updated, we add it with new weight
			res = append(res, newVal)
		}

		// consume pre-existing validators
		delete(prevValsSet, key)
	}

	// prevValsSet now contains all evicted validators (and only those)
	for pkBytes := range prevValsSet {
		//#nosec:G703 // bytes comes from a pk
		pk, _ := bytes.ToBytes48([]byte(pkBytes))
		res = append(res, &transition.ValidatorUpdate{
			Pubkey:           pk,
			EffectiveBalance: 0, // signal val eviction to consensus
		})
	}
	return res
}

// nextEpochValidatorSet returns the current estimation of what next epoch
// validator set would be.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) getActiveVals(st BeaconStateT, epoch math.Epoch) ([]ValidatorT, error) {
	vals, err := st.GetValidators()
	if err != nil {
		return nil, err
	}
	activeVals := make([]ValidatorT, 0, len(vals))
	for _, val := range vals {
		if val.IsActive(epoch) {
			activeVals = append(activeVals, val)
		}
	}

	return activeVals, nil
}
