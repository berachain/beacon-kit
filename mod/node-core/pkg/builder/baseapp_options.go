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

package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/store"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cast"
)

// This file contains Options that extend our default baseapp options to be
// called by cosmos when building the app.
// TODO: refactor into consensus_options for serverv2 migration.

// WithCometParamStore sets the param store to the comet consensus engine.
func WithCometParamStore(
	chainSpec common.ChainSpec,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetParamStore(comet.NewConsensusParamsStore(chainSpec))
	}
}

// WithPrepareProposal sets the prepare proposal handler to the baseapp.
func WithPrepareProposal(
	handler sdk.PrepareProposalHandler,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPrepareProposal(handler)
	}
}

// WithProcessProposal sets the process proposal handler to the baseapp.
func WithProcessProposal(
	handler sdk.ProcessProposalHandler,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetProcessProposal(handler)
	}
}

// WithPreBlocker sets the pre-blocker to the baseapp.
func WithPreBlocker(
	preBlocker sdk.PreBlocker,
) func(bApp *baseapp.BaseApp) {
	return func(bApp *baseapp.BaseApp) {
		bApp.SetPreBlocker(preBlocker)
	}
}

// DefaultBaseappOptions returns the default baseapp options provided by the Cosmos SDK
func DefaultBaseappOptions(appOpts servertypes.AppOptions) []func(*baseapp.BaseApp) {
	var cache storetypes.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		reader, err := os.Open(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}
		defer reader.Close()

		chainID, err = genutiltypes.ParseChainIDFromGenesis(reader)
		if err != nil {
			panic(fmt.Errorf("failed to parse chain-id from genesis file: %w", err))
		}
	}

	snapshotStore, err := server.GetSnapshotStore(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshottypes.NewSnapshotOptions(
		cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	defaultMempool := baseapp.SetMempool(mempool.NoOpMempool{})
	if maxTxs := cast.ToInt(appOpts.Get(server.FlagMempoolMaxTxs)); maxTxs >= 0 {
		defaultMempool = baseapp.SetMempool(
			mempool.NewSenderNonceMempool(
				mempool.SenderNonceMaxTxOpt(maxTxs),
			),
		)
	}

	return []func(*baseapp.BaseApp){
		baseapp.SetPruning(pruningOpts),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(server.FlagDisableIAVLFastNode))),
		defaultMempool,
		baseapp.SetChainID(chainID),
	}
}
