package consensus

import (
	"context"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/server/v2/cometbft/handlers"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	"github.com/cosmos/gogoproto/proto"
)

// Engine is the interface that must be implemented by all consensus
// engines.
type Engine[T transaction.Tx, ValidatorUpdateT any] interface {
	InitGenesis(ctx context.Context, genesisBz []byte) ([]ValidatorUpdateT, error)
	Prepare(context.Context, handlers.AppManager[T], []T, proto.Message) ([]T, error)
	Process(context.Context, handlers.AppManager[T], []T, proto.Message) error
	PreBlock(ctx context.Context, msg proto.Message) error
	EndBlock(ctx context.Context) ([]ValidatorUpdateT, error)
}

var _ Engine[transaction.Tx, any] = (*cometbft.ConsensusEngine[transaction.Tx, any])(nil)

func NewEngine[T transaction.Tx, ValidatorUpdateT any](
	version types.ConsensusVersion,
	txCodec transaction.Codec[T],
	m types.Middleware,
) Engine[T, ValidatorUpdateT] {
	switch version {
	case types.CometBFTConsensus:
		return cometbft.NewConsensusEngine[T, ValidatorUpdateT](
			txCodec, m,
		)
	case types.RollKitConsensus:
		panic("not implemented")
	default:
		panic("unknown consensus version")
	}
}
