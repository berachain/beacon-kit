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

package sha256

import (
	sha256 "github.com/minio/sha256-simd"
	"github.com/sourcegraph/conc/iter"
)

// This hashing library provides multiple ways to utilize the Hash function:
// 1. Directly pass a raw []byte slice to the Hash function.
// 2. Implement the `Hashable` interface and invoke the HashTreeRoot method
//    for a customized hashing approach.

// Hash returns the SHA256 hash of the input bytes.
func Hash(bz []byte) [32]byte {
	return sha256.Sum256(bz)
}

// HashArray returns the SHA256 hash of each input byte slice.
func HashArray(input [][]byte) [][32]byte {
	roots := make([][32]byte, len(input))
	iter.ForEachIdx[[]byte](
		input,
		func(i int, el *[]byte) {
			roots[i] = Hash(*el)
		},
	)
	return roots
}

// HashRoot returns the SHA256 merkle root of the input bytes.
func HashRootBz(input [][]byte) [32]byte {
	return UnsafeMerkleizeVector(HashArray(input), uint64(len(input)))
}

// HashRootBzBytes returns the SHA256 merkle root of the input bytes as a byte slice.
func HashRootBzAsSlice(input [][]byte) []byte {
	bz := HashRootBz(input)
	return bz[:]
}

// HashRootAndMixinLength returns the SHA256 merkle root of the input bytes with
// the length mixed in.
func HashRootAndMixinLengthBz(input [][]byte) [32]byte {
	return UnsafeMerkleizeVectorAndMixinLength(HashArray(input), uint64(len(input)))
}

// HashRootAndMixinLength returns the SHA256 merkle root of the input bytes with
// the length mixed in.
func HashRootAndMixinLengthAsBzSlice(input [][]byte) []byte {
	bz := HashRootAndMixinLengthBz(input)
	return bz[:]
}

// Hashable is an interface for objects that can be hashed.
func HashElement[H Hashable](el H) ([32]byte, error) {
	return el.HashTreeRoot()
}

// HashElements hashes each element in the list and then returns each item as a
// 32 byte buffer. Each element is individually hashed to produce a corresponding
// root. This process is applied to all elements in the input list, resulting in a
// list of roots.
func HashElements[H Hashable](input []H) ([][32]byte, error) {
	var (
		err   error
		roots = make([][32]byte, len(input))
	)

	// Hash each element in the list.
	iter.ForEachIdx[H](
		input,
		func(i int, el *H) {
			var localErr error
			roots[i], localErr = (*el).HashTreeRoot()
			if err != nil {
				err = localErr
			}
		},
	)

	// Return the list of roots and any error encountered.
	return roots, err
}

// HashRoot returns the SHA256 merkle root of the input bytes.
func HashRoot[H Hashable](input []H) [32]byte {
	b, err := BuildMerkleRoot(input, uint64(len(input)))
	if err != nil {
		panic(err)
	}
	return b
}

// HashRootAsSlice returns the SHA256 merkle root of the input bytes as a byte slice.
func HashRootAsSlice[H Hashable](input []H) []byte {
	bz := HashRoot(input)
	return bz[:]
}

// HashRootAndMixinLength returns the SHA256 merkle root of the input bytes with
// the length mixed in.
func HashRootAndMixinLength[H Hashable](input []H) [32]byte {
	b, err := BuildMerkleRootAndMixinLength(input, uint64(len(input)))
	if err != nil {
		panic(err)
	}
	return b
}

// HashRootAndMixinLengthAsSlice returns the SHA256 merkle root of the input bytes with
// the length mixed in as a byte slice.
func HashRootAndMixinLengthAsSlice[H Hashable](input []H) []byte {
	bz := HashRootAndMixinLength(input)
	return bz[:]
}
