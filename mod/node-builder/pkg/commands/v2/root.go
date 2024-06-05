// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/deposit"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/genesis"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/jwt"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/v2/client"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/commands/v2/cometbft"
	beaconconfig "github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-builder/pkg/node-builder/v2"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

// DefaultRootCommandSetup sets up the default commands for the root command.
func DefaultRootCommandSetup[T types.Tx](
	rootCmd *cobra.Command,
	mm *module.Manager,
	newApp nodebuilder.AppCreator[T],
	chainSpec primitives.ChainSpec,
) {
	// Add the ToS Flag to the root command.
	beaconconfig.AddToSFlag(rootCmd)

	// Add all the commands to the root command.
	rootCmd.AddCommand(
		// `comet`
		cometbft.Commands(),
		// `client`
		client.Commands(),
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
		// pruning.Cmd(newApp),
		// `rollback`
		// serverv2.NewRollbackCmd(newApp),
		// `snapshots`
		// snapshot.Cmd(newApp),
		// `start`
		StartCmd(newApp, beaconconfig.AddBeaconKitFlags),
		// `status`
		server.StatusCommand(),
		// `version`
		version.NewVersionCommand(),
	)
}
