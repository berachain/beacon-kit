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

package baseapp

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"cosmossdk.io/core/header"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v1"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/exp/maps"
)

type (
	execMode uint8

	// StoreLoader defines a customizable function to control how we load the
	// CommitMultiStore from disk. This is useful for state migration, when
	// loading a datastore written with an older version of the software. In
	// particular, if a module changed the substore key name (or removed a
	// substore)
	// between two versions of the software.
	StoreLoader func(ms storetypes.CommitMultiStore) error
)

const (
	execModeCheck               execMode = iota // Check a transaction
	execModeReCheck                             // Recheck a (pending) transaction after a commit
	execModeSimulate                            // Simulate a transaction
	execModePrepareProposal                     // Prepare a block proposal
	execModeProcessProposal                     // Process a block proposal
	execModeVoteExtension                       // Extend or verify a pre-commit vote
	execModeVerifyVoteExtension                 // Verify a vote extension
	execModeFinalize                            // Finalize a block proposal
)

var _ servertypes.ABCI = (*BaseApp)(nil)

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	logger      log.Logger
	name        string                      // application name from abci.BlockInfo
	db          dbm.DB                      // common DB backend
	cms         storetypes.CommitMultiStore // Main (uncached) state
	qms         storetypes.MultiStore       // Optional alternative multistore for querying only.
	storeLoader StoreLoader                 // function to handle store loading, may be overridden with SetStoreLoader()

	initChainer        sdk.InitChainer            // ABCI InitChain handler
	preBlocker         sdk.PreBlocker             // logic to run before BeginBlocker
	beginBlocker       sdk.BeginBlocker           // (legacy ABCI) BeginBlock handler
	endBlocker         sdk.EndBlocker             // (legacy ABCI) EndBlock handler
	processProposal    sdk.ProcessProposalHandler // ABCI ProcessProposal handler
	prepareProposal    sdk.PrepareProposalHandler // ABCI PrepareProposal handler
	prepareCheckStater sdk.PrepareCheckStater     // logic to run during commit using the checkState
	precommiter        sdk.Precommiter            // logic to run during commit using the deliverState

	// volatile states:
	//
	// - checkState is set on InitChain and reset on Commit
	// - finalizeBlockState is set on InitChain and FinalizeBlock and set to nil
	// on Commit.
	//
	// - checkState: Used for CheckTx, which is set based on the previous
	// block's
	// state. This state is never committed.
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
	checkState           *state
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

// NewBaseApp returns a reference to an initialized BaseApp. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
func NewBaseApp(
	name string,
	logger log.Logger,
	db dbm.DB,
	options ...func(*BaseApp),
) *BaseApp {
	app := &BaseApp{
		logger: logger.With(log.ModuleKey, "baseapp"),
		name:   name,
		db:     db,
		cms: store.NewCommitMultiStore(
			db,
			logger,
			storemetrics.NewNoOpMetrics(),
		), // by default we use a no-op metric gather in store
		storeLoader: DefaultStoreLoader,
	}

	for _, option := range options {
		option(app)
	}

	if app.interBlockCache != nil {
		app.cms.SetInterBlockCache(app.interBlockCache)
	}

	return app
}

// Name returns the name of the BaseApp.
func (app *BaseApp) Name() string {
	return app.name
}

// AppVersion returns the application's protocol version.
func (app *BaseApp) AppVersion(ctx context.Context) (uint64, error) {
	if app.paramStore == nil {
		return 0, errors.New("app.paramStore is nil")
	}

	cp, err := app.paramStore.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get consensus params: %w", err)
	}
	if cp.Version == nil {
		return 0, nil
	}
	return cp.Version.App, nil
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...storetypes.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *storetypes.KVStoreKey:
			app.MountStore(key, storetypes.StoreTypeIAVL)
		case *storetypes.TransientStoreKey:
			app.MountStore(key, storetypes.StoreTypeTransient)

		case *storetypes.MemoryStoreKey:
			app.MountStore(key, storetypes.StoreTypeMemory)

		default:
			panic(fmt.Sprintf("Unrecognized store key type :%T", key))
		}
	}
}

// MountKVStores mounts all IAVL or DB stores to the provided keys in the
// BaseApp multistore.
func (app *BaseApp) MountKVStores(keys map[string]*storetypes.KVStoreKey) {
	for _, key := range keys {
		app.MountStore(key, storetypes.StoreTypeIAVL)
	}
}

// MountTransientStores mounts all transient stores to the provided keys in
// the BaseApp multistore.
func (app *BaseApp) MountTransientStores(
	keys map[string]*storetypes.TransientStoreKey,
) {
	for _, key := range keys {
		app.MountStore(key, storetypes.StoreTypeTransient)
	}
}

