package simulated

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

// SimComet is normal Comet under the hood, but we override the Start method to avoid starting the actual
// CometBFT core loop so that we can orchestrate it ourselves.
type SimComet struct {
	//  We are forced to stutter here as we want to override the implementations of the original comet service.
	Comet *cometbft.Service
}

func ProvideSimComet(
	logger *phuslu.Logger,
	blockchain blockchain.BlockchainI,
	blockBuilder validator.BlockBuilderI,
	db dbm.DB,
	cmtCfg *cmtcfg.Config,
	appOpts config.AppOptions,
	telemetrySink *metrics.TelemetrySink) *SimComet {
	return &SimComet{
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

func (s *SimComet) Start(_ context.Context) error {
	return nil
}

func (s *SimComet) Stop() error {
	return nil
}

func (s *SimComet) Name() string {
	return s.Comet.Name()
}

func (s *SimComet) CreateQueryContext(height int64, prove bool) (sdk.Context, error) {
	return s.Comet.CreateQueryContext(height, prove)
}

func (s *SimComet) LastBlockHeight() int64 {
	panic("unimplemented")
}
