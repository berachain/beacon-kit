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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

type (
	execMode uint8
)

const (
	execModePrepareProposal     execMode = iota // Check a transaction
	execModeProcessProposal                     // Process a block proposal
	execModeVoteExtension                       // Extend or verify a pre-commit vote
	execModeVerifyVoteExtension                 // Verify a vote extension
	execModeFinalize                            // Finalize a block proposal
)

var _ servertypes.ABCI = (*BaseApp)(nil)

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	logger     log.Logger
	name       string                      // application name from abci.BlockInfo
	db         dbm.DB                      // common DB backend
	cms        storetypes.CommitMultiStore // Main
	Middleware MiddlewareI

	// volatile states:
	//
	// - prepareProposalState: Used for PrepareProposal, which is set based on
	// the previous block's state. This state is never committed. In case of
	// multiple consensus rounds, the state is always reset to the previous
	// block's state.
	//
	// - processProposalState: Used for ProcessProposal, which is set based on
	// the previous block's state. This state is never committed. In case of
	// multiple consensus rounds, the state is always reset to the previous
	// block's state.
	//
	// - finalizeBlockState: Used for FinalizeBlock, which is set based on the
	// previous block's state. This state is committed.
	prepareProposalState *state
	processProposalState *state
	finalizeBlockState   *state

	// An inter-block write-through cache provided to the context during the
	// ABCI
	// FinalizeBlock call.
	interBlockCache storetypes.MultiStorePersistentCache

	// paramStore is used to query for ABCI consensus parameters from an
	// application parameter store.
	paramStore ParamStore

	// initialHeight is the initial height at which we start the BaseApp
	initialHeight int64

	// minRetainBlocks defines the minimum block height offset from the current
	// block being committed, such that all blocks past this offset are pruned
	// from CometBFT. It is used as part of the process of determining the
	// ResponseCommit.RetainHeight value during ABCI Commit. A value of 0
	// indicates
	// that no blocks should be pruned.
	//
	// Note: CometBFT block pruning is dependent on this parameter in
	// conjunction with the unbonding (safety threshold) period, state pruning
	// and state sync
	// snapshot parameters to determine the correct minimum value of
	// ResponseCommit.RetainHeight.
	minRetainBlocks uint64

	// application's version string
	version string

	chainID string
}

// NewBaseApp returns a reference to an initialized cometbft. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
func NewBaseApp(
	storeKey *storetypes.KVStoreKey,
	logger log.Logger,
	db dbm.DB,
	middleware MiddlewareI,
	loadLatest bool,
	options ...func(*BaseApp),
) *BaseApp {
	app := &BaseApp{
		logger: logger.With(log.ModuleKey, "baseapp"),
		name:   "BeaconKit",
		db:     db,
		cms: store.NewCommitMultiStore(
			db,
			logger,
			storemetrics.NewNoOpMetrics(),
		), // by default we use a no-op metric gather in store
		Middleware: middleware,
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

// Name returns the name of the cometbft.
func (app *BaseApp) Name() string {
	return app.name
}

// CommitMultiStore returns the CommitMultiStore of the cometbft.
func (app *BaseApp) CommitMultiStore() storetypes.CommitMultiStore {
	return app.cms
}

// AppVersion returns the application's protocol version.
func (app *BaseApp) AppVersion(ctx context.Context) (uint64, error) {
	cp, err := app.paramStore.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get consensus params: %w", err)
	}
	if cp.Version == nil {
		return 0, nil
	}
	return cp.Version.App, nil
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(
	key storetypes.StoreKey,
	typ storetypes.StoreType,
) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running cometbft.
func (app *BaseApp) LoadLatestVersion() error {
	if err := app.cms.LoadLatestVersion(); err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	// Validator pruning settings.
	return app.cms.GetPruning().Validate()
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running cometbft.
func (app *BaseApp) LoadVersion(version int64) error {
	app.logger.Info(
		"NOTICE: this could take a long time to migrate IAVL store to fastnode if you enable Fast Node.\n",
	)
	err := app.cms.LoadVersion(version)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}

	// Validate Pruning settings.
	return app.cms.GetPruning().Validate()
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() storetypes.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

func (app *BaseApp) setMinRetainBlocks(minRetainBlocks uint64) {
	app.minRetainBlocks = minRetainBlocks
}

func (app *BaseApp) setInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) {
	app.interBlockCache = cache
}

// setState sets the BaseApp's state for the corresponding mode with a branched
// multi-store (i.e. a CacheMultiStore) and a new Context with the same
// multi-store branch, and provided header.
func (app *BaseApp) setState(mode execMode) {
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
// BaseApp's
// ParamStore. If the BaseApp has no ParamStore defined, nil is returned.
func (app *BaseApp) GetConsensusParams(
	ctx context.Context,
) cmtproto.ConsensusParams {
	//#nosec:G703 // bet.
	cp, _ := app.paramStore.Get(ctx)
	return cp
}

func (app *BaseApp) validateFinalizeBlockHeight(
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

// Close is called in start cmd to gracefully cleanup resources.
func (app *BaseApp) Close() error {
	var errs []error

	// Close app.db (opened by cosmos-sdk/server/start.go call to openDB)
	if app.db != nil {
		app.logger.Info("Closing application.db")
		if err := app.db.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
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
