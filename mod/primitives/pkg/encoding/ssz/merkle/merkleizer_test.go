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

package merkle_test

import (
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/stretchr/testify/require"
)

// BasicItem represents a basic item in the SSZ Spec.
type BasicItem uint64

// SizeSSZ returns the size of the U64 in bytes.
func (u BasicItem) SizeSSZ() int {
	return 8
}

// MarshalSSZ marshals the U64 into a byte slice.
func (u BasicItem) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// HashTreeRoot computes the Merkle root of the U64 using SSZ hashing rules.
func (u BasicItem) HashTreeRoot() ([32]byte, error) {
	// In practice we can use a simpler function.
	return merkle.
		NewMerkleizer[[32]byte, BasicItem]().MerkleizeBasic(u)
}

// BasicContainer represents a container of two basic items.
type BasicContainer[SpecT any] struct {
	Item1 BasicItem
	Item2 BasicItem
}

// SizeSSZ returns the size of the container in bytes.
func (c *BasicContainer[SpecT]) SizeSSZ() int {
	return c.Item1.SizeSSZ() + c.Item2.SizeSSZ()
}

// HashTreeRoot computes the Merkle root of the container using SSZ hashing
// rules.
func (c *BasicContainer[SpecT]) HashTreeRoot() ([32]byte, error) {
	return merkle.NewMerkleizer[[32]byte, BasicItem]().
		MerkleizeVectorCompositeOrContainer(
			[]BasicItem{c.Item1, c.Item2},
		)
}

// MarshalSSZ marshals the container into a byte slice.
func (c *BasicContainer[SpecT]) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	item1Bytes, err := c.Item1.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	item2Bytes, err := c.Item2.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	copy(buf[:8], item1Bytes)
	copy(buf[8:], item2Bytes)
	return buf, nil
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
