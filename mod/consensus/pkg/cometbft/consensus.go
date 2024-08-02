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

package cometbft

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	cmtcfg "github.com/cometbft/cometbft/config"
	bls12381 "github.com/cometbft/cometbft/crypto/bls12381"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
)

// Consensus is a wrapper around the CometBFT node and client-side application
// which serves the responsibilty of receiving and routing ABCI requests to the
// node, and returning the responses to the consensus engine.
type Consensus[
	LoggerT log.AdvancedLogger[any, LoggerT],
	ClientT engine.Client,
] struct {
	Logger LoggerT

	// CometBFT node
	CometBFTNode *node.Node
	// Client-side application to route
	// Comet calls to the Node
	App *Application[ClientT]

	// Config
	config Config
}

func NewConsensus[
	LoggerT log.AdvancedLogger[any, LoggerT],
	ClientT engine.Client,
](
	cfg Config,
	logger LoggerT,
	client ClientT,
	chainSpec common.ChainSpec,
) *Consensus[LoggerT, ClientT] {
	return &Consensus[LoggerT, ClientT]{
		Logger: logger,
		config: cfg,
		App:    NewApplication(logger, client, chainSpec),
	}
}

func (c *Consensus[LoggerT, ClientT]) Init(homeDir string) (*privval.FilePV, error) {
	// Ensure the config subdirectory exists
	configDir := filepath.Join(homeDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	dataDir := filepath.Join(homeDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	c.config.Config = c.config.SetRoot(homeDir)

	pvKeyFile := c.config.PrivValidatorKeyFile()
	pvStateFile := c.config.PrivValidatorStateFile()

	var pv *privval.FilePV
	if _, err := os.Stat(pvKeyFile); os.IsNotExist(err) {
		privKey, err := bls12381.GenPrivKey()
		if err != nil {
			return nil, err
		}
		pv = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
		pv.Save()
	} else if err != nil {
		return nil, err
	} else {
		pv = privval.LoadFilePV(pvKeyFile, pvStateFile)
	}

	return pv, nil
}

func (c *Consensus[LoggerT, ClientT]) Start(ctx context.Context) error {
	// Should this generate a key if it doesn't exist?
	nodeKey, err := p2p.LoadOrGenNodeKey(c.config.NodeKeyFile())
	if err != nil {
		return err
	}

	if c.CometBFTNode, err = node.NewNode(
		ctx,
		c.config.Config,
		privval.LoadFilePV(c.config.PrivValidatorKeyFile(), c.config.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewConsensusSyncLocalClientCreator(c.App),
		node.DefaultGenesisDocProviderFunc(c.config.Config),
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(c.config.Instrumentation),
		cometLoggerFromLogger(c.Logger),
	); err != nil {
		return err
	}

	return c.CometBFTNode.Start()
}

func (c *Consensus[LoggerT, ClientT]) Stop(context.Context) error {
	if c.CometBFTNode != nil && c.CometBFTNode.IsRunning() {
		return c.CometBFTNode.Stop()
	}

	return nil
}

type CometLogger[LoggerT log.AdvancedLogger[any, LoggerT]] struct {
	Logger LoggerT
}

func cometLoggerFromLogger[LoggerT log.AdvancedLogger[any, LoggerT]](
	logger LoggerT,
) cmtlog.Logger {
	return &CometLogger[LoggerT]{Logger: logger}
}

func (cl *CometLogger[LoggerT]) Debug(msg string, keyVals ...interface{}) {
	cl.Logger.Debug(msg, keyVals...)
}

func (cl *CometLogger[LoggerT]) Info(msg string, keyVals ...interface{}) {
	cl.Logger.Info(msg, keyVals...)
}

func (cl *CometLogger[LoggerT]) Error(msg string, keyVals ...interface{}) {
	cl.Logger.Error(msg, keyVals...)
}

func (cl *CometLogger[LoggerT]) With(keyVals ...interface{}) cmtlog.Logger {
	return &CometLogger[LoggerT]{Logger: cl.Logger.With(keyVals...)}
}
