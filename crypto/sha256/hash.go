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
	"golang.org/x/sync/errgroup"
)

// This hashing library provides multiple ways to utilize the Hash function:
// 1. Directly pass a raw []byte slice to the Hash function.
// 2. Implement the `Hashable` interface and invoke the HashTreeRoot method
//    for a customized hashing approach.

// Hash returns the SHA256 hash of the input bytes.
func Hash(bz []byte) [32]byte {
	return sha256.Sum256(bz)
}

// HashArray returns the SHA256 hashes of the input bytes in parallel while maintaining the order.
func HashArray(input [][]byte) ([][32]byte, error) {
	roots := make([][32]byte, len(input))

	// Create a channel to receive the results
	resultCh := make(chan struct {
		index int
		hash  [32]byte
	}, len(input))

	// Iterate over the input bytes and hash them in parallel while maintaining order
	for i, el := range input {
		go func(index int, data []byte) {
			resultCh <- struct {
				index int
				hash  [32]byte
			}{index, Hash(data)}
		}(i, el)
	}

	// Collect the results from the channel and place them in the correct order
	for range input {
		res := <-resultCh
		roots[res.index] = res.hash
	}

	close(resultCh)

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

// Hashable is an interface for objects that can be hashed.
func HashElement[H Hashable](el H) ([32]byte, error) {
	return el.HashTreeRoot()
}

// HashElements hashes each element in the list and then returns each item as a 32 byte buffer.
// Where each Element is hashed individually to produce a corresponding Root.
// This process is applied to all elements in the input list, resulting in a list of roots.
func HashElements[H Hashable](input []H) ([][32]byte, error) {
	roots := make([][32]byte, len(input))
	// Create a channel to receive the results
	resultCh := make(chan struct {
		index int
		hash  [32]byte
	}, len(input))

	// Use error group to handle errors from goroutines
	var eg errgroup.Group

	// Iterate over the input bytes and hash them in parallel while maintaining order
	for i, el := range input {
		i, el := i, el // Capture loop variables
		eg.Go(func() error {
			data, err := HashElement(el)
			if err != nil {
				return err
			}
			resultCh <- struct {
				index int
				hash  [32]byte
			}{i, data}
			return nil
		})
	}

	// Collect the results from the channel and place them in the correct order
	for range input {
		res := <-resultCh
		roots[res.index] = res.hash
	}

	close(resultCh)

	// Check for any errors from the goroutines
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return roots, nil
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
