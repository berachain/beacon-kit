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
	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/cockroachdb/errors"
)

const (
	MaxBlobCommitmentsPerBlock = 16
	// KZGMerkleIndex is the merkle index of BlobKzgCommitments' root
	// in the merkle tree built from the block body.
	KZGMerkleIndex        = 24
	KZGOffset      uint64 = KZGMerkleIndex * MaxBlobCommitmentsPerBlock
)

// MerkleProofKZGCommitment generates a Merkle proof for a given index in a list
// of commitments using the KZG algorithm.
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
		LogBodyLengthDeneb,
	)
	if err != nil {
		return nil, err
	}

	topProof, err := tree.MerkleProof(KZGPositionDeneb)
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
	bodyTree, err := merkle.NewTreeWithMaxLeaves(
		LeavesFromCommitments(commitments),
		MaxBlobCommitmentsPerBlock,
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
