package components

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/cometbft"
	"cosmossdk.io/server/v2/cometbft/mempool"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/comet"
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
	*cometbft.CometBFTServer[NodeT, T]

	consensus *comet.Consensus[T, ValidatorUpdateT]
}

type CometServerInput[T transaction.Tx, ValidatorUpdateT any] struct {
	depinject.In

	Consensus *comet.Consensus[T, ValidatorUpdateT]
	TxDecoder *TxDecoder[T]
}

// ProvideCometServer provides a CometServer.
func ProvideCometServer[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	in CometServerInput[T, ValidatorUpdateT],
) *CometBFTServer[NodeT, T, ValidatorUpdateT] {
	return &CometBFTServer[NodeT, T, ValidatorUpdateT]{
		CometBFTServer: cometbft.New[NodeT, T](in.TxDecoder),
		consensus:      in.Consensus,
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
	s.CometBFTServer.App.SetMempool(mempool.NoOpMempool[T]{})
	s.CometBFTServer.App.SetPrepareProposalHandler(s.consensus.PrepareProposal)
	s.CometBFTServer.App.SetProcessProposalHandler(s.consensus.ProcessProposal)
	return nil
}
