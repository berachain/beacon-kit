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
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

// GetExecutionFeeRecipient returns the fee recipient from the latest execution
// payload header for the given block id, along with the proof that can be
// verified against the beacon block root.
func (h *Handler[
	ContextT, BeaconBlockHeaderT, _, _, _, _,
]) GetExecutionFeeRecipient(c ContextT) (any, error) {
	params, err := utils.BindAndValidate[types.ExecutionFeeRecipientRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, beaconState, blockHeader, err := h.resolveBlockID(params.BlockID)
	if err != nil {
		return nil, err
	}

	// Generate the proof (along with the "correct" beacon block root to
	// verify against) for the execution payload fee recipient.
	h.Logger().Info("Generating execution fee recipient proof", "slot", slot)
	proof, beaconBlockRoot, err := merkle.ProveExecutionFeeRecipientInBlock(
		blockHeader, beaconState,
	)
	if err != nil {
		return nil, err
	}

	// Get the block number from the latest execution payload header.
	leph, err := beaconState.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	return types.ExecutionFeeRecipientResponse[BeaconBlockHeaderT]{
		BeaconBlockHeader:          blockHeader,
		BeaconBlockRoot:            beaconBlockRoot,
		ExecutionFeeRecipient:      leph.GetFeeRecipient(),
		ExecutionFeeRecipientProof: proof,
	}, nil
}
