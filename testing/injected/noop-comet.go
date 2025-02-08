package injected

import (
	"context"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	"github.com/berachain/beacon-kit/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/builder"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NoopCometService is normal Comet under the hood, but we override the Start method to avoid starting the actual
// CometBFT core loop so that we can orchestrate it ourselves.
type NoopCometService struct {
	comet *cometbft.Service
}

func ProvideNoopCometService(
	logger *phuslu.Logger,
	blockchain blockchain.BlockchainI,
	blockBuilder validator.BlockBuilderI,
	db dbm.DB,
	cmtCfg *cmtcfg.Config,
	appOpts config.AppOptions,
	telemetrySink *metrics.TelemetrySink) *NoopCometService {
	return &NoopCometService{
		cometbft.NewService(
			logger,
			db,
			blockchain,
			blockBuilder,
			cmtCfg,
			telemetrySink,
			builder.DefaultServiceOptions(appOpts)...,
		)}
}

func (n *NoopCometService) Start(_ context.Context) error {
	return nil
}

func (n *NoopCometService) Stop() error {
	return nil
}

func (n *NoopCometService) Name() string {
	return n.comet.Name()
}

func (n *NoopCometService) CreateQueryContext(_ int64, _ bool) (sdk.Context, error) {
	panic("unimplemented")
}

func (n *NoopCometService) LastBlockHeight() int64 {
	panic("unimplemented")
}
