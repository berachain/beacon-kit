package cometbft

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// TODO: We must rid this of comet bft types.
type Middleware interface {
	InitGenesis(
		ctx context.Context,
		bz []byte,
	) (transition.ValidatorUpdates, error)

	PrepareProposal(
		context.Context,
		math.Slot,
	) ([]byte, []byte, error)

	ProcessProposal(
		ctx context.Context,
		req *cmtabci.ProcessProposalRequest,
	) error

	PreBlock(
		_ context.Context, req *cmtabci.FinalizeBlockRequest,
	) error

	EndBlock(
		ctx context.Context,
	) (transition.ValidatorUpdates, error)
}
