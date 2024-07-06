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

package tree

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

// New returns a Merkle tree of the given leaves.
// As defined in the Ethereum 2.0 Spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#generalized-merkle-tree-index
//
//nolint:lll // link.
func New[LeafT ~[32]byte](
	leaves []LeafT,
	hashFn func([]byte) LeafT,
) []LeafT {
	/*
	   Return an array representing the tree nodes by generalized index:
	   [0, 1, 2, 3, 4, 5, 6, 7], where each layer is a power of 2. The 0 index is ignored. The 1 index is the root.
	   The result will be twice the size as the padded bottom layer for the input leaves.
	*/
	bottomLength := math.U64(len(leaves)).NextPowerOfTwo()
	//nolint:mnd // 2 is okay.
	o := make([]LeafT, bottomLength*2)
	copy(o[bottomLength:], leaves)
	for i := bottomLength - 1; i > 0; i-- {
		o[i] = hashFn(append(o[i*2][:], o[i*2+1][:]...))
	}
	return o
}

// Tree represents a Merkle tree structure.
type Tree[RootT ~[32]byte] struct {
	root   *Node[RootT]
	hasher *merkle.RootHasher[RootT]
}

// NewTreeFromLeaves constructs a Merkle tree, with the minimum
// depth required to support the number of leaves.
func NewTreeFromLeaves[RootT ~[32]byte](
	leaves []RootT,
	maxLeaves uint64,
) (*Tree[RootT], error) {
	return NewTreeFromLeavesWithDepth(
		leaves,
		math.U64(len(leaves)).NextPowerOfTwo().ILog2Ceil(),
		math.U64(maxLeaves).NextPowerOfTwo().ILog2Ceil(),
	)
}

// NewTreeFromLeavesWithDepth constructs a Merkle tree
// from a sequence of byte slices.
// It will fill the tree with zero hashes to create the required depth.
func NewTreeFromLeavesWithDepth[RootT ~[32]byte](
	chunks []RootT,
	depth uint8,
	limitDepth uint8,
) (*Tree[RootT], error) {
	rh := merkle.NewRootHasher(
		crypto.NewHasher[RootT](sha256.Hash),
		merkle.BuildParentTreeRoots,
	)

	// Handle the case where the tree is not full
	if len(chunks) == 0 {
		return &Tree[RootT]{
			root:   &Node[RootT]{},
			hasher: rh,
		}, nil
	}

	if err := VerifySufficientDepth(len(chunks), limitDepth); err != nil {
		return &Tree[RootT]{}, err
	}

	// Create the root node
	currentLayer := make([]*Node[RootT], len(chunks))
	// Create leaf nodes
	for i, leaf := range chunks {
		currentLayer[i] = &Node[RootT]{value: leaf}
	}

	// Build the tree bottom-up
	for d := range depth {
		//nolint:mnd // its okay.
		nextLayer := make([]*Node[RootT], (len(currentLayer)+1)/2)
		for i := 0; i < len(currentLayer); i += 2 {
			left := currentLayer[i]
			var right *Node[RootT]
			if i+1 < len(currentLayer) {
				right = currentLayer[i+1]
			} else {
				right = NewZeroNodeAtDepth[RootT](d)
			}
			parent := NewNodeFromChildren(left, right, rh.Combi)
			nextLayer[i/2] = parent
		}
		currentLayer = nextLayer
	}

	// If we need to extend the tree to be deeper, we do it virtually.
	currentNode := currentLayer[0]
	for j := depth; j < limitDepth; j++ {
		currentNode = NewNodeFromChildren(
			currentNode, NewZeroNodeAtDepth[RootT](j), rh.Combi,
		)
	}
	h := currentNode.value

	root := currentLayer[0]
	root.value = h

	return &Tree[RootT]{
		root:   root,
		hasher: rh,
	}, nil
}

// Root returns the root hash of the Merkle tree.
func (t *Tree[RootT]) Root() RootT {
	if t.root == nil {
		return RootT{}
	}
	return t.root.value
}

// VerifySufficientDepth ensures that the depth is sufficient to build a tree.
func VerifySufficientDepth(numLeaves int, depth uint8) error {
	switch {
	case depth > merkle.MaxTreeDepth:
		return merkle.ErrExceededDepth
	case numLeaves > (1 << depth):
		return errors.Wrap(
			merkle.ErrInsufficientDepthForLeaves,
			fmt.Sprintf(
				"attempted to build tree/root with %d leaves at depth %d",
				numLeaves, depth),
		)
	}
	return nil
}
