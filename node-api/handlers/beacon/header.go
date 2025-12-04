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

package beacon

import (
	"errors"
	"fmt"
	stdmath "math"

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

var ErrMismatchedSlotAndParentBlock = errors.New("slot does not match with parent block")

func (h *Handler) GetBlockHeaders(c handlers.Context) (any, error) {
	req, errReq := utils.BindAndValidate[beacontypes.GetBlockHeadersRequest](c, h.Logger())
	if errReq != nil {
		return nil, errReq
	}

	switch {
	case len(req.Slot) == 0 && len(req.ParentRoot) == 0:
		// no parameter specified, pick chain HEAD
		// by requesting special height -1.
		height := utils.Head
		return h.makeBlockHeaderResponse(height, true /*resultsInList*/)

	case len(req.Slot) != 0 && len(req.ParentRoot) == 0:
		slot, errSlot := math.U64FromString(req.Slot)
		// errSlot should always be nil since we validated slots in BindAndValidate.
		if errSlot != nil {
			return nil, fmt.Errorf("failed retrieving slot from input parameters: %w", errSlot)
		}
		if slot > stdmath.MaxInt64 { // appease linters
			return 0, fmt.Errorf("%w: slot %d", utils.ErrFailedMappingHeightTooHigh, slot)
		}
		return h.makeBlockHeaderResponse(int64(slot), true /*resultsInList*/) //#nosec: G115 // safe

	case len(req.Slot) == 0 && len(req.ParentRoot) != 0:
		parentHeight, errParent := utils.BlockIDToHeight(req.ParentRoot, h.backend)
		if errParent != nil {
			return nil, fmt.Errorf("%w, failed retrieving parent root with error: %w", handlertypes.ErrNotFound, errParent)
		}
		if parentHeight == utils.Head {
			return nil, fmt.Errorf("%w, requested header of tip's child", handlertypes.ErrNotFound)
		}
		height := parentHeight + 1
		return h.makeBlockHeaderResponse(height, true /*resultsInList*/)

	default:
		var (
			slot, errSlot         = math.U64FromString(req.Slot)
			parentSlot, errParent = utils.BlockIDToHeight(req.ParentRoot, h.backend)
		)
		if err := errors.Join(errSlot, errParent); err != nil {
			return nil, err
		}

		if slot > stdmath.MaxInt64 { // appease linters
			return 0, fmt.Errorf("%w: slot %d", utils.ErrFailedMappingHeightTooHigh, slot)
		}
		height := int64(slot) //#nosec: G115 // safe
		if height != parentSlot+1 {
			return nil, fmt.Errorf("%w: request slot %d, parent block slot %d", ErrMismatchedSlotAndParentBlock, slot, parentSlot)
		}
		return h.makeBlockHeaderResponse(height, true /*resultsInList*/)
	}
}

func (h *Handler) GetBlockHeaderByID(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeaderRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}
	slot, err := utils.BlockIDToHeight(req.BlockID, h.backend)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving slot from block ID %s: %w", req.BlockID, err)
	}
	return h.makeBlockHeaderResponse(slot, false /*resultsInList*/)
}

func (h *Handler) makeBlockHeaderResponse(height int64, resultsInList bool) (any, error) {
	st, slot, err := h.backend.StateAndSlotFromHeight(height)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get state from height %d, %s", handlertypes.ErrNotFound, height, err.Error())
	}
	// Return after updating the state root in the block header.
	header, err := st.GetLatestBlockHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block header: %w", err)
	}
	header.SetStateRoot(st.HashTreeRoot())

	// Retrieve the block signature from the block store. The signature may not be available
	// if the block is outside the availability window or if querying genesis.
	var signatureStr string
	if slot > 0 {
		signature, sigErr := h.backend.GetSignatureBySlot(slot)
		if sigErr == nil {
			signatureStr = signature.String()
		}
	}

	// While an Ethereum node may have multiple blocks per slot, BeaconKit
	// will access only one, given single slot finality and the fact that we only
	// access finalized blocks in this APIs. Still we may return a list of responses
	// to be compliant with API specs https://ethereum.github.io/beacon-APIs/?urls.primaryName=v3.1.0#/Beacon/getBlockHeader
	headerResp := beacontypes.BlockHeaderResponse{
		Root:      header.GetBodyRoot(),
		Canonical: true,
		Header: &beacontypes.SignedBeaconBlockHeader{
			Message:   beacontypes.BeaconBlockHeaderFromConsensus(header),
			Signature: signatureStr,
		},
	}

	if !resultsInList {
		return beacontypes.NewResponse(&headerResp), nil
	}

	res := []beacontypes.BlockHeaderResponse{
		headerResp,
	}
	return beacontypes.NewResponse(res), nil
}
