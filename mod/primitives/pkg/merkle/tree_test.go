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
	"strconv"
	"testing"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/stretchr/testify/require"
)

func TestNewTreeFromLeavesWithDepth_NoItemsProvided(t *testing.T) {
	treeDepth := uint8(32)
	_, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		nil,
		treeDepth,
	)
	require.ErrorIs(t, err, merkle.ErrEmptyLeaves)
}

func TestNewTreeFromLeavesWithDepth_DepthSupport(t *testing.T) {
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("BB")),
		byteslib.ToBytes32([]byte("CCC")),
		byteslib.ToBytes32([]byte("DDDD")),
		byteslib.ToBytes32([]byte("EEEEE")),
		byteslib.ToBytes32([]byte("FFFFFF")),
		byteslib.ToBytes32([]byte("GGGGGGG")),
	}
	// Supported depth
	m1, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		merkle.MaxTreeDepth,
	)
	require.NoError(t, err)
	proof, err := m1.MerkleProofWithMixin(2)
	require.NoError(t, err)
	require.Len(t, proof, int(merkle.MaxTreeDepth)+1)
	// Unsupported depth
	_, err = merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		merkle.MaxTreeDepth+1,
	)
	require.ErrorIs(t, err, merkle.ErrExceededDepth)
}

func TestMerkleTree_IsValidMerkleBranch(t *testing.T) {
	treeDepth := uint8(32)
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("B")),
		byteslib.ToBytes32([]byte("C")),
		byteslib.ToBytes32([]byte("D")),
		byteslib.ToBytes32([]byte("E")),
		byteslib.ToBytes32([]byte("F")),
		byteslib.ToBytes32([]byte("G")),
		byteslib.ToBytes32([]byte("H")),
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(t, err)

	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(t, err)
	require.Len(
		t,
		proof,
		int(treeDepth)+1,
	)

	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	require.True(t, merkle.VerifyProof(
		root, items[0], 0, proof,
	), "First Merkle proof did not verify")

	proof, err = m.MerkleProofWithMixin(3)
	require.NoError(t, err)
	require.True(
		t,
		merkle.VerifyProof(
			root,
			items[3],
			3,
			proof,
		),
	)
	require.False(
		t,
		merkle.IsValidMerkleBranch(
			byteslib.ToBytes32([]byte("buzz")),
			proof,
			treeDepth,
			3,
			root,
		),
	)
}

func TestMerkleTree_VerifyProof(t *testing.T) {
	treeDepth := uint8(32)
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("B")),
		byteslib.ToBytes32([]byte("C")),
		byteslib.ToBytes32([]byte("D")),
		byteslib.ToBytes32([]byte("E")),
		byteslib.ToBytes32([]byte("F")),
		byteslib.ToBytes32([]byte("G")),
		byteslib.ToBytes32([]byte("H")),
	}

	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(t, err)
	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(t, err)
	require.Len(
		t,
		proof,
		int(treeDepth)+1,
	)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	if ok := merkle.VerifyProof(root, items[0], 0, proof); !ok {
		t.Error("First Merkle proof did not verify")
	}
	proof, err = m.MerkleProofWithMixin(3)
	require.NoError(t, err)
	require.True(t, merkle.VerifyProof(root, items[3], 3, proof))
	require.False(
		t,
		merkle.VerifyProof(
			root,
			byteslib.ToBytes32([]byte("buzz")),
			3,
			proof,
		),
	)
}

func TestMerkleTree_NegativeIndexes(t *testing.T) {
	treeDepth := uint8(32)
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("B")),
		byteslib.ToBytes32([]byte("C")),
		byteslib.ToBytes32([]byte("D")),
		byteslib.ToBytes32([]byte("E")),
		byteslib.ToBytes32([]byte("F")),
		byteslib.ToBytes32([]byte("G")),
		byteslib.ToBytes32([]byte("H")),
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(t, err)
	err = m.Insert(byteslib.ToBytes32([]byte{'J'}), -1)
	require.ErrorIs(t, err, merkle.ErrNegativeIndex)
}

