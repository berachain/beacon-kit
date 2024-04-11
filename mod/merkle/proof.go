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
	sha256 "github.com/minio/sha256-simd"
)

// VerifyMerkleProof given a tree root, a leaf, the generalized merkle index
// of the leaf in the tree, and the proof itself.
func VerifyMerkleProof(
	root, leaf [32]byte,
	merkleIndex uint64,
	proof [][32]byte,
) bool {
	return VerifyMerkleProofWithDepth(
		root,
		leaf,
		merkleIndex,
		proof,
		uint64(len(proof)),
	)
}

// VerifyMerkleProofWithDepth verifies a Merkle branch against a root of a tree.
func VerifyMerkleProofWithDepth(
	root, leaf [32]byte,
	index uint64,
	proof [][32]byte,
	depth uint64,
) bool {
	if uint64(len(proof)) != depth {
		return false
	}
	return merkleRootFromBranch(leaf, proof, depth, index) == root
}

// Compute a root hash from a leaf and a Merkle proof.
func merkleRootFromBranch(
	leaf [32]byte,
	branch [][32]byte,
	depth uint64,
	index uint64,
) [32]byte {
	if uint64(len(branch)) != depth {
		panic("proof length should equal depth")
	}
	merkleRoot := leaf
	var hashInput [64]byte
	for i := uint64(0); i < depth; i++ {
		//nolint:gomnd // from spec.
		ithBit := (index >> i) & 0x01
		if ithBit == 1 {
			copy(hashInput[:32], branch[i][:])
			copy(hashInput[32:], merkleRoot[:])
		} else {
			copy(hashInput[:32], merkleRoot[:])
			copy(hashInput[32:], branch[i][:])
		}
		merkleRoot = sha256.Sum256(hashInput[:])
	}
	return merkleRoot
}

// def is_valid_merkle_branch(leaf: Bytes32, branch: Sequence[Bytes32], depth:
// uint64, index: uint64, root: Root) -> bool:
//     """
// Check if ``leaf`` at ``index`` verifies against the Merkle ``root`` and
// ``branch``.
//     """
//     value = leaf
//     for i in range(depth):
//         if index // (2**i) % 2:
//             value = hash(branch[i] + value)
//         else:
//             value = hash(value + branch[i])
//     return value == root
