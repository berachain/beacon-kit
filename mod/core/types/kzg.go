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
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/cockroachdb/errors"
)

const (
	Two                        = 2
	RootLength                 = 32
	MaxBlobCommitmentsPerBlock = 16
	// LogMaxBlobCommitments is the Log_2 of MaxBlobCommitmentsPerBlock (16).
	LogMaxBlobCommitments uint64 = 4
	// If you are adding values to the BeaconBlockBodyDeneb struct,
	// the body length must be increased and GetTopLevelRoots updated.
	BodyLength = 5
	// LogBodyLength is the Log_2 of BodyLength (6).
	LogBodyLength = 3
	// KZGPosition is the position of BlobKzgCommitments in the block body.
	KZGPosition = 4
	// KZGMerkleIndex is the merkle index of BlobKzgCommitments' root
	// in the merkle tree built from the block body.
	KZGMerkleIndex        = 24
	KZGOffset      uint64 = KZGMerkleIndex * MaxBlobCommitmentsPerBlock
)

// VerifyKZGInclusionProof verifies the inclusion proof for a commitment in a
// Merkle tree. It takes the commitment, root hash, inclusion proof, and index
// as input parameters.
// The commitment is the value being proven to be included in the Merkle tree.
// The root is the root hash of the Merkle tree.
// The proof is a list of intermediate hashes that prove the inclusion of the
// commitment in the Merkle tree.
// The index is the position of the commitment in the Merkle tree.
// If the inclusion proof is valid, the function returns nil.
// Otherwise, it returns an error indicating an invalid inclusion proof.
func VerifyKZGInclusionProof(
	blob *datypes.BlobSidecar,
) error { // TODO: add wrapped type with inclusion proofs
	verified := merkle.VerifyMerkleProof(
		blob.BeaconBlockHeader.BodyRoot,
		blob.KzgCommitment.ToHashChunks()[0],
		blob.Index+KZGOffset,
		blob.InclusionProof,
	)
	if !verified {
		return errInvalidInclusionProof
	}
	return nil
}

// MerkleProofKZGCommitment generates a Merkle proof for a given index in a list
// of commitments using the KZG algorithm. It takes a 2D byte slice of
// commitments and an index as input, and returns a 2D byte slice representing
// the Merkle proof. If an error occurs during the generation of the proof, it
// returns nil and the error. The function internally calls the `BodyProof`
// function to generate the body proof, and the `topLevelRoots` function to
// obtain the top level roots. It then uses the `merkle.NewTreeFromLeaves`
// function to generate a sparse Merkle tree from the top level roots. Finally,
// it calls the `MerkleProof` method on the sparse Merkle tree to obtain the top
// proof, and appends it to the body proof. Note that the last element of the
// top proof is removed before returning the final proof, as it is not needed.
func MerkleProofKZGCommitment(
	body BeaconBlockBody,
	index uint64,
) ([][32]byte, error) {
	commitments := body.GetBlobKzgCommitments()

	proof, err := BodyProof(commitments, index)
	if err != nil {
		return nil, err
	}

	membersRoots, err := GetTopLevelRoots(body)
	if err != nil {
		return nil, err
	}

	tree, err := merkle.NewTreeFromLeavesWithDepth(
		membersRoots,
		LogBodyLength,
	)
	if err != nil {
		return nil, err
	}

	topProof, err := tree.MerkleProof(KZGPosition)
	if err != nil {
		return nil, err
	}
	return append(proof, topProof...), nil
}

// BodyProof returns the Merkle proof of the subtree up to the root of the KZG
// commitment list.
func BodyProof(commitments kzg.Commitments, index uint64) ([][32]byte, error) {
	if index >= uint64(len(commitments)) {
		return nil, errors.New("index out of range")
	}
	leaves := LeavesFromCommitments(commitments)
	bodyTree, err := merkle.NewTreeFromLeavesWithDepth(
		leaves,
		LogMaxBlobCommitments,
	)
	if err != nil {
		return nil, err
	}

	return bodyTree.MerkleProofWithMixin(index)
}

// LeavesFromCommitments hashes each commitment to construct a slice of roots.
func LeavesFromCommitments(commitments kzg.Commitments) [][32]byte {
	leaves := make([][32]byte, len(commitments))
	for i, commitment := range commitments {
		leaves[i] = commitment.ToHashChunks()[0]
	}
	return leaves
}
