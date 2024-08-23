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
	"errors"
	"fmt"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	servercmtlog "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/log"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/params"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

type (
	execMode uint8
)

const (
	execModePrepareProposal execMode = iota
	execModeProcessProposal
	execModeFinalize
)

const InitialAppVersion uint64 = 0

type Service struct {
	node   *node.Node
	cmtCfg *cmtcfg.Config
	// initialized on creation
	logger     log.Logger
	name       string
	db         dbm.DB
	cms        storetypes.CommitMultiStore
	Middleware MiddlewareI

	prepareProposalState *state
	processProposalState *state
	finalizeBlockState   *state
	interBlockCache      storetypes.MultiStorePersistentCache
	paramStore           *params.ConsensusParamsStore
	initialHeight        int64
	minRetainBlocks      uint64
	// application's version string
	version string
	chainID string
}

func NewService(
	storeKey *storetypes.KVStoreKey,
	logger log.Logger,
	db dbm.DB,
	middleware MiddlewareI,
	loadLatest bool,
	cmtCfg *cmtcfg.Config,
	cs common.ChainSpec,
	options ...func(*Service),
) *Service {
	app := &Service{
		logger: logger.With(log.ModuleKey, "cometbft"),
		name:   "beacond",
		db:     db,
		cms: store.NewCommitMultiStore(
			db,
			logger,
			storemetrics.NewNoOpMetrics(),
		),
		Middleware: middleware,
		cmtCfg:     cmtCfg,
		paramStore: params.NewConsensusParamsStore(cs),
	}

	app.SetVersion(version.Version)
	app.MountStore(storeKey, storetypes.StoreTypeIAVL)

	for _, option := range options {
		option(app)
	}

	if app.interBlockCache != nil {
		app.cms.SetInterBlockCache(app.interBlockCache)
	}

	// Load the app.
	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			panic(err)
		}
	}

	return app
}

// TODO: Move nodeKey into being created within the function.
func (app *Service) Start(
	ctx context.Context,
) error {
	cfg := app.cmtCfg
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return err
	}

	app.node, err = node.NewNode(
		ctx,
		cfg,
		pvm.LoadOrGenFilePV(
			cfg.PrivValidatorKeyFile(),
			cfg.PrivValidatorStateFile(),
		),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		GetGenDocProvider(cfg),
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		servercmtlog.CometLoggerWrapper{Logger: app.logger},
	)
	if err != nil {
		return err
	}

	return app.node.Start()
}

// Close is called in start cmd to gracefully cleanup resources.
func (app *Service) Close() error {
	var errs []error

	if app.node != nil && app.node.IsRunning() {
		app.logger.Info("Stopping CometBFT Node")
		//#nosec:G703 // its a bet.
		_ = app.node.Stop()
	}

	// Close app.db (opened by cosmos-sdk/server/start.go call to openDB)
	if app.db != nil {
		app.logger.Info("Closing application.db")
		if err := app.db.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Name returns the name of the cometbft.
func (app *Service) Name() string {
	return app.name
}

// CommitMultiStore returns the CommitMultiStore of the cometbft.
func (app *Service) CommitMultiStore() storetypes.CommitMultiStore {
	return app.cms
}

// AppVersion returns the application's protocol version.
func (app *Service) AppVersion(ctx context.Context) (uint64, error) {
	cp, err := app.paramStore.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get consensus params: %w", err)
	}
	if cp.Version == nil {
		return 0, nil
	}
	return cp.Version.App, nil
}

// MountStore mounts a store to the provided key in the Service multistore,
// using the default DB.
func (app *Service) MountStore(
	key storetypes.StoreKey,
	typ storetypes.StoreType,
) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

func (app *Service) LoadLatestVersion() error {
	if err := app.cms.LoadLatestVersion(); err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	// Validator pruning settings.
	return app.cms.GetPruning().Validate()
}

func (app *Service) LoadVersion(version int64) error {
	err := app.cms.LoadVersion(version)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}

	// Validate Pruning settings.
	return app.cms.GetPruning().Validate()
}

// LastCommitID returns the last CommitID of the multistore.
func (app *Service) LastCommitID() storetypes.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *Service) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

func (app *Service) setMinRetainBlocks(minRetainBlocks uint64) {
	app.minRetainBlocks = minRetainBlocks
}

func (app *Service) setInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) {
	app.interBlockCache = cache
}

func (app *Service) setState(mode execMode) {
	ms := app.cms.CacheMultiStore()
	baseState := &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, false, app.logger),
	}

	switch mode {
	case execModePrepareProposal:
		app.prepareProposalState = baseState

	case execModeProcessProposal:
		app.processProposalState = baseState

	case execModeFinalize:
		app.finalizeBlockState = baseState

	default:
		panic(fmt.Sprintf("invalid runTxMode for setState: %d", mode))
	}
}

// GetConsensusParams returns the current consensus parameters from the
// Service's
// ParamStore. If the Service has no ParamStore defined, nil is returned.
func (app *Service) GetConsensusParams(
	ctx context.Context,
) cmtproto.ConsensusParams {
	//#nosec:G703 // bet.
	cp, _ := app.paramStore.Get(ctx)
	return cp
}

func (app *Service) validateFinalizeBlockHeight(
	req *abci.FinalizeBlockRequest,
) error {
	if req.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Height)
	}

	lastBlockHeight := app.LastBlockHeight()

	// expectedHeight holds the expected height to validate
	var expectedHeight int64
	if lastBlockHeight == 0 && app.initialHeight > 1 {
		// In this case, we're validating the first block of the chain, i.e no
		// previous commit. The height we're expecting is the initial height.
		expectedHeight = app.initialHeight
	} else {
		// This case can mean two things:
		//
		// - Either there was already a previous commit in the store, in which
		// case we increment the version from there.
		// - Or there was no previous commit, in which case we start at version
		// 1.
		expectedHeight = lastBlockHeight + 1
	}

	if req.Height != expectedHeight {
		return fmt.Errorf(
			"invalid height: %d; expected: %d",
			req.Height,
			expectedHeight,
		)
	}

	return nil
}

// convertValidatorUpdate abstracts the conversion of a
// transition.ValidatorUpdate to an appmodulev2.ValidatorUpdate.
// TODO: this is so hood, bktypes -> sdktypes -> generic is crazy
// maybe make this some kind of codec/func that can be passed in?
func convertValidatorUpdate[ValidatorUpdateT any](
	u **transition.ValidatorUpdate,
) (ValidatorUpdateT, error) {
	var valUpdate ValidatorUpdateT
	update := *u
	if update == nil {
		return valUpdate, errors.New("undefined validator update")
	}
	return any(abci.ValidatorUpdate{
		PubKeyBytes: update.Pubkey[:],
		PubKeyType:  crypto.CometBLSType,
		//#nosec:G701 // this is safe.
		Power: int64(update.EffectiveBalance.Unwrap()),
	}).(ValidatorUpdateT), nil
}
