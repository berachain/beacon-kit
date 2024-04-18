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
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/prysmaticlabs/gohashtree"
)

type Basic interface {
	// TODO: add 128.
	~uint8 | ~uint16 | ~uint32 | ~uint64 | primitives.U256L | bool
}

type Composite[RootT ~[32]byte] interface {
	SizeSSZ() int
	HashTreeRoot() (RootT, error)
}

type Container interface {
	Marshallable
}

// SizeOfBasic returns the size of a basic type.
func SizeOfBasic[B Basic](b B) uint64 {
	// TODO: Boolean maybe this doesnt work.
	return uint64(reflect.TypeOf(b).Size())
}

// SizeOfComposite returns the size of a composite type.
func SizeOfComposite[RootT ~[32]byte, C Composite[RootT]](c C) uint64 {
	return uint64(c.SizeSSZ())
}

// ChunkCount returns the number of chunks required to store a value.
func ChunkCountBasic[B Basic](B) uint64 {
	return 1
}

// ChunkCountBitListVec returns the number of chunks required to store a bitlist
// or bitvector.
func ChunkCountBitListVec[T any](t []T) uint64 {
	return (uint64(len(t)) + 255) / 256
}

// ChunkCountBasicListVec returns the number of chunks required to store a list
// or vector of basic types.
func ChunkCountBasicListVec[B Basic](b []B) uint64 {
	if len(b) == 0 {
		return 0
	}
	return (uint64(len(b))*SizeOfBasic[B](b[0]) + 31) / 32
}

// ChunkCountCompositeList returns the number of chunks required to store a
// list or vector of composite types.
func ChunkCountCompositeList[RootT ~[32]byte, C Composite[RootT]](c []C) uint64 {
	return uint64(len(c))
}

// ChunkCountContainer returns the number of chunks required to store a
// container.
func ChunkCountContainer[C Container](c C) uint64 {
	return uint64(reflect.ValueOf(c).NumField())
}

// PadTo function to pad the chunks to the effective limit with zeroed chunks.
func PadTo[ChunkT ~[32]byte](
	chunks []ChunkT,
	effectiveLimit primitives.U64,
) []ChunkT {
	paddedChunks := make([]ChunkT, effectiveLimit)
	copy(paddedChunks, chunks)
	for i := uint64(len(chunks)); i < uint64(effectiveLimit); i++ {
		paddedChunks[i] = ChunkT{}
	}
	return paddedChunks
}

// Pack packs a list of SSZ-marshallable elements into a single byte slice.
func Pack[B Basic, RootT ~[32]byte](b []B) ([]RootT, error) {
	// Pack each element into separate buffers.
	var packed []byte
	for _, el := range b {
		switch el := reflect.ValueOf(el).Interface().(type) {
		case uint8:
			var buffer [1]byte
			buffer[0] = el
			packed = append(packed, buffer[:]...)
		case uint16:
			var buffer [2]byte
			binary.LittleEndian.PutUint16(buffer[:], el)
			packed = append(packed, buffer[:]...)
		case uint32:
			var buffer [4]byte
			binary.LittleEndian.PutUint32(buffer[:], el)
			packed = append(packed, buffer[:]...)
		case uint64:
			var buffer [8]byte
			binary.LittleEndian.PutUint64(buffer[:], el)
			packed = append(packed, buffer[:]...)
		case primitives.U256L:
			var buffer [32]byte
			copy(buffer[:], el[:])
			packed = append(packed, buffer[:]...)
		case bool:
			var buffer [1]byte
			if el {
				buffer[0] = 1
			}
			packed = append(packed, buffer[:]...)
		default:
			return nil, fmt.Errorf("unsupported type %T", el)
		}
	}

	root, _, err := PartitionBytes[RootT](packed)
	return root, err
}

func PartitionBytes[RootT ~[32]byte](input []byte) ([]RootT, uint64, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / constants.RootLength
	if numChunks == 0 {
		return nil, 0, ErrInvalidNilSlice
	}
	chunks := make([]RootT, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return chunks, numChunks, nil
}

// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSlice[RootT ~[32]byte](input []byte) (RootT, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	chunks, numChunks, err := PartitionBytes[RootT](input)
	if err != nil {
		return RootT{}, err
	}
	return Merkleize[RootT, RootT](
		chunks,
		numChunks,
	)
}

// Merkleize hashes a list of chunks and returns the HTR of the list of.
//
// merkleize(chunks, limit=None): Given ordered BYTES_PER_CHUNK-byte chunks,
// merkleize the chunks, and return the root: The merkleization depends on the
// effective input, which must be padded/limited:
//
//	if no limit:
//		pad the chunks with zeroed chunks to next_pow_of_two(len(chunks))
//
// (virtually for memory efficiency).
//
//	if limit >= len(chunks):
//		pad the chunks with zeroed chunks to next_pow_of_two(limit) (virtually for
//
// memory efficiency).
//
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

	return merkle.NewRootWithMaxLeaves[ChunkT, RootT](
		effectiveChunks,
		effectiveLimit.Unwrap(),
	)
}

// MixinLength takes a root element and mixes in the length of the elements
// that were hashed to produce it.
func MixinLength[RootT ~[32]byte](element RootT, length uint64) RootT {
	//nolint:gomnd // 2 is okay.
	chunks := make([][32]byte, 2)
	chunks[0] = element
	binary.LittleEndian.PutUint64(chunks[1][:], length)
	if err := gohashtree.Hash(chunks, chunks); err != nil {
		return RootT{}
	}
	return chunks[0]
}
