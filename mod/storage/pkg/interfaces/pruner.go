package interfaces

import "context"

type Pruner interface {
	// TODO - Add methods
	Prune(ctx context.Context)
}

type Prunable interface {
	DeleteRange(from uint64, to uint64) error
}
