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
var _ schema.SSZEnumerable[Byte] = (Vector[Byte])(nil)

// Vector represents a vector of elements.
type Vector[T schema.MinimalSSZObject] []T

// VectorBasicFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func VectorFromElements[T schema.MinimalSSZObject](elements ...T) Vector[T] {
	return elements
}

// ByteVectorFromBytes creates a new Vector[Byte]ß from bytes.
func ByteVectorFromBytes(bytes []byte) Vector[Byte] {
	//#nosec:G103 // its fine, but we should find  abetter solution.
	v := *(*Vector[Byte])(unsafe.Pointer(&bytes))
	return v
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (v Vector[T]) SizeSSZ() int {
	var b T
	return b.SizeSSZ() * len(v)
}

// isFixed returns true if the VectorBasic is fixed size.
func (Vector[T]) IsFixed() bool {
	// If the element in the vector is fixed size, then
	// the vector is fixed size.
	var b T
	return b.IsFixed()
}

// Type returns the type of the VectorBasic.
func (v Vector[T]) Type() schema.SSZType {
	var t T
	return schema.DefineVector(t.Type(), uint64(len(v)))
}

// ChunkCount returns the number of chunks in the VectorBasic.
func (v Vector[T]) ChunkCount() uint64 {
	var b T
	switch t := b.Type().ID(); {
	case t.IsBasic():
		//#nosec:G701 // its fine.
		//nolint:mnd // 31 is okay.
		return (v.N()*uint64(b.SizeSSZ()) + 31) / constants.BytesPerChunk
	default:
		return v.N()
	}
}

// N returns the N value as defined in the SSZ specification.
func (v Vector[T]) N() uint64 {
	// vector: ordered fixed-length homogeneous collection, with N values
	// notation Vector[type, N], e.g. Vector[uint64, N]
	return uint64(len(v))
}

// Elements returns the elements of the VectorBasic.
func (v Vector[T]) Elements() []T {
	return v
}

/* -------------------------------------------------------------------------- */
/*                                Merkleization                               */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith returns the Merkle root of the VectorBasic
// with a given merkle.
func (v Vector[T]) HashTreeRootWith(
	merkleizer *merkle.Merkleizer[[32]byte, T],
) ([32]byte, error) {
	var b T
	switch t := b.Type().ID(); {
	case t.IsBasic():
		return merkleizer.MerkleizeVectorBasic(v)
	case t.IsComposite():
		return merkleizer.MerkleizeVectorCompositeOrContainer(v)
	default:
		return [32]byte{}, errors.Wrapf(ErrUnknownType, "%v", b.Type())
	}
}

// HashTreeRoot returns the Merkle root of the VectorBasic.
func (v Vector[T]) HashTreeRoot() ([32]byte, error) {
	return v.HashTreeRootWith(merkle.NewMerkleizer[[32]byte, T]())
}

/* -------------------------------------------------------------------------- */
/*                                Serialization                               */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the VectorBasic into SSZ format.
func (v Vector[T]) MarshalSSZTo(_ []byte) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}

// MarshalSSZ marshals the VectorBasic into SSZ format.
func (v Vector[T]) MarshalSSZ() ([]byte, error) {
	return v.MarshalSSZTo(make([]byte, 0, v.SizeSSZ()))
}

// NewFromSSZ creates a new VectorBasic from SSZ format.
func (v Vector[T]) NewFromSSZ(_ []byte) (Vector[T], error) {
	return nil, errors.New("not implemented yet")
}
