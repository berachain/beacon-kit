// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
	seeds = "c28827cb96c14c905b127b92065a3fb4cd77d7f6@testnet-seeds.whispernode.com:25456," +
		"8a0fbd4a06050519b6bce88c03932bd0a57060bd@beacond-testnet.blacknodes.net:26656," +
		"d9903c2f9902243c88f2758fe2e81e305e737fb3@bera-testnet-seeds.nodeinfra.com:26656," +
		"9c50cc419880131ea02b6e2b92027cefe17941b9@139.59.151.125:26656," +
		"cf44af098820f50a1e96d51b7af6861bc961e706@berav2-seeds.staketab.org:5001," +
		"6b5040a0e1b29a2cca224b64829d8f3d8796a3e3@berachain-testnet-v2-2.seed.l0vd.com:21656"
		// "4f93da5553f0dfaafb620532901e082255ec3ad3@berachain-testnet-v2-1.seed.l0vd.com:61656,"
		// +
		// "a62eefaa284eaede7460315d2f1d1f92988e01f1@135.125.188.10:26656"

	chainID           = "beacond-2061"
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
	consensus := cometbft.NewConsensus(
		cfg.CometBFT,
		logger,
		appBuilder.App(),
		chainSpec,
	)

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
		genesis := cometbft.NewGenesis(
			chainID,
			genesisBz,
			cometbft.DefaultConsensusParams().ConsensusParams,
		)
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

	appBuilder.WithComponents(
		components.DefaultComponentsWithStandardTypes()...)
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
