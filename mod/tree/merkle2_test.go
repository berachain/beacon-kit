package tree_test

import (
	"bytes"
	"crypto/sha256"
	"testing"

	treezz "github.com/berachain/beacon-kit/mod/tree"
)

func TestNewMerkleTree(t *testing.T) {
	leafHash := sha256.Sum256([]byte("leaf"))
	emptyTree := treezz.NewMerkleTree(nil, 0)
	if emptyTree == nil {
		t.Error("Expected non-nil tree for zero depth and no leaves")
	}
	if emptyTree.Depth != 0 {
		t.Errorf("Expected depth 0, got %d", emptyTree.Depth)
	}

	singleLeafTree := treezz.NewMerkleTree([][32]byte{leafHash}, 0)
	if !bytes.Equal(singleLeafTree.Hash[:], leafHash[:]) {
		t.Error("Hash of single leaf tree does not match expected hash")
	}
}

func TestPushLeaf(t *testing.T) {
	leafHash := sha256.Sum256([]byte("leaf"))
	tree := treezz.NewMerkleTree(nil, 1)

	err := tree.PushLeaf(leafHash, 1)
	if err != nil {
		t.Errorf("Failed to push leaf: %s", err)
	}

	if !bytes.Equal(tree.Left.Hash[:], leafHash[:]) {
		t.Error("Leaf hash does not match pushed leaf hash")
	}
}

func TestGenerateProof(t *testing.T) {
	leafHash := sha256.Sum256([]byte("leaf"))
	tree := treezz.NewMerkleTree([][32]byte{leafHash}, 1)
	_, proof, err := tree.GenerateProof(0, 1)
	if err != nil {
		t.Errorf("Failed to generate proof: %s", err)
	}

	rootHash := tree.HashValue()
	if !treezz.VerifyMerkleProof2(leafHash, proof, 1, 0, rootHash) {
		t.Error("Failed to verify Merkle proof")
	}

}

func TestHashValue(t *testing.T) {
	leafHash := sha256.Sum256([]byte("leaf"))
	tree := treezz.NewMerkleTree([][32]byte{leafHash}, 0)
	if tree.HashValue() != leafHash {
		t.Error("Hash value does not match expected leaf hash")
	}
}
