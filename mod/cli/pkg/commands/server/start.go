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

package server

import (
	"context"

	pruningtypes "cosmossdk.io/store/pruning/types"
	types "github.com/berachain/beacon-kit/mod/cli/pkg/commands/server/types"
	"github.com/berachain/beacon-kit/mod/storage/pkg/db"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

const (
	// CometBFT full-node start flags.
	flagAddress         = "address"
	flagTransport       = "transport"
	FlagHaltHeight      = "halt-height"
	FlagHaltTime        = "halt-time"
	FlagInterBlockCache = "inter-block-cache"

	FlagPruning             = "pruning"
	FlagPruningKeepRecent   = "pruning-keep-recent"
	FlagPruningInterval     = "pruning-interval"
	FlagMinRetainBlocks     = "min-retain-blocks"
	FlagIAVLCacheSize       = "iavl-cache-size"
	FlagDisableIAVLFastNode = "iavl-disable-fastnode"
)

// StartCmdOptions defines options that can be customized in
// `StartCmdWithOptions`,.
type StartCmdOptions[T interface {
	Start(context.Context) error
}] struct {
	// AddFlags allows adding custom flags to the start command.
	AddFlags func(cmd *cobra.Command)
}

// StartCmd runs the service passed in, either stand-alone or in-process with
// CometBFT.
func StartCmd[T interface {
	Start(context.Context) error
}](
	appCreator types.AppCreator[T],
) *cobra.Command {
	return StartCmdWithOptions(appCreator, StartCmdOptions[T]{})
}

// StartCmdWithOptions runs the service passed in, either stand-alone or
// in-process with
// CometBFT.
func StartCmdWithOptions[T interface {
	Start(context.Context) error
}](
	appCreator types.AppCreator[T],
	opts StartCmdOptions[T],
) *cobra.Command {
	//nolint:lll // its okay.
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with CometBFT in process. By
default, the application will run with CometBFT in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent', and
'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 362880 states are kept, pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: 2 latest states will be kept; pruning at 10 block intervals.
custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'

`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := client.GetLoggerFromCmd(cmd)
			cfg := client.GetConfigFromCmd(cmd)
			v := client.GetViperFromCmd(cmd)
			_, err := GetPruningOptionsFromFlags(v)
			if err != nil {
				return err
			}

			// Open the Database
			db, err := db.OpenDB(cfg.RootDir, dbm.PebbleDBBackend)
			if err != nil {
				return err
			}

			// Create the application.
			return appCreator(logger, db, nil, cfg, v).
				Start(cmd.Context())
		},
	}

	addStartNodeFlags(cmd, opts)
	return cmd
}

// addStartNodeFlags should be added to any CLI commands that start the network.
//
//nolint:lll // todo fix.
func addStartNodeFlags[T interface {
	Start(context.Context) error
}](
	cmd *cobra.Command,
	opts StartCmdOptions[T],
) {
	cmd.Flags().String(
		flagAddress, "tcp://127.0.0.1:26658", "Listen address")
	cmd.Flags().
		String(
			flagTransport,
			"socket",
			"Transport protocol: socket, grpc")
	cmd.Flags().
		Uint64(
			FlagHaltHeight,
			0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().
		Uint64(
			FlagHaltTime,
			0,
			"Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(
		FlagInterBlockCache,
		true,
		"Enable inter-block caching")
	cmd.Flags().
		String(
			FlagPruning,
			pruningtypes.PruningOptionDefault,
			"Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().
		Uint64(
			FlagPruningKeepRecent,
			0,
			"Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().
		Uint64(FlagPruningInterval,
			0,
			"Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().
		Uint64(
			FlagMinRetainBlocks,
			0,
			"Minimum block height offset during ABCI commit to prune CometBFT blocks")
	cmd.Flags().
		Bool(FlagDisableIAVLFastNode, false, "Disable fast node for IAVL tree")

	// add support for all CometBFT-specific command line options
	cmtcmd.AddNodeFlags(cmd)

	if opts.AddFlags != nil {
		opts.AddFlags(cmd)
	}
}
