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

// Package trie defines utilities for sparse merkle tries for Ethereum
// consensus.
package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/crypto/sha256"
	byteslib "github.com/berachain/beacon-kit/lib/bytes"
	"github.com/cockroachdb/errors"
	"github.com/protolambda/ztyp/tree"
)

const (
	// 2^63 would overflow.
	MaxTrieDepth = 62
)

// SparseMerkleTrie implements a sparse, general purpose Merkle trie
// to be used across Ethereum consensus functionality.
type SparseMerkleTrie struct {
	depth    uint64
	branches [][][]byte
	// list of provided items before hashing them into leaves.
	originalItems [][]byte
}

// NewFromItems constructs a Merkle trie
// from a sequence of byte slices.
func NewFromItems(
	items [][]byte,
	depth uint64,
) (*SparseMerkleTrie, error) {
	if len(items) == 0 {
		return &SparseMerkleTrie{}, errors.New(
			"no items provided to generate Merkle trie",
		)
	}
	if depth == 0 {
		return &SparseMerkleTrie{}, errors.New("depth must be greater than 0")
	}
	if depth > MaxTrieDepth {
		// PowerOf2 would overflow
		return &SparseMerkleTrie{}, errors.New(
			"supported merkle trie depth exceeded (max uint64 depth is 63, " +
				"theoretical max sparse merkle trie depth is 64)")
	}

	leaves := items
	layers := make([][][]byte, depth+1)
	transformedLeaves := make([][]byte, len(leaves))
	for i := range leaves {
		arr := byteslib.ToBytes32(leaves[i])
		transformedLeaves[i] = arr[:]
	}
	layers[0] = transformedLeaves
	for i := uint64(0); i < depth; i++ {
		if len(layers[i])%2 == 1 {
			layers[i] = append(layers[i], tree.ZeroHashes[i][:])
		}
		updatedValues := make([][]byte, 0)
		for j := 0; j < len(layers[i]); j += 2 {
			concat := sha256.Hash(append(layers[i][j], layers[i][j+1]...))
			updatedValues = append(updatedValues, concat[:])
		}
		layers[i+1] = updatedValues
	}
	return &SparseMerkleTrie{
		branches:      layers,
		originalItems: items,
		depth:         depth,
	}, nil
}

// Items returns the original items passed in when creating the Merkle trie.
func (m *SparseMerkleTrie) Items() [][]byte {
	return m.originalItems
}

// HashTreeRoot returns the hash root of the Merkle trie
// defined in the deposit contract.
func (m *SparseMerkleTrie) HashTreeRoot() ([32]byte, error) {
	var enc [32]byte
	numItems := uint64(len(m.originalItems))
	if len(m.originalItems) == 1 &&
		bytes.Equal(m.originalItems[0], tree.ZeroHashes[0][:]) {
		// Accounting for empty tries
		numItems = 0
	}
	binary.LittleEndian.PutUint64(enc[:], numItems)
	return sha256.Hash(
		append(m.branches[len(m.branches)-1][0], enc[:]...),
	), nil
}

// Insert an item into the trie.
func (m *SparseMerkleTrie) Insert(item []byte, index int) error {
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
			parentHash := sha256.Hash(append(root[:], neighbor...))
			root = parentHash
		} else {
			parentHash := sha256.Hash(append(neighbor, root[:]...))
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

// Copy performs a deep copy of the trie.
func (m *SparseMerkleTrie) Copy() *SparseMerkleTrie {
	dstBranches := make([][][]byte, len(m.branches))
	for i1, srcB1 := range m.branches {
		dstBranches[i1] = byteslib.SafeCopy2D(srcB1)
	}

	return &SparseMerkleTrie{
		depth:         m.depth,
		branches:      dstBranches,
		originalItems: byteslib.SafeCopy2D(m.originalItems),
	}
}

// NumOfItems returns the num of items stored in
// the sparse merkle trie. We handle a special case
// where if there is only one item stored and it is an
// empty 32-byte root.
func (m *SparseMerkleTrie) NumOfItems() int {
	var zeroBytes [32]byte
	if len(m.originalItems) == 1 &&
		bytes.Equal(m.originalItems[0], zeroBytes[:]) {
		return 0
	}
	return len(m.originalItems)
}
