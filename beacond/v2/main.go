package main

import (
	"context"
	"encoding/json"
	"os"

	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
)

func run() error {
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Engine.JWTSecretPath = "./testing/files/jwt.hex"
	appOpts := &components.AppOptions{
		HomeDir: ".tmp/beacond",
	}
	logger := clicomponents.ProvideLogger(clicomponents.LoggerInput{
		Cfg: cfg,
		Out: os.Stdout,
	})
	chainSpec := components.ProvideChainSpec()

	var (
		storageBackend  = &StorageBackend{}
		stateProcessor  = &StateProcessor{}
		consensusClient = &RuntimeApp{}
	)
	appBuilder := app.NewBuilder[*StorageBackend, *StateProcessor]()
	consensus := cometbft.NewConsensus(cfg.CometBFT, logger, appBuilder.App(), chainSpec)
	privKey, err := consensus.Init(appOpts.HomeDir)
	if err != nil {
		return err
	}
	pubKey, err := privKey.GetPubKey()
	if err != nil {
		return err
	}

	genesisBz, err := json.Marshal(genesis.DefaultGenesisDeneb())
	if err != nil {
		return err
	}
	genesis := cometbft.NewGenesis(genesisBz, cometbft.DefaultConsensusParams().ConsensusParams)
	if err := genesis.SaveAs(cfg.CometBFT.GenesisFile()); err != nil {
		return err
	}
	if err := AddDepositToGenesis(genesis, chainSpec, cfg, appOpts, bytes.B48(pubKey.Bytes())); err != nil {
		return err
	}
	if err := AddExecutionPayloadToGenesis(genesis, chainSpec, &cfg.CometBFT, "testing/files/eth-genesis.json"); err != nil {
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

	if err := consensus.Start(ctx); err != nil {
		return err
	}

	// Keep the process alive until ctx.Done()
	<-ctx.Done()
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
