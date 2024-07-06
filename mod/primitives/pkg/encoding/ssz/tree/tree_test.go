package tree_test

import (
	"testing"

	bytes "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/tree"
	"github.com/stretchr/testify/require"
)

func TestNewTreeFromLeaves(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tr, err := tree.NewTreeFromLeaves[[32]byte](nil, 0)
		require.NoError(t, err)
		require.NotNil(t, tr)
	})

	t.Run("single leaf", func(t *testing.T) {
		leaves := [][32]byte{{1, 2, 3}}
		tr, err := tree.NewTreeFromLeaves[[32]byte](leaves, 1)
		require.NoError(t, err)
		require.NotNil(t, tr)
	})

	t.Run("multiple leaves", func(t *testing.T) {
		leaves := [][32]byte{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}
		tr, err := tree.NewTreeFromLeaves(leaves, 4)
		require.NoError(t, err)
		require.NotNil(t, tr)
	})

	t.Run("exceeding max leaves", func(t *testing.T) {
		leaves := [][32]byte{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
			{10, 11, 12},
		}
		_, err := tree.NewTreeFromLeaves(leaves, 1)
		require.Error(t, err)
	})
}

func TestNew(t *testing.T) {
	leaves := [][32]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	result := tree.New(leaves, sha256.Hash)
	require.Len(t, result, 8)
	require.NotEqual(t, [32]byte{}, result[1]) // Root should not be empty
}

func TestCompareTreeToBytes96HashTreeRoot(t *testing.T) {
	t.Run("compare tree to Bytes96 hash tree root", func(t *testing.T) {
		// Create a sample Bytes96 (3 x 32 bytes)
		bytes96 := [96]byte{}
		for i := 0; i < 96; i++ {
			bytes96[i] = byte(i)
		}

		// Split Bytes96 into 3 chunks of 32 bytes each
		var leaves [][32]byte
		for i := 0; i < 3; i++ {
			var leaf [32]byte
			copy(leaf[:], bytes96[i*32:(i+1)*32])
			leaves = append(leaves, leaf)
		}

		// Create a tree using our NewTreeFromLeaves function
		tr, err := tree.NewTreeFromLeaves(leaves, 4)
		require.NoError(t, err)
		require.NotNil(t, tr)

		// Get the root of our tree
		treeRoot := tr.Root()

		// Calculate the hash tree root of Bytes96 using ssz.HashTreeRoot
		sszRoot, err := bytes.B96(bytes96).HashTreeRoot()
		require.NoError(t, err)

		// Compare the roots
		require.Equal(t, sszRoot, treeRoot, "Tree root should match SSZ hash tree root")
	})
}
