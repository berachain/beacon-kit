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
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/runtime/v2"
	serverv2 "cosmossdk.io/server/v2"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/client"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/cometbft"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/deposit"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/genesis"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/jwt"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

// Commands sets up the default commands for the root command.
func Commands[
	NodeT types.Node[T],
	T transaction.Tx,
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	root *Root,
	mm *runtime.MM[T],
	appCreator serverv2.AppCreator[NodeT, T],
	chainSpec common.ChainSpec,
) []*cobra.Command {
	// Setup the custom start command options.
	// startCmdOptions := server.StartCmdOptions[NodeT]{
	// 	AddFlags: flags.AddBeaconKitFlags,
	// }

	cmds := []*cobra.Command{
		// `comet`
		cometbft.Commands[NodeT](appCreator),
		// `client`
		client.Commands(), // we don't need this anymore once cometbftserver
		// adheres to the HasStartCmd flag.
		// `config`
		confixcmd.ConfigCommand(),
		// `init`
		genutilcli.InitCmd(mm),
		// `genesis`
		genesis.Commands(chainSpec),
		// `deposit`
		deposit.Commands[ExecutionPayloadT](chainSpec),
		// `jwt`
		jwt.Commands(),
		// `keys`
		keys.Commands(),

		// Not yet implemented on SimappV2
		// `prune`
		// pruning.Cmd(appCreator),
		// // `rollback`
		// server.NewRollbackCmd(appCreator),
		// // `snapshots`
		// snapshot.Cmd(appCreator),

		// `status`
		server.StatusCommand(),
		// `version`
		version.NewVersionCommand(),
	}

	return cmds
}
