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

package ssz

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/serializer"
)

/* -------------------------------------------------------------------------- */
/*                                    Basic                                   */
/* -------------------------------------------------------------------------- */

// ListBasic is a list of basic types.
type ListBasic[B Basic[B]] struct {
	elements []B
	limit    uint64
}

// ListBasicFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func ListBasicFromElements[B Basic[B]](
	limit uint64,
	elements ...B,
) *ListBasic[B] {
	return &ListBasic[B]{
		elements: elements,
		limit:    limit,
	}
}

// IsFixed returns true if the ListBasic is fixed size.
func (l ListBasic[B]) IsFixed() bool {
	return false
}

// N returns the N value as defined in the SSZ specification.
func (l ListBasic[B]) N() uint64 {
	// list: ordered variable-length homogeneous collection, limited to N values
	// notation List[type, N], e.g. List[uint64, N]
	return l.limit
}

// ChunkCount returns the number of chunks in the ListBasic.
func (l ListBasic[B]) ChunkCount() uint64 {
	// List[B, N] and Vector[B, N], where B is a basic type:
	// (N * size_of(B) + 31) // 32 (dividing by chunk size, rounding up)
	var b B
	//nolint:mnd // 31 is okay.
	return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.RootLength
}

// SizeSSZ returns the size of the list in bytes.
func (l ListBasic[B]) SizeSSZ() int {
	// The same for ListBasic as for VectorBasic.
	return VectorBasic[B](l.elements).SizeSSZ()
}

// HashTreeRootWith returns the Merkle root of the ListBasic
// with a given merkleizer.
func (l ListBasic[B]) HashTreeRootWith(
	merkleizer BasicMerkleizer[[32]byte, B],
) ([32]byte, error) {
	return merkleizer.MerkleizeListBasic(l.elements, l.limit)
}

// HashTreeRoot returns the Merkle root of the ListBasic.
func (l ListBasic[B]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, B]())
}

// MarshalSSZTo marshals the ListBasic into SSZ format.
func (l ListBasic[B]) MarshalSSZTo(out []byte) ([]byte, error) {
	return VectorBasic[B](l.elements).MarshalSSZTo(out)
}

// MarshalSSZ marshals the ListBasic into SSZ format.
func (l ListBasic[B]) MarshalSSZ() ([]byte, error) {
	// The same for ListBasic as for VectorBasic.
	return VectorBasic[B](l.elements).MarshalSSZ()
}

// NewFromSSZ creates a new ListBasic from SSZ format.
func (l ListBasic[B]) NewFromSSZ(buf []byte) (*ListBasic[B], error) {
	// The same for ListBasic as for VectorBasic
	var (
		elements = make(VectorBasic[B], 0)
		err      error
	)

	elements, err = elements.NewFromSSZ(buf)
	return &ListBasic[B]{
		elements: elements,
	}, err
}

/* -------------------------------------------------------------------------- */
/*                                  Composite                                 */
/* -------------------------------------------------------------------------- */

// ListComposite is a list of Composite types.
type ListComposite[T Composite[T]] struct {
	t     []T
	limit uint64
}

// ListCompositeFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func ListCompositeFromElements[T Composite[T]](
	limit uint64, elements ...T,
) *ListComposite[T] {
	return &ListComposite[T]{
		t:     elements,
		limit: limit,
	}
}

// IsFixed returns true if the ListBasic is fixed size.
func (l ListComposite[T]) IsFixed() bool {
	return false
}

// N returns the N value as defined in the SSZ specification.
func (l ListComposite[T]) N() uint64 {
	return l.limit
}

// ChunkCount returns the number of chunks in the VectorComposite.
func (l ListComposite[T]) ChunkCount() uint64 {
	// List[C, N] and Vector[C, N], where C is a composite type: N
	return (l.N())
}

// SizeSSZ returns the size of the list in bytes.
func (l ListComposite[T]) SizeSSZ() int {
	// The same for ListComposite as for VectorComposite.
	return VectorComposite[T](l.t).SizeSSZ()
}

// HashTreeRootWith returns the Merkle root of the ListComposite
// with a given merkleizer.
func (l ListComposite[T]) HashTreeRootWith(
	merkleizer CompositeMerkleizer[common.ChainSpec, [32]byte, T],
) ([32]byte, error) {
	return merkleizer.MerkleizeListComposite(l.t)
}

// HashTreeRoot returns the Merkle root of the ListComposite.
func (l ListComposite[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, T]())
}

// MarshalSSZTo marshals the ListComposite into SSZ format.
func (l ListComposite[T]) MarshalSSZTo(out []byte) ([]byte, error) {
	var t T
	if !t.IsFixed() {
		panic("not implemented yet")
	}

	// Safe to use Vector helper for a list here.
	return serializer.MarshalVectorFixed(out, l.t)
}

// MarshalSSZ marshals the ListComposite into SSZ format.
func (l ListComposite[T]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new ListComposite from SSZ format.
func (ListComposite[T]) NewFromSSZ(
	buf []byte,
	limit uint64,
) (*ListComposite[T], error) {
	var t T
	if !t.IsFixed() {
		panic("not implemented yet")
	}

	// We can use Vector helper for a list here, it is safe.
	elems, err := serializer.UnmarshalVectorFixed[T](buf)
	if err != nil {
		return nil, err
	}

	return &ListComposite[T]{
		t:     elems,
		limit: limit,
	}, nil
}
