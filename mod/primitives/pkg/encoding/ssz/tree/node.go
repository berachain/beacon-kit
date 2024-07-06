package tree

// Node represents a node in the tree backing of an SSZ object.
type Node struct {
	// left is the left child node.
	left *Node

	// right is the right child node.
	right *Node

	// value holds the node's data.
	value []byte
}
