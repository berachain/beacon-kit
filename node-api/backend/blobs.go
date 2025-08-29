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
	"errors"
	"fmt"

	apitypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BlobSidecarsByIndices is the backend helper function that will query the
// data availability store for all sidecars for a slot, returning only those
// sidecars specified by the indices, or all sidecars if left unspecified.
func (b *Backend) BlobSidecarsByIndices(slot math.Slot, indices []uint64) ([]*apitypes.Sidecar, error) {
	currentSlot, _ := b.node.GetSyncData()
	if currentSlot < 0 {
		return nil, errors.New("invalid negative block height")
	}

	// If the requested slot is 0 (head, finalized, justified), use the current slot.
	if slot == 0 {
		slot = math.Slot(currentSlot)
	}

	// Validate the requested slot is within the Data Availability Period.
	if !b.cs.WithinDAPeriod(slot, math.Slot(currentSlot)) {
		return nil, fmt.Errorf(
			"requested slot (%d) is not within Data Availability Period (previous %d epochs)",
			slot, b.cs.MinEpochsForBlobsSidecarsRequest(),
		)
	}

	// Validate request indices.
	if uint64(len(indices)) >= b.cs.MaxBlobsPerBlock() {
		return nil, errors.New("too many indices requested")
	}
	for _, index := range indices {
		if index >= b.cs.MaxBlobsPerBlock() {
			return nil, errors.New("blob index out of range")
		}
	}

	blobSidecars, err := b.sb.AvailabilityStore().GetBlobSidecars(slot)
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
	blobSidecarsResponse := make([]*apitypes.Sidecar, 0, responseCap)

	for _, blobSidecar := range blobSidecars {
		// Skip if indices specified and this index not requested.
		if len(indices) > 0 && !isRequestIndex[blobSidecar.GetIndex()] {
			continue
		}
		// Craft and append the blob sidecar serialized data to the response.
		blobSidecarsResponse = append(blobSidecarsResponse,
			apitypes.SidecarFromConsensus(blobSidecar),
		)
	}
	return blobSidecarsResponse, nil
}
