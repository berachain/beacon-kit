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
	"github.com/berachain/beacon-kit/primitives/merkle/zero"
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

// FinalizedNode represents a finalized node and satisfies the TreeNode
// interface.
type FinalizedNode struct {
	depositCount uint64
	hash         [32]byte
}

// GetRoot returns the root of the Merkle tree.
func (f *FinalizedNode) GetRoot() [32]byte {
	return f.hash
}

// IsFull returns whether there is space left for deposits.
// A FinalizedNode will always return true as by definition it
// is full and deposits can't be added to it.
func (*FinalizedNode) IsFull() bool {
	return true
}

// Finalize marks deposits of the Merkle tree as finalized.
func (f *FinalizedNode) Finalize(_, _ uint64) (TreeNode, error) {
	return f, nil
}

// GetFinalized returns a list of hashes of all the finalized nodes and the
// number of deposits.
func (f *FinalizedNode) GetFinalized(result [][32]byte) (uint64, [][32]byte) {
	return f.depositCount, append(result, f.hash)
}

// PushLeaf adds a new leaf node at the next available zero node.
func (*FinalizedNode) PushLeaf([32]byte, uint64) (TreeNode, error) {
	return nil, ErrFinalizedNodeCannotPushLeaf
}

// Right returns nil as a finalized node can't have any children.
func (*FinalizedNode) Right() TreeNode {
	return nil
}

// Left returns nil as a finalized node can't have any children.
func (*FinalizedNode) Left() TreeNode {
	return nil
}

// LeafNode represents a leaf node holding a deposit and satisfies the TreeNode
// interface.
type LeafNode struct {
	hash [32]byte
}

// GetRoot returns the root of the Merkle tree.
func (l *LeafNode) GetRoot() [32]byte {
	return l.hash
}

// IsFull returns whether there is space left for deposits.
// A LeafNode will always return true as it is the last node
// in the tree and therefore can't have any deposits added to it.
func (*LeafNode) IsFull() bool {
	return true
}

// Finalize marks deposits of the Merkle tree as finalized.
func (l *LeafNode) Finalize(_, _ uint64) (TreeNode, error) {
	return &FinalizedNode{1, l.hash}, nil
}

// GetFinalized returns a list of hashes of all the finalized nodes and the
// number of deposits.
func (*LeafNode) GetFinalized(result [][32]byte) (uint64, [][32]byte) {
	return 0, result
}

// PushLeaf adds a new leaf node at the next available zero node.
func (*LeafNode) PushLeaf([32]byte, uint64) (TreeNode, error) {
	return nil, ErrLeafNodeCannotPushLeaf
}

// Right returns nil as a leaf node is the last node and can't have any
// children.
func (*LeafNode) Right() TreeNode {
	return nil
}

// Left returns nil as a leaf node is the last node and can't have any children.
func (*LeafNode) Left() TreeNode {
	return nil
}

// InnerNode represents an inner node with two children and satisfies the
// TreeNode interface.
type InnerNode struct {
	left, right TreeNode
	hasher      merkle.Hasher[[32]byte]
}

// GetRoot returns the root of the Merkle tree.
func (n *InnerNode) GetRoot() [32]byte {
	left := n.left.GetRoot()
	right := n.right.GetRoot()
	return n.hasher.Combi(left, right)
}

// IsFull returns whether there is space left for deposits.
func (n *InnerNode) IsFull() bool {
	return n.right.IsFull()
}

// Finalize marks deposits of the Merkle tree as finalized.
func (n *InnerNode) Finalize(
	depositsToFinalize uint64,
	depth uint64,
) (TreeNode, error) {
	var err error
	deposits := pow.TwoToThePowerOf(depth)
	if deposits <= depositsToFinalize {
		return &FinalizedNode{deposits, n.GetRoot()}, nil
	}
	if depth == 0 {
		return nil, ErrZeroDepth
	}
	n.left, err = n.left.Finalize(depositsToFinalize, depth-1)
	if err != nil {
		return nil, err
	}

	//nolint:mnd // spec.
	if depositsToFinalize > deposits/2 {
		remaining := depositsToFinalize - deposits/2
		n.right, err = n.right.Finalize(remaining, depth-1)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

// GetFinalized returns a list of hashes of all the finalized nodes and the
// number of deposits.
func (n *InnerNode) GetFinalized(result [][32]byte) (uint64, [][32]byte) {
	leftDeposits, result := n.left.GetFinalized(result)
	rightDeposits, result := n.right.GetFinalized(result)
	return leftDeposits + rightDeposits, result
}

// PushLeaf adds a new leaf node at the next available zero node.
//
//nolint:nestif // recursion.
func (n *InnerNode) PushLeaf(leaf [32]byte, depth uint64) (TreeNode, error) {
	if !n.left.IsFull() {
		left, err := n.left.PushLeaf(leaf, depth-1)
		if err == nil {
			n.left = left
		} else {
			return n, err
		}
	} else {
		right, err := n.right.PushLeaf(leaf, depth-1)
		if err == nil {
			n.right = right
		} else {
			return n, err
		}
	}
	return n, nil
}

// Right returns the child node on the right.
func (n *InnerNode) Right() TreeNode {
	return n.right
}

// Left returns the child node on the left.
func (n *InnerNode) Left() TreeNode {
	return n.left
}

// ZeroNode represents an empty node without a deposit and satisfies the
// TreeNode interface.
type ZeroNode struct {
	depth  uint64
	hasher merkle.Hasher[[32]byte]
}

// GetRoot returns the root of the Merkle tree.
func (z *ZeroNode) GetRoot() [32]byte {
	if z.depth == DepositContractDepth {
		return z.hasher.Combi(zero.Hashes[z.depth-1], zero.Hashes[z.depth-1])
	}
	return zero.Hashes[z.depth]
}

// IsFull returns wh   ether there is space left for deposits.
// A ZeroNode will always return false as a ZeroNode is an empty node
// that gets replaced by a deposit.
func (*ZeroNode) IsFull() bool {
	return false
}

// Finalize marks deposits of the Merkle tree as finalized.
func (*ZeroNode) Finalize(_, _ uint64) (TreeNode, error) {
	return nil, nil //nolint:nilnil // spec.
}

// GetFinalized returns a list of hashes of all the finalized nodes and the
// number of deposits.
func (*ZeroNode) GetFinalized(result [][32]byte) (uint64, [][32]byte) {
	return 0, result
}

// PushLeaf adds a new leaf node at the next available zero node.
func (z *ZeroNode) PushLeaf(leaf [32]byte, depth uint64) (TreeNode, error) {
	return create(z.hasher, [][32]byte{leaf}, depth), nil
}

// Right returns nil as a zero node can't have any children.
func (*ZeroNode) Right() TreeNode {
	return nil
}

// Left returns nil as a zero node can't have any children.
func (*ZeroNode) Left() TreeNode {
	return nil
}
