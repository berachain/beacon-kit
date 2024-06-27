package components

import (
	"cosmossdk.io/core/transaction"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
)

type ConsensusEngineInput[
	T transaction.Tx, ValidatorUpdateT any,
] struct {
	TxCodec    *components.TxCodec[T]
	Middleware *components.ABCIMiddleware
}

func ProvideConsensusEngine[T transaction.Tx, ValidatorUpdateT any](
	in ConsensusEngineInput[T, ValidatorUpdateT],
) *cometbft.ConsensusEngine[T, ValidatorUpdateT] {
	return cometbft.NewConsensusEngine[T, ValidatorUpdateT](
		in.TxCodec,
		in.Middleware,
	)
}
