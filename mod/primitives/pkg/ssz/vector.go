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

// VectorBasic is a vector of basic types.
type VectorBasic[B Basic[B]] []B

// VectorBasicFromElements creates a new ListComposite from elements.
// TODO: Deprecate once off of Fastssz
func VectorBasicFromElements[B Basic[B]](elements ...B) VectorBasic[B] {
	return elements
}

// SizeSSZ returns the size of the list in bytes.
func (l VectorBasic[B]) SizeSSZ() int {
	var b B
	return b.SizeSSZ() * len(l)
}

// isFixed returns true if the VectorBasic is fixed size.
func (VectorBasic[B]) IsFixed() bool {
	return true
}

// ChunkCount returns the number of chunks in the VectorBasic.
func (l VectorBasic[B]) ChunkCount() uint64 {
	// List[B, N] and Vector[B, N], where B is a basic type:
	// (N * size_of(B) + 31) // 32 (dividing by chunk size, rounding up)
	var b B
	//#nosec:G701 // its fine.
	//nolint:mnd // 31 is okay.
	return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.RootLength
}

// N returns the N value as defined in the SSZ specification.
func (l VectorBasic[B]) N() uint64 {
	// vector: ordered fixed-length homogeneous collection, with N values
	// notation Vector[type, N], e.g. Vector[uint64, N]
	return uint64(len(l))
}

// HashTreeRootWith returns the Merkle root of the VectorBasic
// with a given merkleizer.
func (l VectorBasic[B]) HashTreeRootWith(
	merkleizer BasicMerkleizer[[32]byte, B],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorBasic(l)
}

// HashTreeRoot returns the Merkle root of the VectorBasic.
func (l VectorBasic[B]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, B]())
}

// MarshalSSZToBytes marshals the VectorBasic into SSZ format.
func (l VectorBasic[B]) MarshalSSZTo(out []byte) ([]byte, error) {
	return serializer.MarshalVectorFixed(out, l)
}

// MarshalSSZ marshals the VectorBasic into SSZ format.
func (l VectorBasic[B]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new VectorBasic from SSZ format.
func (VectorBasic[B]) NewFromSSZ(buf []byte) (VectorBasic[B], error) {
	return serializer.UnmarshalVectorFixed[B](buf)
}

/* -------------------------------------------------------------------------- */
/*                                  Composite                                 */
/* -------------------------------------------------------------------------- */

// VectorComposite is a vector of Composite types.
type VectorComposite[C Composite[C]] []C

// VectorCompositeFromElements creates a new VectorComposite from elements.
// TODO: Deprecate once off of Fastssz
func VectorCompositeFromElements[C Composite[C]](
	elements ...C,
) VectorComposite[C] {
	return elements
}

// isFixed returns true if the VectorBasic is fixed size.
func (VectorComposite[C]) IsFixed() bool {
	var c C
	return c.IsFixed()
}

// N returns the N value as defined in the SSZ specification.
func (l VectorComposite[C]) N() uint64 {
	// vector: ordered fixed-length homogeneous collection, with N values
	// notation Vector[type, N], e.g. Vector[uint64, N]
	return uint64(len(l))
}

// SizeSSZ returns the size of the list in bytes.
func (l VectorComposite[C]) SizeSSZ() int {
	var c C
	return c.SizeSSZ() * len(l)
}

// ChunkCount returns the number of chunks in the VectorComposite.
func (l VectorComposite[C]) ChunkCount() uint64 {
	// List[C, N] and Vector[C, N], where C is a composite type: N
	return (l.N())
}

// HashTreeRootWith returns the Merkle root of the VectorComposite
// with a given merkleizer.
func (l VectorComposite[C]) HashTreeRootWith(
	merkleizer CompositeMerkleizer[common.ChainSpec, [32]byte, C],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorComposite(l)
}

// HashTreeRoot returns the Merkle root of the VectorComposite.
func (l VectorComposite[C]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, C]())
}

// MarshalSSZToBytes marshals the VectorComposite into SSZ format.
func (l VectorComposite[C]) MarshalSSZTo(out []byte) ([]byte, error) {
	var c C
	if !c.IsFixed() {
		panic("not implemented yet")
	}

	return serializer.MarshalVectorFixed(out, l)
}

// MarshalSSZ marshals the VectorComposite into SSZ format.
func (l VectorComposite[C]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new VectorComposite from SSZ format.
func (VectorComposite[C]) NewFromSSZ(
	buf []byte,
) (VectorComposite[C], error) {
	var c C
	if !c.IsFixed() {
		panic("not implemented yet")
	}

	return serializer.UnmarshalVectorFixed[C](buf)
}
