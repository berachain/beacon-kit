// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle

import (
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
	"github.com/prysmaticlabs/gohashtree"
)

const (
	// 2^63 would overflow.
	MaxTreeDepth = 62
)

// Tree[LeafT, RootT] implements a Merkle tree that has been optimized to
// handle leaves that are 32 bytes in size.
type Tree[LeafT, RootT ~[32]byte] struct {
	depth    uint8
	branches [][]LeafT
	leaves   []LeafT
}

// NewTreeFromLeaves constructs a Merkle tree, with the minimum
// depth required to support the number of leaves.
func NewTreeFromLeaves[LeafT, RootT ~[32]byte](
	leaves []LeafT,
) (*Tree[LeafT, RootT], error) {
	return NewTreeFromLeavesWithDepth[LeafT, RootT](
		leaves,
		math.U64(len(leaves)).NextPowerOfTwo().ILog2Ceil(),
	)
}

// NewTreeWithMaxLeaves constructs a Merkle tree with a maximum number of
// leaves.
func NewTreeWithMaxLeaves[LeafT, RootT ~[32]byte](
	leaves []LeafT,
	maxLeaves uint64,
) (*Tree[LeafT, RootT], error) {
	return NewTreeFromLeavesWithDepth[LeafT, RootT](
		leaves,
		math.U64(maxLeaves).NextPowerOfTwo().ILog2Ceil(),
	)
}

// NewTreeFromLeaves constructs a Merkle tree from a sequence of byte slices.
// It will fill the tree with zero hashes to create the required depth.
func NewTreeFromLeavesWithDepth[LeafT, RootT ~[32]byte](
	leaves []LeafT,
	depth uint8,
) (*Tree[LeafT, RootT], error) {
	if err := verifySufficientDepth(len(leaves), depth); err != nil {
		return &Tree[LeafT, RootT]{}, err
	}

	layers := make([][]LeafT, depth+1)
	layers[0] = leaves

	var err error
	for d := range depth {
		currentLayer := layers[d]
		if len(currentLayer)%2 == 1 {
			currentLayer = append(currentLayer, zero.Hashes[d])
		}
		layers[d+1], err = BuildParentTreeRoots[LeafT, LeafT](currentLayer)
		if err != nil {
			return &Tree[LeafT, RootT]{}, err
		}
	}

	return &Tree[LeafT, RootT]{
		branches: layers,
		leaves:   leaves,
		depth:    depth,
	}, nil
}

// Insert an item into the tree.
func (m *Tree[LeafT, RootT]) Insert(item [32]byte, index int) error {
	if index < 0 {
		return errors.Wrap(ErrNegativeIndex, fmt.Sprintf("index: %d", index))
	}
	for index >= len(m.branches[0]) {
		m.branches[0] = append(m.branches[0], zero.Hashes[0])
	}
	m.branches[0][index] = item
	if index >= len(m.leaves) {
		m.leaves = append(m.leaves, item)
	} else {
		m.leaves[index] = item
	}

	var (
		hashFn       func([]byte) [32]byte
		neighbor     = [32]byte{}
		input        = [64]byte{}
		currentIndex = index
		root         = item
	)

	//nolint:mnd // 5 as defined by the library.
	if m.depth > 5 {
		hashFn = sha256.CustomSHA256Hasher()
	} else {
		hashFn = sha256.Sum256
	}

	for i := range m.depth {
		if neighborIdx := currentIndex ^ 1; neighborIdx >= len(m.branches[i]) {
			neighbor = zero.Hashes[i]
		} else {
			neighbor = m.branches[i][neighborIdx]
		}

		//nolint:mnd // 2 is allowed.
		if isLeft := currentIndex%2 == 0; isLeft {
			copy(input[0:32], root[:])
			copy(input[32:64], neighbor[:])
		} else {
			copy(input[0:32], neighbor[:])
			copy(input[32:64], root[:])
		}
		root = hashFn(input[:])

		//nolint:mnd // 2 is allowed.
		parentIdx := currentIndex / 2
		if len(m.branches[i+1]) == 0 || parentIdx >= len(m.branches[i+1]) {
			m.branches[i+1] = append(m.branches[i+1], root)
		} else {
			m.branches[i+1][parentIdx] = root
		}
		currentIndex = parentIdx
	}
	return nil
}

// Root returns the root of the Merkle tree.
func (m *Tree[LeafT, RootT]) Root() [32]byte {
	return m.branches[len(m.branches)-1][0]
}

// HashTreeRoot returns the Root of the Merkle tree with the
// number of leaves mixed in.
func (m *Tree[LeafT, RootT]) HashTreeRoot() ([32]byte, error) {
	numItems := uint64(len(m.leaves))
	if len(m.leaves) == 1 &&
		m.leaves[0] == zero.Hashes[0] {
		numItems = 0
	}
	return MixinLength(m.Root(), numItems), nil
}

// MerkleProof computes a proof from a tree's branches using a Merkle index.
func (m *Tree[LeafT, RootT]) MerkleProof(leafIndex uint64) ([][32]byte, error) {
	numLeaves := uint64(len(m.branches[0]))
	if leafIndex >= numLeaves {
		return nil, errors.Newf(
			"merkle index out of range in tree, max range: %d, received: %d",
			numLeaves,
			leafIndex,
		)
	}
	proof := make([][32]byte, m.depth)
	for i := range m.depth {
		subIndex := (leafIndex >> i) ^ 1
		if subIndex < uint64(len(m.branches[i])) {
			proof[i] = m.branches[i][subIndex]
		} else {
			proof[i] = zero.Hashes[i]
		}
	}
	return proof, nil
}

// MerkleProofWithMixin computes a proof from a tree's branches using a Merkle
// index.
func (m *Tree[LeafT, RootT]) MerkleProofWithMixin(
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

// MixinLength takes a root element and mixes in the length of the elements
// that were hashed to produce it.
func MixinLength[RootT ~[32]byte](element RootT, length uint64) RootT {
	chunks := make([][32]byte, two)
	chunks[0] = element
	binary.LittleEndian.PutUint64(chunks[1][:], length)
	if err := gohashtree.Hash(chunks, chunks); err != nil {
		return [32]byte{}
	}
	return chunks[0]
}
