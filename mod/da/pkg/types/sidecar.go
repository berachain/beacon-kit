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

package types

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

// BlobSidecar as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/p2p-interface.md?ref=bankless.ghost.io#blobsidecar
//
//nolint:lll
type BlobSidecar struct {
	// Index represents the index of the blob in the block.
	Index uint64
	// Blob represents the blob data.
	Blob eip4844.Blob
	// KzgCommitment is the KZG commitment of the blob.
	KzgCommitment eip4844.KZGCommitment
	// Kzg proof allows folr the verification of the KZG commitment.
	KzgProof eip4844.KZGProof
	// BeaconBlockHeader represents the beacon block header for which this blob
	// is being included.
	BeaconBlockHeader *types.BeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block
	// body.
	InclusionProof [][32]byte `ssz-size:"8,32"`
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	index math.U64,
	header *types.BeaconBlockHeader,
	blob *eip4844.Blob,
	commitment eip4844.KZGCommitment,
	proof eip4844.KZGProof,
	inclusionProof [][32]byte,
) *BlobSidecar {
	return &BlobSidecar{
		Index:             index.Unwrap(),
		Blob:              *blob,
		KzgCommitment:     commitment,
		KzgProof:          proof,
		BeaconBlockHeader: header,
		InclusionProof:    inclusionProof,
	}
}

// HasValidInclusionProof verifies the inclusion proof of the
// blob in the beacon body.
func (b *BlobSidecar) HasValidInclusionProof(
	kzgOffset uint64,
) bool {
	// Calculate the hash tree root of the KZG commitment.
	leaf, err := b.KzgCommitment.HashTreeRoot()
	if err != nil {
		return false
	}

	gIndex := kzgOffset + b.Index

	// Verify the inclusion proof.
	return merkle.IsValidMerkleBranch(
		leaf,
		b.InclusionProof,
		//#nosec:G701 // safe.
		uint8(
			len(b.InclusionProof),
		), // TODO: use KZG_INCLUSION_PROOF_DEPTH calculation.
		gIndex,
		b.BeaconBlockHeader.BodyRoot,
	)
}
