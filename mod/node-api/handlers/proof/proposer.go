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
	"unsafe"

	ptypes "github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/mod/node-api/types"
)

// GetBlockProposer returns the block proposer for the given block id along
// with a merkle proof that can be verified against the beacon block root.
func (h *Handler[
	ContextT, BeaconBlockHeaderT, _, _, _, _, ValidatorT,
]) GetBlockProposer(c ContextT) (any, error) {
	params, err := utils.BindAndValidate[types.BlockIDRequest](c)
	if err != nil {
		return nil, err
	}

	// Get the slot from the given input of block id, block header, and beacon
	// state for the desired slot.
	slot, err := utils.SlotFromBlockID(params.BlockID)
	if err != nil {
		return nil, err
	}
	blockHeader, err := h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, err
	}
	beaconState, err := h.backend.StateFromSlot(slot)
	if err != nil {
		return nil, err
	}

	// Get the beacon state struct for the proving the proposer validator
	// exists within the state.
	beaconStateForValidatorProof, err := ptypes.NewBeaconStateForValidator(
		beaconState, h.backend.ChainSpec(),
	)
	if err != nil {
		return nil, err
	}

	// Set the "correct" state root on the block header, which yields the exact
	// beacon block root that is set on execution clients for EIP-4788.
	//
	// TODO: Determine why the backend's block header from beacon state has an
	// incorrect state root (also retreivable from backend via StateRootAtSlot).
	stateRootCorrection, err := beaconStateForValidatorProof.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	blockHeader.SetStateRoot(stateRootCorrection)

	// Now get the beacon block struct for proving the proposer validator exists
	// within the state in this block.
	beaconBlockForValidatorProof, err := ptypes.NewBeaconBlockForValidator(
		blockHeader, beaconStateForValidatorProof,
	)
	if err != nil {
		return nil, err
	}

	// Generate the proof (along with the "correct" beacon block root to verify
	// against) for the proposer validator.
	validatorProof, beaconBlockRoot, err := ptypes.GetProofForProposer_FastSSZ(
		beaconBlockForValidatorProof,
	)
	if err != nil {
		return nil, err
	}

	//#nosec:G103 // on purpose.
	validators := *(*[]ValidatorT)(
		unsafe.Pointer(&beaconBlockForValidatorProof.StateRoot.Validators),
	)

	return ptypes.BlockProposerProofResponse[
		BeaconBlockHeaderT, ValidatorT,
	]{
		BeaconBlockHeader: blockHeader,
		BeaconBlockRoot:   beaconBlockRoot,
		Validator:         validators[blockHeader.GetProposerIndex()],
		ValidatorProof:    validatorProof,
	}, nil
}
