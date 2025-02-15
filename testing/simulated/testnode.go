package simulated

import (
	"io"
	"os"
	"testing"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/log/phuslu"
	nodetypes "github.com/berachain/beacon-kit/node-core/types"
	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/berachain/beacon-kit/storage/db"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// TestNodeInput takes the input for building and starting a node
type TestNodeInput struct {
	TempHomeDir string
	CometConfig *cmtcfg.Config
	AuthRPC     *url.ConnectionURL
	Logger      *phuslu.Logger
	AppOpts     *viper.Viper
	Components  []any
}

type TestNode struct {
	nodetypes.Node
	StorageBackend blockchain.StorageBackend
}

// NewTestNode Uses the testnet chainspec.
func NewTestNode(
	t *testing.T,
	input TestNodeInput,
) TestNode {
	t.Helper()

	beaconKitConfig := createBeaconKitConfig(t)
	beaconKitConfig.Engine.RPCDialURL = input.AuthRPC
	appOpts := getAppOptions(t, input.AppOpts, beaconKitConfig, input.TempHomeDir)

	// Create a database
	database, err := db.OpenDB(input.TempHomeDir, dbm.PebbleDBBackend)
	require.NoError(t, err)

	// Build a node
	node := buildNode(
		input.Logger,
		database,
		os.Stdout, // or some other writer
		input.CometConfig,
		appOpts,
		input.Components,
	)
	return node
}

// buildNode run the same logic as primary build, but it returns the components allowing us to query them.
func buildNode(
	logger *phuslu.Logger,
	db dbm.DB,
	_ io.Writer,
	cmtCfg *cmtcfg.Config,
	appOpts servertypes.AppOptions,
	components []any,
) TestNode {
	// variables to hold the components needed to set up BeaconApp
	var (
		apiBackend interface {
			AttachQueryBackend(nodetypes.ConsensusService)
		}
		beaconNode     nodetypes.Node
		cmtService     nodetypes.ConsensusService
		config         *config.Config
		storageBackend blockchain.StorageBackend
	)

	// build all node components using depinject
	if err := depinject.Inject(
		depinject.Configs(
			depinject.Provide(
				components...,
			),
			depinject.Supply(
				appOpts,
				logger,
				db,
				cmtCfg,
			),
		),
		&apiBackend,
		&beaconNode,
		&cmtService,
		&config,
		&storageBackend,
	); err != nil {
		panic(err)
	}
	if config == nil {
		panic("config is nil")
	}
	if apiBackend == nil {
		panic("node or api backend is nil")
	}

	logger.WithConfig(config.GetLogger())
	apiBackend.AttachQueryBackend(cmtService)
	return TestNode{beaconNode, storageBackend}
}
