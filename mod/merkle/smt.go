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

package merkle

import (
	"encoding/binary"
	"fmt"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	sha256 "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"
)

const (
	// 2^63 would overflow.
	MaxTreeDepth = 62
)

// SparseMerkleTree implements a sparse, general purpose Merkle tree
// to be used across Ethereum consensus functionality.
type SparseMerkleTree struct {
	depth         uint64
	branches      [][][32]byte
	originalItems [][32]byte
}

// NewTreeFromItems constructs a Merkle tree from a sequence of byte slices.
func NewTreeFromItems(
	items [][32]byte,
	depth uint64,
) (*SparseMerkleTree, error) {
	switch {
	case len(items) == 0:
		return nil, ErrEmptyItems
	case depth == 0:
		return nil, ErrZeroDepth
	case depth > MaxTreeDepth:
		return nil, ErrExceededDepth
	}

	layers := make([][][32]byte, depth+1)
	layers[0] = items

	for i := uint64(0); i < depth; i++ {
		currentLayer := layers[i]
		//nolint:gomnd // div 2.
		nextLayerSize := (len(currentLayer) + 1) / 2
		nextLayer := make([][32]byte, nextLayerSize)
		for j := 0; j < len(currentLayer); j += 2 {
			left := currentLayer[j]
			right := tree.ZeroHashes[i]
			if j+1 < len(currentLayer) {
				right = currentLayer[j+1]
			}
			hashInput := append(left[:], right[:]...)
			nextLayer[j/2] = sha256.Sum256(hashInput)
		}
		layers[i+1] = nextLayer
	}

	return &SparseMerkleTree{
		branches:      layers,
		originalItems: items,
		depth:         depth,
	}, nil
}

// Root returns the root of the Merkle tree.
func (m *SparseMerkleTree) Root() ([32]byte, error) {
	return sha256.Sum256(m.branches[len(m.branches)-1][0][:]), nil
}

// HashTreeRoot returns the Root of the Merkle tree with the
// number of leaves mixed in.
func (m *SparseMerkleTree) HashTreeRoot() ([32]byte, error) {
	var enc [32]byte
	numItems := uint64(len(m.originalItems))
	if len(m.originalItems) == 1 &&
		m.originalItems[0] == tree.ZeroHashes[0] {
		numItems = 0
	}
	binary.LittleEndian.PutUint64(enc[:], numItems)
	hashInput := append(m.branches[len(m.branches)-1][0][:], enc[:]...)
	return sha256.Sum256(hashInput), nil
}

// Items returns the original items passed in when creating the Merkle tree.
func (m *SparseMerkleTree) Items() [][32]byte {
	return m.originalItems
}

// Insert an item into the tree.
func (m *SparseMerkleTree) Insert(item []byte, index int) error {
	if index < 0 {
		return fmt.Errorf("negative index provided: %d", index)
	}
	for index >= len(m.branches[0]) {
		m.branches[0] = append(m.branches[0], tree.ZeroHashes[0])
	}
	someItem := byteslib.ToBytes32(item)
	m.branches[0][index] = someItem
	if index >= len(m.originalItems) {
		m.originalItems = append(m.originalItems, someItem)
	} else {
		m.originalItems[index] = someItem
	}
	neighbor := [32]byte{}
	input := [64]byte{}
	currentIndex := index
	root := byteslib.ToBytes32(item)
	for i := uint64(0); i < m.depth; i++ {
		if neighborIdx := currentIndex ^ 1; neighborIdx >= len(m.branches[i]) {
			neighbor = tree.ZeroHashes[i]
		} else {
			neighbor = m.branches[i][neighborIdx]
		}

		//nolint:gomnd
		if isLeft := currentIndex%2 == 0; isLeft {
			copy(input[0:32], root[:])
			copy(input[32:64], neighbor[:])
		} else {
			copy(input[0:32], neighbor[:])
			copy(input[32:64], root[:])
		}
		root = sha256.Sum256(input[:])

		//nolint:gomnd
		parentIdx := currentIndex / 2
		if len(m.branches[i+1]) == 0 || parentIdx >= len(m.branches[i+1]) {
			newItem := root
			m.branches[i+1] = append(m.branches[i+1], newItem)
		} else {
			newItem := root
			m.branches[i+1][parentIdx] = newItem
		}
		currentIndex = parentIdx
	}
	return nil
}

// MerkleProof computes a proof from a tree's branches using a Merkle index.
func (m *SparseMerkleTree) MerkleProof(index uint64) ([][32]byte, error) {
	numLeaves := uint64(len(m.branches[0]))
	if index >= numLeaves {
		return nil, fmt.Errorf(
			"merkle index out of range in tree, max range: %d, received: %d",
			numLeaves,
			index,
		)
	}
	proof := make([][32]byte, m.depth)
	for i := uint64(0); i < m.depth; i++ {
		subIndex := (index >> i) ^ 1
		if subIndex < uint64(len(m.branches[i])) {
			proof[i] = m.branches[i][subIndex]
		} else {
			proof[i] = tree.ZeroHashes[i]
		}
	}
	return proof, nil
}

// MerkleProofWithMixin computes a proof from a tree's branches using a Merkle
// index.
func (m *SparseMerkleTree) MerkleProofWithMixin(
	index uint64,
) ([][32]byte, error) {
	proof, err := m.MerkleProof(index)
	if err != nil {
		return nil, err
	}

	mixin := [32]byte{}
	binary.LittleEndian.PutUint64(mixin[:8], uint64(len(m.originalItems)))
	return append(proof, mixin), nil
}
