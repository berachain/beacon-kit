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

package ssz_test

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/stretchr/testify/require"
)

type RelevantChainSpecFunctionsToThisType interface{}

// BasicItem represnets a basic item in the SSZ Spec.
type BasicItem uint64

// SizeSSZ returns the size of the U64 in bytes.
func (u BasicItem) SizeSSZ() int {
	return 8
}

// HashTreeRoot computes the Merkle root of the U64 using SSZ hashing rules.
func (u BasicItem) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	var hashRoot [32]byte
	copy(hashRoot[:], buf)
	return hashRoot, nil
}

// BasiContainer represents a container of two basic items.
type BasicContainer[SpecT any] struct {
	Item1 BasicItem
	Item2 BasicItem
}

// SizeSSZ returns the size of the container in bytes.
func (c *BasicContainer[SpecT]) SizeSSZ() int {
	// TODO: We should be able to generalize SizeSSZ() as well.
	return c.Item1.SizeSSZ() + c.Item2.SizeSSZ()
}

// HashTreeRoot computes the Merkle root of the container using SSZ hashing
// rules.
func (c *BasicContainer[SpecT]) HashTreeRoot() ([32]byte, error) {
	return ssz.MerkleizeContainer[any, math.U64](c)
}

// TestBasicItemMerkleization tests the Merkleization of a basic item.
func TestBasicContainerMerkleization(t *testing.T) {
	container := BasicContainer[any]{
		Item1: BasicItem(1),
		Item2: BasicItem(2),
	}

	// Merkleize the container.
	actualRoot, err := container.HashTreeRoot()
	require.NoError(t, err)

	// Manually compute our own root, using our merkle tree knowledge.
	htr1, err := container.Item1.HashTreeRoot()
	require.NoError(t, err)
	htr2, err := container.Item2.HashTreeRoot()
	require.NoError(t, err)
	expectedRoot := sha256.Sum256(append(htr1[:], htr2[:]...))

	// Should match
	require.Equal(t, expectedRoot, actualRoot)
}
