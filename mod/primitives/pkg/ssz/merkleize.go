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
	"reflect"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

// Merkleize hashes the packed value and returns the HTR.
func MerkleizeBasic[
	SpecT any, U64T U64[U64T], U256L U256LT,
	RootT ~[32]byte, B Basic[SpecT, RootT],
](
	value B,
) (RootT, error) {
	return MerkleizeVecBasic[U64T, U256L, RootT, SpecT]([]B{value})
}

// MerkleizeVec implements the SSZ merkleization algorithm for a vector of basic
// types.
func MerkleizeVecBasic[
	U64T U64[U64T], U256L U256LT, RootT ~[32]byte,
	SpecT any, B Basic[SpecT, RootT],
](
	value []B,
) (RootT, error) {
	packed, err := Pack[U64T, U256L, SpecT](value)
	if err != nil {
		return [32]byte{}, err
	}
	return Merkleize[U64T, RootT](packed)
}

// MerkleizeListBasic implements the SSZ merkleization algorithm for a list of
// basic types.
func MerkleizeListBasic[
	SpecT any, U64T U64[U64T], U256L U256LT, RootT ~[32]byte,
	B Basic[SpecT, RootT],
](
	value []B,
	limit uint64,
) (RootT, error) {
	packed, err := Pack[U64T, U256L, SpecT](value)
	if err != nil {
		return [32]byte{}, err
	}
	root, err := Merkleize[U64T, RootT](
		packed,
		ChunkCountBasicList[SpecT](value, limit),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return merkle.MixinLength(root, uint64(len(value))), nil
}

// TODO: MerkleizeBitlist

// MerkleizeContainer implements the SSZ merkleization algorithm for a
// container.
func MerkleizeContainer[
	SpecT any, U64T U64[U64T], RootT ~[32]byte,
	C Container[SpecT, RootT],
](
	value C, _ ...SpecT,
) (RootT, error) {
	rValue := reflect.ValueOf(value)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}
	numFields := rValue.NumField()
	htrs := make([]RootT, numFields)
	var err error
	for i := range numFields {
		fieldValue := rValue.Field(i)
		if !fieldValue.CanInterface() {
			return RootT{}, fmt.Errorf(
				"cannot interface with field %v",
				fieldValue,
			)
		}

		// TODO: handle different types.
		field, ok := fieldValue.Interface().(Basic[SpecT, RootT])
		if !ok {
			return RootT{}, fmt.Errorf(
				"field %d does not implement Hashable",
				i,
			)
		}
		htrs[i], err = field.HashTreeRoot( /*args...*/ )
		if err != nil {
			return RootT{}, err
		}
	}
	return Merkleize[U64T, RootT](htrs)
}

// MerkleizeVecComposite implements the SSZ merkleization algorithm for a vector
// of composite types.
func MerkleizeVecComposite[
	SpecT any, U64T U64[U64T], RootT ~[32]byte, C Composite[SpecT, RootT],
](
	value []C,
) (RootT, error) {
	htrs := make([]RootT, len(value))
	var err error
	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}
	return Merkleize[U64T, RootT](htrs)
}

// MerkleizeListComposite implements the SSZ merkleization algorithm for a list
// of composite types.
func MerkleizeListComposite[
	SpecT any, U64T U64[U64T], RootT ~[32]byte,
	C Composite[SpecT, RootT],
](
	value []C,
	limit uint64,
) (RootT, error) {
	htrs := make([]RootT, len(value))
	var err error
	for i, el := range value {
		htrs[i], err = el.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}
	root, err := Merkleize[U64T, RootT](
		htrs,
		ChunkCountCompositeList[C](value, limit),
	)
	if err != nil {
		return RootT{}, err
	}
	return merkle.MixinLength(root, uint64(len(value))), nil
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
func Merkleize[U64T U64[U64T], RootT, ChunkT ~[32]byte](
	chunks []ChunkT,
	limit ...uint64,
) (RootT, error) {
	var (
		effectiveLimit  U64T
		effectiveChunks []ChunkT
		lenChunks       = uint64(len(chunks))
	)

	//#nosec:G701 // This is a safe operation.
	switch {
	case len(limit) == 0:
		//#nosec:G701 // This is a safe operation.
		effectiveLimit = U64T(lenChunks).NextPowerOfTwo()
	case limit[0] >= lenChunks:
		//#nosec:G701 // This is a safe operation.
		effectiveLimit = U64T(limit[0]).NextPowerOfTwo()
	default:
		//#nosec:G701 // This is a safe operation.
		if limit[0] < lenChunks {
			return RootT{}, errors.New("input exceeds limit")
		}
		//#nosec:G701 // This is a safe operation.
		effectiveLimit = U64T(limit[0])
	}

	effectiveChunks = PadTo(chunks, effectiveLimit)
	if len(effectiveChunks) == 1 {
		return RootT(effectiveChunks[0]), nil
	}

	return merkle.NewRootWithMaxLeaves[U64T, ChunkT, RootT](
		effectiveChunks,
		//#nosec:G701 // This is a safe operation.
		uint64(effectiveLimit),
	)
}
