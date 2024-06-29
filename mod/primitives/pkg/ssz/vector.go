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
	var b B
	return (l.N()*uint64(b.SizeSSZ()) + 31) / constants.RootLength
}

// N returns the N value as defined in the SSZ specification
func (l VectorBasic[B]) N() uint64 {
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
type VectorComposite[T Composite[T]] []T

// VectorCompositeFromElements creates a new VectorComposite from elements.
// TODO: Deprecate once off of Fastssz
func VectorCompositeFromElements[T Composite[T]](
	elements ...T,
) VectorComposite[T] {
	return elements
}

// isFixed returns true if the VectorBasic is fixed size.
func (VectorComposite[T]) IsFixed() bool {
	var t T
	return t.IsFixed()
}

// N returns the N value as defined in the SSZ specification
func (l VectorComposite[T]) N() uint64 {
	return uint64(len(l))
}

// SizeSSZ returns the size of the list in bytes.
func (l VectorComposite[T]) SizeSSZ() int {
	var t T
	return t.SizeSSZ() * len(l)
}

// HashTreeRootWith returns the Merkle root of the VectorComposite
// with a given merkleizer.
func (l VectorComposite[T]) HashTreeRootWith(
	merkleizer CompositeMerkleizer[common.ChainSpec, [32]byte, T],
) ([32]byte, error) {
	return merkleizer.MerkleizeVectorComposite(l)
}

// HashTreeRoot returns the Merkle root of the VectorComposite.
func (l VectorComposite[T]) HashTreeRoot() ([32]byte, error) {
	// Create a merkleizer
	return l.HashTreeRootWith(merkleizer.New[[32]byte, T]())
}

// MarshalSSZToBytes marshals the VectorComposite into SSZ format.
func (l VectorComposite[T]) MarshalSSZTo(out []byte) ([]byte, error) {
	var t T
	if !t.IsFixed() {
		panic("not implemented yet")
	}

	return serializer.MarshalVectorFixed(out, l)
}

// MarshalSSZ marshals the VectorComposite into SSZ format.
func (l VectorComposite[T]) MarshalSSZ() ([]byte, error) {
	return l.MarshalSSZTo(make([]byte, 0, l.SizeSSZ()))
}

// NewFromSSZ creates a new VectorComposite from SSZ format.
func (VectorComposite[T]) NewFromSSZ(
	buf []byte,
) (VectorComposite[T], error) {
	var t T
	if !t.IsFixed() {
		panic("not implemented yet")
	}

	return serializer.UnmarshalVectorFixed[T](buf)
}
