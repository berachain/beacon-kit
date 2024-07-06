package tree

// Node represents a node in the tree backing of an SSZ object.
type Node struct {
	// left is the left child node.
	left *Node
	// right is the right child node.
	right *Node
	// value holds the node's serialized data and/or hash.
	value []byte
}

// Left returns the left child node.
func (n *Node) Left() *Node {
	return n.left
}

// Right returns the right child node.
func (n *Node) Right() *Node {
	return n.right
}

// Value returns the node's data.
func (n *Node) Value() []byte {
	return n.value
}
