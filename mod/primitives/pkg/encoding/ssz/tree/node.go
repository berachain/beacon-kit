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

// Node represents a node in the tree backing of an SSZ object.
type Node[RootT ~[32]byte] struct {
	// left is the left child node.
	left *Node[RootT]
	// right is the right child node.
	right *Node[RootT]
	// value holds the node's serialized data and/or hash.
	value RootT
	// GeneralizedIndex is the generalized index of the node.
	gIndex GeneralizedIndex[RootT]
}

// NewNodeAtDepth creates a new Node with the given value and depth.
func NewNodeAtDepth[RootT ~[32]byte](
	value RootT,
	d uint8,
	layerIndex uint64,
) *Node[RootT] {
	return &Node[RootT]{
		value:  value,
		gIndex: NewGeneralizedIndex[RootT](d, layerIndex),
	}
}

// NewNodeFromChildren creates a new Node from left and right child nodes.
// It calculates the value of the new node by combining the values of its
// children.
func NewNodeFromChildren[RootT ~[32]byte](
	left, right *Node[RootT], hasher func(RootT, RootT) RootT,
) *Node[RootT] {
	return &Node[RootT]{
		left:   left,
		right:  right,
		value:  hasher(left.value, right.value),
		gIndex: left.gIndex.Parent(),
	}
}

// Left returns the left child node.
func (n *Node[RootT]) Left() *Node[RootT] {
	return n.left
}

// Right returns the right child node.
func (n *Node[RootT]) Right() *Node[RootT] {
	return n.right
}

// Value returns the node's data.
func (n *Node[RootT]) Value() RootT {
	return n.value
}
