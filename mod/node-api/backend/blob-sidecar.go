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

	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BlobSidecarsAtSlot returns the blob sidecars at a given slot.
func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlobSidecarsAtSlot(
	slot math.Slot,
) ([]*types.BlobSidecarData[BeaconBlockHeaderT], error) {
	blockHeader, err := b.BlockHeaderAtSlot(slot)
	if err != nil {
		return nil, err
	}

	blobSidecars, err := b.sb.AvailabilityStore().GetBlobSideCars(slot)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob sidecars: %w", err)
	}

	// Convert the returned blobSidecars to
	// []*beacontypes.BlobSidecarData[BeaconBlockHeaderT]
	result := make(
		[]*types.BlobSidecarData[BeaconBlockHeaderT],
		len(*blobSidecars),
	)
	for i, sidecar := range *blobSidecars {
		result[i] = &types.BlobSidecarData[BeaconBlockHeaderT]{
			Index:         sidecar.GetIndex(),
			Blob:          sidecar.GetBlob(),
			KzgCommitment: sidecar.GetKzgCommitment(),
			KzgProof:      sidecar.GetKzgProof(),
			BeaconBlockHeader: &types.BlockHeader[BeaconBlockHeaderT]{
				Message:   blockHeader,
				Signature: crypto.BLSSignature{}, // TODO: Implement signature.
			},
			// sidecar.GetBeaconBlockHeader(),
			KzgCommitmentInclusionProof: sidecar.GetInclusionProof(),
		}
	}

	return result, nil
}
