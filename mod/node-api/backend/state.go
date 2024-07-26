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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// StateFromSlot returns the state at the given slot using query context.
func (b *Backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateFromSlot(
	slot uint64,
) (BeaconStateT, error) {
	var state BeaconStateT
	//#nosec:G701 // not an issue in practice.
	queryCtx, err := b.node.CreateQueryContext(int64(slot), false)
	if err != nil {
		return state, err
	}

	return b.sb.StateFromContext(queryCtx), nil
}

// GetStateRoot returns the root of the state at the given stateID.
//
// TODO: fix https://github.com/berachain/beacon-kit/issues/1777.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateRootAtSlot(
	slot uint64,
) (common.Root, error) {
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}

	// This is required to handle the semantical expectation that
	// 0 -> latest despite 0 != latest.
	if slot == 0 {
		var latestSlot math.U64
		latestSlot, err = st.GetSlot()
		if err != nil {
			return common.Root{}, err
		}
		slot = latestSlot.Unwrap()
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
	st, err := b.StateFromSlot(slot)
	if err != nil {
		return fork, err
	}
	return st.GetFork()
}
