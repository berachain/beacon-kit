// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/berachain/beacon-kit/mod/primitives/merkle"
)

// BlobSidecar as per the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/deneb/p2p-interface.md?ref=bankless.ghost.io#blobsidecar
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen -path ./sidecar.go -objs BlobSidecar -include ../../primitives/math,../../primitives/kzg,../../primitives,../../primitives,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output sidecar.ssz.go
type BlobSidecar struct {
	// Index represents the index of the blob in the block.
	Index uint64
	// Blob represents the blob data.
	Blob kzg.Blob
	// KzgCommitment is the KZG commitment of the blob.
	KzgCommitment kzg.Commitment
	// Kzg proof allows folr the verification of the KZG commitment.
	KzgProof kzg.Proof
	// BeaconBlockHeader represents the beacon block header for which this blob
	// is being included.
	BeaconBlockHeader *primitives.BeaconBlockHeader
	// InclusionProof is the inclusion proof of the blob in the beacon block
	// body.
	InclusionProof [][32]byte `ssz-size:"8,32"`
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	index uint64,
	header *primitives.BeaconBlockHeader,
	blob *kzg.Blob,
	commitment kzg.Commitment,
	proof kzg.Proof,
	inclusionProof [][32]byte,
) *BlobSidecar {
	return &BlobSidecar{
		Index:             index,
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
