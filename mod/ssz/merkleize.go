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
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
)

// Merkleize hashes a list of chunks and returns the HTR of the list of.
//
// merkleize(chunks, limit=None): Given ordered BYTES_PER_CHUNK-byte chunks, merkleize the chunks, and return the root:
// The merkleization depends on the effective input, which must be padded/limited:
//
//	if no limit:
//		pad the chunks with zeroed chunks to next_pow_of_two(len(chunks)) (virtually for memory efficiency).
//	if limit >= len(chunks):
//		pad the chunks with zeroed chunks to next_pow_of_two(limit) (virtually for memory efficiency).
//	if limit < len(chunks):
//		do not merkleize, input exceeds limit. Raise an error instead.
//	  Then, merkleize the chunks (empty input is padded to 1 zero chunk):
//	 If 1 chunk: the root is the chunk itself.
//	If > 1 chunks: merkleize as binary tree.
func Merkleize[ChunkT, RootT ~[32]byte](
	chunks []ChunkT,
	limit ...uint64,
) (RootT, error) {
	var (
		effectiveLimit  primitives.U64
		effectiveChunks []ChunkT
		lenChunks       = uint64(len(chunks))
	)

	if len(limit) == 0 {
		effectiveLimit = primitives.U64(lenChunks).NextPowerOfTwo()
	} else if limit[0] >= lenChunks {
		effectiveLimit = primitives.U64(limit[0]).NextPowerOfTwo()
		effectiveChunks = PadTo(chunks, effectiveLimit)
	} else {
		limit := limit[0]
		if limit < uint64(lenChunks) {
			return RootT{}, fmt.Errorf("input exceeds limit")
		}
		effectiveLimit = primitives.U64(limit)
	}

	if len(effectiveChunks) == 0 {
		effectiveChunks = PadTo(chunks, 1)
	}

	if len(effectiveChunks) == 1 {
		return RootT(effectiveChunks[0]), nil
	}

	return merkle.NewRootWithMaxLeaves[ChunkT, RootT](effectiveChunks, effectiveLimit.Unwrap())
}

// PadTo function to pad the chunks to the effective limit with zeroed chunks.
func PadTo[ChunkT ~[32]byte](chunks []ChunkT, effectiveLimit primitives.U64) []ChunkT {
	paddedChunks := make([]ChunkT, effectiveLimit)
	copy(paddedChunks, chunks)
	for i := uint64(len(chunks)); i < uint64(effectiveLimit); i++ {
		paddedChunks[i] = ChunkT{}
	}
	return paddedChunks
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
	return Merkleize[[32]byte, [32]byte](
		chunks,
		numChunks,
	)
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
	return merkle.NewRootWithMaxLeaves[[32]byte, [32]byte](roots, length)
}
