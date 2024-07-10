package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Fork[T any] interface {
	// New creates a new fork.
	New(
		previousVersion common.Version,
		currentVersion common.Version,
		epoch math.Epoch,
	) T
}
