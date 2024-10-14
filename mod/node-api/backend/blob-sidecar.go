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
	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
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

	// blobSidecars, err := b.getBlobSidecarsFromBlockRoot(
	//	blockHeader.GetBodyRoot(),
	//	slot,
	// )
	// if err != nil {
	//	return nil, err
	// }

	// TODO: Implement with real data.
	blobSidecars := []*types.BlobSidecarData[BeaconBlockHeaderT]{
		{
			Index: 0,
			Blob: eip4844.Blob{
				0x62, 0x6c, 0x6f, 0x62, 0x31, // "blob1" in hex
			},
			KzgCommitment: eip4844.KZGCommitment{
				0x6b, 0x7a, 0x67, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
				0x6d, 0x65, 0x6e, 0x74, // "kzg_commitment" in hex
			},
			KzgProof: eip4844.KZGProof{
				0x6b, 0x7a, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, // "kzg_proof" in hex
			},
			BeaconBlockHeader: types.BlockHeader[BeaconBlockHeaderT]{
				Message:   blockHeader,
				Signature: crypto.BLSSignature{}, // TODO: Implement signature.
			},
			KzgCommitmentInclusionProof: []common.Root{
				{
					0x69, 0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x6f, 0x6e, // "inclusion" in hex
				},
			},
		},
	}
	return blobSidecars, nil
}
