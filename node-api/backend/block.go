// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	types "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BlockHeaderAtSlot returns the block header at the given slot.
func (b *Backend) BlockHeaderAtSlot(slot math.Slot) (*ctypes.BeaconBlockHeader, error) {
	st, _, err := b.StateAtSlot(slot)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	blockHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get latest block header")
	}

	// Return after updating the state root in the block header.
	blockHeader.SetStateRoot(st.HashTreeRoot())
	return blockHeader, nil
}

// GetBlockRoot returns the root of the block at the given stateID.
func (b *Backend) BlockRootAtSlot(slot math.Slot) (common.Root, error) {
	st, _, err := b.StateAtSlot(slot)
	if err != nil {
		return common.Root{}, errors.Wrapf(err, "failed to get state from slot %d", slot)
	}

	// Get the latest block header, which has the same hash tree root as the requested block.
	blockHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return common.Root{}, errors.Wrapf(err, "failed to get latest block header")
	}

	// Return hash tree root of the block header after updating the state root in it.
	blockHeader.SetStateRoot(st.HashTreeRoot())
	return blockHeader.HashTreeRoot(), nil
}

// TODO: Implement this.
func (b *Backend) BlockRewardsAtSlot(_ math.Slot) (*types.BlockRewardsData, error) {
	return &types.BlockRewardsData{
		ProposerIndex:     1,
		Total:             1,
		Attestations:      1,
		SyncAggregate:     1,
		ProposerSlashings: 1,
		AttesterSlashings: 1,
	}, nil
}
