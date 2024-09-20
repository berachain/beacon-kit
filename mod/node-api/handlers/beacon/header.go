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

package beacon

import (
	"github.com/berachain/beacon-kit/mod/errors"
	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func (h *Handler[
	BeaconBlockHeaderT, ContextT, _, _,
]) GetBlockHeaders(c ContextT) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlockHeadersRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}
	var slot math.Slot
	// If slot is being passed in the request, fetch block header at that slot.
	if req.Slot != "" {
		slot, err = utils.U64FromString(req.Slot)
		if err != nil {
			return nil, err
		}
	} else if req.ParentRoot != "" {
		// TODO: If parent_root is being passed in the request,
		// fetch block header at that parent_root.
		// Convert the string to common.Root
		var parentRoot common.Root
		err = parentRoot.UnmarshalText([]byte(req.ParentRoot))
		if err != nil {
			return nil, errors.Wrapf(err, "invalid parent root: %v", req.ParentRoot)
		}

		slot, err = h.backend.GetSlotByStateRoot(parentRoot)
		if err != nil {
			return nil, err
		}
	} else {
		// If neither slot nor parent_root is provided, fetch current head slot.
		slot, err = h.backend.GetHeadSlot()
		if err != nil {
			return nil, err
		}
	}
	header, err := h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, err
	}

	root, err := h.backend.BlockRootAtSlot(slot)
	if err != nil {
		return nil, err
	}

	return types.Wrap(&beacontypes.BlockHeaderResponse[BeaconBlockHeaderT]{
		Root:      root,
		Canonical: true,
		Header: &beacontypes.BlockHeader[BeaconBlockHeaderT]{
			Message:   header,
			Signature: crypto.BLSSignature{}, // TODO: implement
		},
	}), nil
}

func (h *Handler[
	BeaconBlockHeaderT, ContextT, _, _,
]) GetBlockHeaderByID(c ContextT) (any, error) {
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

	root, err := h.backend.BlockRootAtSlot(slot)
	if err != nil {
		return nil, err
	}
	return types.Wrap(&beacontypes.BlockHeaderResponse[BeaconBlockHeaderT]{
		Root:      root, // This is root hash of entire beacon block.
		Canonical: true,
		Header: &beacontypes.BlockHeader[BeaconBlockHeaderT]{
			Message:   header,
			Signature: crypto.BLSSignature{}, // TODO: implement
		},
	}), nil
}
