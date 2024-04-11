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

package merkle

import (
	"encoding/binary"
	"fmt"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	sha256 "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"
)

// VerifyMerkleProof given a tree root, a leaf, the generalized merkle index
// of the leaf in the tree, and the proof itself.
func VerifyMerkleProof(
	root, leaf [32]byte,
	merkleIndex uint64,
	proof [][32]byte,
) bool {
	if len(proof) == 0 {
		return false
	}
	return VerifyMerkleProofWithDepth(
		root,
		leaf,
		merkleIndex,
		proof,
		uint64(len(proof))-1,
	)
}

// VerifyMerkleProofWithDepth verifies a Merkle branch against a root of a tree.
func VerifyMerkleProofWithDepth(
	root, item [32]byte,
	merkleIndex uint64,
	proof [][32]byte,
	depth uint64,
) bool {
	if uint64(len(proof)) != depth+1 {
		return false
	}
	for i := uint64(0); i <= depth; i++ {
		if (merkleIndex & 1) == 1 {
			item = sha256.Sum256(append(proof[i][:], item[:]...))
		} else {
			item = sha256.Sum256(append(item[:], proof[i][:]...))
		}
		merkleIndex /= 2
	}
	return root == item
}

// MerkleProof computes a proof from a tree's branches using a Merkle index.
func (m *SparseMerkleTree) MerkleProof(index uint64) ([][32]byte, error) {
	numLeaves := uint64(len(m.branches[0]))
	if index >= numLeaves {
		return nil, fmt.Errorf(
			"merkle index out of range in tree, max range: %d, received: %d",
			numLeaves,
			index,
		)
	}
	proof := make([][32]byte, m.depth)
	for i := uint64(0); i < m.depth; i++ {
		subIndex := (index >> i) ^ 1
		layer := m.branches[i]
		if subIndex < uint64(len(layer)) {
			proof[i] = byteslib.ToBytes32(layer[subIndex])
		} else {
			proof[i] = tree.ZeroHashes[i]
		}
	}
	return proof, nil
}

// MerkleProofWithMixin computes a proof from a tree's branches using a Merkle
// index.
func (m *SparseMerkleTree) MerkleProofWithMixin(
	index uint64,
) ([][32]byte, error) {
	proof, err := m.MerkleProof(index)
	if err != nil {
		return nil, err
	}

	var mixin [32]byte
	binary.LittleEndian.PutUint64(mixin[:], uint64(len(m.originalItems)))
	return append(proof, mixin), nil
}
