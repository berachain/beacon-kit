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

const (
	two = 2
)

// MerkleizeChunksSSZ hashes a list of chunks by building
// a merkle tree and returning the root.
func MerkleizeChunksSSZ(
	chunks [][32]byte,
	limit uint64,
) ([32]byte, error) {
	trie, err := NewFromChunks(chunks, limit)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}

// MerkleizeByteSliceSSZ hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSliceSSZ(input []byte) ([32]byte, error) {
	trie, err := NewFromByteSlice(input)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}

// MerkleizeListSSZ hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
func MerkleizeListSSZ[T Hashable](
	elements []T,
) ([32]byte, error) {
	trie, err := NewFromList(elements)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}

// MerkleizeVectorSSZ hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func MerkleizeVectorSSZ[T Hashable](
	elements []T,
) ([32]byte, error) {
	trie, err := NewFromVector(elements)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}

func MerkleizeContainerSSZ[C Container](
	container C,
) ([32]byte, error) {
	trie, err := NewFromContainer(container)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}
