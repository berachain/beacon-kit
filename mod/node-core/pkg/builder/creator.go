package builder

import (
	"io"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// AppCreator is a function that creates an application.
// It is necessary to adhere to the types.AppCreator[T] interface.
func (b *Builder[NodeT]) AppCreator(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) NodeT {
	// Check for goleveldb cause bad project.
	if appOpts.Get("app-db-backend") == "goleveldb" {
		panic("goleveldb is not supported")
	}

	b.node.SetApplication(app.NewBeaconKitApp(
		logger, db, traceStore, true,
		appOpts,
		b.depInjectCfg,
		b.chainSpec,
		append(
			server.DefaultBaseappOptions(appOpts),
			func(bApp *baseapp.BaseApp) {
				bApp.SetParamStore(comet.NewConsensusParamsStore(b.chainSpec))
			})...,
	))
	return b.node
}
