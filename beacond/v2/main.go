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

const (
	seeds = "2f8ce8462cddc9ae865ab8ec1f05cc286f07c671@34.152.0.40:26656,3037b09eaa2eed5cd1b1d3d733ab8468bf4910ee@35.203.36.128:26656,add35d414bee9c0be3b10bcf8fbc12a059eb9a3b@35.246.180.53:26656,925221ce669017eb2fd386bc134f13c03c5471d4@34.159.151.132:26656,ae50b817fcb2f35da803aa0190a5e37f4f8bcdb5@34.64.62.166:26656"

	chainID           = "bartio-beacon-80084"
	homeDir           = ".tmp/beacond"
	jwtSecretPath     = "testing/files/jwt.hex"
	ethGenesisPath    = "testing/files/eth-genesis.json"
	bartioGenesisPath = "testing/networks/80084/genesis.json"
)

func run() error {
	sync := os.Getenv("CHAIN_SPEC") == "testnet"
	ctx := context.Background()
	cfg := config.DefaultConfig()
	cfg.Engine.JWTSecretPath = jwtSecretPath
	appOpts := &components.AppOptions{
		HomeDir: homeDir,
	}
	logger := clicomponents.ProvideLogger(clicomponents.LoggerInput{
		Cfg: cfg,
		Out: os.Stdout,
	})
	chainSpec := components.ProvideChainSpec()

	if sync {
		// Add testnet syncing seeds
		cfg.CometBFT.P2P.Seeds = seeds
	}

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
	if !sync {
		pubKey, err := privKey.GetPubKey()
		if err != nil {
			return err
		}
		genesisBz, err := json.Marshal(genesis.DefaultGenesisDeneb())
		if err != nil {
			return err
		}
		genesis := cometbft.NewGenesis(chainID, genesisBz, cometbft.DefaultConsensusParams().ConsensusParams)
		if err := genesis.SaveAs(cfg.CometBFT.GenesisFile()); err != nil {
			return err
		}
		if err := AddDepositToGenesis(genesis, chainSpec, cfg, appOpts, bytes.B48(pubKey.Bytes())); err != nil {
			return err
		}
		if err := AddExecutionPayloadToGenesis(genesis, chainSpec, &cfg.CometBFT, ethGenesisPath); err != nil {
			return err
		}
	} else {
		if err := ConvertBartioGenesis(&cfg.CometBFT, bartioGenesisPath); err != nil {
			return err
		}
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
