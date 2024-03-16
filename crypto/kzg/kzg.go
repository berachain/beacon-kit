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

package kzg

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/beacon/core/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/merkle"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/sourcegraph/conc/iter"
)

var (
	errInvalidIndex          = errors.New("index out of bounds")
	errInvalidBodyRoot       = errors.New("invalid Beacon Block Body root")
	errInvalidInclusionProof = errors.New("invalid KZG commitment inclusion proof")
	errNilBlockHeader        = errors.New("received nil beacon block header")
)

// ConvertCommitmentToVersionedHash computes a SHA-256 hash of the given
// commitment and prefixes it with the BlobCommitmentVersion. This function is
// used to generate
// a versioned hash for KZG commitments.
//
// The restulting hash is intended for use in contexts where a versioned
// identifier
// for the commitment is required.
func ConvertCommitmentToVersionedHash(
	commitment [48]byte,
) primitives.ExecutionHash {
	hash := sha256.Sum256(commitment[:])
	// Prefix the hash with the BlobCommitmentVersion to create a versioned
	// hash.
	hash[0] = BlobCommitmentVersion
	return hash
}

// ConvertCommitmentsToVersionedHashes converts a slice of commitments to a
// slice of versioned hashes. This function is used to generate versioned hashes
// for KZG commitments.
//
// The resulting hashes are intended for use in contexts where a versioned
// identifier
// for the commitments is required.
func ConvertCommitmentsToVersionedHashes(
	commitments [][48]byte,
) []primitives.ExecutionHash {
	return iter.Map(commitments, func(bz *[48]byte) primitives.ExecutionHash {
		return ConvertCommitmentToVersionedHash(*bz)
	})
}

// VerifyKZGInclusionProof verifies the inclusion proof for a commitment in a Merkle tree.
// It takes the commitment, root hash, inclusion proof, and index as input parameters.
// The commitment is the value being proven to be included in the Merkle tree.
// The root is the root hash of the Merkle tree.
// The proof is a list of intermediate hashes that prove the inclusion of the commitment in the Merkle tree.
// The index is the position of the commitment in the Merkle tree.
// If the inclusion proof is valid, the function returns nil.
// Otherwise, it returns an error indicating an invalid inclusion proof.
func VerifyKZGInclusionProof(root []byte, blob *types.BlobSidecar, index uint64) error { // TODO: add wrapped type with inclusion proofs
	if len(root) != rootLength {
		return errInvalidBodyRoot
	}
	chunks := make([][32]byte, 2)
	copy(chunks[0][:], blob.KzgCommitment)
	copy(chunks[1][:], blob.KzgCommitment[rootLength:])
	gohashtree.HashChunks(chunks, chunks)
	verified := merkle.VerifyMerkleProof(root, chunks[0][:], index+KZGOffset, blob.InclusionProof)
	if !verified {
		return errInvalidInclusionProof
	}
	return nil
}

// MerkleProofKZGCommitment generates a Merkle proof for a given index in a list of commitments using the KZG algorithm.
// It takes a 2D byte slice of commitments and an index as input, and returns a 2D byte slice representing the Merkle proof.
// If an error occurs during the generation of the proof, it returns nil and the error.
// The function internally calls the `bodyProof` function to generate the body proof, and the `topLevelRoots` function to obtain the top level roots.
// It then uses the `merkle.GenerateTrieFromItems` function to generate a sparse Merkle tree from the top level roots.
// Finally, it calls the `MerkleProof` method on the sparse Merkle tree to obtain the top proof, and appends it to the body proof.
// Note that the last element of the top proof is removed before returning the final proof, as it is not needed.
func MerkleProofKZGCommitment(blk beacontypes.BeaconBlock, index int) ([][]byte, error) {
	commitments := blk.GetBody().GetBlobKzgCommitments()

	cmts := make([][]byte, len(commitments))
	fmt.Println("PRE COMMITMENTS", len(commitments))
	for i, c := range commitments {
		cmts[i] = c[:]
	}

	fmt.Println("PRE BODY PROOF")

	fmt.Println("LENGTH COMMITMENTS", len(cmts), "INDEX", index)
	proof, err := bodyProof(cmts, index)
	if err != nil {
		return nil, err
	}

	fmt.Println("PRE TOP LEVEL ROOTS")

	membersRoots, err := blk.GetBody().GetTopLevelRoots()
	if err != nil {
		return nil, err
	}

	fmt.Println("PAST TOP LEVEL ROOTS")

	sparse, err := merkle.GenerateTrieFromItems(membersRoots, logBodyLength)
	if err != nil {
		return nil, err
	}

	fmt.Println("PRE MERKLE PROOF")
	topProof, err := sparse.MerkleProof(kzgPosition)
	if err != nil {
		return nil, err
	}
	// sparse.MerkleProof always includes the length of the slice this is
	// why we remove the last element that is not needed in topProof
	proof = append(proof, topProof[:len(topProof)-1]...)
	return proof, nil
}

// bodyProof returns the Merkle proof of the subtree up to the root of the KZG
// commitment list.
func bodyProof(commitments [][]byte, index int) ([][]byte, error) {
	if index < 0 || index >= len(commitments) {
		return nil, errors.New("index out of range")
	}
	fmt.Println("PRE LEAVES FROM COMMITMENTS")
	leaves := leavesFromCommitments(commitments)
	fmt.Println("PRE MERKLE GENERATE TRIE FROM ITEMS")
	sparse, err := merkle.GenerateTrieFromItems(leaves, logMaxBlobCommitments)
	if err != nil {
		return nil, err
	}

	fmt.Println("PRE MERKLE PROOF")
	proof, err := sparse.MerkleProof(index)
	if err != nil {
		return nil, err
	}
	return proof, err
}

// leavesFromCommitments hashes each commitment to construct a slice of roots
func leavesFromCommitments(commitments [][]byte) [][]byte {
	leaves := make([][]byte, len(commitments))
	for i, kzg := range commitments {
		chunk := make([][32]byte, 2)
		copy(chunk[0][:], kzg)
		copy(chunk[1][:], kzg[rootLength:])
		gohashtree.HashChunks(chunk, chunk)
		leaves[i] = chunk[0][:]
	}
	return leaves
}
