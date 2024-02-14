// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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
	"github.com/minio/sha256-simd"
)

// Hash returns the SHA256 hash of the input bytes.
func Hash(bz []byte) [32]byte {
	return sha256.Sum256(bz)
}

// HashArray returns the SHA256 hashes of the input bytes.
func HashArray(input [][]byte) ([][32]byte, error) {
	roots := make([][32]byte, len(input))
	for i, el := range input {
		roots[i] = Hash(el)
	}
	return roots, nil
}

// HashRoot returns the SHA256 merkle root of the input bytes.
func HashRootBz(input [][]byte) [32]byte {
	// Hash all the elements in the array to create the
	// leaves of the merkle tree.
	roots, err := HashArray(input)
	if err != nil {
		return [32]byte{}
	}

	// Then build the merkle tree from the leaves.
	return UnsafeMerkleizeVector(roots, uint64(len(input)))
}

// HashRootBzBytes returns the SHA256 merkle root of the input bytes as a byte slice.
func HashRootBzAsSlice(input [][]byte) []byte {
	bz := HashRootBz(input)
	return bz[:]
}

// HashRootAndMixinLength returns the SHA256 merkle root of the input bytes with
// the length mixed in.
func HashRootAndMixinLengthBz(input [][]byte) [32]byte {
	roots, err := HashArray(input)
	if err != nil {
		return [32]byte{}
	}
	return UnsafeMerkleizeVectorAndMixinLength(roots, uint64(len(input)))
}

// HashRootAndMixinLength returns the SHA256 merkle root of the input bytes with
// the length mixed in.
func HashRootAndMixinLengthAsBzSlice(input [][]byte) []byte {
	bz := HashRootAndMixinLengthBz(input)
	return bz[:]
}

// HashRoot returns the SHA256 merkle root of the input bytes.
func HashRoot[H Hashable](input []H) [32]byte {
	b, _ := BuildMerkleRoot(input, uint64(len(input)))
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
	b, _ := BuildMerkleRootAndMixinLength(input, uint64(len(input)))
	return b
}

// HashRootAndMixinLengthAsSlice returns the SHA256 merkle root of the input bytes with
// the length mixed in as a byte slice.
func HashRootAndMixinLengthAsSlice[H Hashable](input []H) []byte {
	bz := HashRootAndMixinLength(input)
	return bz[:]
}
