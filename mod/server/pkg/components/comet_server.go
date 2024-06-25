package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	sdkcomet "cosmossdk.io/server/v2/cometbft"
	"cosmossdk.io/server/v2/cometbft/mempool"
	"github.com/berachain/beacon-kit/mod/consensus/pkg"
	consensustypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/spf13/viper"
)

var _ serverv2.ServerComponent[
	types.Node[transaction.Tx], transaction.Tx,
] = (*CometBFTServer[
	types.Node[transaction.Tx], transaction.Tx, any,
])(nil)

type CometBFTServer[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
] struct {
	*sdkcomet.CometBFTServer[NodeT, T]

	TxCodec transaction.Codec[T]
}

type CometServerInput[T transaction.Tx, ValidatorUpdateT any] struct {
	depinject.In

	TxCodec *nodecomponents.TxCodec[T]
}

// ProvideCometServer provides a CometServer.
func ProvideCometServer[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	in CometServerInput[T, ValidatorUpdateT],
) *CometBFTServer[NodeT, T, ValidatorUpdateT] {
	return &CometBFTServer[NodeT, T, ValidatorUpdateT]{
		CometBFTServer: sdkcomet.New[NodeT, T](in.TxCodec),
		TxCodec:        in.TxCodec,
	}
}

// Init wraps the default Init method and sets the PrepareProposal and
// ProcessProposal handlers.
func (s *CometBFTServer[NodeT, T, ValidatorUpdateT]) Init(
	node NodeT, v *viper.Viper, logger log.Logger,
) error {
	if err := s.CometBFTServer.Init(node, v, logger); err != nil {
		return err
	}
	var middleware nodecomponents.ABCIMiddleware
	registry := node.GetServiceRegistry()
	if err := registry.FetchService(&middleware); err != nil {
		return err
	}

	engine := consensus.NewEngine[T, ValidatorUpdateT](
		consensustypes.CometBFTConsensus,
		s.TxCodec,
		&middleware,
	)

	s.CometBFTServer.App.SetMempool(mempool.NoOpMempool[T]{})
	s.CometBFTServer.App.SetPrepareProposalHandler(engine.Prepare)
	s.CometBFTServer.App.SetProcessProposalHandler(engine.Process)
	return nil
}
