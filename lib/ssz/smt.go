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
	"github.com/berachain/beacon-kit/lib/ssz/common"
	"github.com/protolambda/ztyp/tree"
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

// BuildTreeFromChunks uses our optimized routine to build
// a SparseMerkleTrie from a vector of 32-byte chunks.
// limit=0 means no limit.
func BuildTreeFromChunks(
	chunks [][32]byte,
	limit uint64,
) (*SparseMerkleTrie, error) {
	length := uint64(len(chunks))
	// No limit, set the limit to the length of the chunks.
	if limit == 0 {
		limit = length
	}

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

	if length == 1 {
		return &SparseMerkleTrie{
			depth:    depth,
			branches: [][][32]byte{{chunks[0]}},
		}, nil
	}

	// Build the binary tree from the chunks.
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

// ChunkifyBytes pads the input byte slice to the next multiple of 32,
// if needed, and chunkifies it into 32-byte chunks.
func ChunkifyBytes(input []byte) ([][32]byte, error) {
	// Add (chunkSize - 1) to round up the division
	// to the next multiple of chunkSize.
	var chunkSize = common.BytesPerChunk
	numChunks := (len(input) + (chunkSize - 1)) / chunkSize
	if numChunks == 0 {
		return nil, ErrInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[i*chunkSize:])
	}
	return chunks, nil
}

// BuildTreeFromBasic builds a SparseMerkleTrie from a basic object,
// or a vector of basic objects (fixedSize=true)
// or from a list of basic objects (fixedSize=false).
func BuildTreeFromBasic(
	value common.SSZObject,
	fixedSized bool,
) (*SparseMerkleTrie, error) {
	bz, err := value.Marshal()
	if err != nil {
		return nil, err
	}
	chunks, err := ChunkifyBytes(bz)
	if err != nil {
		return nil, err
	}
	var limit uint64
	if !fixedSized {
		// TODO: Limit for the list of basic objects.
		limit = 0
	}
	body, err := BuildTreeFromChunks(chunks, limit)
	if err != nil {
		return nil, err
	}
	if fixedSized {
		return body, nil
	}
	// Mix in the length of the list of basic objects.
	return mixInLength(body, uint64(len(chunks))), nil
}

// BuildTreeFromComposite builds a SparseMerkleTrie from a composite type.
func BuildTreeFromComposite(
	value common.Composite,
	fixedSize bool,
) (*SparseMerkleTrie, error) {
	elements := value.Elements()
	length := uint64(len(elements))
	elemRoots := make([][32]byte, length)
	var err error
	for i, el := range elements {
		elemRoots[i], err = el.HashTreeRoot()
		if err != nil {
			return nil, err
		}
	}
	// Per the spec, if fixedSize, limit=None.
	// Otherwise, limit=(32*length + 31) / 32 = length.
	// In both cases, the effect is the same as limit=length.
	body, err := BuildTreeFromChunks(elemRoots, length)
	if err != nil {
		return nil, err
	}
	if fixedSize {
		return body, nil
	}
	// Mix in the length of the variable-length composite type.
	return mixInLength(body, length), nil
}

func mixInLength(t *SparseMerkleTrie, length uint64) *SparseMerkleTrie {
	bodyRoot := t.HashTreeRoot()
	root := tree.GetHashFn().Mixin(bodyRoot, length)
	return &SparseMerkleTrie{
		depth:    t.depth + 1,
		branches: append(t.branches, [][32]byte{root}),
	}
}
