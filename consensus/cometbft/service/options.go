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
	"time"

	pruningtypes "cosmossdk.io/store/pruning/types"
	storetypes "cosmossdk.io/store/types"
)

// File for storing in-package cometbft optional functions,
// for options that need access to non-exported fields of the Service

// SetPruning sets a pruning option on the multistore associated with the s.
func SetPruning(opts pruningtypes.PruningOptions) func(*Service) {
	return func(bs *Service) { bs.sm.GetCommitMultiStore().SetPruning(opts) }
}

// SetMinRetainBlocks returns a Service option function that sets the minimum
// block retention height value when determining which heights to prune during
// ABCI Commit.
func SetMinRetainBlocks(minRetainBlocks uint64) func(*Service) {
	return func(bs *Service) { bs.setMinRetainBlocks(minRetainBlocks) }
}

// SetIAVLCacheSize provides a Service option function that sets the size of
// IAVL cache.
func SetIAVLCacheSize(size int) func(*Service) {
	return func(bs *Service) {
		bs.sm.GetCommitMultiStore().SetIAVLCacheSize(size)
	}
}

// SetIAVLDisableFastNode enables(false)/disables(true) fast node usage from the
// IAVL store.
func SetIAVLDisableFastNode(disable bool) func(*Service) {
	return func(bs *Service) {
		bs.sm.GetCommitMultiStore().SetIAVLDisableFastNode(disable)
	}
}

// SetInterBlockCache provides a Service option function that sets the
// inter-block cache.
func SetInterBlockCache(cache storetypes.MultiStorePersistentCache) func(*Service) {
	return func(s *Service) {
		s.setInterBlockCache(cache)
	}
}

// SetChainID sets the chain ID in cometbft.
func SetChainID(chainID string) func(*Service) {
	return func(s *Service) { s.chainID = chainID }
}

// SetInterBlockCache provides a Service option function that sets the stable
// block time upgrade height and, optionally, time if the upgrade happened in
// the past.
//
// If the network starts from genesis, you don't need to set this option.
func SetSBTUpgradeHeightAndTime(height int64, time time.Time) func(*Service) {
	return func(bs *Service) { bs.setSBTUpgradeHeightAndTime(height, time) }
}
