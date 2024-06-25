package components

import (
	"cosmossdk.io/core/transaction"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

func DefaultServerComponents[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
]() []any {
	return []any{
		ProvideCometServer[NodeT, T, ValidatorUpdateT],
		ProvideConsensus[T, ValidatorUpdateT],
		ProvideTxDecoder[T],
	}
}
