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
	"encoding/binary"
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/prysmaticlabs/gohashtree"
)

// SizeOfBasic returns the size of a basic type.
func SizeOfBasic[RootT ~[32]byte, B Basic[SpecT, RootT], SpecT any](
	b B,
) uint64 {
	// TODO: Boolean maybe this doesnt work.
	return uint64(reflect.TypeOf(b).Size())
}

// SizeOfComposite returns the size of a composite type.
func SizeOfComposite[RootT ~[32]byte, C Composite[SpecT, RootT], SpecT any](
	c C,
) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(c.SizeSSZ())
}

// SizeOfContainer returns the size of a container type.
func SizeOfContainer[RootT ~[32]byte, C Container[SpecT, RootT], SpecT any](
	c C,
) int {
	size := 0
	rValue := reflect.ValueOf(c)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}
	for i := range rValue.NumField() {
		fieldValue := rValue.Field(i)
		if !fieldValue.CanInterface() {
			return -1
		}

		// TODO: handle different types.
		field, ok := fieldValue.Interface().(Basic[SpecT, RootT])
		if !ok {
			return -1
		}
		size += field.SizeSSZ()

		// TODO: handle the offset calculation.
	}

	// TODO: This doesn't yet handle anything to do with offset calculation.
	return size
}

// ChunkCount returns the number of chunks required to store a value.
func ChunkCountBasic[RootT ~[32]byte, B Basic[SpecT, RootT], SpecT any](
	B,
) uint64 {
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
func ChunkCountBasicList[SpecT any, RootT ~[32]byte, B Basic[SpecT, RootT]](
	b []B,
	maxCapacity uint64,
) uint64 {
	numItems := uint64(len(b))
	if numItems == 0 {
		return 1
	}
	size := SizeOfBasic[RootT, B, SpecT](b[0])
	//nolint:mnd // 32 is okay.
	limit := (maxCapacity*size + 31) / 32
	if limit != 0 {
		return limit
	}

	return numItems
}

// ChunkCountCompositeList returns the number of chunks required to store a
// list or vector of composite types.
func ChunkCountCompositeList[
	SpecT any, RootT ~[32]byte, C Composite[SpecT, RootT],
](
	c []C,
	limit uint64,
) uint64 {
	return max(uint64(len(c)), limit)
}

// ChunkCountContainer returns the number of chunks required to store a
// container.
func ChunkCountContainer[SpecT any, RootT ~[32]byte, C Container[SpecT, RootT]](
	c C,
) uint64 {
	//#nosec:G701 // This is a safe operation.
	return uint64(reflect.ValueOf(c).NumField())
}

// PadTo function to pad the chunks to the effective limit with zeroed chunks.
func PadTo[U64T ~uint64, ChunkT ~[32]byte](
	chunks []ChunkT,
	size U64T,
) []ChunkT {
	switch numChunks := U64T(len(chunks)); {
	case numChunks == size:
		return chunks
	case numChunks > size:
		return chunks[:size]
	default:
		return append(chunks, make([]ChunkT, size-numChunks)...)
	}
}

// Pack packs a list of SSZ-marshallable elements into a single byte slice.
func Pack[
	U64T U64[U64T],
	U256L U256LT,
	SpecT any,
	RootT ~[32]byte,
	B Basic[SpecT, RootT],
](b []B) ([]RootT, error) {
	// Pack each element into separate buffers.
	var packed []byte
	for _, el := range b {
		fieldValue := reflect.ValueOf(el)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if !fieldValue.CanInterface() {
			return nil, errors.Newf(
				"cannot interface with field %v",
				fieldValue,
			)
		}

		// TODO: Do we need a safety check for Basic only here?
		// TODO: use a real interface instead of hood inline.
		el, ok := reflect.ValueOf(el).
			Interface().(interface{ MarshalSSZ() ([]byte, error) })
		if !ok {
			return nil, errors.Newf("unsupported type %T", el)
		}

		// TODO: Do we need a safety check for Basic only here?
		buf, err := el.MarshalSSZ()
		if err != nil {
			return nil, err
		}
		packed = append(packed, buf...)
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
	return Merkleize[U64T, RootT](
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
