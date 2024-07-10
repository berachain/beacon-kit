package types

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Eth1Data[T any] interface {
	// New creates a new eth1 data.
	New(
		depositRoot common.Root,
		depositCount math.U64,
		blockHash gethprimitives.ExecutionHash,
	) T
	// GetDepositCount returns the deposit count of the eth1 data.
	GetDepositCount() math.U64
}
