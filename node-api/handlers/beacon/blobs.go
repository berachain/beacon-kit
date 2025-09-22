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
	"github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/mapping"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetBlobSidecars provides an implementation for the
// "/eth/v1/beacon/blob_sidecars/:block_id" API endpoint.
//
//nolint:gocognit // TODO: fix
func (h *Handler) GetBlobSidecars(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[types.GetBlobSidecarsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	// Map requested blockID to slot
	slotID, err := mapping.BlockIDToHeight(req.BlockID, h.backend)
	if err != nil {
		return nil, err
	}

	var slot math.Slot
	if slotID == mapping.Head {
		latestHeight, _ := h.backend.GetSyncData()
		if latestHeight < 0 {
			return nil, errors.New("invalid negative block height")
		}
		slot = math.Slot(latestHeight)
	} else {
		slot = math.Slot(slotID) //#nosec: G115 // practically safe
	}

	// Convert indices to uint64.
	indices := make([]uint64, len(req.Indices))
	for i, idxS := range req.Indices {
		var idx math.U64
		idx, err = math.U64FromString(idxS)
		if err != nil {
			return nil, err
		}
		indices[i] = idx.Unwrap()
	}

	// Validate the requested slot is within the Data Availability Period.
	if !h.cs.WithinDAPeriod(slot, slot) {
		return nil, fmt.Errorf(
			"requested slot (%d) is not within Data Availability Period (previous %d epochs)",
			slotID, h.cs.MinEpochsForBlobsSidecarsRequest(),
		)
	}

	// Validate request indices.
	if uint64(len(indices)) >= h.cs.MaxBlobsPerBlock() {
		return nil, errors.New("too many indices requested")
	}
	for _, index := range indices {
		if index >= h.cs.MaxBlobsPerBlock() {
			return nil, errors.New("blob index out of range")
		}
	}

	blobSidecars, err := h.backend.GetBlobSidecarsAtSlot(slot)
	if err != nil {
		return nil, err
	}

	// Create a map of requested indices for O(1) index lookups.
	isRequestIndex := make(map[uint64]bool)
	for _, idx := range indices {
		isRequestIndex[idx] = true
	}

	// Preallocate response slice - if indices specified, size will be len(indices),
	// otherwise size will be all sidecars.
	responseCap := len(blobSidecars)
	if len(indices) > 0 {
		responseCap = len(indices)
	}
	blobSidecarsResponse := make([]*types.Sidecar, 0, responseCap)

	for _, blobSidecar := range blobSidecars {
		// Skip if indices specified and this index not requested.
		if len(indices) > 0 && !isRequestIndex[blobSidecar.GetIndex()] {
			continue
		}
		// Craft and append the blob sidecar serialized data to the response.
		blobSidecarsResponse = append(blobSidecarsResponse,
			types.SidecarFromConsensus(blobSidecar),
		)
	}

	return types.SidecarsResponse{
		Data: blobSidecarsResponse,
	}, nil
}
