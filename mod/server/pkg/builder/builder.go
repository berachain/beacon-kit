package serverbuilder

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

type ServerBuilder[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
] struct {
	server *serverv2.Server[NodeT, T]

	depInjectCfg depinject.Config
	Components   []any
}

func New[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	opts ...Opt[NodeT, T, ValidatorUpdateT],
) *ServerBuilder[NodeT, T, ValidatorUpdateT] {
	b := &ServerBuilder[NodeT, T, ValidatorUpdateT]{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (sb *ServerBuilder[
	NodeT, T, ValidatorUpdateT,
]) Build(logger log.Logger) *serverv2.Server[NodeT, T] {
	// var (
	// logger    log.Logger
	// cmtServer *components.CometBFTServer[NodeT, T, ValidatorUpdateT]
	// )
	// if err := depinject.Inject(
	// 	depinject.Configs(
	// 		sb.depInjectCfg,
	// 		depinject.Provide(
	// 			sb.components...,
	// 		),
	// 	),
	// 	&logger,
	// 	&cmtServer,
	// ); err != nil {
	// 	panic(err)
	// }

	// sb.server = serverv2.NewServer(
	// 	logger,
	// 	sb.components...,
	// )

	return sb.server
}
