package types

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type ProofProvider interface {
	GetBlockProposerProof(ctx context.Context, timestamp uint64) (math.U64, bytes.B48, [][32]byte, error)
	GetExecutionNumberProof(ctx context.Context, timestamp uint64) ([][32]byte, error)
}
