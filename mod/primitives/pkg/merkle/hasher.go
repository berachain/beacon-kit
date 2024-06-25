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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
	"github.com/prysmaticlabs/gohashtree"
)

// Hasher can be re-used for constructing Merkle tree roots.
type Hasher[RootT ~[32]byte] struct {
	buffer Buffer[RootT]
}

// NewHasher creates a new merkle Hasher.
func NewHasher[RootT ~[32]byte](buffer Buffer[RootT]) *Hasher[RootT] {
	return &Hasher[RootT]{
		buffer: buffer,
	}
}

// NewRootWithMaxLeaves constructs a Merkle tree root from a set of.
func (m *Hasher[RootT]) NewRootWithMaxLeaves(
	leaves []RootT,
	length uint64,
) (RootT, error) {
	return m.NewRootWithDepth(
		leaves, math.U64(length).NextPowerOfTwo().ILog2Ceil(),
	)
}

// NewRootWithDepth constructs a Merkle tree root from a set of leaves.
func (m *Hasher[RootT]) NewRootWithDepth(
	leaves []RootT,
	depth uint8,
) (RootT, error) {
	// Return zerohash at depth
	if len(leaves) == 0 {
		return zero.Hashes[depth], nil
	}

	// Preallocate a single buffer large enough for the maximum layer size
	// TODO: It seems that BuildParentTreeRoots has different behaviour
	// when we pass leaves in directly.
	buf := m.buffer.Get((len(leaves) + 1) / two)

	var err error
	for i := range depth {
		layerLen := len(leaves)
		if layerLen%two == 1 {
			leaves = append(leaves, zero.Hashes[i])
		}

		newLayerSize := (layerLen + 1) / two
		if err = BuildParentTreeRoots(buf[:newLayerSize], leaves); err != nil {
			return zero.Hashes[depth], err
		}
		leaves, buf = buf[:newLayerSize], leaves
	}
	if len(leaves) != 1 {
		return zero.Hashes[depth], nil
	}
	return leaves[0], nil
}

// MixinLength takes a root element and mixes in the length of the elements
// that were hashed to produce it.
//
// TODO: move to ssz package.
func MixinLength[RootT ~[32]byte](element RootT, length uint64) RootT {
	chunks := make([][32]byte, two)
	chunks[0] = element
	binary.LittleEndian.PutUint64(chunks[1][:], length)
	if err := gohashtree.Hash(chunks, chunks); err != nil {
		return [32]byte{}
	}
	return chunks[0]
}
