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

	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/sourcegraph/conc/iter"
)

//nolint:lll
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

	var idx math.ValidatorIndex
	for si, val := range vals {
		if val.IsEligibleForActivationQueue(minEffectiveBalance) {
			val.SetActivationEligibilityEpoch(nextEpoch)
			idx, err = st.ValidatorIndexByPubkey(val.GetPubkey())
			if err != nil {
				return fmt.Errorf("registry update, failed loading validator index, state index %d: %w", si, err)
			}
			if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
				return fmt.Errorf("registry update, failed updating validator idx %d: %w", idx, err)
			}
		}
	}

	return nil
}

// processValidatorsSetUpdates returns the validators set updates that
// will be used by consensus.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) processValidatorsSetUpdates(
	st BeaconStateT,
) (transition.ValidatorUpdates, error) {
	// at this state slot has not been updated yet so
	// we pick nextEpochValidatorSet
	activeVals, err := sp.nextEpochValidatorSet(st)
	if err != nil {
		return nil, err
	}

	// pick prev epoch validators
	slot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	sp.valSetMu.Lock()
	defer sp.valSetMu.Unlock()

	// prevEpoch is calculated assuming current block
	// will turn epoch but we have not update slot yet
	prevEpoch := sp.cs.SlotToEpoch(slot)
	currEpoch := prevEpoch + 1
	if slot == 0 {
		currEpoch = 0 // prevEpoch for genesis is zero
	}
	prevEpochVals := sp.valSetByEpoch[prevEpoch] // picks nil if it's genesis

	// calculate diff
	res := sp.validatorSetsDiffs(prevEpochVals, activeVals)

	// clear up sets we won't lookup to anymore
	sp.valSetByEpoch[currEpoch] = activeVals
	if prevEpoch >= 1 {
		delete(sp.valSetByEpoch, prevEpoch-1)
	}
	return res, nil
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
