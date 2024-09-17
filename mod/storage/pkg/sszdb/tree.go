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

package sszdb

import (
	"reflect"
	"unsafe"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	fastssz "github.com/ferranbt/fastssz"
)

type Node struct {
	Left  *Node
	Right *Node
	Value []byte
}

type Treeable interface {
	GetTree() (*fastssz.Node, error)
	DefineSchema(*schema.Codec)
}

func NewTreeFromFastSSZ(tr Treeable) (*Node, error) {
	root, err := tr.GetTree()
	if err != nil {
		return nil, err
	}
	return copyTree(root), nil
}

func (n *Node) CachedHash() []byte {
	if (n.Left == nil && n.Right == nil) || n.Value != nil {
		return n.Value
	}
	n.Value = hashFn(append(n.Left.CachedHash(), n.Right.CachedHash()...))
	return n.Value
}

// TODO this is a big hack to speed up development
// to be replaced with either a custom walker or simply ssz/v2
// It can also be used for regression testing against the fastssz
// implementation.
func copyTree(node *fastssz.Node) *Node {
	if node == nil {
		return nil
	}
	reflectNode := reflect.Indirect(reflect.ValueOf(node))

	f := reflectNode.FieldByIndex([]int{0})
	left := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*fastssz.Node)

	f = reflectNode.FieldByIndex([]int{1})
	right := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface().(*fastssz.Node)

	f = reflectNode.FieldByIndex([]int{3})
	value := f.Bytes()

	return &Node{
		Left:  copyTree(left),
		Right: copyTree(right),
		Value: value,
	}
}
