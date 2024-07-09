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

package merkleizer

import (
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes/buffer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

// merkleizer can be used for merkleizing SSZ types.
type Merkleizer[
	RootT ~[32]byte, T types.MerkleizableSSZObject[RootT],
] struct {
	rootHasher  *merkle.RootHasher[RootT]
	bytesBuffer Buffer[RootT]
}

// New creates a new merkleizer with a reusable hasher and bytes buffer.
func New[
	RootT ~[32]byte, T types.MerkleizableSSZObject[RootT],
]() *Merkleizer[RootT, T] {
	return &Merkleizer[RootT, T]{
		rootHasher: merkle.NewRootHasher(
			crypto.NewHasher[RootT](sha256.CustomHashFn()),
			merkle.BuildParentTreeRoots,
		),
		bytesBuffer: buffer.NewReusableBuffer[RootT](),
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Vector                                   */
/* -------------------------------------------------------------------------- */

// MerkleizeBasic hashes the packed value and returns the HTR.
func (m *Merkleizer[RootT, T]) MerkleizeBasic(
	value T,
) (RootT, error) {
	return m.MerkleizeVectorBasic([]T{value})
}

// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
//
// TODO: Deprecate in favor of Merkelize(List/Vector)Basic with type T = Byte.
func (m *Merkleizer[RootT, T]) MerkleizeByteSlice(
	input []byte,
) (RootT, error) {
	chunks, numChunks := chunkifyBytes[RootT](input)
	return m.Merkleize(chunks, numChunks)
}

// MerkleizeVectorBasic implements the SSZ merkleization algorithm
// for a vector of basic types.
func (m *Merkleizer[RootT, T]) MerkleizeVectorBasic(
	value []T,
) (RootT, error) {
	// merkleize(pack(value))
	// if value is a basic object or a vector of basic objects.
	packed, _, err := pack[RootT](value)
	if err != nil {
		return [32]byte{}, err
	}
	return m.Merkleize(packed)
}

// MerkleizeVectorComposite implements the SSZ merkleization algorithm for a
// vector of composite types or a container.
func (m *Merkleizer[RootT, T]) MerkleizeVectorCompositeOrContainer(
	value []T,
) (RootT, error) {
	var (
		err  error
		htrs = m.bytesBuffer.Get(len(value))
	)

	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}
	return m.Merkleize(htrs)
}

/* -------------------------------------------------------------------------- */
/*                                    List                                    */
/* -------------------------------------------------------------------------- */

// MerkleizeListBasic implements the SSZ merkleization algorithm for a list of
// basic types.
func (m *Merkleizer[RootT, T]) MerkleizeListBasic(
	value []T,
	chunkCount uint64,
) (RootT, error) {
	// mix_in_length(
	// 		merkleize(
	// 			pack(value),
	// 			limit=chunk_count(type),
	// 		),
	//      len(value),
	// )
	// if value is a list of basic objects.
	packed, _, err := pack[RootT](value)
	if err != nil {
		return [32]byte{}, err
	}

	root, err := m.Merkleize(
		packed, chunkCount,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return m.rootHasher.MixIn(root, uint64(len(value))), nil
}

// MerkleizeListComposite implements the SSZ merkleization algorithm for a list
// of composite types.
func (m *Merkleizer[RootT, T]) MerkleizeListComposite(
	value []T,
	chunkCount uint64,
) (RootT, error) {
	var (
		err  error
		htrs = m.bytesBuffer.Get(len(value))
	)

	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}

	root, err := m.Merkleize(
		htrs, chunkCount,
	)
	if err != nil {
		return RootT{}, err
	}

	return m.rootHasher.MixIn(root, uint64(len(value))), nil
}

/* -------------------------------------------------------------------------- */
/*                                  Merkleize                                 */
/* -------------------------------------------------------------------------- */

// Merkleize hashes a list of chunks and returns the HTR of the list of.
//
// From Spec:
//
// merkleize(chunks, limit=None): Given ordered BYTES_PER_CHUNK-byte chunks,
// merkleize the chunks, and return the root: The merkleization depends on the
// effective input, which must be padded/limited.
func (m *Merkleizer[RootT, T]) Merkleize(
	chunks []RootT,
	limit ...uint64,
) (RootT, error) {
	var (
		// effectiveLimit is used to track the "virtual padding of"
		effectiveLimit math.U64
		lenChunks      = uint64(len(chunks))
	)

	// The merkleization depends on the effective input, which must be
	// padded/limited
	switch {
	// From Spec:
	//
	// if no limit: pad the chunks with zeroed chunks to
	// next_pow_of_two(len(chunks)) (virtually for memory efficiency).
	case len(limit) == 0:
		effectiveLimit = math.U64(lenChunks).NextPowerOfTwo()

	// From Spec:
	//
	// limit >= len(chunks), pad the chunks with zeroed chunks to
	// next_pow_of_two(limit) (virtually for memory efficiency).
	case limit[0] >= lenChunks:
		effectiveLimit = math.U64(limit[0]).NextPowerOfTwo()

	// From Spec:
	//
	// if limit < len(chunks): do not merkleize,
	// input exceeds limit. Raise an error instead.
	default:
		if limit[0] < lenChunks {
			return RootT{}, errors.New("input exceeds limit")
		}
		effectiveLimit = math.U64(limit[0])
	}

	// From Spec:
	//
	// If 1 chunk: the root is the chunk itself.
	if lenChunks == 1 && effectiveLimit == 1 {
		return chunks[0], nil
	}

	// If > 1 chunks: merkleize as binary tree.
	return m.rootHasher.NewRootWithMaxLeaves(chunks, effectiveLimit)
}
