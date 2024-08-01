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

	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	cmtcfg "github.com/cometbft/cometbft/config"
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
	LoggerT log.Logger[any],
	NodeT engine.Client,
] struct {
	Logger LoggerT

	// CometBFT node
	CometBFTNode *node.Node
	// Client-side application to route
	// Comet calls to the Node
	App *Application[NodeT]

	// Config
	config Config
}

func NewConsensus[
	LoggerT log.Logger[any],
	NodeT engine.Client,
](
	cfg Config,
	logger LoggerT,
	node NodeT,
	chainSpec common.ChainSpec,
) *Consensus[LoggerT, NodeT] {
	return &Consensus[LoggerT, NodeT]{
		Logger: logger,
		config: cfg,
		App:    NewApplication(logger, node, chainSpec),
	}
}

func (c *Consensus[LoggerT, NodeT]) Init() error {
	// This function needs to build the validator files
	// and the config file??? (if that's necessary)
	return nil
}

func (c *Consensus[LoggerT, NodeT]) Start(ctx context.Context) error {
	// Should this generate a key if it doesn't exist?
	nodeKey, err := p2p.LoadOrGenNodeKey(c.config.NodeKeyFile)
	if err != nil {
		return err
	}

	if c.CometBFTNode, err = node.NewNode(
		ctx,
		c.config.CometConfig,
		privval.LoadFilePV(c.config.PrivValidatorKeyFile, c.config.PrivValidatorStateFile),
		nodeKey,
		proxy.NewConsensusSyncLocalClientCreator(c.App),
		nil, // TODO: implement the genesis doc provider
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(c.config.CometConfig.Instrumentation),
		// cometLoggerFromLogger(c.Logger),
		cmtlog.NewNopLogger(), // TODO: make adapter for our logger
	); err != nil {
		return err
	}

	return c.CometBFTNode.Start()
}

func (c *Consensus[LoggerT, NodeT]) Stop(context.Context) error {
	if c.CometBFTNode != nil && c.CometBFTNode.IsRunning() {
		return c.CometBFTNode.Stop()
	}

	return nil
}

// func cometLoggerFromLogger[LoggerT log.AdvancedLogger[any, LoggerT]](
// 	logger LoggerT,
// ) cmtlog.Logger {
// 	return logger
// }
