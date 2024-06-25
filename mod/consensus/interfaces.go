package consensus

import (
	"context"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/server/v2/cometbft/handlers"
	"github.com/cosmos/gogoproto/proto"
)

// Consensus is the interface that must be implemented by all consensus
// engines.
type Consensus[T transaction.Tx, ValidatorUpdateT any] interface {
	InitGenesis(
		ctx context.Context, genesisBz []byte,
	) ([]ValidatorUpdateT, error)
	Prepare(context.Context, handlers.AppManager[T], []T, proto.Message) ([]T, error)
	Process(context.Context, handlers.AppManager[T], []T, proto.Message) error
	PreBlock(ctx context.Context, msg proto.Message) error
	EndBlock(ctx context.Context) ([]ValidatorUpdateT, error)
}
