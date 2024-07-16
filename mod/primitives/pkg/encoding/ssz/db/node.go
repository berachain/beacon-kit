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

package db

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

// TODO: Figure out what the best way to do our DB types is?
// Should this even be a DB package, should we just have a tree that
// is completely decoupled fromt he DB?
// Putting this here for now just to make the deps nicely structured
// but long term we need to figure out the play.

// Node represents a node in the SSZ merkle tree.
type Node[GIndexT ~uint64, RootT ~[32]byte] struct {
	// SSZType is the SSZ type of the node.
	schema.SSZType
	// gIndex is the generalized index of the node in the Merkle tree.
	gIndex GIndexT
	// offset is the byte offset within the 32-byte chunk where the node's data
	// begins.
	offset uint8
}

// NewTreeNode locates a node in the SSZ merkle tree by its path and a root
// schema node to begin traversal from with gindex 1.
func NewTreeNode[GIndexT ~uint64, RootT ~[32]byte](
	root schema.SSZType, path merkle.ObjectPath[GIndexT, RootT],
) (Node[GIndexT, RootT], error) {
	found, gindex, offset, err := path.GetGeneralizedIndex(root)
	return Node[GIndexT, RootT]{
		SSZType: found,
		gIndex:  gindex,
		offset:  offset,
	}, err
}

// GeIndex returns the generalized index of the node in the Merkle tree.
func (n Node[GIndexT, RootT]) GIndex() GIndexT {
	return n.gIndex
}

// Offset returns the byte offset within the 32-byte chunk where the node's data
// begins.
func (n Node[_, _]) Offset() uint8 {
	return n.offset
}
