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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	server "github.com/berachain/beacon-kit/cli/commands/server"
	"github.com/berachain/beacon-kit/config"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log"
	"github.com/cosmos/cosmos-sdk/client/flags"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cast"
)

// This file contains Options that extend our default Service options to be
// called by cosmos when building the app.
// TODO: refactor into consensus_options for serverv2 migration.

// DefaultServiceOptions returns the default Service options provided by the
// Cosmos SDK.
func DefaultServiceOptions[
	LoggerT log.AdvancedLogger[LoggerT],
](
	appOpts config.AppOptions,
) []func(*cometbft.Service[LoggerT]) {
	var cache storetypes.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	// get chainID, possibly falling back to genesis if flag is not set
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		chainID, err = loadChainIDFromGenesis(appOpts)
		if err != nil {
			panic(err)
		}
	}

	return []func(*cometbft.Service[LoggerT]){
		cometbft.SetPruning[LoggerT](pruningOpts),
		cometbft.SetMinRetainBlocks[LoggerT](
			cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks)),
		),
		cometbft.SetInterBlockCache[LoggerT](cache),
		cometbft.SetIAVLCacheSize[LoggerT](
			cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize)),
		),
		cometbft.SetIAVLDisableFastNode[LoggerT](
			// default to true
			true,
		),
		cometbft.SetChainID[LoggerT](chainID),
	}
}

func loadChainIDFromGenesis(appOpts config.AppOptions) (string, error) {
	var (
		homeDir = cast.ToString(appOpts.Get(flags.FlagHome))
		fp      = filepath.Join(homeDir, "config", "genesis.json")
	)

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return "", err
	}

	chainID, err := genutiltypes.ParseChainIDFromGenesis(f)
	if err != nil {
		return "",
			errors.Join(
				f.Close(),
				fmt.Errorf(
					"failed to parse chain-id from genesis file: %w",
					err,
				),
			)
	}
	return chainID, f.Close()
}
