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

package commands

import (
	"os"
	"path/filepath"

	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/start"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	cliflags "github.com/berachain/beacon-kit/mod/cli/pkg/flags"
)

// DefaultCommandConfig adds a start command to the root command.
func DefaultCommandConfig[NodeT types.Node[T], T transaction.Tx](
	rootCmd *cobra.Command,
	appCreator serverv2.AppCreator[NodeT, T],
	logger log.Logger,
	bkCommands []*cobra.Command,
	server *serverv2.Server[NodeT, T],
) (serverv2.CLIConfig, error) {
	// TODO: this is weird, but is how cosmos does it. pls fix later
	flags := server.StartFlags()
	startCmd := start.NewStartCmd[NodeT, T](
		appCreator,
		server,
		flags,
	)
	cliflags.AddBeaconKitFlags(startCmd)
	cmds := server.CLICommands()
	cmds.Commands = append(cmds.Commands, bkCommands...)
	cmds.Commands = append(cmds.Commands, startCmd)

	return cmds, nil
}

// AddCommands adds the start command to the root command and sets the
// server context
func AddCommands[NodeT types.Node[T], T transaction.Tx](
	rootCmd *cobra.Command,
	newApp serverv2.AppCreator[NodeT, T],
	logger log.Logger,
	cmdConfig serverv2.CLIConfig,
	server *serverv2.Server[NodeT, T],
) error {
	originalPersistentPreRunE := rootCmd.PersistentPreRunE
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// set the default command outputs
		cmd.SetOut(cmd.OutOrStdout())
		cmd.SetErr(cmd.ErrOrStderr())

		if err := configHandle(server, cmd); err != nil {
			return err
		}

		if rootCmd.PersistentPreRun != nil {
			rootCmd.PersistentPreRun(cmd, args)
			return nil
		}

		return originalPersistentPreRunE(cmd, args)
	}

	rootCmd.AddCommand(cmdConfig.Commands...)
	return nil
}

// configHandle writes the default config to the home directory if it does not exist and sets the server context
func configHandle[NodeT types.Node[T], T transaction.Tx](
	s *serverv2.Server[NodeT, T],
	cmd *cobra.Command,
) error {
	home, err := cmd.Flags().GetString(serverv2.FlagHome)
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, "config")

	// we need to check app.toml as the config folder can already exist for the client.toml
	if _, err := os.Stat(filepath.Join(configDir, "app.toml")); os.IsNotExist(err) {
		if err = s.WriteConfig(configDir); err != nil {
			return err
		}
	}

	v, err := serverv2.ReadConfig(configDir)
	if err != nil {
		return err
	}

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	log, err := serverv2.NewLogger(v, cmd.OutOrStdout())
	if err != nil {
		return err
	}

	return serverv2.SetCmdServerContext(cmd, v, log)
}
