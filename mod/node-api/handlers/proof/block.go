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

package proof

import (
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Get the slot from the given input of block id, beacon state, and beacon
// block header for the resolved slot.
func (h *Handler[
	ContextT, BeaconBlockHeaderT, BeaconStateT, _, _, _,
]) resolveBlockID(blockID string) (
	uint64, BeaconStateT, BeaconBlockHeaderT, error,
) {
	var (
		beaconState BeaconStateT
		blockHeader BeaconBlockHeaderT
	)

	slot, err := utils.SlotFromBlockID(blockID, h.backend)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	beaconState, slot, err = h.backend.StateFromSlot(slot)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	// Get the latest block header from the state itself and not the backend
	// to avoid querying state from disk (query context) again.
	blockHeader, err = beaconState.GetLatestBlockHeader()
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	// For proofs, we need to patch the latest block header on the beacon state
	// with the version that was used to calculate the parent beacon block root,
	// which has the empty state root in the latest block header. Check
	// EIP-4788 and the spec for more details.
	blockHeaderForProofInState := blockHeader.New(
		math.Slot(slot),
		blockHeader.GetProposerIndex(),
		blockHeader.GetParentBlockRoot(),
		common.Root{},
		blockHeader.GetBodyRoot(),
	)
	if err = beaconState.SetLatestBlockHeader(
		blockHeaderForProofInState,
	); err != nil {
		return 0, beaconState, blockHeader, err
	}

	return slot, beaconState, blockHeader, nil
}
