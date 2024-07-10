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
	"unsafe"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Vector conforms to the SSZEenumerable interface.
var _ schema.SSZEnumerable[Byte] = (*List[Byte])(nil)

// List is a list of basic types.
type List[T schema.MinimalSSZObject] struct {
	// elements is the list of elements.
	elements []T
	// limit is the maximum number of elements in the list.
	limit uint64
}

// ListFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func ListFromElements[T schema.MinimalSSZObject](
	limit uint64,
	elements ...T,
) *List[T] {
	return &List[T]{
		elements: elements,
		limit:    limit,
	}
}

// ByteList from Bytes creates a new List from bytes.
func ByteListFromBytes(bytes []byte, limit uint64) *List[Byte] {
	//#nosec:G103 // its fine, but we should find  abetter solution.
	elements := *(*[]Byte)(unsafe.Pointer(&bytes))
	return &List[Byte]{
		elements: elements,
		limit:    limit,
	}
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (l *List[T]) SizeSSZ() int {
	// The same for List as for Vector.
	return Vector[T](l.elements).SizeSSZ()
}

// IsFixed returns true if the List is fixed size.
func (l *List[T]) IsFixed() bool {
	// We recursively define "variable-size" types to be lists, unions,
	// Bitlists.
	// Therefore all Lists are NOT fixed.
	return false
}

// N returns the N value as defined in the SSZ specification.
func (l *List[T]) N() uint64 {
	// list: ordered variable-length homogeneous collection, limited to N values
	// notation List[type, N], e.g. List[uint64, N]
	return l.limit
}

// ChunkCount returns the number of chunks in the List.
func (l *List[T]) ChunkCount() uint64 {
	var b T
	switch t := b.Type().ID(); {
	case t.IsBasic():
		//#nosec:G701 // its fine.
		//nolint:mnd // 31 is okay.
		return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.BytesPerChunk
	default:
		return l.N()
	}
}

// Type returns the type of the List.
func (l *List[T]) Type() schema.SSZType {
	var t T
	// TODO: Fix this is a bad hack.
	if l == nil {
		return schema.DefineList(t.Type(), 0)
	}
	return schema.DefineList(t.Type(), l.limit)
}

// Elements returns the elements of the List.
func (l *List[T]) Elements() []T {
	return l.elements
}

// HashTreeRootWith returns the Merkle root of the List
// with a given merkle.
func (l *List[T]) HashTreeRootWith(
	merkleizer *merkle.Merkleizer[[32]byte, T],
) ([32]byte, error) {
	var b T
	switch t := b.Type().ID(); {
	case t.IsBasic():
		return merkleizer.MerkleizeListBasic(l.elements, l.ChunkCount())
	case t.IsComposite():
		return merkleizer.MerkleizeListComposite(l.elements, l.ChunkCount())
	default:
		return [32]byte{}, errors.Wrapf(ErrUnknownType, "%v", b.Type())
	}
}

// HashTreeRoot returns the Merkle root of the List.
func (l *List[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkle.NewMerkleizer[[32]byte, T]())
}

// MarshalSSZTo marshals the List into SSZ format.
func (l *List[T]) MarshalSSZTo(out []byte) ([]byte, error) {
	return Vector[T](l.elements).MarshalSSZTo(out)
}

// MarshalSSZ marshals the List into SSZ format.
func (l *List[T]) MarshalSSZ() ([]byte, error) {
	// The same for List as for Vector.
	return Vector[T](l.elements).MarshalSSZ()
}

// NewFromSSZ creates a new List from SSZ format.
func (l *List[T]) NewFromSSZ(buf []byte, limit uint64) (*List[T], error) {
	// We can use Vector helper for a list here, it is safe.
	elements, err := Vector[T](l.elements).NewFromSSZ(buf)
	if err != nil {
		return nil, err
	}

	return &List[T]{
		elements: elements,
		limit:    limit,
	}, nil
}
