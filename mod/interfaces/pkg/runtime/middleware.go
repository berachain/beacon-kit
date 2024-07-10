package runtime

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/cosmos/gogoproto/proto"
)

// Middleware is a middleware between ABCI and the validator logic.
type Middleware interface {
	// InitGenesis initializes the validators with the given genesis state.
	InitGenesis(
		ctx context.Context, bz []byte,
	) (transition.ValidatorUpdates, error)
	// PrepareProposal prepares a proposal for the given slot.
	PrepareProposal(ctx context.Context, slot math.Slot) ([]byte, []byte, error)
	// ProcessProposal processes a proposal.
	ProcessProposal(ctx context.Context, req proto.Message) (proto.Message, error)
	// PreBlock is called before processing a block.
	PreBlock(ctx context.Context, req proto.Message) error
	// EndBlock is called after processing a block.
	// It returns the validator updates from the beacon state.
	EndBlock(ctx context.Context) (transition.ValidatorUpdates, error)
}
