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
	"slices"

	"github.com/berachain/beacon-kit/primitives/pkg/math/log"
	"github.com/berachain/beacon-kit/primitives/pkg/math/pow"
)

// Inspired by the Ethereum 2.0 spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#helpers-for-generalized-indices
//
//nolint:lll // link.

type (
	// GeneralizedIndex is a generalized index.
	GeneralizedIndex uint64

	// GeneralizedIndices is a list of generalized indices.
	GeneralizedIndices []GeneralizedIndex
)

// NewGeneralizedIndex calculates the generalized index from the depth and
// index. Inspired by:
// https://github.com/protolambda/remerkleable/blob/master/remerkleable/tree.py#L20
//
//nolint:lll // link.
func NewGeneralizedIndex(
	depth uint8,
	index uint64,
) GeneralizedIndex {
	return GeneralizedIndex((1 << depth) | index)
}

// Unwrap returns the underlying uint64 value of the GeneralizedIndex.
func (g GeneralizedIndex) Unwrap() uint64 {
	return uint64(g)
}

// Length returns the length of the generalized index.
func (g GeneralizedIndex) Length() int {
	//#nosec:G701 // uint8 cannot overflow int.
	return int(log.ILog2Floor(g))
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
func (g GeneralizedIndex) GetBranchIndices() GeneralizedIndices {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the given tree index to the root.
	o := GeneralizedIndices{g.Sibling()}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent().Sibling())
	}
	return o[:len(o)-1]
}

// GetPathIndices returns the generalized indices of the nodes on the path from
// the leaf to the root.
func (g GeneralizedIndex) GetPathIndices() GeneralizedIndices {
	// Get the generalized indices of the sister chunks along the path from the
	// chunk with the given tree index to the root.
	o := GeneralizedIndices{g}
	for o[len(o)-1] > 1 {
		o = append(o, o[len(o)-1].Parent())
	}
	return o[:len(o)-1]
}

// Concat multiple generalized indices into a single generalized index
// representing the path from the first to the last node.
func (gs GeneralizedIndices) Concat() GeneralizedIndex {
	o := GeneralizedIndex(1)
	for _, i := range gs {
		floorPower := pow.PrevPowerOfTwo(i)
		o *= floorPower
		o += i - floorPower
	}
	return o
}

// GetHelperIndices returns the generalized indices of all "extra" chunks in the
// tree needed to prove the chunks with the given generalized indices. The
// decreasing order is chosen deliberately to ensure equivalence to the order of
// hashes in a regular single-item Merkle proof in the single-item case.
func (gs GeneralizedIndices) GetHelperIndices() GeneralizedIndices {
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

	difference := make(GeneralizedIndices, 0, len(allHelperIndices))
	for helperIndex := range allHelperIndices {
		if _, exists := allPathIndices[helperIndex]; !exists {
			difference = append(difference, helperIndex)
		}
	}

	// Sort in decreasing order.
	slices.SortFunc(difference, GeneralizedIndexReverseComparator)

	return difference
}

// Comparator function used to sort generalized indices in reverse order.
func GeneralizedIndexReverseComparator(i, j GeneralizedIndex) int {
	switch {
	case i < j:
		return 1
	case i == j:
		return 0
	case i > j:
		return -1
	default:
		panic("unreachable")
	}
}
