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
//

package cometbft

import (
	"fmt"
	"time"

	cmtcfg "github.com/cometbft/cometbft/config"
	cmttypes "github.com/cometbft/cometbft/types"
)

const ( // appeases mnd
	timeoutPropose   = 1000 * time.Millisecond
	timeoutPrevote   = 1000 * time.Millisecond
	timeoutPrecommit = 1000 * time.Millisecond
	timeoutCommit    = 500 * time.Millisecond

	maxBlockSize = 100 * 1024 * 1024

	precision    = 505 * time.Millisecond
	messageDelay = 15 * time.Second
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
	cfg.P2P.MaxNumInboundPeers = 120
	cfg.P2P.MaxNumOutboundPeers = 40

	cfg.Mempool.Type = "nop"
	cfg.Mempool.Recheck = false
	cfg.Mempool.RecheckTimeout = 0
	cfg.Mempool.Broadcast = false
	cfg.Mempool.Size = 0
	cfg.Mempool.MaxTxBytes = 0
	cfg.Mempool.MaxTxsBytes = 0
	cfg.Mempool.CacheSize = 0

	consensus := cfg.Consensus
	consensus.TimeoutPropose = timeoutPropose
	consensus.TimeoutPrevote = timeoutPrevote
	consensus.TimeoutPrecommit = timeoutPrecommit
	consensus.TimeoutCommit = timeoutCommit

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

// extractConsensusParams pull consensus parameters (not config) set in
// genesis. They are mostly used to (not) update consensus parameters once
// a block is finalized.
func extractConsensusParams(cmtCfg *cmtcfg.Config) (*cmttypes.ConsensusParams, error) {
	genFunc := GetGenDocProvider(cmtCfg)
	genDoc, err := genFunc()
	if err != nil {
		return nil, err
	}

	// Todo: add validation for genesis params by chainID
	cmtConsensusParams := genDoc.GenesisDoc.ConsensusParams
	return cmtConsensusParams, nil
}
