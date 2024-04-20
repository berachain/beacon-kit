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

package tree

import (
	"crypto/sha256"
	"sort"

	"github.com/berachain/beacon-kit/mod/primitives"
)

// GetBranchIndices returns the generalized indices of the sister chunks along
// the path from the chunk with the
// given tree index to the root.
func GetBranchIndices(treeIndex GeneralizedIndex) []GeneralizedIndex {
	var o []GeneralizedIndex
	o = append(o, GeneralizedIndexSibling(treeIndex))
	for o[len(o)-1] > 1 {
		o = append(
			o,
			GeneralizedIndexSibling(GeneralizedIndexParent(o[len(o)-1])),
		)
	}
	if len(o) > 1 {
		return o[:len(o)-1]
	}
	return o
}

// GetPathIndices returns the generalized indices of the chunks along the path
// from the chunk with the
// given tree index to the root.
func GetPathIndices(treeIndex GeneralizedIndex) []GeneralizedIndex {
	var o []GeneralizedIndex
	o = append(o, treeIndex)
	for o[len(o)-1] > 1 {
		o = append(o, GeneralizedIndexParent(o[len(o)-1]))
	}
	if len(o) > 1 {
		return o[:len(o)-1]
	}
	return o
}

// GetHelperIndices returns the generalized indices of all "extra" chunks in the
// tree needed to prove the chunks with the given generalized indices. Note that
// the decreasing order is chosen deliberately to ensure equivalence to the
// order of hashes in a regular single-item Merkle proof in the single-item
// case.
func GetHelperIndices(indices []GeneralizedIndex) []GeneralizedIndex {
	allHelperIndices := make(map[GeneralizedIndex]struct{})
	allPathIndices := make(map[GeneralizedIndex]struct{})

	for _, index := range indices {
		for _, idx := range GetBranchIndices(index) {
			allHelperIndices[idx] = struct{}{}
		}
		for _, idx := range GetPathIndices(index) {
			allPathIndices[idx] = struct{}{}
		}
	}

	var diff []GeneralizedIndex
	for idx := range allHelperIndices {
		if _, exists := allPathIndices[idx]; !exists {
			diff = append(diff, idx)
		}
	}

	sort.Slice(diff, func(i, j int) bool {
		return diff[i] > diff[j]
	})

	return diff
}

// CalculateMerkleRoot calculates the Merkle root from a leaf and a proof based
// on the generalized index.
func CalculateMerkleRoot(
	leaf primitives.Bytes32,
	proof []primitives.Bytes32,
	index GeneralizedIndex,
) primitives.Bytes32 {
	if len(proof) != GetGeneralizedIndexLength(index) {
		panic("proof length does not match the expected length from index")
	}
	for i, h := range proof {
		if GetGeneralizedIndexBit(index, i) {
			leaf = sha256.Sum256(append(h[:], leaf[:]...))
		} else {
			leaf = sha256.Sum256(append(leaf[:], h[:]...))
		}
	}
	return leaf
}

// VerifyMerkleProof verifies a Merkle proof for a single leaf and a given root.
func VerifyMerkleProof(
	leaf primitives.Bytes32,
	proof []primitives.Bytes32,
	index GeneralizedIndex,
	root primitives.Bytes32,
) bool {
	return CalculateMerkleRoot(leaf, proof, index) == root
}

// CalculateMultiMerkleRoot calculates the Merkle root for multiple leaves with
// a given proof and indices.
func CalculateMultiMerkleRoot(
	leaves []primitives.Bytes32,
	proof []primitives.Bytes32,
	indices []GeneralizedIndex,
) primitives.Bytes32 {
	if len(leaves) != len(indices) {
		panic("number of leaves does not match number of indices")
	}
	helperIndices := GetHelperIndices(indices)
	if len(proof) != len(helperIndices) {
		panic("proof length does not match helper indices length")
	}

	objects := make(map[GeneralizedIndex]primitives.Bytes32)
	for i, leaf := range leaves {
		objects[indices[i]] = leaf
	}
	for i, h := range proof {
		objects[helperIndices[i]] = h
	}

	keys := make([]GeneralizedIndex, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	pos := 0
	for pos < len(keys) {
		k := keys[pos]
		if _, exists := objects[k]; exists &&
			objects[GeneralizedIndex(k^1)] != (primitives.Bytes32{}) &&
			objects[GeneralizedIndex(k/2)] == (primitives.Bytes32{}) {
			leftChild := objects[GeneralizedIndex((k|1)^1)]
			rightChild := objects[GeneralizedIndex(k|1)]
			hashed := sha256.Sum256(append(leftChild[:], rightChild[:]...))
			parentIndex := GeneralizedIndex(k / 2)
			objects[parentIndex] = hashed
			keys = append(keys, parentIndex)
		}
		pos++
	}
	return objects[GeneralizedIndex(1)]
}

// VerifyMerkleMultiproof verifies a Merkle multiproof for multiple leaves and a
// given root.
func VerifyMerkleMultiproof(
	leaves []primitives.Bytes32,
	proof []primitives.Bytes32,
	indices []GeneralizedIndex,
	root primitives.Bytes32,
) bool {
	return CalculateMultiMerkleRoot(leaves, proof, indices) == root
}
