// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package runtime

import (
	"context"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	blocksync "github.com/itsdevbear/bolaris/beacon/execution/block-sync"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/cosmos/abci/commit"
	"github.com/itsdevbear/bolaris/cosmos/abci/preblock"
	proposal "github.com/itsdevbear/bolaris/cosmos/abci/proposal"
	beaconkeeper "github.com/itsdevbear/bolaris/cosmos/x/beacon/keeper"
	"github.com/itsdevbear/bolaris/types/config"
)

// BeaconKeeper is an interface that defines the methods needed for the EVM setup.
type BeaconKeeper interface {
	// Setup initializes the EVM keeper.
	Setup(engine.Caller) error
}

// CosmosApp is an interface that defines the methods needed for the Cosmos setup.
type CosmosApp interface {
	SetPrepareProposal(sdk.PrepareProposalHandler)
	baseapp.ProposalTxVerifier
	SetMempool(mempool.Mempool)
	SetAnteHandler(sdk.AnteHandler)
	SetExtendVoteHandler(sdk.ExtendVoteHandler)
	SetProcessProposal(sdk.ProcessProposalHandler)
	SetVerifyVoteExtensionHandler(sdk.VerifyVoteExtensionHandler)
	PreBlocker() sdk.PreBlocker
	SetPreBlocker(sdk.PreBlocker)
	SetPrepareCheckStater(sdk.PrepareCheckStater)
	ChainID() string
}

// Polaris is a struct that wraps the Polaris struct from the polar package.
type Polaris struct {
	cfg *config.Config
	engine.Caller
	blocksyncer *blocksync.BlockSync
	// logger is the underlying logger supplied by the sdk.
	logger log.Logger
}

// New creates a new Polaris runtime from the provided
// dependencies.
func New(
	appOpts servertypes.AppOptions,
	logger log.Logger,
) (*Polaris, error) {
	var err error
	p := &Polaris{
		logger: logger,
	}

	// Read the configuration from the cosmos app options
	p.cfg, err = config.ReadConfigFromAppOpts(appOpts)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// New creates a new Polaris runtime from the provided
// dependencies, panics on error.
func MustNew(appOpts servertypes.AppOptions, logger log.Logger) *Polaris {
	p, err := New(appOpts, logger)
	if err != nil {
		panic(err)
	}
	return p
}

// Build is a function that sets up the Polaris struct.
// It takes a BaseApp and an BeaconKeeper as arguments.
// It returns an error if the setup fails.
func (p *Polaris) Build(app CosmosApp, bk *beaconkeeper.Keeper) error {
	mp := mempool.NewSenderNonceMempool()
	app.SetMempool(mp)

	jwtSecret, err := eth.LoadJWTSecret(p.cfg.ExecutionClient.JWTSecretPath, p.logger)
	if err != nil {
		return err
	}

	// Create the eth1 client that will be used to interact with the execution client.
	opts := []eth.Option{
		eth.WithHTTPEndpointAndJWTSecret(p.cfg.ExecutionClient.RPCDialURL, jwtSecret),
		eth.WithLogger(p.logger),
		eth.WithRequiredChainID(p.cfg.ExecutionClient.RequiredChainID),
	}
	eth1Client, err := eth.NewEth1Client(context.Background(), opts...)
	if err != nil {
		return err
	}

	// Engine Caller wraps the eth1 client and provides the interface for the
	// blockchain service to interact with the execution client.
	engineCallerOpts := []engine.Option{
		engine.WithBeaconConfig(&p.cfg.BeaconConfig),
		engine.WithLogger(p.logger),
	}
	p.Caller = engine.NewCaller(eth1Client, engineCallerOpts...)

	// Create the blockchain service that will be used to process blocks.
	chainOpts := []blockchain.Option{
		blockchain.WithBeaconConfig(&p.cfg.BeaconConfig),
		blockchain.WithLogger(p.logger),
		blockchain.WithForkChoiceStoreProvider(bk),
		blockchain.WithEngineCaller(p.Caller),
	}
	blkChain := blockchain.NewService(chainOpts...)

	// Block Syncer
	blockSyncOpts := []blocksync.Option{
		blocksync.WithBeaconConfig(&p.cfg.BeaconConfig),
		blocksync.WithLogger(p.logger),
		blocksync.WithHeadSubscriber(p.Caller),
		blocksync.WithForkChoiceStoreProvider(bk),
	}
	p.blocksyncer = blocksync.New(blockSyncOpts...)
	p.blocksyncer.Start(context.TODO())

	// Build Proposal Handler
	defaultProposalHandler := baseapp.NewDefaultProposalHandler(mp, app)
	proposalHandler := proposal.NewHandler(blkChain,
		defaultProposalHandler.PrepareProposalHandler(), defaultProposalHandler.ProcessProposalHandler(),
		p.blocksyncer)
	app.SetPrepareProposal(proposalHandler.PrepareProposalHandler)
	app.SetProcessProposal(proposalHandler.ProcessProposalHandler)

	// Build PreBlock Handler
	app.SetPreBlocker(
		preblock.NewBeaconPreBlockHandler(p.logger, bk, nil).PreBlocker(),
	)

	fn := func(ctx sdk.Context) {}
	// Build PrepareCheckStater
	app.SetPrepareCheckStater(
		commit.NewBeaconPrepareCheckStateHandler(
			p.logger, bk, blkChain, fn,
			// func(ctx sdk.Context) { _ = app.ModuleManager.PrepareCheckState },
		).PrepareCheckStater(),
	)

	return nil
}

func (p *Polaris) SyncEL(ctx context.Context) error {
	return p.blocksyncer.WaitforExecutionClientSync(ctx)
}