// MountMemoryStores mounts all in-memory KVStores with the BaseApp's internal
// commit multi-store.
func (app *BaseApp) MountMemoryStores(
	keys map[string]*storetypes.MemoryStoreKey,
) {
	skeys := maps.Keys(keys)
	sort.Strings(skeys)
	for _, key := range skeys {
		memKey := keys[key]
		app.MountStore(memKey, storetypes.StoreTypeMemory)
	}
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
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion() error {
	err := app.storeLoader(app.cms)
	if err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	return app.Init()
}

// DefaultStoreLoader will be used by default and loads the latest version.
func DefaultStoreLoader(ms storetypes.CommitMultiStore) error {
	return ms.LoadLatestVersion()
}

// CommitMultiStore returns the root multi-store.
// App constructor can use this to access the `cms`.
// UNSAFE: must not be used during the abci life cycle.
func (app *BaseApp) CommitMultiStore() storetypes.CommitMultiStore {
	return app.cms
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running baseapp.
func (app *BaseApp) LoadVersion(version int64) error {
	app.logger.Info(
		"NOTICE: this could take a long time to migrate IAVL store to fastnode if you enable Fast Node.\n",
	)
	err := app.cms.LoadVersion(version)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}

	return app.Init()
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() storetypes.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// ChainID returns the chainID of the app.
func (app *BaseApp) ChainID() string {
	return app.chainID
}

// Init initializes the app. It seals the app, preventing any
// further modifications. In addition, it validates the app against
// the earlier provided settings. Returns an error if validation fails.
// nil otherwise. Panics if the app is already sealed.
func (app *BaseApp) Init() error {
	if app.cms == nil {
		return errors.New("commit multi-store must not be nil")
	}

	// needed for the export command which inits from store but never calls
	// initchain
	app.setState(execModeCheck, cmtproto.Header{ChainID: app.chainID})

	return app.cms.GetPruning().Validate()
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
func (app *BaseApp) setState(mode execMode, h cmtproto.Header) {
	ms := app.cms.CacheMultiStore()
	headerInfo := header.Info{
		Height:  h.Height,
		Time:    h.Time,
		ChainID: h.ChainID,
		AppHash: h.AppHash,
	}
	baseState := &state{
		ms: ms,
		ctx: sdk.NewContext(ms, false, app.logger).
			WithBlockHeader(h).
			WithHeaderInfo(headerInfo),
	}

	switch mode {
	case execModeCheck:
		baseState.SetContext(
			baseState.Context().WithIsCheckTx(true),
		)
		app.checkState = baseState

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

// StoreConsensusParams sets the consensus parameters to the BaseApp's param
// store.
func (app *BaseApp) StoreConsensusParams(
	ctx context.Context,
	cp cmtproto.ConsensusParams,
) error {
	if app.paramStore == nil {
		return errors.New(
			"cannot store consensus params with no params store set",
		)
	}

	return app.paramStore.Set(ctx, cp)
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

func (app *BaseApp) preBlock(
	req *abci.FinalizeBlockRequest,
) error {
	if app.preBlocker != nil {
		ctx := app.finalizeBlockState.Context()
		if err := app.preBlocker(ctx, req); err != nil {
			return err
		}
		// ConsensusParams can change in preblocker, so we need to
		// write the consensus parameters in store to context
		ctx = ctx.WithConsensusParams(app.GetConsensusParams(ctx))
		app.finalizeBlockState.SetContext(ctx)
	}
	return nil
}

func (app *BaseApp) beginBlock(
	_ *abci.FinalizeBlockRequest,
) (sdk.BeginBlock, error) {
	if app.beginBlocker != nil {
		return app.beginBlocker(app.finalizeBlockState.Context())
	}

	return sdk.BeginBlock{}, nil
}

func (app *BaseApp) deliverTx(tx []byte) *abci.ExecTxResult {
	gInfo, result, err := app.runTx(execModeFinalize, tx)
	if err != nil {
		space, code, log := errorsmod.ABCIInfo(err, false)
		return &abci.ExecTxResult{
			Codespace: space,
			Code:      code,
			Log:       log,
			//#nosec:G701 // bet.
			GasWanted: int64(gInfo.GasWanted),
			//#nosec:G701 // bet.
			GasUsed: int64(gInfo.GasUsed),
		}
	}

	return &abci.ExecTxResult{
		//#nosec:G701 // bet.
		GasWanted: int64(gInfo.GasWanted),
		//#nosec:G701 // bet.
		GasUsed: int64(gInfo.GasUsed),
		Log:     result.Log,
		Data:    result.Data,
	}
}

// endBlock is an application-defined function that is called after transactions
// have been processed in FinalizeBlock.
func (app *BaseApp) endBlock(_ context.Context) (sdk.EndBlock, error) {
	if app.endBlocker != nil {
		return app.endBlocker(app.finalizeBlockState.Context())
	}

	return sdk.EndBlock{}, nil
}

// runTx processes a transaction within a given execution mode, encoded
// transaction bytes, and the decoded transaction itself. All state transitions
// occur through
// a cached Context depending on the mode provided. State only gets persisted
// if all messages get executed successfully and the execution mode is
// DeliverTx.
// Note, gas execution info is always returned. A reference to a Result is
// returned if the tx does not run out of gas and if all the messages are valid
// and execute successfully. An error is returned otherwise.
func (app *BaseApp) runTx(
	_ execMode,
	_ []byte,
) (gInfo sdk.GasInfo, result *sdk.Result, err error) {
	return sdk.GasInfo{
			GasUsed:   0,
			GasWanted: 0,
		}, nil, sdkerrors.ErrTxDecode.Wrap(
			errors.New("skip decoding").Error(),
		)
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
