// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/validator"
	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	statem "github.com/berachain/beacon-kit/consensus/cometbft/service/state"
	errorsmod "github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/storage"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	initialAppVersion uint64 = 0
	AppName           string = "beacond"
)

type Service struct {
	node *node.Node

	// cmtConsensusParams are part of the blockchain state and
	// are agreed upon by all validators in the network.
	cmtConsensusParams *cmttypes.ConsensusParams

	// cmtCgf are node-specific settings that influence how
	// the consensus engine operates on a particular node.
	// Loaded from config file (config.toml), not part of state.
	cmtCfg *cmtcfg.Config

	telemetrySink TelemetrySink

	logger       *phuslu.Logger
	sm           *statem.Manager
	Blockchain   blockchain.BlockchainI
	BlockBuilder validator.BlockBuilderI

	// prepareProposalState is used for PrepareProposal, which is set based on
	// the previous block's state. This state is never committed. In case of
	// multiple consensus rounds, the state is always reset to the previous
	// block's state.
	prepareProposalState *state

	// processProposalState is used for ProcessProposal, which is set based on
	// the previous block's state. This state is never committed. In case of
	// multiple consensus rounds, the state is always reset to the previous
	// block's state.
	processProposalState *state

	// finalizeBlockState is used for FinalizeBlock, which is set based on the
	// previous block's state. This state is committed. finalizeBlockState is
	// set
	// on InitChain and FinalizeBlock and set to nil on Commit.
	finalizeBlockState *state

	interBlockCache storetypes.MultiStorePersistentCache

	// initialHeight is the initial height at which we start the node
	initialHeight   int64
	minRetainBlocks uint64

	chainID string

	// ctx is a cancellable context based on the parent ctx passed in Start().
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewService(
	logger *phuslu.Logger,
	db dbm.DB,
	blockchain blockchain.BlockchainI,
	blockBuilder validator.BlockBuilderI,
	cmtCfg *cmtcfg.Config,
	telemetrySink TelemetrySink,
	options ...func(*Service),
) *Service {
	if err := validateConfig(cmtCfg); err != nil {
		panic(err)
	}
	cmtConsensusParams, err := extractConsensusParams(cmtCfg)
	if err != nil {
		panic(err)
	}

	s := &Service{
		logger: logger,
		sm: statem.NewManager(
			db,
			servercmtlog.WrapSDKLogger(logger),
		),
		Blockchain:         blockchain,
		BlockBuilder:       blockBuilder,
		cmtConsensusParams: cmtConsensusParams,
		cmtCfg:             cmtCfg,
		telemetrySink:      telemetrySink,
	}

	s.MountStore(storage.StoreKey, storetypes.StoreTypeIAVL)

	for _, option := range options {
		option(s)
	}

	if s.interBlockCache != nil {
		s.sm.GetCommitMultiStore().SetInterBlockCache(s.interBlockCache)
	}

	// Load latest height, once all stores have been set
	if err = s.sm.LoadLatestVersion(); err != nil {
		panic(fmt.Errorf("failed loading latest version: %w", err))
	}

	return s
}

// TODO: Move nodeKey into being created within the function.
func (s *Service) Start(
	ctx context.Context,
) error {
	cfg := s.cmtCfg
	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return err
	}

	privVal, err := pvm.LoadOrGenFilePV(
		cfg.PrivValidatorKeyFile(),
		cfg.PrivValidatorStateFile(),
		nil,
	)
	if err != nil {
		return err
	}

	s.ctx = ctx
	s.node, err = node.NewNode(
		ctx,
		cfg,
		privVal,
		nodeKey,
		proxy.NewLocalClientCreator(s),
		GetGenDocProvider(cfg),
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		servercmtlog.WrapCometLogger(s.logger),
	)
	if err != nil {
		return err
	}

	started := make(chan struct{})

	// we start the node in a goroutine since calling Start() can block if genesis
	// time is in the future causing us not to handle signals gracefully.
	go func() {
		err = s.node.Start()
		started <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		s.logger.Warn("DEBUG: Start() ctx.Done()")
		return ctx.Err()
	case <-s.ctx.Done():
		s.logger.Warn("DEBUG: Start() s.ctx.Done()")
		return s.ctx.Err()
	case <-started:
	}

	close(started)

	return err
}

