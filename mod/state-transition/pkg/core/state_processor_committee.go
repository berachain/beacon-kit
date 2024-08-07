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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/sourcegraph/conc/iter"
)

// processSyncCommitteeUpdates processes the sync committee updates.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, ValidatorT, _, _, _, _,
]) processSyncCommitteeUpdates(
	st BeaconStateT,
) (transition.ValidatorUpdates, error) {
	vals, err := st.GetValidatorsByEffectiveBalance()
	if err != nil {
		return nil, err
	}

	return iter.MapErr(
		vals,
		func(val *ValidatorT) (*transition.ValidatorUpdate, error) {
			v := (*val)
			return &transition.ValidatorUpdate{
				Pubkey:           v.GetPubkey(),
				EffectiveBalance: v.GetEffectiveBalance(),
			}, nil
		},
	)
}
