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
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math/pow"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/berachain/beacon-kit/primitives/merkle/zero"
)

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

// Equals returns true if this finalized node is equal to the other node.
func (f *FinalizedNode) Equals(other TreeNode) bool {
	if f == nil && other == nil {
		return true
	}
	if f == nil || other == nil {
		return false
	}
	fn, ok := other.(*FinalizedNode)
	if !ok {
		return false
	}
	return f.depositCount == fn.depositCount && f.hash == fn.hash
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

// Equals returns true if this leaf node is equal to the other node.
func (l *LeafNode) Equals(other TreeNode) bool {
	if l == nil && other == nil {
		return true
	}
	if l == nil || other == nil {
		return false
	}
	ln, ok := other.(*LeafNode)
	if !ok {
		return false
	}
	return l.hash == ln.hash
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

// Equals returns true if this inner node is equal to the other node.
func (n *InnerNode) Equals(other TreeNode) bool {
	if n == nil && other == nil {
		return true
	}
	if n == nil || other == nil {
		return false
	}
	in, ok := other.(*InnerNode)
	if !ok {
		return false
	}
	return n.left.Equals(in.left) && n.right.Equals(in.right)
}

// ZeroNode represents an empty node without a deposit and satisfies the
// TreeNode interface.
type ZeroNode struct {
	depth  uint64
	hasher merkle.Hasher[[32]byte]
}

// GetRoot returns the root of the Merkle tree.
func (z *ZeroNode) GetRoot() [32]byte {
	if z.depth == constants.DepositContractDepth {
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

// Equals returns true if this zero node is equal to the other node.
func (z *ZeroNode) Equals(other TreeNode) bool {
	if z == nil && other == nil {
		return true
	}
	if z == nil || other == nil {
		return false
	}
	zn, ok := other.(*ZeroNode)
	if !ok {
		return false
	}
	return z.depth == zn.depth
}
