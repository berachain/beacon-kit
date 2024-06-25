package serverbuilder

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

type Opt[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
] func(*ServerBuilder[NodeT, T, ValidatorUpdateT])

func WithDepInjectConfig[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	cfg depinject.Config,
) Opt[NodeT, T, ValidatorUpdateT] {
	return func(b *ServerBuilder[NodeT, T, ValidatorUpdateT]) {
		b.depInjectCfg = cfg
	}
}

func WithComponents[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	components ...any,
) Opt[NodeT, T, ValidatorUpdateT] {
	return func(b *ServerBuilder[NodeT, T, ValidatorUpdateT]) {
		b.components = append(b.components, components...)
	}
}
