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

package backend

import (
	"fmt"
	datypes "github.com/berachain/beacon-kit/da/types"
	apitypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (b *Backend[
	_, _, _, _,
]) BlobSidecarsByIndices(slot math.Slot, indices []uint64) ([]*apitypes.Sidecar, error) {
	var blobSidecars datypes.BlobSidecars
	// TODO: Check if we are WithinDAPeriod(). Have to get current head slot somehow.

	fmt.Println("Fetcing blobs from da store at slot", slot)
	blobSidecars, err := b.sb.AvailabilityStore().GetBlobSidecars(slot)
	fmt.Println("Fetched", len(blobSidecars), "blobSidecars")
	if err != nil {
		return nil, err
	}

	// Create a map of requested indices for O(1) lookup if indices are specified.
	indexMap := make(map[uint64]bool)
	if len(indices) > 0 {
		for _, idx := range indices {
			indexMap[idx] = true
		}
	}

	// Preallocate response slice - if indices specified, size will be len(indices),
	// otherwise size will be all sidecars>
	responseSize := len(blobSidecars)
	if len(indices) > 0 {
		responseSize = len(indices)
	}
	blobSidecarsResponse := make([]*apitypes.Sidecar, 0, responseSize)

	for _, blobSidecar := range blobSidecars {
		// Skip if indices specified and this index not requested.
		if len(indices) > 0 && !indexMap[blobSidecar.GetIndex()] {
			continue
		}
		// Craft and append the blob sidecar serialized data to the response.
		blobSidecarsResponse = append(blobSidecarsResponse,
			apitypes.SidecarFromConsensus(blobSidecar),
		)
	}
	return blobSidecarsResponse, nil
}
