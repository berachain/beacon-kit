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

	"github.com/berachain/beacon-kit/primitives/math/pow"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// create builds a new merkle tree.
func create(
	hasher merkle.Hasher[[32]byte],
	leaves [][32]byte,
	depth uint64,
) TreeNode {
	length := uint64(len(leaves))
	if length == 0 {
		return &ZeroNode{
			depth:  depth,
			hasher: hasher,
		}
	}

	if depth == 0 {
		return &LeafNode{hash: leaves[0]}
	}

	split := min(pow.TwoToThePowerOf(depth-1), length)
	left := create(hasher, leaves[0:split], depth-1)
	right := create(hasher, leaves[split:], depth-1)
	return &InnerNode{
		left:   left,
		right:  right,
		hasher: hasher,
	}
}

// fromSnapshotParts creates a new Merkle tree from a list of finalized leaves,
// number of deposits and specified depth.
func fromSnapshotParts(
	hasher merkle.Hasher[[32]byte],
	finalized [][32]byte,
	deposits uint64,
	level uint64,
) (TreeNode, error) {
	var err error

	if len(finalized) < 1 || deposits == 0 {
		return &ZeroNode{
			depth:  level,
			hasher: hasher,
		}, nil
	}
	if deposits == pow.TwoToThePowerOf(level) {
		return &FinalizedNode{
			depositCount: deposits,
			hash:         finalized[0],
		}, nil
	}
	if level == 0 {
		return nil, ErrZeroLevel
	}
	node := InnerNode{
		hasher: hasher,
	}
	if leftSubtree := pow.TwoToThePowerOf(level - 1); deposits <= leftSubtree {
		node.left, err = fromSnapshotParts(hasher, finalized, deposits, level-1)
		if err != nil {
			return nil, err
		}
		node.right = &ZeroNode{
			depth:  level - 1,
			hasher: hasher,
		}
	} else {
		node.left = &FinalizedNode{
			depositCount: leftSubtree,
			hash:         finalized[0],
		}
		node.right, err = fromSnapshotParts(hasher, finalized[1:], deposits-leftSubtree, level-1)
		if err != nil {
			return nil, err
		}
	}
	return &node, nil
}

// generateProof returns a merkle proof and root.
func generateProof(
	tree TreeNode,
	index uint64,
	depth uint64,
) ([32]byte, [][32]byte) {
	var proof [][32]byte
	node := tree
	for depth > 0 {
		ithBit := (index >> (depth - 1)) & 0x1 //nolint:mnd // spec.
		if ithBit == 1 {
			proof = append(proof, node.Left().GetRoot())
			node = node.Right()
		} else {
			proof = append(proof, node.Right().GetRoot())
			node = node.Left()
		}
		depth--
	}

	slices.Reverse(proof)
	return node.GetRoot(), proof
}
