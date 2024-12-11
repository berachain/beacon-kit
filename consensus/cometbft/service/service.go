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

	storetypes "cosmossdk.io/store/types"
	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/params"
	statem "github.com/berachain/beacon-kit/consensus/cometbft/service/state"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/transition"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	initialAppVersion uint64 = 0
	appName           string = "beacond"
)

type Service[
	LoggerT log.AdvancedLogger[LoggerT],
] struct {
	node   *node.Node
	cmtCfg *cmtcfg.Config

	logger     LoggerT
	sm         *statem.Manager
	Middleware MiddlewareI

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
	paramStore      *params.ConsensusParamsStore

	// initialHeight is the initial height at which we start the node
	initialHeight   int64
	minRetainBlocks uint64

	chainID string
}

func NewService[
	LoggerT log.AdvancedLogger[LoggerT],
](
	storeKey *storetypes.KVStoreKey,
	logger LoggerT,
	db dbm.DB,
	middleware MiddlewareI,
	cmtCfg *cmtcfg.Config,
	cs common.ChainSpec,
	options ...func(*Service[LoggerT]),
) *Service[LoggerT] {
	s := &Service[LoggerT]{
		logger: logger,
		sm: statem.NewManager(
			db,
			servercmtlog.WrapSDKLogger(logger),
		),
		Middleware: middleware,
		cmtCfg:     cmtCfg,
		paramStore: params.NewConsensusParamsStore(cs),
	}

	s.MountStore(storeKey, storetypes.StoreTypeIAVL)

	for _, option := range options {
		option(s)
	}

	if s.interBlockCache != nil {
		s.sm.CommitMultiStore().SetInterBlockCache(s.interBlockCache)
	}

	// Load latest height, once all stores have been set
	if err := s.sm.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return s
}

// TODO: Move nodeKey into being created within the function.
func (s *Service[_]) Start(
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

	return s.node.Start()
}

// Close is called in start cmd to gracefully cleanup resources.
func (s *Service[_]) Close() error {
	var errs []error

	if s.node != nil && s.node.IsRunning() {
		s.logger.Info("Stopping CometBFT Node")
		//#nosec:G703 // its a bet.
		_ = s.node.Stop()
	}

	s.logger.Info("Closing application.db")
	if err := s.sm.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// Name returns the name of the cometbft.
func (s *Service[_]) Name() string {
	return appName
}

// CommitMultiStore returns the CommitMultiStore of the cometbft.
func (s *Service[_]) CommitMultiStore() storetypes.CommitMultiStore {
	return s.sm.CommitMultiStore()
}

// AppVersion returns the application's protocol version.
func (s *Service[_]) AppVersion(_ context.Context) (uint64, error) {
	return s.appVersion()
}

func (s *Service[_]) appVersion() (uint64, error) {
	cp := s.paramStore.Get()
	return cp.Version.App, nil
}

// MountStore mounts a store to the provided key in the Service multistore,
// using the default DB.
func (s *Service[_]) MountStore(
	key storetypes.StoreKey,
	typ storetypes.StoreType,
) {
	s.sm.CommitMultiStore().MountStoreWithDB(key, typ, nil)
}

// LastBlockHeight returns the last committed block height.
func (s *Service[_]) LastBlockHeight() int64 {
	return s.sm.CommitMultiStore().LastCommitID().Version
}

func (s *Service[_]) setMinRetainBlocks(minRetainBlocks uint64) {
	s.minRetainBlocks = minRetainBlocks
}

func (s *Service[_]) setInterBlockCache(
	cache storetypes.MultiStorePersistentCache,
) {
	s.interBlockCache = cache
}

// resetState provides a fresh state which can be used to reset
// prepareProposal/processProposal/finalizeBlock State.
// A state is explicitly returned to avoid false positives from
// nilaway tool.
func (s *Service[LoggerT]) resetState() *state {
	ms := s.sm.CommitMultiStore().CacheMultiStore()
	return &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, false, servercmtlog.WrapSDKLogger(s.logger)),
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
		//#nosec:G701 // this is safe.
		Power: int64(update.EffectiveBalance.Unwrap()),
	}).(ValidatorUpdateT), nil
}
