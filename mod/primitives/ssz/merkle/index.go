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
	"crypto/sha256"
	"errors"
	"sort"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
)

type (
	// GeneralizedIndex is a generalized index.
	GeneralizedIndex math.U64

	// GeneralizedIndicies is a list of generalized indices.
	GeneralizedIndicies []GeneralizedIndex
)

// NewGeneralizedIndex calculates the generalized index from the depth and
// index.
func NewGeneralizedIndex(depth uint8, index uint64) GeneralizedIndex {
	return GeneralizedIndex((1 << depth) + index)
}

// Length returns the length of the generalized index.
func (g GeneralizedIndex) Length() uint64 {
	return uint64(math.U64(g).ILog2Ceil())
}

// IndexBit returns the bit at the specified position in a generalized index.
func (g GeneralizedIndex) IndexBit(position int) bool {
	return (g & (1 << position)) > 0
}

// Sibling returns the sibling index of the current generalized index.
func (g GeneralizedIndex) Sibling() GeneralizedIndex {
	return g ^ 1
}

// LeftChild returns the left child index of the current generalized index.
//
//nolint:mnd // from spec.
func (g GeneralizedIndex) LeftChild() GeneralizedIndex {
	return 2 * g
}

// RightChild returns the right child index of the current generalized index.
func (g GeneralizedIndex) RightChild() GeneralizedIndex {
	return 2*g + 1
}

// Parent returns the parent index of the current generalized index.
//
//nolint:mnd // from spec.
func (g GeneralizedIndex) Parent() GeneralizedIndex {
	return g / 2
}

// GetBranchIndices returns the generalized indices of the nodes on the path
// from the root to the leaf.
func (g GeneralizedIndex) GetBranchIndices() GeneralizedIndicies {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the
	// given tree index to the root.
	o := []GeneralizedIndex{g.Sibling()}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent().Sibling())
	}
	return o[:len(o)-1]
}

// GetPathIndices returns the generalized indices of the nodes on the path from
// the leaf to the root.
func (g GeneralizedIndex) GetPathIndices() GeneralizedIndicies {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the
	// given tree index to the root.
	o := []GeneralizedIndex{g}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent())
	}
	return o[:len(o)-1]
}

// CalculateMerkleRoot calculates the Merkle root from the leaf and proof.
func (g GeneralizedIndex) CalculateMerkleRoot(
	leaf primitives.Bytes32,
	proof []primitives.Bytes32,
) primitives.Root {
	if len(proof) != int(g.Length()) {
		panic("proof length does not match index length")
	}
	for i, h := range proof {
		if g.IndexBit(i) {
			leaf = sha256.Sum256(append(h[:], leaf[:]...))
		} else {
			leaf = sha256.Sum256(append(leaf[:], h[:]...))
		}
	}
	return leaf
}

// VerifyMerkleProof verifies the Merkle proof for the given
// leaf, proof, and root.
func (g GeneralizedIndex) VerifyMerkleProof(
	leaf primitives.Bytes32,
	proof []primitives.Bytes32,
	root primitives.Root,
) bool {
	return g.CalculateMerkleRoot(leaf, proof) == root
}

// Concatenates multiple generalized indices into a single generalized index
// representing the path from the first to the last node.
func (gs GeneralizedIndicies) Concat() GeneralizedIndex {
	o := GeneralizedIndex(1)
	for _, i := range gs {
		floorPower := math.U64(i).PrevPowerOfTwo()
		o = GeneralizedIndex(
			math.U64(o)*floorPower + (math.U64(i) - floorPower),
		)
	}
	return o
}

// GetHelperIndices returns the generalized indices of all "extra" chunks in the
// tree needed to prove the chunks with the given generalized indices. The
// decreasing order is chosen deliberately to ensure equivalence to the order of
// hashes in a regular single-item Merkle proof in the single-item case.
func (gs GeneralizedIndicies) GetHelperIndices() GeneralizedIndicies {
	allHelperIndices := make(map[GeneralizedIndex]struct{})
	allPathIndices := make(map[GeneralizedIndex]struct{})

	for _, index := range gs {
		for _, helperIndex := range index.GetBranchIndices() {
			allHelperIndices[helperIndex] = struct{}{}
		}
		for _, pathIndex := range index.GetPathIndices() {
			allPathIndices[pathIndex] = struct{}{}
		}
	}

	difference := make([]GeneralizedIndex, 0, len(allHelperIndices))
	for helperIndex := range allHelperIndices {
		if _, exists := allPathIndices[helperIndex]; !exists {
			difference = append(difference, helperIndex)
		}
	}

	// Sort in decreasing order
	sort.Slice(difference, func(i, j int) bool {
		return difference[i] > difference[j]
	})

	return difference
}

// CalculateMultiMerkleRoot calculates the Merkle root for multiple leaves with
// their corresponding proofs and indices.
func (gs GeneralizedIndicies) CalculateMultiMerkleRoot(
	leaves []primitives.Bytes32,
	proof []primitives.Bytes32,
) (primitives.Root, error) {
	if len(leaves) != len(gs) {
		return primitives.Root{}, errors.New(
			"mismatched leaves and indices length",
		)
	}

	helperIndices := gs.GetHelperIndices()
	if len(proof) != len(helperIndices) {
		return primitives.Root{}, errors.New(
			"mismatched proof and helper indices length",
		)
	}

	objects := make(map[GeneralizedIndex]primitives.Bytes32)
	for i, index := range gs {
		objects[index] = leaves[i]
	}
	for i, index := range helperIndices {
		objects[index] = proof[i]
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
		if _, ok := objects[k]; ok {
			if sibling, ok := objects[k^1]; ok {
				if _, ok := objects[k/2]; !ok {
					obj := objects[(k|1)^1]
					objects[k/2] = sha256.Sum256(append(obj[:], sibling[:]...))
					keys = append(keys, k/2)
				}
			}
		}
		pos++
	}
	return objects[GeneralizedIndex(1)], nil
}

// VerifyMerkleMultiproof verifies the Merkle multiproof by comparing the
// calculated root with the provided root.
func (gs GeneralizedIndicies) VerifyMerkleMultiproof(
	leaves []primitives.Bytes32,
	proof []primitives.Bytes32,
	indices []GeneralizedIndex,
	root primitives.Root,
) bool {
	calculatedRoot, err := gs.CalculateMultiMerkleRoot(leaves, proof)
	if err != nil {
		return false
	}
	return calculatedRoot == root
}
