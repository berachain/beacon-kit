package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	consensus "github.com/berachain/beacon-kit/mod/consensus/pkg"
	consensustypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
)

// this probably doesnt belong here but doing it for now

// lol
type ConsensusEngineInput[T transaction.Tx] struct {
	depinject.In
	TxCodec    *nodecomponents.TxCodec[T]
	Middleware *nodecomponents.ABCIMiddleware
}

// lolol
func ProvideConsensusEngine[T transaction.Tx, ValidatorUpdateT any](
	in ConsensusEngineInput[T],
) consensus.Engine[T, ValidatorUpdateT] {
	return consensus.NewEngine[T, ValidatorUpdateT](
		consensustypes.CometBFTConsensus,
		in.TxCodec,
		in.Middleware,
	)
}
