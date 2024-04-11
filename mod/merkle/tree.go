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

	"github.com/berachain/beacon-kit/mod/merkle/htr"
	"github.com/cockroachdb/errors"
	sha256 "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"
)

const (
	// 2^63 would overflow.
	MaxTreeDepth = 62
)

// Tree implements a Merkle tree that has been optimized to
// handle leaves that are 32 bytes in size.
type Tree struct {
	depth    uint64
	branches [][][32]byte
	leaves   [][32]byte
}

// NewTreeFromLeaves constructs a Merkle tree, with the minimum
// depth required to support the number of leaves.
func NewTreeFromLeaves(
	leaves [][32]byte,
) (*Tree, error) {
	return NewTreeFromLeavesWithDepth(leaves, uint64(len(leaves)))
}

// NewTreeFromLeaves constructs a Merkle tree from a sequence of byte slices.
// It will fill the tree with zero hashes to create the required depth.
func NewTreeFromLeavesWithDepth(
	leaves [][32]byte,
	depth uint64,
) (*Tree, error) {
	numLeaves := len(leaves)
	switch {
	case numLeaves == 0:
		return &Tree{}, ErrEmptyLeaves
	case depth == 0:
		return &Tree{}, ErrZeroDepth
	case depth > MaxTreeDepth:
		return &Tree{}, ErrExceededDepth
	case numLeaves > (1 << depth):
		return &Tree{}, errors.Wrap(
			ErrInsufficientDepthForLeaves,
			fmt.Sprintf("attempted to store %d leaves with depth %d",
				numLeaves, depth))
	}

	layers := make([][][32]byte, depth+1)
	layers[0] = leaves

	var err error
	for i := uint64(0); i < depth; i++ {
		currentLayer := layers[i]
		if len(currentLayer)%2 == 1 {
			currentLayer = append(currentLayer, tree.ZeroHashes[i])
		}
		layers[i+1], err = htr.BuildParentTreeRoots(currentLayer)
		if err != nil {
			return &Tree{}, err
		}
	}

	return &Tree{
		branches: layers,
		leaves:   leaves,
		depth:    depth,
	}, nil
}

// Insert an item into the tree.
func (m *Tree) Insert(item [32]byte, index int) error {
	if index < 0 {
		return errors.Wrap(ErrNegativeIndex, fmt.Sprintf("index: %d", index))
	}
	for index >= len(m.branches[0]) {
		m.branches[0] = append(m.branches[0], tree.ZeroHashes[0])
	}
	m.branches[0][index] = item
	if index >= len(m.leaves) {
		m.leaves = append(m.leaves, item)
	} else {
		m.leaves[index] = item
	}
	neighbor := [32]byte{}
	input := [64]byte{}
	currentIndex := index
	root := item
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
			m.branches[i+1] = append(m.branches[i+1], root)
		} else {
			copy(m.branches[i+1][parentIdx][:], root[:])
		}
		currentIndex = parentIdx
	}
	return nil
}

// Root returns the root of the Merkle tree.
func (m *Tree) Root() ([32]byte, error) {
	return sha256.Sum256(m.branches[len(m.branches)-1][0][:]), nil
}

// HashTreeRoot returns the Root of the Merkle tree with the
// number of leaves mixed in.
func (m *Tree) HashTreeRoot() ([32]byte, error) {
	var enc [32]byte
	numItems := uint64(len(m.leaves))
	if len(m.leaves) == 1 &&
		m.leaves[0] == tree.ZeroHashes[0] {
		numItems = 0
	}
	binary.LittleEndian.PutUint64(enc[:], numItems)
	hashInput := append(m.branches[len(m.branches)-1][0][:], enc[:]...)
	return sha256.Sum256(hashInput), nil
}

// MerkleProof computes a proof from a tree's branches using a Merkle index.
func (m *Tree) MerkleProof(leafIndex uint64) ([][32]byte, error) {
	numLeaves := uint64(len(m.branches[0]))
	if leafIndex >= numLeaves {
		return nil, fmt.Errorf(
			"merkle index out of range in tree, max range: %d, received: %d",
			numLeaves,
			leafIndex,
		)
	}
	proof := make([][32]byte, m.depth)
	for i := uint64(0); i < m.depth; i++ {
		subIndex := (leafIndex >> i) ^ 1
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
func (m *Tree) MerkleProofWithMixin(
	index uint64,
) ([][32]byte, error) {
	proof, err := m.MerkleProof(index)
	if err != nil {
		return nil, err
	}

	mixin := [32]byte{}
	binary.LittleEndian.PutUint64(mixin[:8], uint64(len(m.leaves)))
	return append(proof, mixin), nil
}
