// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package commands

import (
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/client"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/cometbft"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/deposit"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/genesis"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/jwt"
	beaconconfig "github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
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
