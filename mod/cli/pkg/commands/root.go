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
	"context"
	"os"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/client"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/cometbft"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/deposit"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/genesis"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/jwt"
	beaconconfig "github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	"github.com/berachain/beacon-kit/mod/primitives"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// DefaultRootCommandSetup sets up the default commands for the root command.
func DefaultRootCommandSetup[T servertypes.Application](
	rootCmd *cobra.Command,
	mm *module.Manager,
	newApp servertypes.AppCreator[T],
	chainSpec primitives.ChainSpec,
) {
	// Add the ToS Flag to the root command.
	beaconconfig.AddToSFlag(rootCmd)

	// Setup the custom start command options.
	startCmdOptions := server.StartCmdOptions[T]{
		AddFlags: beaconconfig.AddBeaconKitFlags,
		PostSetup: func(app T, svrCtx *server.Context, clientCtx sdkclient.Context, ctx context.Context, g *errgroup.Group) error {
			svrCtx.Logger = log.NewCustomLogger(zerolog.New(os.Stdout).With().Timestamp().Logger()).With("module", "HENLO OOGA")
			return nil
		},
	}

	// Add all the commands to the root command.
	rootCmd.AddCommand(
		// `comet`
		cometbft.Commands(newApp),
		// `client`
		client.Commands[T](),
		// `config`
		confixcmd.ConfigCommand(),
		// `init`
		genutilcli.InitCmd(mm),
		// `genesis`
		genesis.Commands(chainSpec),
		// `deposit`
		deposit.Commands(chainSpec),
		// `jwt`
		jwt.Commands(),
		// `keys`
		keys.Commands(),
		// `prune`
		pruning.Cmd(newApp),
		// `rollback`
		server.NewRollbackCmd(newApp),
		// `snapshots`
		snapshot.Cmd(newApp),
		// `start`
		server.StartCmdWithOptions(newApp, startCmdOptions),
		// `status`
		server.StatusCommand(),
		// `version`
		version.NewVersionCommand(),
	)
}
