package tree

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives"
)

func TestMerkleTreeAndRoot(t *testing.T) {
	// Create a sample list of leaves
	leaves := []primitives.Bytes32{
		{0x6c, 0x65, 0x61, 0x66, 0x31}, // "leaf1"
		{0x6c, 0x65, 0x61, 0x66, 0x32}, // "leaf2"
		{0x6c, 0x65, 0x61, 0x66, 0x33}, // "leaf3"
		{0x6c, 0x65, 0x61, 0x66, 0x34}, // "leaf4"
	}

	// Generate the Merkle tree from the leaves
	tree := MerkleTree(leaves)

	// Calculate the Merkle root using the tree
	merkleRoot := tree[1] // Index 1 is the root of the tree

	_ = merkleRoot
	t.Log("Merkle root:", merkleRoot)
}