func (s *Service) Stop() error {
	var errs []error

	if s.node != nil && s.node.IsRunning() {
		s.logger.Info("Stopping CometBFT Node")
		err := s.node.Stop()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to stop CometBFT Node: %w", err))
		}
		s.logger.Info("Waiting for CometBFT Node to stop")
		s.node.Wait()
	}

	// cancel the service context
	if s.cancelFunc != nil {
		s.logger.Warn("DEBUG: Stop() is cancelling the service context")
		//s.cancelFunc()
	}

	s.logger.Info("Closing application.db")
	if err := s.sm.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close application.id: %w", err))
	}
	return errors.Join(errs...)
}

// Name returns the name of the cometbft.
func (s *Service) Name() string {
	return AppName
}

// CommitMultiStore returns the CommitMultiStore of the cometbft.
// This needs to be exposed because it is used by commands like Rollback.
func (s *Service) CommitMultiStore() storetypes.CommitMultiStore {
	return s.sm.GetCommitMultiStore()
}

// AppVersion returns the application's protocol version.
func (s *Service) AppVersion(_ context.Context) (uint64, error) {
	return s.appVersion()
}

func (s *Service) appVersion() (uint64, error) {
	cp := s.cmtConsensusParams.ToProto()
	return cp.Version.App, nil
}

// MountStore mounts a store to the provided key in the Service multistore,
// using the default DB.
func (s *Service) MountStore(
	key storetypes.StoreKey,
	typ storetypes.StoreType,
) {
	s.sm.GetCommitMultiStore().MountStoreWithDB(key, typ, nil)
}

// LastBlockHeight returns the last committed block height.
func (s *Service) LastBlockHeight() int64 {
	return s.sm.GetCommitMultiStore().LastCommitID().Version
}

func (s *Service) setMinRetainBlocks(minRetainBlocks uint64) {
	s.minRetainBlocks = minRetainBlocks
}

func (s *Service) setInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) {
	s.interBlockCache = cache
}

// resetState provides a fresh state which can be used to reset
// prepareProposal/processProposal/finalizeBlock State.
// A state is explicitly returned to avoid false positives from
// nilaway tool.
func (s *Service) resetState(ctx context.Context) *state {
	ms := s.sm.GetCommitMultiStore().CacheMultiStore()

	newCtx := sdk.NewContext(
		ms,
		false,
		servercmtlog.WrapSDKLogger(s.logger),
	).WithContext(ctx)

	return &state{
		ms:  ms,
		ctx: newCtx,
	}
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
	//nolint:errcheck // should be safe
	return any(abci.ValidatorUpdate{
		PubKeyBytes: update.Pubkey[:],
		PubKeyType:  crypto.CometBLSType,
		Power:       int64(update.EffectiveBalance.Unwrap()), // #nosec G115 -- this is safe.
	}).(ValidatorUpdateT), nil
}

// getContextForProposal returns the correct Context for PrepareProposal and
// ProcessProposal. We use finalizeBlockState on the first block to be able to
// access any state changes made in InitChain.
func (s *Service) getContextForProposal(
	ctx sdk.Context,
	height int64,
) sdk.Context {
	if height != s.initialHeight {
		return ctx
	}

	if s.finalizeBlockState == nil {
		// this is unexpected since cometBFT won't call PrepareProposal
		// on initialHeight. Panic appeases nilaway.
		panic(fmt.Errorf("getContextForProposal: %w", errNilFinalizeBlockState))
	}
	ctx, _ = s.finalizeBlockState.Context().CacheContext()
	return ctx
}

// CreateQueryContext creates a new sdk.Context for a query, taking as args
// the block height and whether the query needs a proof or not.
func (s *Service) CreateQueryContext(
	height int64,
	prove bool,
) (sdk.Context, error) {
	// use custom query multi-store if provided
	lastBlockHeight := s.sm.GetCommitMultiStore().LatestVersion()
	if lastBlockHeight == 0 {
		return sdk.Context{}, errorsmod.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"%s is not ready; please wait for first block",
			AppName,
		)
	}

	if height > lastBlockHeight {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidHeight,
				"cannot query with height in the future; please provide a valid height",
			)
	}

	// when a client did not provide a query height, manually inject the latest
	if height == 0 {
		height = lastBlockHeight
	}

	if height <= 1 && prove {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			)
	}

	cacheMS, err := s.sm.GetCommitMultiStore().CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{},
			errorsmod.Wrapf(
				sdkerrors.ErrNotFound,
				"failed to load state at height %d; %s (latest height: %d)",
				height,
				err,
				lastBlockHeight,
			)
	}

	return sdk.NewContext(
		cacheMS,
		true,
		servercmtlog.WrapSDKLogger(s.logger),
	), nil
}
