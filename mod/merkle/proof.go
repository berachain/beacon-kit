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

// VerifyProof given a tree root, a leaf, the generalized merkle index
// of the leaf in the tree, and the proof itself.
func VerifyProof(
	root, leaf [32]byte,
	merkleIndex uint64,
	proof [][32]byte,
) bool {
	//#nosec:G701 `int`` is at minimum 32-bits and thus a
	// uint8 will always fit.
	if len(proof) > int(^uint8(0)) {
		return false
	}
	return IsValidMerkleBranch(
		leaf,
		proof,
		//#nosec:G701 // we check the length of the proof above.
		uint8(len(proof)),
		merkleIndex,
		root,
	)
}

// IsValidMerkleBranch as per the Ethereum 2.0 spec:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_valid_merkle_branch
//
//nolint:lll
func IsValidMerkleBranch(
	leaf [32]byte, branch [][32]byte, depth uint8, index uint64, root [32]byte,
) bool {
	//#nosec:G701 `int`` is at minimum 32-bits and thus a
	// uint8 will always fit.
	if len(branch) != int(depth) {
		return false
	}
	return RootFromBranch(leaf, branch, depth, index) == root
}

// RootFromBranch calculates the Merkle root from a leaf and a branch.
// Inspired by:
// https://github.com/sigp/lighthouse/blob/2cd0e609f59391692b4c8e989e26e0dac61ff801/consensus/merkle_proof/src/lib.rs#L357
//
//nolint:lll
func RootFromBranch(
	leaf [32]byte,
	branch [][32]byte,
	depth uint8,
	index uint64,
) [32]byte {
	merkleRoot := leaf
	var hashInput [64]byte
	for i := range depth {
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
