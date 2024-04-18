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

	"github.com/berachain/beacon-kit/mod/primitives/merkle"
)

// Merkleize hashes the packed value and returns the HTR.
func MerkleizeBasic[
	U64T U64[U64T], U256L U256LT,
	B Basic[RootT], RootT ~[32]byte](
	value B,
) (RootT, error) {
	return MerkleizeVecBasic[U64T, U256L, B, RootT]([]B{value})
}

// MerkleizeVec implements the SSZ merkleization algorithm for a vector of basic
// types.
func MerkleizeVecBasic[U64T U64[U64T], U256L U256LT, B Basic[RootT], RootT ~[32]byte](
	value []B,
) (RootT, error) {
	packed, err := Pack[U64T, U256L, B, RootT](value)
	if err != nil {
		return [32]byte{}, err
	}
	return Merkleize[U64T, RootT, RootT](packed)
}

// MerkleizeListBasic implements the SSZ merkleization algorithm for a list of
// basic types.
func MerkleizeListBasic[U64T U64[U64T], U256L U256LT, B Basic[RootT], RootT ~[32]byte](
	value []B,
	limit uint64,
) (RootT, error) {
	packed, err := Pack[U64T, U256L](value)
	if err != nil {
		return [32]byte{}, err
	}
	root, err := Merkleize[U64T, RootT, RootT](
		packed,
		ChunkCountBasicList[RootT, B](value, limit),
	)
	if err != nil {
		return [32]byte{}, err
	}
	return merkle.MixinLength(root, uint64(len(value))), nil
}

// TODO: MerkleizeBitlist

// MerkleizeContainer implements the SSZ merkleization algorithm for a
// container.
func MerkleizeContainer[U64T U64[U64T], C Composite[RootT], RootT ~[32]byte](
	value C,
) (RootT, error) {
	rValue := reflect.ValueOf(value)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}
	htrs := make([]RootT, rValue.NumField())
	var err error
	for i := range rValue.NumField() {
		field, ok := rValue.Field(i).Interface().(Hashable[RootT])
		if !ok {
			return RootT{}, fmt.Errorf(
				"field %d does not implement Hashable",
				i,
			)
		}
		htrs[i], err = field.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}
	return Merkleize[U64T, RootT, RootT](htrs)
}

// MerkleizeVecComposite implements the SSZ merkleization algorithm for a vector
// of composite types.
func MerkleizeVecComposite[U64T U64[U64T], C Composite[RootT], RootT ~[32]byte](
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
	return Merkleize[U64T, RootT, RootT](htrs)
}

// MerkleizeListComposite implements the SSZ merkleization algorithm for a list
// of composite types.
func MerkleizeListComposite[
	U64T U64[U64T], C Composite[RootT], RootT ~[32]byte,
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
	root, err := Merkleize[U64T, RootT, RootT](
		htrs,
		ChunkCountCompositeList(value, limit),
	)
	if err != nil {
		return RootT{}, err
	}
	return merkle.MixinLength(root, uint64(len(value))), nil
}
