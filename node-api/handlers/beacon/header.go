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
		// by requesting special slot 0.
		return makeBlockHeaderResponse(h.backend, 0, true /*resultsInList*/)

	case len(req.Slot) != 0 && len(req.ParentRoot) == 0:
		slot, errSlot := math.U64FromString(req.Slot)
		// errSlot should always be nil since we validated slots in BindAndValidate.
		if errSlot != nil {
			return nil, fmt.Errorf("failed retrieving slot from input parameters: %w", errSlot)
		}
		return makeBlockHeaderResponse(h.backend, slot, true /*resultsInList*/)

	case len(req.Slot) == 0 && len(req.ParentRoot) != 0:
		parentSlot, errParent := utils.SlotFromBlockID(req.ParentRoot, h.backend)
		if errParent != nil {
			return nil, fmt.Errorf("%w, failed retrieving parent root with error: %w", handlertypes.ErrNotFound, errParent)
		}
		return makeBlockHeaderResponse(h.backend, parentSlot+1, true /*resultsInList*/)

	default:
		var (
			slot, errSlot         = math.U64FromString(req.Slot)
			parentSlot, errParent = utils.SlotFromBlockID(req.ParentRoot, h.backend)
		)
		if err := errors.Join(errSlot, errParent); err != nil {
			return nil, err
		}
		if slot != parentSlot+1 {
			return nil, fmt.Errorf("%w: request slot %d, parent block slot %d", ErrMismatchedSlotAndParentBlock, slot, parentSlot)
		}
		return makeBlockHeaderResponse(h.backend, slot, true /*resultsInList*/)
	}
}

func (h *Handler) GetBlockHeaderByID(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeaderRequest](c, h.Logger())
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving slot from block ID %s: %w", req.BlockID, err)
	}

	return makeBlockHeaderResponse(h.backend, slot, false /*resultsInList*/)
}

func makeBlockHeaderResponse(backend Backend, slot math.Slot, resultsInList bool) (any, error) {
	header, err := backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, fmt.Errorf("%w: failed retrieving header at slot %d: %w", handlertypes.ErrNotFound, slot, err)
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
			Signature: "", // TODO: implement
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
