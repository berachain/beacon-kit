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
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processRewardsAndPenalties(st BeaconStateT) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// processRewardsAndPenalties does not really do anything right now.
	// However we cannot simply drop it because appHash accounts
	// for the list of operations carried out over the state
	// even if the operations does not affect the final state
	// (rewards and penalties are always zero at this stage of beaconKit)

	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// go head doing the processing, eve
	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork3Height):
	default:
		// no real need to perform hollowProcessRewardsAndPenalties
		return nil
	}

	if sp.cs.SlotToEpoch(slot) == math.U64(constants.GenesisEpoch) {
		return nil
	}

	// this has been simplified to make clear that
	// we are not really doing anything here
	valCount, err := st.GetTotalValidators()
	if err != nil {
		return err
	}

	for i := range valCount {
		// Increase the balance of the validator.
		if err = st.IncreaseBalance(math.ValidatorIndex(i), 0); err != nil {
			return err
		}

		// Decrease the balance of the validator.
		if err = st.DecreaseBalance(math.ValidatorIndex(i), 0); err != nil {
			return err
		}
	}

	return nil
}
