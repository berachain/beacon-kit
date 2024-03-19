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
	"github.com/berachain/beacon-kit/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
)

const (
	chunkSize = 32
)

// SparseMerkleTrie implements a sparse, general purpose Merkle trie
// to be used across Ethereum consensus functionality.
type SparseMerkleTrie struct {
	depth uint8
	// List of branches for each level of the trie.
	// The first level is the leaves, and the last level is the root.
	branches [][][32]byte
}

// HashTreeRoot returns the hash root of the trie.
func (m *SparseMerkleTrie) HashTreeRoot() [32]byte {
	return m.branches[len(m.branches)-1][0]
}

// NewFromChunks uses our optimized routine to build
// a SparseMerkleTrie from a vector of 32-byte chunks.
func NewFromChunks(
	chunks [][32]byte,
	limit uint64,
) (*SparseMerkleTrie, error) {
	length := uint64(len(chunks))
	if limit < length {
		return nil, ErrInputExceedsLimit
	}

	// The chunks was virtually padded with zeroed chunks
	// to the next_pow_of_two of limit.
	depth := tree.CoverDepth(limit)
	// Return zerohash at depth
	if length == 0 {
		return &SparseMerkleTrie{
			depth:    depth,
			branches: [][][32]byte{{tree.ZeroHashes[depth]}},
		}, nil
	}

	layers := make([][][32]byte, depth+1)
	layers[0] = make([][32]byte, length)
	copy(layers[0], chunks)
	layerLen := len(chunks)
	for i := uint8(0); i < depth; i++ {
		oddNodeLength := layerLen%two == 1
		if oddNodeLength {
			chunks = append(chunks, tree.ZeroHashes[i])
		}
		var err error
		chunks, err = sha256.BuildParentTreeRoots(chunks)
		if err != nil {
			return nil, err
		}
		layerLen = len(chunks)
		layers[i+1] = make([][32]byte, layerLen)
		copy(layers[i+1], chunks)
	}
	// At the end of the loop, elements will only
	// contain the root of the trie.
	if len(chunks) != 1 {
		return nil, ErrBuildMerkleTree
	}
	return &SparseMerkleTrie{
		depth:    depth,
		branches: layers,
	}, nil
}

// NewFromByteSlice builds a SparseMerkleTrie from a byte slice.
// The byte slice can be considered as a vector of 32-byte chunks.
// The leaves of the trie are built by chunkifying the byte slice
// into 32-byte chunks.
func NewFromByteSlice(input []byte) (*SparseMerkleTrie, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / chunkSize
	if numChunks == 0 {
		return nil, ErrInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[chunkSize*i:])
	}
	return NewFromChunks(chunks, numChunks)
}

// NewFromVector builds a SparseMerkleTrie from
// a fixed-length vector of elements.
func NewFromVector[T Hashable](
	elements []T,
) (*SparseMerkleTrie, error) {
	length := uint64(len(elements))
	elemRoots := make([][32]byte, length)
	var err error
	for i, el := range elements {
		elemRoots[i], err = el.HashTreeRoot()
		if err != nil {
			return nil, err
		}
	}
	// Per the spec, limit=None, whose behavior
	// is the same as limit=length.
	return NewFromChunks(elemRoots, length)
}

// NewFromList builds a SparseMerkleTrie from a variable-length
// list of elements, with the length mixed in.
func NewFromList[T Hashable](
	elements []T,
) (*SparseMerkleTrie, error) {
	body, err := NewFromVector(elements)
	if err != nil {
		return nil, err
	}
	bodyRoot := body.HashTreeRoot()
	root := tree.GetHashFn().Mixin(bodyRoot, uint64(len(elements)))
	return &SparseMerkleTrie{
		depth:    body.depth + 1,
		branches: append(body.branches, [][32]byte{root}),
	}, nil
}

// NewFromContainer builds a SparseMerkleTrie from a container.
func NewFromContainer[C Container](
	container C,
) (*SparseMerkleTrie, error) {
	return NewFromVector(container.Fields())
}
