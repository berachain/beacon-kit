package consensus

import (
	"context"

	"github.com/cosmos/gogoproto/proto"
)

// Consensus is the interface that must be implemented by all consensus
// engines.
type Consensus[ValidatorUpdateT any] interface {
	InitGenesis(
		ctx context.Context, genesisBz []byte,
	) ([]ValidatorUpdateT, error)
	Prepare(ctx context.Context, msg proto.Message) (proto.Message, error)
	Process(ctx context.Context, msg proto.Message) (proto.Message, error)
	PreBlock(ctx context.Context, msg proto.Message) error
	EndBlock(ctx context.Context) ([]ValidatorUpdateT, error)
}
