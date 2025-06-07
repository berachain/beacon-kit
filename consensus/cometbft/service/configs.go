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
//

package cometbft

import (
	"errors"
	"fmt"
	"time"

	"github.com/berachain/beacon-kit/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmttypes "github.com/cometbft/cometbft/types"
)

const ( // appeases mnd
	// These timeouts are the ones we tested are necessary
	// at minimum to have a smooth network. We enforce that
	// these minima are respected.
	minTimeoutPropose   = 2000 * time.Millisecond
	minTimeoutPrevote   = 2000 * time.Millisecond
	minTimeoutPrecommit = 2000 * time.Millisecond
	minTimeoutCommit    = 500 * time.Millisecond

	maxBlockSize = 100 * 1024 * 1024

	precision    = 505 * time.Millisecond
	messageDelay = 15 * time.Second

	defaultMaxNumInboundPeers  = 40
	defaultMaxNumOutboundPeers = 10
)

var (
	ErrInvalidaConfig          = errors.New("invalid comet config for BeaconKit")
	ErrInvalidaConsensusParams = errors.New("invalid comet consensus params for BeaconKit")
)

// DefaultConfig returns the default configuration for the CometBFT
// consensus engine. It overrides a few values based on our own measurements
// and development level of BeaconKit. Recall that these are node-specific
// values (although they influence consensus).
// This should be the only place in the entire BeaconKit codebase where
// cmtcfg.DefaultConfig() is called.
func DefaultConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// BeaconKit forces PebbleDB as the database backend.
	cfg.BaseConfig.DBBackend = "pebbledb"

	// These settings are set by default for performance reasons.
	cfg.P2P.MaxNumInboundPeers = defaultMaxNumInboundPeers
	cfg.P2P.MaxNumOutboundPeers = defaultMaxNumOutboundPeers

	cfg.Mempool.Type = "nop"
	cfg.Mempool.Recheck = false
	cfg.Mempool.RecheckTimeout = 0
	cfg.Mempool.Broadcast = false
	cfg.Mempool.Size = 0
	cfg.Mempool.MaxTxBytes = 0
	cfg.Mempool.MaxTxsBytes = 0
	cfg.Mempool.CacheSize = 0

	// By default, we set timeouts to the minima we tested
	consensus := cfg.Consensus
	consensus.TimeoutPropose = minTimeoutPropose
	consensus.TimeoutPrevote = minTimeoutPrevote
	consensus.TimeoutPrecommit = minTimeoutPrecommit
	consensus.TimeoutCommit = minTimeoutCommit

	cfg.Storage.DiscardABCIResponses = true

	cfg.TxIndex.Indexer = "null"

	cfg.Instrumentation.Prometheus = true
	cfg.Instrumentation.MaxOpenConnections = 800

	// Disable profiling by default
	// cfg.RPC.PprofListenAddress = "localhost:6060"

	if err := cfg.ValidateBasic(); err != nil {
		panic(fmt.Errorf("invalid comet config: %w", err))
	}

	return cfg
}

// DefaultConsensusParams returns the default consensus parameters
// shared by every node in the network. Consensus parameters are
// inscripted in genesis.
func DefaultConsensusParams(consensusKeyAlgo string) *cmttypes.ConsensusParams {
	res := cmttypes.DefaultConsensusParams()
	res.Validator.PubKeyTypes = []string{consensusKeyAlgo}

	// set max block size in order to accommodate max blobs size
	// This matches current cmttypes.MaxBlockSizeBytes but it's
	// explicitly hard coded for safety across deps upgrades.
	res.Block.MaxBytes = maxBlockSize

	// activate pbst and hard code values to
	// be safe across dependencies upgrades
	res.Feature.PbtsEnableHeight = 1
	res.Synchrony.Precision = precision
	res.Synchrony.MessageDelay = messageDelay

	if err := res.ValidateBasic(); err != nil {
		panic(fmt.Errorf("invalid default consensus parameters: %w", err))
	}

	return res
}

func validateConfig(cfg *cmtcfg.Config) error {
	if cfg.Consensus.TimeoutPropose < minTimeoutPropose {
		return fmt.Errorf("%w, config timeout propose %v, min requested %v",
			ErrInvalidaConfig,
			cfg.Consensus.TimeoutPropose,
			minTimeoutPropose,
		)
	}

	if cfg.Consensus.TimeoutPrevote < minTimeoutPrevote {
		return fmt.Errorf("%w, config timeout prevote %v, min requested %v",
			ErrInvalidaConfig,
			cfg.Consensus.TimeoutPrevote,
			minTimeoutPrevote,
		)
	}

	if cfg.Consensus.TimeoutPrecommit < minTimeoutPrecommit {
		return fmt.Errorf("%w, config timeout propose %v, min requested %v",
			ErrInvalidaConfig,
			cfg.Consensus.TimeoutPrecommit,
			minTimeoutPrecommit,
		)
	}

	if cfg.Consensus.TimeoutCommit < minTimeoutCommit {
		return fmt.Errorf("%w, config timeout propose %v, min requested %v",
			ErrInvalidaConfig,
			cfg.Consensus.TimeoutCommit,
			minTimeoutCommit,
		)
	}

	return nil
}

func warnAboutConfigs(
	cmtCfg *cmtcfg.Config,
	logger log.Logger,
) {
	connectionsCap := defaultMaxNumInboundPeers + defaultMaxNumOutboundPeers
	if cmtCfg.P2P.MaxNumInboundPeers+cmtCfg.P2P.MaxNumOutboundPeers > connectionsCap {
		logger.Warn(
			"excessive peering",
			"max_num_inbound_peers", cmtCfg.P2P.MaxNumInboundPeers,
			"recommended max_num_inbound_peers", defaultMaxNumInboundPeers,
			"max_num_outbound_peers", cmtCfg.P2P.MaxNumOutboundPeers,
			"recommended max_num_outbound_peers", defaultMaxNumOutboundPeers,
			"recommended connections cap (inbound + outbound)", connectionsCap,
		)
	}
}

// extractConsensusParams pull consensus parameters (not config) set in
// genesis. They are mostly used to (not) update consensus parameters once
// a block is finalized.
func extractConsensusParams(cmtCfg *cmtcfg.Config) (*cmttypes.ConsensusParams, error) {
	genFunc := GetGenDocProvider(cmtCfg)
	genDoc, err := genFunc()
	if err != nil {
		return nil, fmt.Errorf("failed pulling consensus params: %w", err)
	}

	// Todo: add validation for genesis params by chainID
	cmtConsensusParams := genDoc.GenesisDoc.ConsensusParams
	return cmtConsensusParams, validateConsensusParams(cmtConsensusParams)
}

func validateConsensusParams(params *cmttypes.ConsensusParams) error {
	if params.Block.MaxBytes < maxBlockSize {
		return fmt.Errorf("%w, param max size %v, requested size %v",
			ErrInvalidaConsensusParams,
			params.Block.MaxBytes,
			maxBlockSize,
		)
	}
	return nil
}
