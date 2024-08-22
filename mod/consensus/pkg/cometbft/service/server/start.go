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
//nolint:mnd // its okay.
package server

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	pruningtypes "cosmossdk.io/store/pruning/types"
	types "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/node"
	dbm "github.com/cosmos/cosmos-db"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
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
	FlagIndexEvents         = "index-events"
	FlagMinRetainBlocks     = "min-retain-blocks"
	FlagIAVLCacheSize       = "iavl-cache-size"
	FlagDisableIAVLFastNode = "iavl-disable-fastnode"
	FlagShutdownGrace       = "shutdown-grace"
)

// StartCmdOptions defines options that can be customized in
// `StartCmdWithOptions`,.
type StartCmdOptions[T types.Application] struct {
	// AddFlags allows adding custom flags to the start command.
	AddFlags func(cmd *cobra.Command)
}

// StartCmd runs the service passed in, either stand-alone or in-process with
// CometBFT.
func StartCmd[T types.Application](
	appCreator types.AppCreator[T],
) *cobra.Command {
	return StartCmdWithOptions(appCreator, StartCmdOptions[T]{})
}

// StartCmdWithOptions runs the service passed in, either stand-alone or
// in-process with
// CometBFT.
func StartCmdWithOptions[T types.Application](
	appCreator types.AppCreator[T],
	opts StartCmdOptions[T],
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with CometBFT in or out of process. By
default, the application will run with CometBFT in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent', and
'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 362880 states are kept, pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: 2 latest states will be kept; pruning at 10 block intervals.
custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.

The node may be started in a 'query only' mode where only the gRPC and JSON HTTP
API services are enabled via the 'grpc-only' flag. In this mode, CometBFT is
bypassed and can be used when legacy queries are needed after an on-chain upgrade
is performed. Note, when enabled, gRPC will also be automatically enabled.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := GetServerContextFromCmd(cmd)
			_, err := GetPruningOptionsFromFlags(serverCtx.Viper)
			if err != nil {
				return err
			}

			err = start(serverCtx, appCreator)

			serverCtx.Logger.Debug("received quit signal")
			//#nosec:G703 // its a bet.
			graceDuration, _ := cmd.Flags().GetDuration(FlagShutdownGrace)
			if graceDuration > 0 {
				serverCtx.Logger.Info(
					"graceful shutdown start",
					FlagShutdownGrace,
					graceDuration,
				)
				<-time.After(graceDuration)
				serverCtx.Logger.Info("graceful shutdown complete")
			}

			return err
		},
	}

	addStartNodeFlags(cmd, opts)
	return cmd
}

func start[T types.Application](
	svrCtx *Context,
	appCreator types.AppCreator[T],
) error {
	home := svrCtx.Config.RootDir
	db, err := OpenDB(home, dbm.PebbleDBBackend)
	if err != nil {
		return err
	}

	_ = appCreator(svrCtx.Logger, db, nil, svrCtx.Config, svrCtx.Viper)
	return nil
}

// GetGenDocProvider returns a function which returns the genesis doc from the
// genesis file.
func GetGenDocProvider(
	cfg *cmtcfg.Config,
) func() (node.ChecksummedGenesisDoc, error) {
	return func() (node.ChecksummedGenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(cfg.GenesisFile())
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}

		gen, err := appGenesis.ToGenesisDoc()
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}
		genbz, err := gen.AppState.MarshalJSON()
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}

		bz, err := json.Marshal(genbz)
		if err != nil {
			return node.ChecksummedGenesisDoc{
				Sha256Checksum: []byte{},
			}, err
		}
		sum := sha256.Sum256(bz)

		return node.ChecksummedGenesisDoc{
			GenesisDoc:     gen,
			Sha256Checksum: sum[:],
		}, nil
	}
}

// addStartNodeFlags should be added to any CLI commands that start the network.
func addStartNodeFlags[T types.Application](
	cmd *cobra.Command,
	opts StartCmdOptions[T],
) {
	cmd.Flags().String(flagAddress, "tcp://127.0.0.1:26658", "Listen address")
	cmd.Flags().
		String(flagTransport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().
		Uint64(FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().
		Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().
		String(FlagPruning, pruningtypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().
		Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().
		Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().
		Uint64(FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune CometBFT blocks")
	cmd.Flags().
		Bool(FlagDisableIAVLFastNode, false, "Disable fast node for IAVL tree")
	cmd.Flags().
		Duration(FlagShutdownGrace, 0*time.Second, "On Shutdown, duration to wait for resource clean up")

	// add support for all CometBFT-specific command line options
	cmtcmd.AddNodeFlags(cmd)

	if opts.AddFlags != nil {
		opts.AddFlags(cmd)
	}
}
