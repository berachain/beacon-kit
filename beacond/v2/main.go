package main

import (
	"context"
	"os"

	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components"
)

// TODO: inject the genesis bytes
func run() error {
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Engine.JWTSecretPath = "./testing/files/jwt.hex"
	appOpts := &components.AppOptions{
		HomeDir: ".tmp/testingd",
	}
	logger := clicomponents.ProvideLogger(clicomponents.LoggerInput{
		Cfg: cfg,
		Out: os.Stdout,
	})
	chainSpec := components.ProvideChainSpec()

	// genesisBz, err := json.Marshal(genesis.DefaultGenesisDeneb())
	// if err != nil {
	// 	return err
	// }

	var (
		storageBackend  = &StorageBackend{}
		stateProcessor  = &StateProcessor{}
		consensusClient = &RuntimeApp{}
	)
	appBuilder := app.NewBuilder[*StorageBackend, *StateProcessor]()
	consensus := cometbft.NewConsensus(cfg.CometBFT, logger, appBuilder.App(), chainSpec)
	if err := consensus.Init(); err != nil {
		return err
	}

	appBuilder.WithComponents(components.DefaultComponentsWithStandardTypes()...)
	appBuilder.WithStateProcessor(stateProcessor)
	appBuilder.WithStorageBackend(storageBackend)
	appBuilder.WithConsensusClient(consensusClient)
	app, err := appBuilder.Build(logger, appOpts, cfg)
	if err != nil {
		return err
	}

	if err := app.Start(ctx); err != nil {
		return err
	}

	return consensus.Start(ctx)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
