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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/serializer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Vector conforms to the SSZEenumerable interface.
var _ types.SSZEnumerable[U64] = (*List[U64])(nil)

// List is a list of basic types.
type List[B types.SSZType[B]] struct {
	elements []B
	limit    uint64
}

// ListFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func ListFromElements[B types.SSZType[B]](
	limit uint64,
	elements ...B,
) *List[B] {
	return &List[B]{
		elements: elements,
		limit:    limit,
	}
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (l List[B]) SizeSSZ() int {
	// The same for List as for Vector.
	return Vector[B](l.elements).SizeSSZ()
}

// IsFixed returns true if the List is fixed size.
func (l List[B]) IsFixed() bool {
	// We recursively define "variable-size" types to be lists, unions,
	// Bitlists.
	// Therefore all Lists are NOT fixed.
	return false
}

// N returns the N value as defined in the SSZ specification.
func (l List[B]) N() uint64 {
	// list: ordered variable-length homogeneous collection, limited to N values
	// notation List[type, N], e.g. List[uint64, N]
	return l.limit
}

// ChunkCount returns the number of chunks in the List.
func (l List[B]) ChunkCount() uint64 {
	var b B
	switch b.Type() {
	case types.Basic:
		//#nosec:G701 // its fine.
		//nolint:mnd // 31 is okay.
		return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.BytesPerChunk
	default:
		return l.N()
	}
}

// Type returns the type of the List.
func (l List[B]) Type() types.Type {
	return types.Composite
}

// Elements returns the elements of the List.
func (l List[B]) Elements() []B {
	return l.elements
}

// HashTreeRootWith returns the Merkle root of the List
// with a given merkleizer.
func (l List[B]) HashTreeRootWith(
	merkleizer ListMerkleizer[[32]byte, B],
) ([32]byte, error) {
	var b B
	switch b.Type() {
	case types.Basic:
		return merkleizer.MerkleizeListBasic(l.elements, l.ChunkCount())
	case types.Composite:
		return merkleizer.MerkleizeListComposite(l.elements, l.ChunkCount())
	default:
		return [32]byte{}, errors.Wrapf(ErrUnknownType, "%v", b.Type())
	}
}

// HashTreeRoot returns the Merkle root of the List.
func (l List[B]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, B]())
}

// MarshalSSZTo marshals the List into SSZ format.
func (l List[B]) MarshalSSZTo(out []byte) ([]byte, error) {
	return Vector[B](l.elements).MarshalSSZTo(out)
}

// MarshalSSZ marshals the List into SSZ format.
func (l List[B]) MarshalSSZ() ([]byte, error) {
	// The same for List as for Vector.
	return Vector[B](l.elements).MarshalSSZ()
}

// NewFromSSZ creates a new List from SSZ format.
func (l List[B]) NewFromSSZ(buf []byte, limit uint64) (*List[B], error) {
	var b B
	if !b.IsFixed() {
		panic("not implemented yet")
	}

	// We can use Vector helper for a list here, it is safe.
	elements, err := serializer.UnmarshalVectorFixed[B](buf)
	if err != nil {
		return nil, err
	}

	return &List[B]{
		elements: elements,
		limit:    limit,
	}, nil
}
