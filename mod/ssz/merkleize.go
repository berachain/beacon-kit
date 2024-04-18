// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package ssz

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/merkle/zero"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
)

// Merkleize hashes a list of chunks and returns the HTR of the list of.
// As per the Ethereum 2.0 SSZ Specifcation:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/simple-serialize.md#merkleization
//
//nolint:lll
func Merkleize[ChunkT, RootT ~[32]byte](
	chunks []ChunkT,
	limit ...uint64,
) (RootT, error) {
	var effectiveLimit primitives.U64
	if len(limit) == 0 {
		effectiveLimit = primitives.U64(len(chunks))
	} else {
		limit := limit[0]
		if limit < uint64(len(chunks)) {
			return RootT{}, fmt.Errorf("input exceeds limit")
		}
		effectiveLimit = primitives.U64(limit)
	}

	if effectiveLimit == 0 {
		return zero.Hashes[0], nil
	} else if effectiveLimit == 1 {
		if len(chunks) == 0 {
			return zero.Hashes[0], nil
		}
		return RootT(chunks[0]), nil
	}

	paddedChunks := make([]ChunkT, effectiveLimit)
	copy(paddedChunks, chunks)
	for i := len(chunks); i < len(paddedChunks); i++ {
		paddedChunks[i] = zero.Hashes[0]
	}

	tree, err := merkle.NewTreeWithMaxLeaves[ChunkT, RootT](paddedChunks, uint64(effectiveLimit))
	if err != nil {
		return RootT{}, err
	}
	return tree.Root(), nil
}

// Pack packs a list of SSZ-marshallable elements into a single byte slice.
func Pack[S Marshallable](s []S) ([]byte, error) {
	// Pack each element into a single buffer.
	buf := make([]byte, 0)
	for _, el := range s {
		packed, err := el.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		buf = append(buf, packed...)
	}

	// Right pad the buffer to ensure it is a multiple of 32 bytes
	paddingSize := 32 - (len(buf) % 32)
	if paddingSize > 0 {
		padding := make([]byte, paddingSize)
		buf = append(buf, padding...)
	}
	return buf, nil
}

// MixinLength mixes the length into the root.
func MixinLength(root [32]byte, length uint64) [32]byte {
	return merkle.MixinLength(root, length)
}

// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSlice(input []byte) ([32]byte, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / constants.RootLength
	if numChunks == 0 {
		return [32]byte{}, ErrInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return Merkleize[[32]byte, [32]byte](chunks)
}

// MerkleizeList hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
func MerkleizeList[T Hashable[[32]byte]](
	elements []T, limit uint64,
) ([32]byte, error) {
	body, err := MerkleizeVector(elements, limit)
	if err != nil {
		return [32]byte{}, err
	}
	return merkle.MixinLength(body, uint64(len(elements))), nil
}

// MerkleizeVector hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func MerkleizeVector[T Hashable[[32]byte]](
	elements []T, length uint64,
) ([32]byte, error) {
	roots := make([][32]byte, len(elements))
	var err error
	for i, el := range elements {
		roots[i], err = el.HashTreeRoot()
		if err != nil {
			return [32]byte{}, err
		}
	}
	return Merkleize[[32]byte, [32]byte](roots, length)
}
