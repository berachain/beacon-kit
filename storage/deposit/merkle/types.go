package merkle

// TreeNode is the interface for a Merkle tree.
type TreeNode interface {
	// GetRoot returns the root of the Merkle tree.
	GetRoot() [32]byte
	// IsFull returns whether there is space left for deposits.
	IsFull() bool
	// Finalize marks deposits of the Merkle tree as finalized.
	Finalize(depositsToFinalize uint64, depth uint64) (TreeNode, error)
	// GetFinalized returns the number of deposits and a list of hashes of all the finalized nodes.
	GetFinalized(result [][32]byte) (uint64, [][32]byte)
	// PushLeaf adds a new leaf node at the next available Zero node.
	PushLeaf(leaf [32]byte, depth uint64) (TreeNode, error)

	// Right represents the right child of a node.
	Right() TreeNode
	// Left represents the left child of a node.
	Left() TreeNode
}
