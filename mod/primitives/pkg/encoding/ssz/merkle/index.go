// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle

import (
	"sort"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math/pow"
)

type (
	// GeneralizedIndex is a generalized index.
	GeneralizedIndex[RootT ~[32]byte] uint64

	// GeneralizedIndices is a list of generalized indices.
	GeneralizedIndices[RootT ~[32]byte] []GeneralizedIndex[RootT]
)

// NewGeneralizedIndex calculates the generalized index from the depth and
// index.
func NewGeneralizedIndex[RootT ~[32]byte](
	depth uint8,
	index uint64,
) GeneralizedIndex[RootT] {
	return GeneralizedIndex[RootT]((1 << depth) + index)
}

// Unwrap returns the underlying uint64 value of the GeneralizedIndex.
func (g GeneralizedIndex[RootT]) Unwrap() uint64 {
	return uint64(g)
}

// Length returns the length of the generalized index.
func (g GeneralizedIndex[RootT]) Length() uint64 {
	return uint64(log.ILog2Floor(uint64(g)))
}

// IndexBit returns the bit at the specified position in a generalized index.
func (g GeneralizedIndex[RootT]) IndexBit(position int) bool {
	return (g & (1 << position)) > 0
}

// Sibling returns the sibling index of the current generalized index.
func (g GeneralizedIndex[RootT]) Sibling() GeneralizedIndex[RootT] {
	return g ^ 1
}

// LeftChild returns the left child index of the current generalized index.
//
//nolint:mnd // from spec.
func (g GeneralizedIndex[RootT]) LeftChild() GeneralizedIndex[RootT] {
	return 2 * g
}

// RightChild returns the right child index of the current generalized index.
func (g GeneralizedIndex[RootT]) RightChild() GeneralizedIndex[RootT] {
	return 2*g + 1
}

// Parent returns the parent index of the current generalized index.
//
//nolint:mnd // from spec.
func (g GeneralizedIndex[RootT]) Parent() GeneralizedIndex[RootT] {
	return g / 2
}

// GetBranchIndices returns the generalized indices of the nodes on the path
// from the root to the leaf.
func (g GeneralizedIndex[RootT]) GetBranchIndices() GeneralizedIndices[RootT] {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the
	// given tree index to the root.
	o := []GeneralizedIndex[RootT]{g.Sibling()}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent().Sibling())
	}
	return o[:len(o)-1]
}

// GetPathIndices returns the generalized indices of the nodes on the path from
// the leaf to the root.
func (g GeneralizedIndex[RootT]) GetPathIndices() GeneralizedIndices[RootT] {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the
	// given tree index to the root.
	o := []GeneralizedIndex[RootT]{g}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent())
	}
	return o[:len(o)-1]
}

// CalculateMerkleRoot calculates the Merkle root from the leaf and proof.
func (g GeneralizedIndex[RootT]) CalculateMerkleRoot(
	leaf RootT,
	proof []RootT,
) (RootT, error) {
	if uint64(len(proof)) != g.Length() {
		return RootT{},
			errors.Newf("expected proof length %d, received %d", g.Length(),
				len(proof))
	}
	for i, h := range proof {
		if g.IndexBit(i) {
			leaf = sha256.Hash(append(h[:], leaf[:]...))
		} else {
			leaf = sha256.Hash(append(leaf[:], h[:]...))
		}
	}
	return leaf, nil
}

// VerifyMerkleProof verifies the Merkle proof for the given
// leaf, proof, and root.
func (g GeneralizedIndex[RootT]) VerifyMerkleProof(
	leaf RootT,
	proof []RootT, // .Bytes32,
	root RootT,
) (bool, error) {
	calculated, err := g.CalculateMerkleRoot(leaf, proof)
	return calculated == root, err
}

// Concat multiple generalized indices into a single generalized index
// representing the path from the first to the last node.
func (gs GeneralizedIndices[RootT]) Concat() GeneralizedIndex[RootT] {
	o := GeneralizedIndex[RootT](1)
	for _, i := range gs {
		floorPower := pow.PrevPowerOfTwo(i)
		o = GeneralizedIndex[RootT](
			uint64(o)*uint64(floorPower) + (uint64(i) - uint64(floorPower)),
		)
	}
	return o
}

// GetHelperIndices returns the generalized indices of all "extra" chunks in the
// tree needed to prove the chunks with the given generalized indices. The
// decreasing order is chosen deliberately to ensure equivalence to the order of
// hashes in a regular single-item Merkle proof in the single-item case.
func (
	gs GeneralizedIndices[RootT],
) GetHelperIndices() GeneralizedIndices[RootT] {
	allHelperIndices := make(map[GeneralizedIndex[RootT]]struct{})
	allPathIndices := make(map[GeneralizedIndex[RootT]]struct{})

	for _, index := range gs {
		for _, helperIndex := range index.GetBranchIndices() {
			allHelperIndices[helperIndex] = struct{}{}
		}
		for _, pathIndex := range index.GetPathIndices() {
			allPathIndices[pathIndex] = struct{}{}
		}
	}

	difference := make([]GeneralizedIndex[RootT], 0, len(allHelperIndices))
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
func (gs GeneralizedIndices[RootT]) CalculateMultiMerkleRoot(
	leaves []RootT,
	proof []RootT,
) (RootT, error) {
	if len(leaves) != len(gs) {
		return RootT{}, errors.New(
			"mismatched leaves and indices length",
		)
	}

	helperIndices := gs.GetHelperIndices()
	if len(proof) != len(helperIndices) {
		return RootT{}, errors.New(
			"mismatched proof and helper indices length",
		)
	}

	objects := make(map[GeneralizedIndex[RootT]]RootT)
	for i, index := range gs {
		objects[index] = leaves[i]
	}
	for i, index := range helperIndices {
		objects[index] = proof[i]
	}

	keys := make([]GeneralizedIndex[RootT], 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	pos := 0
	var sibling RootT
	for pos < len(keys) {
		k := keys[pos]
		if _, ok := objects[k]; ok {
			if sibling, ok = objects[k^1]; ok {
				if _, ok = objects[k/2]; !ok {
					obj := objects[(k|1)^1]
					objects[k/2] = sha256.Hash(append(obj[:], sibling[:]...))
					//nolint:mnd // from spec.
					keys = append(keys, k/2)
				}
			}
		}
		pos++
	}
	return objects[GeneralizedIndex[RootT](1)], nil
}

// VerifyMerkleMultiproof verifies the Merkle multiproof by comparing the
// calculated root with the provided root.
func (gs GeneralizedIndices[RootT]) VerifyMerkleMultiproof(
	leaves []RootT,
	proof []RootT,
	root RootT,
) bool {
	calculatedRoot, err := gs.CalculateMultiMerkleRoot(leaves, proof)
	if err != nil {
		return false
	}
	return calculatedRoot == root
}
