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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
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
	items := make([][32]byte, 0)
	for _, v := range []string{
		"A",
		"BB",
		"CCC",
		"DDDD",
		"EEEEE",
		"FFFFF",
		"GGGGGG",
	} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		items = append(items, item)
	}

	// Supported depth
	m1, err := merkle.NewTreeFromLeavesWithDepth(
		items,
		merkle.MaxTreeDepth,
	)
	require.NoError(t, err)
	proof, err := m1.MerkleProofWithMixin(2)
	require.NoError(t, err)
	require.Len(t, proof, int(merkle.MaxTreeDepth)+1)
	// Unsupported depth
	_, err = merkle.NewTreeFromLeavesWithDepth(
		items,
		merkle.MaxTreeDepth+1,
	)
	require.ErrorIs(t, err, merkle.ErrExceededDepth)
}

func TestMerkleTree_IsValidMerkleBranch(t *testing.T) {
	treeDepth := uint8(32)
	items := make([][32]byte, 0)
	for _, v := range []string{"A", "B", "C", "D", "E", "F", "G", "H"} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		items = append(items, item)
	}
	m, err := merkle.NewTreeFromLeavesWithDepth(
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

	root := m.HashTreeRoot()
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

	it := byteslib.ExtendToSize([]byte("buzz"), byteslib.B32Size)
	item, err := byteslib.ToBytes32(it)
	require.NoError(t, err)
	require.False(
		t,
		merkle.IsValidMerkleBranch(
			common.Root(item),
			proof,
			treeDepth,
			3,
			root,
		),
	)
}

func TestMerkleTree_VerifyProof(t *testing.T) {
	treeDepth := uint8(32)
	items := make([][32]byte, 0)
	for _, v := range []string{"A", "B", "C", "D", "E", "F", "G", "H"} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		items = append(items, item)
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
	root := m.HashTreeRoot()
	if ok := merkle.VerifyProof(root, items[0], 0, proof); !ok {
		t.Error("First Merkle proof did not verify")
	}
	proof, err = m.MerkleProofWithMixin(3)
	require.NoError(t, err)
	require.True(t, merkle.VerifyProof(root, items[3], 3, proof))

	it := byteslib.ExtendToSize([]byte("buzz"), byteslib.B32Size)
	item, err := byteslib.ToBytes32(it)
	require.NoError(t, err)
	require.False(
		t,
		merkle.VerifyProof(
			root,
			common.Root(item),
			3,
			proof,
		),
	)
}

func TestMerkleTree_NegativeIndexes(t *testing.T) {
	treeDepth := uint8(32)
	items := make([][32]byte, 0)
	for _, v := range []string{"A", "B", "C", "D", "E", "F", "G", "H"} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		items = append(items, item)
	}
	m, err := merkle.NewTreeFromLeavesWithDepth(
		items,
		treeDepth,
	)
	require.NoError(t, err)

	it := byteslib.ExtendToSize([]byte("J"), byteslib.B32Size)
	extraItem, err := byteslib.ToBytes32(it)
	require.NoError(t, err)
	err = m.Insert(extraItem, -1)
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
	m, err := merkle.NewTreeFromLeavesWithDepth(
		items,
		treeDepth+1,
	)
	require.NoError(t, err)
	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(t, err)
	root := m.HashTreeRoot()
	require.True(
		t,
		merkle.VerifyProof(
			root,
			items[0],
			0,
			proof,
		),
	)

	// Now we update the merkle
	it := byteslib.ExtendToSize([]byte{5}, byteslib.B32Size)
	item, err := byteslib.ToBytes32(it)

	require.NoError(t, err)
	require.NoError(t, m.Insert(item, 3))
	proof, err = m.MerkleProofWithMixin(3)
	require.NoError(t, err)
	root = m.HashTreeRoot()
	require.True(t, merkle.VerifyProof(
		root, [32]byte{5}, 3, proof,
	), "Second Merkle proof did not verify")
	require.False(t, merkle.VerifyProof(
		root, [32]byte{4}, 3, proof,
	), "Old item should not verify")

	// Now we update the tree at an index larger than the number of items.
	it = byteslib.ExtendToSize([]byte{6}, byteslib.B32Size)
	item, err = byteslib.ToBytes32(it)
	require.NoError(t, err)
	require.NoError(t, m.Insert(item, 15))
}

func BenchmarkNewTreeFromLeavesWithDepth(b *testing.B) {
	treeDepth := uint8(32)
	items := make([][32]byte, 0)
	for _, v := range []string{
		"A",
		"BB",
		"CCC",
		"DDDD",
		"EEEEE",
		"FFFFFF",
		"GGGGGGG",
	} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(b, err)
		items = append(items, item)
	}
	for i := 0; i < b.N; i++ {
		_, err := merkle.NewTreeFromLeavesWithDepth(
			items,
			treeDepth,
		)
		require.NoError(b, err, "Could not generate Merkle tree from items")
	}
}

func BenchmarkInsertTrie_Optimized(b *testing.B) {
	treeDepth := uint8(32)
	b.StopTimer()

	var (
		numDeposits = 16000
		items       = make([][32]byte, numDeposits)
		err         error
	)

	for i := range numDeposits {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		items[i], err = byteslib.ToBytes32(it)
		require.NoError(b, err)
	}
	tr, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)
	require.NoError(b, err)

	it := byteslib.ExtendToSize([]byte("hello-world"), byteslib.B32Size)
	someItem, err := byteslib.ToBytes32(it)
	require.NoError(b, err)

	b.StartTimer()
	for i := range b.N {
		require.NoError(b, tr.Insert(someItem, i%numDeposits))
	}
}

func BenchmarkGenerateProof(b *testing.B) {
	treeDepth := uint8(32)
	b.StopTimer()

	items := make([][32]byte, 0)
	for _, v := range []string{
		"A",
		"BB",
		"CCC",
		"DDDD",
		"EEEEE",
		"FFFFFF",
		"GGGGGGG",
	} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(b, err)
		items = append(items, item)
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

	items := make([][32]byte, 0)
	for _, v := range []string{
		"A",
		"BB",
		"CCC",
		"DDDD",
		"EEEEE",
		"FFFFFF",
		"GGGGGGG",
	} {
		it := byteslib.ExtendToSize([]byte(v), byteslib.B32Size)
		item, err := byteslib.ToBytes32(it)
		require.NoError(b, err)
		items = append(items, item)
	}

	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte](
		items,
		treeDepth,
	)

	require.NoError(b, err)
	proof, err := m.MerkleProofWithMixin(2)
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ok := merkle.IsValidMerkleBranch(
			items[2], proof, treeDepth+1, 2, m.HashTreeRoot(),
		)
		require.True(b, ok, "Merkle proof did not verify")
	}
}
