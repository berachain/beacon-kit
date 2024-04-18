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

	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/prysmaticlabs/gohashtree"
)

// SizeOfBasic returns the size of a basic type.
func SizeOfBasic[RootT ~[32]byte, B Basic[RootT]](b B) uint64 {
	// TODO: Boolean maybe this doesnt work.
	return uint64(reflect.TypeOf(b).Size())
}

// SizeOfComposite returns the size of a composite type.
func SizeOfComposite[RootT ~[32]byte, C Composite[RootT]](c C) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(c.SizeSSZ())
}

// ChunkCount returns the number of chunks required to store a value.
func ChunkCountBasic[RootT ~[32]byte, B Basic[RootT]](B) uint64 {
	return 1
}

// ChunkCountBitListVec returns the number of chunks required to store a bitlist
// or bitvector.
func ChunkCountBitListVec[T any](t []T) uint64 {
	//nolint:mnd // 256 is okay.
	return (uint64(len(t)) + 255) / 256
}

// ChunkCountBasicList returns the number of chunks required to store a list
// or vector of basic types.
func ChunkCountBasicList[B Basic[RootT], RootT ~[32]byte](
	b []B,
	maxCapacity uint64,
) uint64 {
	numItems := uint64(len(b))
	if numItems == 0 {
		return 1
	}
	size := SizeOfBasic[RootT, B](b[0])
	//nolint:mnd // 32 is okay.
	limit := (maxCapacity*size + 31) / 32
	if limit != 0 {
		return limit
	}

	return numItems
}

// ChunkCountCompositeList returns the number of chunks required to store a
// list or vector of composite types.
func ChunkCountCompositeList[C Composite[RootT], RootT ~[32]byte](
	c []C,
	limit uint64,
) uint64 {
	return max(uint64(len(c)), limit)
}

// ChunkCountContainer returns the number of chunks required to store a
// container.
func ChunkCountContainer[C Container[RootT], RootT ~[32]byte](c C) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(reflect.ValueOf(c).NumField())
}

// PadTo function to pad the chunks to the effective limit with zeroed chunks.
func PadTo[U64T U64[U64T], ChunkT ~[32]byte](
	chunks []ChunkT,
	effectiveLimit U64T,
) []ChunkT {
	paddedChunks := make([]ChunkT, effectiveLimit)
	copy(paddedChunks, chunks)
	//#nosec:G701 // This is a safe operation.
	for i := uint64(len(chunks)); i < uint64(effectiveLimit); i++ {
		paddedChunks[i] = ChunkT{}
	}
	return paddedChunks
}

// Pack packs a list of SSZ-marshallable elements into a single byte slice.
func Pack[
	U64T U64[U64T],
	U256L U256LT,
	B Basic[RootT],
	RootT ~[32]byte,
](b []B) ([]RootT, error) {
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
		case U64T:
			var buffer [8]byte
			//#nosec:G701 // This is a safe operation.
			binary.LittleEndian.PutUint64(buffer[:], uint64(el))
			packed = append(packed, buffer[:]...)
		case U256L:
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

// PartitionBytes partitions a byte slice into chunks of a given length.
func PartitionBytes[RootT ~[32]byte](input []byte) ([]RootT, uint64, error) {
	//nolint:mnd // we add 31 in order to round up the division.
	numChunks := max((uint64(len(input))+31)/constants.RootLength, 1)
	chunks := make([]RootT, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return chunks, numChunks, nil
}

// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSlice[U64T U64[U64T], RootT ~[32]byte](
	input []byte,
) (RootT, error) {
	chunks, numChunks, err := PartitionBytes[RootT](input)
	if err != nil {
		return RootT{}, err
	}
	return Merkleize[U64T, RootT, RootT](
		chunks,
		numChunks,
	)
}

// MixinLength takes a root element and mixes in the length of the elements
// that were hashed to produce it.
func MixinLength[RootT ~[32]byte](element RootT, length uint64) RootT {
	//nolint:mnd // 2 is okay.
	chunks := make([][32]byte, 2)
	chunks[0] = element
	binary.LittleEndian.PutUint64(chunks[1][:], length)
	if err := gohashtree.Hash(chunks, chunks); err != nil {
		return RootT{}
	}
	return chunks[0]
}
