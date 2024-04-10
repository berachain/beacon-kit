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

package tree

import (
	"bytes"
	"encoding/binary"
	"fmt"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	sha256 "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"
)

const (
	// 2^63 would overflow.
	MaxDepth = 62
)

// SparseMerkleTree implements a sparse, general purpose Merkle tree
// to be used across Ethereum consensus functionality.
type SparseMerkleTree struct {
	depth    uint64
	branches [][][]byte
	// list of provided items before hashing them into leaves.
	originalItems [][]byte
}

// NewFromItems constructs a Merkle tree from a sequence of byte slices.
func NewFromItems(items [][]byte, depth uint64) (*SparseMerkleTree, error) {
	switch {
	case len(items) == 0:
		return nil, ErrEmptyItems
	case depth == 0:
		return nil, ErrZeroDepth
	case depth > MaxDepth:
		return nil, ErrExceededDepth
	}

	transformedLeaves := make([][]byte, len(items))
	for i, item := range items {
		tl := byteslib.ToBytes32(item)
		transformedLeaves[i] = tl[:]
	}

	layers := make([][][]byte, depth+1)
	layers[0] = transformedLeaves

	for i := uint64(0); i < depth; i++ {
		currentLayer := layers[i]
		//nolint:gomnd // we divide by 2 to get the next layer size.
		nextLayerSize := (len(currentLayer) + 1) / 2
		nextLayer := make([][]byte, nextLayerSize)
		for j := 0; j < len(currentLayer); j += 2 {
			left := currentLayer[j]
			var right []byte
			if j+1 < len(currentLayer) {
				right = currentLayer[j+1]
			} else {
				right = tree.ZeroHashes[i][:]
			}
			h := sha256.Sum256(append(left, right...))
			nextLayer[j/2] = h[:]
		}
		layers[i+1] = nextLayer
	}

	return &SparseMerkleTree{
		branches:      layers,
		originalItems: items,
		depth:         depth,
	}, nil
}

// Items returns the original items passed in when creating the Merkle tree.
func (m *SparseMerkleTree) Items() [][]byte {
	return m.originalItems
}

// HashTreeRoot returns the hash root of the Merkle tree
// defined in the deposit contract.
func (m *SparseMerkleTree) HashTreeRoot() ([32]byte, error) {
	var enc [32]byte
	numItems := uint64(len(m.originalItems))
	if len(m.originalItems) == 1 &&
		bytes.Equal(m.originalItems[0], tree.ZeroHashes[0][:]) {
		numItems = 0
	}
	binary.LittleEndian.PutUint64(enc[:], numItems)
	return sha256.Sum256(
		append(m.branches[len(m.branches)-1][0], enc[:]...),
	), nil
}

// Insert an item into the tree.
func (m *SparseMerkleTree) Insert(item []byte, index int) error {
	if index < 0 {
		return fmt.Errorf("negative index provided: %d", index)
	}
	for index >= len(m.branches[0]) {
		m.branches[0] = append(m.branches[0], tree.ZeroHashes[0][:])
	}
	someItem := byteslib.ToBytes32(item)
	m.branches[0][index] = someItem[:]
	if index >= len(m.originalItems) {
		m.originalItems = append(m.originalItems, someItem[:])
	} else {
		m.originalItems[index] = someItem[:]
	}
	currentIndex := index
	root := byteslib.ToBytes32(item)
	two := 2
	for i := uint64(0); i < m.depth; i++ {
		isLeft := currentIndex%two == 0
		neighborIdx := currentIndex ^ 1
		var neighbor []byte
		if neighborIdx >= len(m.branches[i]) {
			neighbor = tree.ZeroHashes[i][:]
		} else {
			neighbor = m.branches[i][neighborIdx]
		}
		if isLeft {
			parentHash := sha256.Sum256(append(root[:], neighbor...))
			root = parentHash
		} else {
			parentHash := sha256.Sum256(append(neighbor, root[:]...))
			root = parentHash
		}
		parentIdx := currentIndex / two
		if len(m.branches[i+1]) == 0 || parentIdx >= len(m.branches[i+1]) {
			newItem := root
			m.branches[i+1] = append(m.branches[i+1], newItem[:])
		} else {
			newItem := root
			m.branches[i+1][parentIdx] = newItem[:]
		}
		currentIndex = parentIdx
	}
	return nil
}

// Copy performs a deep copy of the tree.
func (m *SparseMerkleTree) Copy() *SparseMerkleTree {
	dstBranches := make([][][]byte, len(m.branches))
	for i1, srcB1 := range m.branches {
		dstBranches[i1] = byteslib.SafeCopy2D(srcB1)
	}

	return &SparseMerkleTree{
		depth:         m.depth,
		branches:      dstBranches,
		originalItems: byteslib.SafeCopy2D(m.originalItems),
	}
}

// NumOfItems returns the num of items stored in
// the sparse merkle tree. We handle a special case
// where if there is only one item stored and it is an
// empty 32-byte root.
func (m *SparseMerkleTree) NumOfItems() int {
	var zeroBytes [32]byte
	if len(m.originalItems) == 1 &&
		bytes.Equal(m.originalItems[0], zeroBytes[:]) {
		return 0
	}
	return len(m.originalItems)
}
