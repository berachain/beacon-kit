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
	"unsafe"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/berachain/beacon-kit/mod/trie"
	"github.com/berachain/beacon-kit/mod/trie/merkleize"
	"github.com/prysmaticlabs/gohashtree"
)

// Commitments represents a slice of KZG commitments.
// TODO: Must be [48]byte for fastssz to work, annoying.
// TODO: Fix this with upstream PR to fastssz.
type Commitments [][48]byte

// HashTreeRoot returns the hash tree root of the commitments.
func (c Commitments) HashTreeRoot() ([32]byte, error) {
	return merkleize.VectorSSZ(
		*(*[]Commitment)(unsafe.Pointer(&c)), uint64(len(c)))
}

// MerkleProof returns the merkle proof for the commitment at the given index.
func (c Commitments) MerkleProof(index uint64) ([][]byte, error) {
	// Generate the leaves for the sparse merkle tree.
	leaves := make([][]byte, len(c))
	for i, commitment := range c {
		leaf, err := Commitment(commitment).HashTreeRoot()
		if err != nil {
			return [][]byte{}, err
		}
		leaves[i] = leaf[:]
	}

	// Build a sparse merkle tree from the leaves.
	sparse, err := trie.NewFromItems(leaves, constants.LogMaxBlobCommitments)
	if err != nil {
		return nil, err
	}

	// Generate the proof for the given index.
	return sparse.MerkleProof(index)
}

// Commitment is a KZG commitment.
type Commitment [48]byte

// HashTreeRoot returns the hash tree root of the commitment.
func (c Commitment) HashTreeRoot() ([32]byte, error) {
	//nolint:gomnd // two is okay.
	chunk := make([][32]byte, 2)
	copy(chunk[0][:], c[:])
	copy(chunk[1][:], c[constants.RootLength:])
	gohashtree.HashChunks(chunk, chunk)
	return chunk[0], nil
}

// ToVersionedHash converts this KZG commitment into a versioned hash.
func (c Commitment) ToVersionedHash() primitives.ExecutionHash {
	hash := sha256.Sum256(c[:])
	// Prefix the hash with the BlobCommitmentVersion
	// to create a versioned hash.
	hash[0] = constants.BlobCommitmentVersion
	return hash
}

// CommitmentsToVersionedHashes converts a slice of commitments to a
// slice of versioned hashes. This function is used to generate versioned hashes
// for KZG commitments.
//
// The resulting hashes are intended for use in contexts where a versioned
// identifier
// for the commitments is required.
func CommitmentsToVersionedHashes(
	commitments Commitments,
) []primitives.ExecutionHash {
	hashes := make([]primitives.ExecutionHash, len(commitments))
	for i, bz := range commitments {
		hashes[i] = Commitment(bz).ToVersionedHash()
	}
	return hashes
}
