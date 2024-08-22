package components

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/cli/pkg/components/log"
	cometbft "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
)

// ProvideCometBFTService provides the CometBFT service component.
func ProvideCometBFTService(
	logger *phuslu.Logger,
	storeKey **storetypes.KVStoreKey,
	abciMiddleware cometbft.MiddlewareI,
	db dbm.DB,
	cmtCfg *cmtcfg.Config,
	appOpts servertypes.AppOptions,
	chainSpec common.ChainSpec,
) *cometbft.Service {
	return cometbft.NewService(
		*storeKey,
		log.WrapSDKLogger(logger),
		db,
		abciMiddleware,
		true,
		cmtCfg,
		append(
			builder.DefaultServiceOptions(appOpts),
			builder.WithCometParamStore(chainSpec),
		)...,
	)
}
