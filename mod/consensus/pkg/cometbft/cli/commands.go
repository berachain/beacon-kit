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

package server

import (
	"context"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server"
	types "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	cmtcfg "github.com/cometbft/cometbft/config"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	cmtversion "github.com/cometbft/cometbft/version"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// AddCommands add server commands.
func AddCommands[T types.Application](
	rootCmd *cobra.Command,
	appCreator types.AppCreator[T],
	opts server.StartCmdOptions[T],
) {
	cometCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		ShowNodeIDCmd(),
		ShowValidatorCmd(),
		ShowAddressCmd(),
		VersionCmd(),
		cmtcmd.ResetAllCmd,
		cmtcmd.ResetStateCmd,
		BootstrapStateCmd[T](appCreator),
	)

	startCmd := server.StartCmdWithOptions(appCreator, opts)
	rootCmd.AddCommand(
		startCmd,
		cometCmd,
		version.NewVersionCommand(),
		server.NewRollbackCmd[T](appCreator),
	)
}

// StatusCommand returns the command to return the status of the network.
func StatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			status, err := cmtservice.GetNodeStatus(
				context.Background(),
				clientCtx,
			)
			if err != nil {
				return err
			}

			output, err := cmtjson.Marshal(status)
			if err != nil {
				return err
			}

			// In order to maintain backwards compatibility, the default json
			// format output
			//#nosec:G703 // its a bet.
			outputFormat, _ := cmd.Flags().GetString(flags.FlagOutput)
			if outputFormat == flags.OutputFormatJSON {
				clientCtx = clientCtx.WithOutputFormat(flags.OutputFormatJSON)
			}

			return clientCtx.PrintRaw(output)
		},
	}

	cmd.Flags().
		StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().
		StringP(flags.FlagOutput, "o", "json", "Output format (text|json)")

	return cmd
}

// ShowNodeIDCmd - ported from CometBFT, dump node ID to stdout.
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := client.GetConfigFromCmd(cmd)
			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}

			cmd.Println(nodeKey.ID())
			return nil
		},
	}
}

// ShowValidatorCmd - ported from CometBFT, show this node's validator info.
func ShowValidatorCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show-validator",
		Short: "Show this node's CometBFT validator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := client.GetConfigFromCmd(cmd)
			privValidator := pvm.LoadFilePV(
				cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile(),
			)
			pk, err := privValidator.GetPubKey()
			if err != nil {
				return err
			}

			sdkPK, err := cryptocodec.FromCmtPubKeyInterface(pk)
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			bz, err := clientCtx.Codec.MarshalInterfaceJSON(sdkPK)
			if err != nil {
				return err
			}

			cmd.Println(string(bz))
			return nil
		},
	}

	return &cmd
}

// ShowAddressCmd - show this node's validator address.
func ShowAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-address",
		Short: "Shows this node's CometBFT validator consensus address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := client.GetConfigFromCmd(cmd)
			privValidator := pvm.LoadFilePV(
				cfg.PrivValidatorKeyFile(),
				cfg.PrivValidatorStateFile(),
			)

			valConsAddr := (sdk.ConsAddress)(privValidator.GetAddress())

			cmd.Println(valConsAddr.String())
			return nil
		},
	}

	return cmd
}

// VersionCmd prints CometBFT and ABCI version numbers.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CometBFT libraries' version",
		Long:  "Print protocols' and libraries' version numbers against which this app has been compiled.",
		RunE: func(cmd *cobra.Command, args []string) error {
			bs, err := yaml.Marshal(&struct {
				CometBFT      string
				ABCI          string
				BlockProtocol uint64
				P2PProtocol   uint64
			}{
				CometBFT:      cmtversion.CMTSemVer,
				ABCI:          cmtversion.ABCIVersion,
				BlockProtocol: cmtversion.BlockProtocol,
				P2PProtocol:   cmtversion.P2PProtocol,
			})
			if err != nil {
				return err
			}

			cmd.Println(string(bs))
			return nil
		},
	}
}

func BootstrapStateCmd[T types.Application](
	appCreator types.AppCreator[T],
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap-state",
		Short: "Bootstrap CometBFT state at an arbitrary block height using a light client",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.NewLogger(cmd.OutOrStdout())
			cfg := client.GetConfigFromCmd(cmd)
			v := client.GetViperFromCmd(cmd)

			height, err := cmd.Flags().GetInt64("height")
			if err != nil {
				return err
			}
			if height == 0 {
				home := v.GetString(flags.FlagHome)
				var db dbm.DB
				db, err = server.OpenDB(home, dbm.PebbleDBBackend)
				if err != nil {
					return err
				}

				app := appCreator(logger, db, nil, cfg, v)
				height = app.CommitMultiStore().LastCommitID().Version
			}

			return node.BootstrapState(
				cmd.Context(),
				cfg,
				cmtcfg.DefaultDBProvider,
				server.GetGenDocProvider(cfg),
				//#nosec:G701 // bet.
				uint64(height),
				nil,
			)
		},
	}

	cmd.Flags().
		Int64("height", 0, "Block height to bootstrap state at, if not provided it uses the latest block height in app state")

	return cmd
}
