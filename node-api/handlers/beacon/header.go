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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// determineSlot determines the slot based on the query params.
// If both slot and parent_root are provided, it verifies that the provided slot
// and parent root match.
// If only slot is provided, it returns the slot.
// If only parent_root is provided, it returns the slot for the parent root.
// If neither slot nor parent_root is provided, it returns 0 so that the latest slot is fetched.
func (h *Handler) determineSlot(req *beacontypes.GetBlockHeadersRequest) (math.Slot, error) {
	var parentRoot common.Root
	switch {
	case req.Slot != "" && req.ParentRoot != "":
		// Both slot and parent_root are provided
		slot, err := math.U64FromString(req.Slot)
		if err != nil {
			return 0, errors.Wrapf(err, "invalid slot: %v", req.Slot)
		}
		err = parentRoot.UnmarshalText([]byte(req.ParentRoot))
		if err != nil {
			return 0, errors.Wrapf(err, "invalid parent root: %v", req.ParentRoot)
		}
		// Verify that the provided slot and parent root match
		verifiedSlot, errVerify := h.backend.GetSlotByParentRoot(parentRoot)
		if errVerify != nil {
			return 0, errVerify
		}
		if verifiedSlot != slot {
			return 0, errors.New(
				"provided slot does not match the slot for the given parent root",
			)
		}
		return slot, nil

	case req.Slot != "":
		// If only slot is provided
		slot, err := math.U64FromString(req.Slot)
		if err != nil {
			return 0, err
		}
		return slot, nil

	case req.ParentRoot != "":
		// If only parent_root is provided

		// Convert the string to common.Root
		err := parentRoot.UnmarshalText([]byte(req.ParentRoot))
		if err != nil {
			return 0, errors.Wrapf(err, "invalid parent root: %v", req.ParentRoot)
		}
		slot, err := h.backend.GetSlotByParentRoot(parentRoot)
		if err != nil {
			return 0, err
		}
		return slot, nil

	default:
		// If neither slot nor parent_root is provided, pass the slot as 0 so that
		// the latest slot is fetched.
		return 0, nil
	}
}

func (h *Handler) GetBlockHeaders(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeadersRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	// Determine the slot based on the query params.
	slot, err := h.determineSlot(&req)
	if err != nil {
		return nil, err
	}

	header, err := h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, err
	}

	root, err := h.backend.BlockRootAtSlot(slot)
	if err != nil {
		return nil, err
	}

	// Create a slice with a single header response
	// This is as per the API spec.
	// https://ethereum.github.io/beacon-APIs/?urls.primaryName=v3.1.0#/Beacon/getBlockHeader
	headers := []beacontypes.BlockHeaderResponse{{
		Root:      root,
		Canonical: true,
		Header: &beacontypes.SignedBeaconBlockHeader{
			Message:   beacontypes.BeaconBlockHeaderFromConsensus(header),
			Signature: "", // TODO: implement
		},
	}}

	return beacontypes.NewResponse(headers), nil
}

func (h *Handler) GetBlockHeaderByID(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeaderRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	slot, err := utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, err
	}
	header, err := h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, err
	}
	return beacontypes.NewResponse(&beacontypes.BlockHeaderResponse{
		Root:      header.GetBodyRoot(),
		Canonical: true,
		Header: &beacontypes.SignedBeaconBlockHeader{
			Message:   beacontypes.BeaconBlockHeaderFromConsensus(header),
			Signature: "", // TODO: implement
		},
	}), nil
}
