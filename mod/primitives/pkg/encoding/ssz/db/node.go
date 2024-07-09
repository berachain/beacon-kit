package db

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
)

// Node represents a node in the SSZ merkle tree.
type Node[RootT ~[32]byte] struct {
	// SSZType is the SSZ type of the node.
	schema.SSZType
	// gIndex is the generalized index of the node in the Merkle tree.
	gIndex uint64
	// offset is the byte offset within the 32-byte chunk where the node's data
	// begins.
	offset uint8
}

// NewTreeNode locates a node in the SSZ merkle tree by its path and a root
// schema node to begin traversal from with gindex 1.
func NewTreeNode[RootT ~[32]byte](
	root schema.SSZType, path merkle.ObjectPath[RootT],
) (Node[RootT], error) {
	found, gindex, offset, err := path.GetGeneralizedIndex(root)
	return Node[RootT]{SSZType: found, gIndex: gindex, offset: offset}, err
}

// GeIndex returns the generalized index of the node in the Merkle tree.
func (n Node[RootT]) GIndex() merkle.GeneralizedIndex[RootT] {
	return n.gIndex
}

// Offset returns the byte offset within the 32-byte chunk where the node's data
// begins.
func (n Node[_]) Offset() uint8 {
	return n.offset
}
