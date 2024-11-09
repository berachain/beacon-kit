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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/sourcegraph/conc/iter"
)

// processValidatorsSetUpdates returns the validators set updates that
// will be used by consensus.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) processValidatorsSetUpdates(
	st BeaconStateT,
) (transition.ValidatorUpdates, error) {
	vals, err := st.GetValidatorsByEffectiveBalance()
	if err != nil {
		return nil, err
	}

	// filter out validators whose effective balance is not sufficient to validate
	activeVals := make([]ValidatorT, 0, len(vals))
	for _, val := range vals {
		if val.GetEffectiveBalance() > math.U64(sp.cs.EjectionBalance()) {
			activeVals = append(activeVals, val)
		}
	}

	// We need to inform consensus of the changes incurred by the validator set
	// We strive to send only diffs (added,updated or removed validators) and
	// avoid re-sending validators that have not changed.
	currentValSet, err := iter.MapErr(
		activeVals,
		func(val *ValidatorT) (*transition.ValidatorUpdate, error) {
			v := (*val)
			return &transition.ValidatorUpdate{
				Pubkey:           v.GetPubkey(),
				EffectiveBalance: v.GetEffectiveBalance(),
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	res := make([]*transition.ValidatorUpdate, 0)

	prevValsSet := make(map[string]math.Gwei, len(sp.prevEpochValidators))
	for _, v := range sp.prevEpochValidators {
		prevValsSet[string(v.Pubkey[:])] = v.EffectiveBalance
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

	// rotate validators set to new epoch ones
	sp.prevEpochValidators = currentValSet
	return res, nil
}
