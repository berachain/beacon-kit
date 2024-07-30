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

package backend

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// StateFromSlotForProof returns the beacon state of the version that was used
// to calculate the parent beacon block root, which has the empty state root in
// the latest block header. Hence we do not process the next slot.
func (b *Backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateFromSlotForProof(slot uint64) (BeaconStateT, uint64, error) {
	return b.stateFromSlotRaw(slot)
}

// GetStateRoot returns the root of the state at the given slot.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateRootAtSlot(
	slot uint64,
) (common.Root, error) {
	st, slot, err := b.stateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}

	// As calculated by the beacon chain. Ideally, this logic
	// should be abstracted by the beacon chain.
	index := slot % b.cs.SlotsPerHistoricalRoot()
	return st.StateRootAtIndex(index)
}

// GetStateFork returns the fork of the state at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, ForkT, _, _, _, _, _, _,
]) StateForkAtSlot(
	slot uint64,
) (ForkT, error) {
	var fork ForkT
	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return fork, err
	}
	return st.GetFork()
}
