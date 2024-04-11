package tree

import (
	"crypto/sha256"
)

const MaxTreeDepth = 32

var zeroHashes [MaxTreeDepth + 1][32]byte

func init() {
	var zero [32]byte
	for i := 0; i <= MaxTreeDepth; i++ {
		zeroHashes[i] = sha256.Sum256(zero[:])
		zero = zeroHashes[i]
	}
}

type MerkleTree struct {
	Hash  [32]byte
	Left  *MerkleTree
	Right *MerkleTree
	Depth int
}

type MerkleTreeError string

func (e MerkleTreeError) Error() string {
	return string(e)
}

const (
	LeafReached                   MerkleTreeError = "leaf reached"
	MerkleTreeFull                MerkleTreeError = "merkle tree full"
	InvalidMerkleTree             MerkleTreeError = "invalid merkle tree"
	DepthTooSmall                 MerkleTreeError = "depth too small"
	ZeroNodeFinalized             MerkleTreeError = "zero node finalized"
	FinalizedNodePushed           MerkleTreeError = "finalized node pushed"
	ProofEncounteredFinalizedNode MerkleTreeError = "proof encountered finalized node"
	PleaseNotifyTheDevs           MerkleTreeError = "please notify the devs"
)

func NewMerkleTree(leaves [][32]byte, depth int) *MerkleTree {
	if len(leaves) == 0 {
		return &MerkleTree{Hash: zeroHashes[depth], Depth: depth}
	}

	if depth == 0 {
		if len(leaves) != 1 {
			panic("depth 0 tree with multiple leaves")
		}
		return &MerkleTree{Hash: leaves[0]}
	}

	subtreeCapacity := 1 << (depth - 1)
	var leftLeaves, rightLeaves [][32]byte

	if len(leaves) <= subtreeCapacity {
		leftLeaves = leaves
		rightLeaves = nil
	} else {
		leftLeaves = leaves[:subtreeCapacity]
		rightLeaves = leaves[subtreeCapacity:]
	}

	leftSubtree := NewMerkleTree(leftLeaves, depth-1)
	rightSubtree := NewMerkleTree(rightLeaves, depth-1)

	return &MerkleTree{
		Hash:  hash32Concat(leftSubtree.Hash[:], rightSubtree.Hash[:]),
		Left:  leftSubtree,
		Right: rightSubtree,
		Depth: depth,
	}
}

func hash32Concat(left, right []byte) [32]byte {
	data := append(left, right...)
	return sha256.Sum256(data)
}

func (t *MerkleTree) PushLeaf(elem [32]byte, depth int) error {
	if depth == 0 {
		return DepthTooSmall
	}

	switch {
	case t.Left == nil && t.Right == nil:
		if t.Depth == 0 {
			return LeafReached
		}
		*t = *NewMerkleTree([][32]byte{elem}, depth)
	case t.Left != nil && t.Right != nil:
		if t.Left.Depth == 0 && t.Right.Depth == 0 {
			return MerkleTreeFull
		}
		err := t.Right.PushLeaf(elem, depth-1)
		if err != nil {
			return err
		}
		t.Hash = hash32Concat(t.Left.Hash[:], t.Right.Hash[:])
	default:
		return InvalidMerkleTree
	}

	return nil
}

func (t *MerkleTree) HashValue() [32]byte {
	return t.Hash
}

func (t *MerkleTree) LeftAndRightBranches() (*MerkleTree, *MerkleTree) {
	if t.Left == nil || t.Right == nil {
		return nil, nil
	}
	return t.Left, t.Right
}

func (t *MerkleTree) IsLeaf() bool {
	return t.Depth == 0
}

func (t *MerkleTree) AppendFinalizedHashes(result *[][32]byte) {
	if t.IsLeaf() || t == nil {
		return
	}
	if t.Left != nil {
		t.Left.AppendFinalizedHashes(result)
	}
	if t.Right != nil {
		t.Right.AppendFinalizedHashes(result)
	}
	*result = append(*result, t.Hash)
}

func (t *MerkleTree) GetFinalizedHashes() [][32]byte {
	var result [][32]byte
	t.AppendFinalizedHashes(&result)
	return result
}

// Return the leaf at `index` and a Merkle proof of its inclusion.
//
// The Merkle proof is in "bottom-up" order, starting with a leaf node
// and moving up the tree. Its length will be exactly equal to `depth`.
func (t *MerkleTree) GenerateProof(index uint, depth uint) ([32]byte, [][32]byte, error) {
	var proof [][32]byte
	currentNode := t
	current_depth := depth
	for current_depth > 0 {
		ith_bit := (index >> (current_depth - 1)) & 0x01
		left, right := currentNode.LeftAndRightBranches()
		// Go right, include the left branch in the proof.
		if ith_bit == 1 {
			proof = append(proof, left.HashValue())
			currentNode = right
		} else {
			proof = append(proof, right.HashValue())
			currentNode = left
		}
		current_depth--
	}

	if len(proof) != int(depth) {
		panic("proof length does not match depth")
	}
	if !currentNode.IsLeaf() {
		panic("current node is not a leaf")
	}

	// Put proof in bottom-up order.
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
	}

	return currentNode.HashValue(), proof, nil
}

// Verify a proof that `leaf` exists at `index` in a Merkle tree rooted at `root`.
//
// The `branch` argument is the main component of the proof: it should be a slice of internal
// node hashes such that the root can be reconstructed (in bottom-up order).
func VerifyMerkleProof2(leaf [32]byte, branch [][32]byte, depth uint, index uint, root [32]byte) bool {
	if len(branch) == int(depth) {
		return MerkleRootFromBranch(leaf, branch, depth, index) == root
	}
	return false
}

// Compute a root hash from a leaf and a Merkle proof.
func MerkleRootFromBranch(leaf [32]byte, branch [][32]byte, depth uint, index uint) [32]byte {
	if len(branch) != int(depth) {
		panic("proof length should equal depth")
	}

	// merkleRoot := make([]byte, 32)
	// copy(merkleRoot, leaf[:])
	var merkleRoot [32]byte
	merkleRoot = leaf
	for i, node := range branch[:depth] {
		ithBit := (index >> i) & 0x01
		if ithBit == 1 {
			merkleRoot = hash32Concat(node[:], merkleRoot[:])
		} else {
			input := append(merkleRoot[:], node[:]...)
			merkleRoot = sha256.Sum256(input)
		}
	}

	return merkleRoot
}
