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
	"github.com/berachain/beacon-kit/primitives/math"
)

// processSlashingsReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#slashings-balances-updates
//
//nolint:lll
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) processSlashingsReset(
	st BeaconStateT,
) error {
	// processSlashingsReset does not really do anything right now.
	// However we cannot simply drop it because appHash accounts
	// for the list of operations carried out over the state
	// even if the operations does not affect the final state
	// (currently no slashing on beaconKit)

	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// go head doing the processing
	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork3Height):
		// go head doing the processing
	default:
		// no real need to perform slashing reset
		return nil
	}

	index := (sp.cs.SlotToEpoch(slot).Unwrap() + 1) % sp.cs.EpochsPerSlashingsVector()
	return st.UpdateSlashingAtIndex(index, 0)
}
