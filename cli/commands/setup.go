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

package commands

import (
	"github.com/berachain/beacon-kit/cli/commands/deposit"
	"github.com/berachain/beacon-kit/cli/commands/genesis"
	"github.com/berachain/beacon-kit/cli/commands/initialize"
	"github.com/berachain/beacon-kit/cli/commands/jwt"
	"github.com/berachain/beacon-kit/cli/commands/server"
	servertypes "github.com/berachain/beacon-kit/cli/commands/server/types"
	"github.com/berachain/beacon-kit/cli/flags"
	cmtcli "github.com/berachain/beacon-kit/consensus/cometbft/cli"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/cosmos/cosmos-sdk/version"
)

// DefaultRootCommandSetup sets up the default commands for the root command.
func DefaultRootCommandSetup(
	root *Root,
	mm *cometbft.Service,
	appCreator servertypes.AppCreator,
	chainSpecCreator servertypes.ChainSpecCreator,
) {
	// Add all the commands to the root command.
	root.cmd.AddCommand(
		// `comet`
		cmtcli.Commands(appCreator),
		// `init`
		initialize.InitCmd(chainSpecCreator, mm),
		// `genesis`
		genesis.Commands(chainSpecCreator),
		// `deposit`
		deposit.Commands(chainSpecCreator, appCreator),
		// `jwt`
		jwt.Commands(),
		// `rollback`
		server.NewRollbackCmd(appCreator),
		// `start`
		server.StartCmdWithOptions(appCreator, server.StartCmdOptions{
			AddFlags: flags.AddBeaconKitFlags,
		}),
		// `status`
		cmtcli.StatusCommand(),
		// `version`
		version.NewVersionCommand(),
	)
}
