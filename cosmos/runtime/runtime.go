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

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/app/contracts"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	eth "github.com/itsdevbear/bolaris/beacon/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
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
		engine.WithEth1Client(eth1Client),
		engine.WithBeaconConfig(&p.cfg.BeaconConfig),
		engine.WithLogger(p.logger),
	}
	p.Caller = engine.NewCaller(engineCallerOpts...)

	// Create the blockchain service that will be used to process blocks.
	chainOpts := []blockchain.Option{
		blockchain.WithBeaconConfig(&p.cfg.BeaconConfig),
		blockchain.WithLogger(p.logger),
		blockchain.WithForkChoiceStoreProvider(bk),
		blockchain.WithEngineCaller(p.Caller),
	}
	blkChain := blockchain.NewService(chainOpts...)

	handlers := make(map[common.Address]callback.LogHandler)

	sc := &contracts.StakingCallbacks{}
	handlers[common.HexToAddress(
		"0x18Df82C7E422A42D47345Ed86B0E935E9718eBda",
	)], _ = callback.NewFrom(sc)
	// Build Log Processor
	logProcessorOpts := []logs.Option{
		logs.WithEthClient(eth1Client),
		logs.WithHandlers(handlers),
	}
	logProcessor, err := logs.NewProcessor(logProcessorOpts...)
	if err != nil {
		return err
	}

	// Build Proposal Handler
	defaultProposalHandler := baseapp.NewDefaultProposalHandler(mp, app)
	proposalHandler := proposal.NewHandler(blkChain,
		defaultProposalHandler.PrepareProposalHandler(), defaultProposalHandler.ProcessProposalHandler())
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
			p.logger, bk, blkChain, fn, logProcessor,
			// func(ctx sdk.Context) { _ = app.ModuleManager.PrepareCheckState },
		).PrepareCheckStater(),
	)

	return nil
}
