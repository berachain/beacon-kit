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
	"strconv"

	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[
	BeaconBlockHeaderT, ContextT, _, _,
]) GetBlobSidecars(c ContextT) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlobSidecarsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	slot, err := utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, err
	}

	// convert indices to uint64
	indices := make([]uint64, len(req.Indices))
	for i, idx := range req.Indices {
		indices[i], err = strconv.ParseUint(idx, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	blobSidecars, err := h.backend.BlobSidecarsAtSlot(slot, indices)
	if err != nil {
		return nil, err
	}

	return beacontypes.BlobSidecarsResponse[BeaconBlockHeaderT]{
		Data: blobSidecars,
	}, nil
}
