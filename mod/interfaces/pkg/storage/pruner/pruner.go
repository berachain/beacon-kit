package pruner

import "context"

// Prunable is an interface representing a store that can be pruned.
type Prunable interface {
	// Prune prunes the store from [start, end).
	Prune(start, end uint64) error
}

// Pruner is an interface for pruning a prunable type.
type Pruner[PrunableT Prunable] interface {
	Name() string
	Start(ctx context.Context)
}
