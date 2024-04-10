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

package tree_test

import (
	"strconv"
	"testing"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	tree "github.com/berachain/beacon-kit/mod/tree"
	"github.com/stretchr/testify/require"
)

const (
	TreeDepth uint64 = 32
)

func TestNewFromItems_NoItemsProvided(t *testing.T) {
	_, err := tree.NewFromItems(nil, TreeDepth)
	require.ErrorContains(t, err, "no items provided to generate Merkle tree")
}

func TestNewFromItems_DepthSupport(t *testing.T) {
	items := [][]byte{
		[]byte("A"),
		[]byte("BB"),
		[]byte("CCC"),
		[]byte("DDDD"),
		[]byte("EEEEE"),
		[]byte("FFFFFF"),
		[]byte("GGGGGGG"),
	}
	// Supported depth
	m1, err := tree.NewFromItems(items, tree.MaxDepth)
	require.NoError(t, err)
	proof, err := m1.MerkleProof(2)
	require.NoError(t, err)
	require.Len(t, proof, int(tree.MaxDepth)+1)
	// Unsupported depth
	_, err = tree.NewFromItems(items, tree.MaxDepth+1)
	require.ErrorIs(t, err, tree.ErrExceededDepth)
}

func TestMerkleTree_VerifyMerkleProofWithDepth(t *testing.T) {
	items := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
		[]byte("E"),
		[]byte("F"),
		[]byte("G"),
		[]byte("H"),
	}
	m, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(t, err)
	proof, err := m.MerkleProof(0)
	require.NoError(t, err)
	require.Len(
		t,
		proof,
		int(TreeDepth)+1,
	)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	if ok := tree.VerifyMerkleProofWithDepth(
		root[:], items[0], 0, proof, TreeDepth,
	); !ok {
		t.Error("First Merkle proof did not verify")
	}
	proof, err = m.MerkleProof(3)
	require.NoError(t, err)
	require.True(
		t,
		tree.VerifyMerkleProofWithDepth(
			root[:],
			items[3],
			3,
			proof,
			TreeDepth,
		),
	)
	require.False(
		t,
		tree.VerifyMerkleProofWithDepth(
			root[:],
			[]byte("buzz"),
			3,
			proof,
			TreeDepth,
		),
	)
}

func TestMerkleTree_VerifyMerkleProof(t *testing.T) {
	items := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
		[]byte("E"),
		[]byte("F"),
		[]byte("G"),
		[]byte("H"),
	}

	m, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(t, err)
	proof, err := m.MerkleProof(0)
	require.NoError(t, err)
	require.Len(
		t,
		proof,
		int(TreeDepth)+1,
	)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	if ok := tree.VerifyMerkleProof(root[:], items[0], 0, proof); !ok {
		t.Error("First Merkle proof did not verify")
	}
	proof, err = m.MerkleProof(3)
	require.NoError(t, err)
	require.True(t, tree.VerifyMerkleProof(root[:], items[3], 3, proof))
	require.False(
		t,
		tree.VerifyMerkleProof(root[:], []byte("buzz"), 3, proof),
	)
}

func TestMerkleTree_NegativeIndexes(t *testing.T) {
	items := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
		[]byte("E"),
		[]byte("F"),
		[]byte("G"),
		[]byte("H"),
	}
	m, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(t, err)
	require.ErrorContains(
		t,
		m.Insert([]byte{'J'}, -1),
		"negative index provided",
	)
}

func TestMerkleTree_VerifyMerkleProof_TrieUpdated(t *testing.T) {
	items := [][]byte{
		{1},
		{2},
		{3},
		{4},
	}
	treeDepth := TreeDepth + 1
	m, err := tree.NewFromItems(items, treeDepth)
	require.NoError(t, err)
	proof, err := m.MerkleProof(0)
	require.NoError(t, err)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	require.True(
		t,
		tree.VerifyMerkleProofWithDepth(root[:], items[0], 0, proof, treeDepth),
	)

	// Now we update the tree.
	require.NoError(t, m.Insert([]byte{5}, 3))
	proof, err = m.MerkleProof(3)
	require.NoError(t, err)
	root, err = m.HashTreeRoot()
	require.NoError(t, err)
	require.True(t, tree.VerifyMerkleProofWithDepth(
		root[:], []byte{5}, 3, proof, treeDepth,
	), "Second Merkle proof did not verify")
	require.False(t, tree.VerifyMerkleProofWithDepth(
		root[:], []byte{4}, 3, proof, treeDepth,
	), "Old item should not verify")

	// Now we update the tree at an index larger than the number of items.
	require.NoError(t, m.Insert([]byte{6}, 15))
}

func TestCopy_OK(t *testing.T) {
	items := [][]byte{
		{1},
		{2},
		{3},
		{4},
	}
	source, err := tree.NewFromItems(
		items,
		TreeDepth+1,
	)
	require.NoError(t, err)
	copiedTree := source.Copy()

	if copiedTree == source {
		t.Errorf("Original tree returned.")
	}
	a, err := copiedTree.HashTreeRoot()
	require.NoError(t, err)
	b, err := source.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, a, b)
}

func BenchmarkNewFromItems(b *testing.B) {
	items := [][]byte{
		[]byte("A"),
		[]byte("BB"),
		[]byte("CCC"),
		[]byte("DDDD"),
		[]byte("EEEEE"),
		[]byte("FFFFFF"),
		[]byte("GGGGGGG"),
	}
	for i := 0; i < b.N; i++ {
		_, err := tree.NewFromItems(
			items,
			TreeDepth,
		)
		require.NoError(b, err, "Could not generate Merkle tree from items")
	}
}

func BenchmarkInsertTrie_Optimized(b *testing.B) {
	b.StopTimer()
	numDeposits := 16000
	items := make([][]byte, numDeposits)
	for i := 0; i < numDeposits; i++ {
		someRoot := byteslib.ToBytes32([]byte(strconv.Itoa(i)))
		items[i] = someRoot[:]
	}
	tr, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(b, err)

	someItem := byteslib.ToBytes32([]byte("hello-world"))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		require.NoError(b, tr.Insert(someItem[:], i%numDeposits))
	}
}

func BenchmarkGenerateProof(b *testing.B) {
	b.StopTimer()
	items := [][]byte{
		[]byte("A"),
		[]byte("BB"),
		[]byte("CCC"),
		[]byte("DDDD"),
		[]byte("EEEEE"),
		[]byte("FFFFFF"),
		[]byte("GGGGGGG"),
	}
	goodTree, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(b, err)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = goodTree.MerkleProof(3)
		require.NoError(b, err)
	}
}

func BenchmarkVerifyMerkleProofWithDepth(b *testing.B) {
	b.StopTimer()
	items := [][]byte{
		[]byte("A"),
		[]byte("BB"),
		[]byte("CCC"),
		[]byte("DDDD"),
		[]byte("EEEEE"),
		[]byte("FFFFFF"),
		[]byte("GGGGGGG"),
	}
	m, err := tree.NewFromItems(
		items,
		TreeDepth,
	)
	require.NoError(b, err)
	proof, err := m.MerkleProof(2)
	require.NoError(b, err)

	root, err := m.HashTreeRoot()
	require.NoError(b, err)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if ok := tree.VerifyMerkleProofWithDepth(
			root[:], items[2], 2, proof, TreeDepth,
		); !ok {
			b.Error("Merkle proof did not verify")
		}
	}
}
