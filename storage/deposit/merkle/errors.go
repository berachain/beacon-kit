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

var (
	// ErrInvalidSnapshotRoot occurs when the snapshot root does not match the calculated root.
	ErrInvalidSnapshotRoot = errors.New("snapshot root is invalid")

	// ErrInvalidDepositCount occurs when the value for mix in length is 0.
	ErrInvalidDepositCount = errors.New("deposit count should be greater than 0")

	// ErrInvalidIndex occurs when the index is less than the number of finalized deposits.
	ErrInvalidIndex = errors.New("index should be greater than finalizedDeposits - 1")

	// ErrTooManyDeposits occurs when the number of deposits exceeds the capacity of the tree.
	ErrTooManyDeposits = errors.New("number of deposits should not be greater than the capacity of the tree")
)
