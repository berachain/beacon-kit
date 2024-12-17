package merkle

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrFinalizedNodeCannotPushLeaf may occur when attempting to push a leaf to a finalized node.
	// When a node is finalized, it cannot be modified or changed.
	ErrFinalizedNodeCannotPushLeaf = errors.New("can't push a leaf to a finalized node")

	// ErrLeafNodeCannotPushLeaf may occur when attempting to push a leaf to a leaf node.
	ErrLeafNodeCannotPushLeaf = errors.New("can't push a leaf to a leaf node")

	// ErrZeroLevel occurs when the value of level is 0.
	ErrZeroLevel = errors.New("level should be greater than 0")

	// ErrZeroDepth occurs when the value of depth is 0.
	ErrZeroDepth = errors.New("depth should be greater than 0")
)