func TestMerkleTree_VerifyProof_TrieUpdated(t *testing.T) {
	treeDepth := uint8(32)
	items := [][32]byte{
		{1},
		{2},
		{3},
		{4},
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth+1,
	)
	require.NoError(t, err)
	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(t, err)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	require.True(
		t,
		merkle.VerifyProof(
			root,
			items[0],
			0,
			proof,
		),
	)

	// Now we update the merkle.
	require.NoError(t, m.Insert(byteslib.ToBytes32([]byte{5}), 3))
	proof, err = m.MerkleProofWithMixin(3)
	require.NoError(t, err)
	root, err = m.HashTreeRoot()
	require.NoError(t, err)
	require.True(t, merkle.VerifyProof(
		root, [32]byte{5}, 3, proof,
	), "Second Merkle proof did not verify")
	require.False(t, merkle.VerifyProof(
		root, [32]byte{4}, 3, proof,
	), "Old item should not verify")

	// Now we update the tree at an index larger than the number of items.
	require.NoError(t, m.Insert(byteslib.ToBytes32([]byte{6}), 15))
}

func BenchmarkNewTreeFromLeavesWithDepth(b *testing.B) {
	treeDepth := uint8(32)
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("BB")),
		byteslib.ToBytes32([]byte("CCC")),
		byteslib.ToBytes32([]byte("DDDD")),
		byteslib.ToBytes32([]byte("EEEEE")),
		byteslib.ToBytes32([]byte("FFFFFF")),
		byteslib.ToBytes32([]byte("GGGGGGG")),
	}
	for i := 0; i < b.N; i++ {
		_, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
			items,
			treeDepth,
		)
		require.NoError(b, err, "Could not generate Merkle tree from items")
	}
}

func BenchmarkInsertTrie_Optimized(b *testing.B) {
	treeDepth := uint8(32)
	b.StopTimer()
	numDeposits := 16000
	items := make([][32]byte, numDeposits)
	for i := range numDeposits {
		items[i] = byteslib.ToBytes32([]byte(strconv.Itoa(i)))
	}
	tr, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(b, err)

	someItem := byteslib.ToBytes32([]byte("hello-world"))
	b.StartTimer()
	for i := range b.N {
		require.NoError(b, tr.Insert(someItem, i%numDeposits))
	}
}

func BenchmarkGenerateProof(b *testing.B) {
	treeDepth := uint8(32)
	b.StopTimer()
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("BB")),
		byteslib.ToBytes32([]byte("CCC")),
		byteslib.ToBytes32([]byte("DDDD")),
		byteslib.ToBytes32([]byte("EEEEE")),
		byteslib.ToBytes32([]byte("FFFFFF")),
		byteslib.ToBytes32([]byte("GGGGGGG")),
	}
	goodTree, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = goodTree.MerkleProofWithMixin(3)
		require.NoError(b, err)
	}
}

func BenchmarkIsValidMerkleBranch(b *testing.B) {
	treeDepth := uint8(4)
	b.StopTimer()
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("BB")),
		byteslib.ToBytes32([]byte("CCC")),
		byteslib.ToBytes32([]byte("DDDD")),
		byteslib.ToBytes32([]byte("EEEEE")),
		byteslib.ToBytes32([]byte("FFFFFF")),
		byteslib.ToBytes32([]byte("GGGGGGG")),
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)

	require.NoError(b, err)
	proof, err := m.MerkleProofWithMixin(2)
	require.NoError(b, err)

	root, err := m.HashTreeRoot()
	require.NoError(b, err)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ok := merkle.IsValidMerkleBranch(
			items[2], proof, treeDepth+1, 2, root,
		)
		require.True(b, ok, "Merkle proof did not verify")
	}
}
