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
	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BlobSidecarsAtSlot returns the blob sidecars at the given slot.
func (b Backend[
	_, _, _, BeaconBlockHeaderT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) BlobSidecarsAtSlot(slot math.Slot, indices []uint64) ([]*beacontypes.BlobSidecarData[BeaconBlockHeaderT], error) {
	blobSidecars, err := b.sb.AvailabilityStore().GetBlobsFromStore(slot)
	if err != nil {
		return nil, err
	}

	// Create a map of requested indices for O(1) lookup if indices are specified
	indexMap := make(map[uint64]bool)
	if len(indices) > 0 {
		for _, idx := range indices {
			indexMap[idx] = true
		}
	}

	// Initialize response slice - if indices specified, size will be len(indices),
	// otherwise size will be all sidecars
	responseSize := blobSidecars.Len()
	if len(indices) > 0 {
		responseSize = len(indices)
	}
	blobSidecarsResponse := make([]*beacontypes.BlobSidecarData[BeaconBlockHeaderT], 0, responseSize)

	// Process each sidecar
	for i := 0; i < blobSidecars.Len(); i++ {
		blobSidecar := blobSidecars.Get(i)

		// Skip if indices specified and this index not requested
		if len(indices) > 0 && !indexMap[blobSidecar.GetIndex()] {
			continue
		}

		blobHex, err := blobSidecar.GetBlob().MarshalText()
		if err != nil {
			return nil, err
		}
		kzgCommitmentHex, err := blobSidecar.GetKzgCommitment().MarshalText()
		if err != nil {
			return nil, err
		}
		kzgProofHex, err := blobSidecar.GetKzgProof().MarshalText()
		if err != nil {
			return nil, err
		}
		inclusionProof := blobSidecar.GetInclusionProof()
		inclusionProofStrings := make([]string, len(inclusionProof))
		for j, proof := range inclusionProof {
			inclusionProofStrings[j] = proof.String()
		}

		blobSidecarsResponse = append(blobSidecarsResponse, &beacontypes.BlobSidecarData[BeaconBlockHeaderT]{
			Index: blobSidecar.GetIndex(),
			Blob:  string(blobHex),
			SignedBlockHeader: &beacontypes.BlockHeader[BeaconBlockHeaderT]{
				Message:   blobSidecar.GetBeaconBlockHeader(),
				Signature: bytes.B48{}, // TODO: implement
			},
			KZGCommitment:               string(kzgCommitmentHex),
			KZGProof:                    string(kzgProofHex),
			KZGCommitmentInclusionProof: inclusionProofStrings,
		})
	}

	return blobSidecarsResponse, nil
}
