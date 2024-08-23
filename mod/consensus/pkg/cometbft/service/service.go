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
	s := &Service{
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

	s.SetVersion(version.Version)
	s.MountStore(storeKey, storetypes.StoreTypeIAVL)

	for _, option := range options {
		option(s)
	}

	if s.interBlockCache != nil {
		s.cms.SetInterBlockCache(s.interBlockCache)
	}

	// Load the s.
	if loadLatest {
		if err := s.LoadLatestVersion(); err != nil {
			panic(err)
		}
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

	s.node, err = node.NewNode(
		ctx,
		cfg,
		pvm.LoadOrGenFilePV(
			cfg.PrivValidatorKeyFile(),
			cfg.PrivValidatorStateFile(),
		),
		nodeKey,
		proxy.NewLocalClientCreator(s),
		GetGenDocProvider(cfg),
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		servercmtlog.CometLoggerWrapper{Logger: s.logger},
	)
	if err != nil {
		return err
	}

	return s.node.Start()
}

// Close is called in start cmd to gracefully cleanup resources.
func (s *Service) Close() error {
	var errs []error

	if s.node != nil && s.node.IsRunning() {
		s.logger.Info("Stopping CometBFT Node")
		//#nosec:G703 // its a bet.
		_ = s.node.Stop()
	}

	// Close s.db (opened by cosmos-sdk/server/start.go call to openDB)
	if s.db != nil {
		s.logger.Info("Closing application.db")
		if err := s.db.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Name returns the name of the cometbft.
func (s *Service) Name() string {
	return s.name
}

// CommitMultiStore returns the CommitMultiStore of the cometbft.
func (s *Service) CommitMultiStore() storetypes.CommitMultiStore {
	return s.cms
}

// AppVersion returns the application's protocol version.
func (s *Service) AppVersion(ctx context.Context) (uint64, error) {
	cp, err := s.paramStore.Get(ctx)
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
func (s *Service) MountStore(
	key storetypes.StoreKey,
	typ storetypes.StoreType,
) {
	s.cms.MountStoreWithDB(key, typ, nil)
}

func (s *Service) LoadLatestVersion() error {
	if err := s.cms.LoadLatestVersion(); err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	// Validator pruning settings.
	return s.cms.GetPruning().Validate()
}

func (s *Service) LoadVersion(version int64) error {
	err := s.cms.LoadVersion(version)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}

	// Validate Pruning settings.
	return s.cms.GetPruning().Validate()
}

// LastCommitID returns the last CommitID of the multistore.
func (s *Service) LastCommitID() storetypes.CommitID {
	return s.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (s *Service) LastBlockHeight() int64 {
	return s.cms.LastCommitID().Version
}

func (s *Service) setMinRetainBlocks(minRetainBlocks uint64) {
	s.minRetainBlocks = minRetainBlocks
}

func (s *Service) setInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) {
	s.interBlockCache = cache
}

func (s *Service) setState(mode execMode) {
	ms := s.cms.CacheMultiStore()
	baseState := &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, false, s.logger),
	}

	switch mode {
	case execModePrepareProposal:
		s.prepareProposalState = baseState

	case execModeProcessProposal:
		s.processProposalState = baseState

	case execModeFinalize:
		s.finalizeBlockState = baseState

	default:
		panic(fmt.Sprintf("invalid runTxMode for setState: %d", mode))
	}
}

// GetConsensusParams returns the current consensus parameters from the
// Service's
// ParamStore. If the Service has no ParamStore defined, nil is returned.
func (s *Service) GetConsensusParams(
	ctx context.Context,
) cmtproto.ConsensusParams {
	//#nosec:G703 // bet.
	cp, _ := s.paramStore.Get(ctx)
	return cp
}

func (s *Service) validateFinalizeBlockHeight(
	req *abci.FinalizeBlockRequest,
) error {
	if req.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Height)
	}

	lastBlockHeight := s.LastBlockHeight()

	// expectedHeight holds the expected height to validate
	var expectedHeight int64
	if lastBlockHeight == 0 && s.initialHeight > 1 {
		// In this case, we're validating the first block of the chain, i.e no
		// previous commit. The height we're expecting is the initial height.
		expectedHeight = s.initialHeight
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
