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

var _ types.SSZEnumerable[Vector[U64], U64] = (Vector[U64])(nil)

type Vector[B interface{ types.SSZType[B] }] []B

// VectorBasicFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func VectorFromElements[B Basic[B]](elements ...B) Vector[B] {
	return elements
}

/* -------------------------------------------------------------------------- */
/*                                 BaseSSZType                                */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the list in bytes.
func (l Vector[B]) SizeSSZ() int {
	var b B
	return b.SizeSSZ() * len(l)
}

// isFixed returns true if the VectorBasic is fixed size.
func (Vector[B]) IsFixed() bool {
	var b B
	return b.IsFixed()
}

// Type returns the type of the VectorBasic.
func (l Vector[B]) Type() types.Type {
	var b B
	return b.Type()
}

// ChunkCount returns the number of chunks in the VectorBasic.
func (l Vector[B]) ChunkCount() uint64 {
	var b B
	switch b.Type() {
	case types.Basic:
		return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.BytesPerChunk
	default:
		return l.N()
	}
}

// N returns the N value as defined in the SSZ specification.
func (l Vector[B]) N() uint64 {
	// vector: ordered fixed-length homogeneous collection, with N values
	// notation Vector[type, N], e.g. Vector[uint64, N]
	return uint64(len(l))
}

// Elements returns the elements of the VectorBasic.
func (l Vector[B]) Elements() []B {
	return l
}

/* -------------------------------------------------------------------------- */
/*                                Merkleization                               */
/* -------------------------------------------------------------------------- */

// HashTreeRootWith returns the Merkle root of the VectorBasic
// with a given merkleizer.
func (l Vector[B]) HashTreeRootWith(
	merkleizer VectorMerkelizer[[32]byte, B],
) ([32]byte, error) {
	var b B
	if b.Type() == types.Basic {
		return merkleizer.MerkleizeVectorBasic(l)
	}
	return merkleizer.MerkleizeVectorComposite(l)
}

// HashTreeRoot returns the Merkle root of the VectorBasic.
func (l Vector[B]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, B]())
}

/* -------------------------------------------------------------------------- */
/*                                Serialization                               */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the VectorBasic into SSZ format.
func (l Vector[B]) MarshalSSZTo(out []byte) ([]byte, error) {
	if !l.IsFixed() {
		// return serializer.MarshalVectorVariable(out, l)
		return nil, errors.New("not implemented yet")
	}
	return serializer.MarshalVectorFixed(out, l)
}

// MarshalSSZ marshals the VectorBasic into SSZ format.
func (l Vector[B]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new VectorBasic from SSZ format.
func (v Vector[B]) NewFromSSZ(buf []byte) (Vector[B], error) {
	if !v.IsFixed() {
		panic("not implemented yet")
	}

	return serializer.UnmarshalVectorFixed[B](buf)
}
