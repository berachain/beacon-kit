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
	types "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// BlockHeader returns the block header at the given slot.
func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockHeaderAtSlot(
	slot uint64,
) (BeaconBlockHeaderT, error) {
	var blockHeader BeaconBlockHeaderT

	st, _, err := b.stateFromSlot(slot)
	if err != nil {
		return blockHeader, err
	}

	blockHeader, err = st.GetLatestBlockHeader()
	return blockHeader, err
}

// GetBlockRoot returns the root of the block at the given stateID.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRootAtSlot(
	slot uint64,
) (common.Root, error) {
	st, slot, err := b.stateFromSlot(slot)
	if err != nil {
		return common.Root{}, err
	}

	// As calculated by the beacon chain. Ideally, this logic
	// should be abstracted by the beacon chain.
	index := slot % b.cs.SlotsPerHistoricalRoot()
	return st.GetBlockRootAtIndex(index)
}

// TODO: Implement this.
func (b Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlockRewardsAtSlot(
	_ uint64,
) (*types.BlockRewardsData, error) {
	return &types.BlockRewardsData{
		ProposerIndex:     1,
		Total:             1,
		Attestations:      1,
		SyncAggregate:     1,
		ProposerSlashings: 1,
		AttesterSlashings: 1,
	}, nil
}
