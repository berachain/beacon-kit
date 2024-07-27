package cometbft

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
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
	LoggerT log.AdvancedLogger[any, LoggerT],
	NodeT engine.Client,
] struct {
	Logger LoggerT

	// CometBFT node
	CometBFTNode *node.Node
	// Client-side application to route
	// Comet calls to the Node
	App *Application[NodeT]

	config Config
}

func (c *Consensus[LoggerT, NodeT]) Start(ctx context.Context) error {
	// Should this generate a key if it doesn't exist?
	nodeKey, err := p2p.LoadNodeKey(c.config.NodeKeyFile)
	if err != nil {
		return err
	}

	if c.CometBFTNode, err = node.NewNode(
		ctx,
		&c.config.cmtConfig,
		privval.LoadFilePV(c.config.PrivValidatorKeyFile, c.config.PrivValidatorStateFile),
		nodeKey,
		proxy.NewConsensusSyncLocalClientCreator(c.App),
		nil, // TODO: implement the genesis doc provider
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(c.config.cmtConfig.Instrumentation),
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
