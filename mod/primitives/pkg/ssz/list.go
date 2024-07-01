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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
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
	// We recursively define "variable-size" types to be lists, unions,
	// Bitlists.
	// Therefore all Lists are NOT fixed.
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
	//#nosec:G701 // its fine.
	//nolint:mnd // 31 is okay.
	return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.BytesPerChunk
}

// SizeSSZ returns the size of the list in bytes.
func (l ListBasic[B]) SizeSSZ() int {
	// The same for ListBasic as for Vector.
	return Vector[B](l.elements).SizeSSZ()
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
	return Vector[B](l.elements).MarshalSSZTo(out)
}

// MarshalSSZ marshals the ListBasic into SSZ format.
func (l ListBasic[B]) MarshalSSZ() ([]byte, error) {
	// The same for ListBasic as for Vector.
	return Vector[B](l.elements).MarshalSSZ()
}

// NewFromSSZ creates a new ListBasic from SSZ format.
func (l ListBasic[B]) NewFromSSZ(buf []byte) (*ListBasic[B], error) {
	// The same for ListBasic as for Vector
	var (
		elements = make(Vector[B], 0)
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
type ListComposite[C Composite[C]] struct {
	elements []C
	limit    uint64
}

// ListCompositeFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func ListCompositeFromElements[C Composite[C]](
	limit uint64, elements ...C,
) *ListComposite[C] {
	return &ListComposite[C]{
		elements: elements,
		limit:    limit,
	}
}

// IsFixed returns true if the ListBasic is fixed size.
func (l ListComposite[C]) IsFixed() bool {
	// We recursively define "variable-size" types to be lists, unions,
	// Bitlists.
	// Therefore all Lists are NOT fixed.
	return false
}

// N returns the N value as defined in the SSZ specification.
func (l ListComposite[C]) N() uint64 {
	// list: ordered variable-length homogeneous collection, limited to N values
	// notation List[type, N], e.g. List[uint64, N]
	return l.limit
}

// ChunkCount returns the number of chunks in the VectorComposite.
func (l ListComposite[C]) ChunkCount() uint64 {
	// List[C, N] and Vector[C, N], where C is a composite type: N
	return (l.N())
}

// SizeSSZ returns the size of the list in bytes.
func (l ListComposite[C]) SizeSSZ() int {
	// The same for ListComposite as for VectorComposite.
	return VectorComposite[C](l.elements).SizeSSZ()
}

// HashTreeRootWith returns the Merkle root of the ListComposite
// with a given merkleizer.
func (l ListComposite[C]) HashTreeRootWith(
	merkleizer CompositeMerkleizer[common.ChainSpec, [32]byte, C],
) ([32]byte, error) {
	return merkleizer.MerkleizeListComposite(l.elements)
}

// HashTreeRoot returns the Merkle root of the ListComposite.
func (l ListComposite[C]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, C]())
}

// MarshalSSZTo marshals the ListComposite into SSZ format.
func (l ListComposite[C]) MarshalSSZTo(out []byte) ([]byte, error) {
	var c C
	if !c.IsFixed() {
		panic("not implemented yet")
	}

	// Safe to use Vector helper for a list here.
	return serializer.MarshalVectorFixed(out, l.elements)
}

// MarshalSSZ marshals the ListComposite into SSZ format.
func (l ListComposite[C]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new ListComposite from SSZ format.
func (ListComposite[C]) NewFromSSZ(
	buf []byte,
	limit uint64,
) (*ListComposite[C], error) {
	var c C
	if !c.IsFixed() {
		panic("not implemented yet")
	}

	// We can use Vector helper for a list here, it is safe.
	elements, err := serializer.UnmarshalVectorFixed[C](buf)
	if err != nil {
		return nil, err
	}

	return &ListComposite[C]{
		elements: elements,
		limit:    limit,
	}, nil
}
