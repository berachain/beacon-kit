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

package core

import (
	"slices"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	crypto "github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// TODO: THIS FUNCTION NEEDS TO BE INTEGRATED BETTER WITH WITHDRAWALS AND
// SLASHING, IT IS A HACKY TEMPORARY WAY TO GET THE VALIDATOR SET UPDATING
// NICELY.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) getNextSyncCommittee(
	st BeaconStateT,
) ([]*transition.ValidatorUpdate, error) {
	syncCommittee, err := st.GetSyncCommittee()
	if err != nil {
		return nil, err
	}

	var idx math.U64
	var val *types.Validator
	validatorUpdates := make([]*transition.ValidatorUpdate, 0)
	aboveEjectionBlance := make([]crypto.BLSPubkey, 0)
	for _, pubkey := range syncCommittee {
		idx, err = st.ValidatorIndexByPubkey(pubkey)
		if err != nil {
			return nil, err
		}

		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return nil, err
		}

		// If the validator is in the committee and above the ejection balance
		// then they get to stay in the committee.
		if val.EffectiveBalance >= math.U64(sp.cs.EjectionBalance()) &&
			!val.Slashed {
			aboveEjectionBlance = append(aboveEjectionBlance, pubkey)
		} else {
			// If the validator is in the committee but below the ejection
			// balance
			// then they get ejected.
			validatorUpdates = append(validatorUpdates, &transition.ValidatorUpdate{
				Pubkey: pubkey,
				Exit:   true,
			})
		}
	}

	// If no validators were ejected / we are at the max committee size we can
	// return early.
	if len(aboveEjectionBlance) == int(sp.cs.SyncCommitteeSize()) {
		return nil, nil
	}

	allValidators, err := st.GetValidators()
	if err != nil {
		return nil, err
	}

	for _, val := range allValidators {
		// We want to stop once we have the max committee size.
		if len(aboveEjectionBlance) == int(sp.cs.SyncCommitteeSize()) {
			break
		}

		// If the validator is not already in the committee and above the
		// ejection balance we can add it to the committee.
		if !slices.Contains(aboveEjectionBlance, val.Pubkey) &&
			val.EffectiveBalance >= math.U64(sp.cs.MaxEffectiveBalance()) &&
			!val.Slashed {
			aboveEjectionBlance = append(aboveEjectionBlance, val.Pubkey)
			validatorUpdates = append(
				validatorUpdates,
				&transition.ValidatorUpdate{
					Pubkey: val.Pubkey,
					Exit:   false,
				},
			)
		}
	}

	// Set the new sync committee.
	if err = st.SetSyncCommittee(aboveEjectionBlance); err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}
