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
	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/deposit"
	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/genesis"
	initcli "github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/init"
	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/jwt"
	"github.com/berachain/beacon-kit/mod/cli/pkg/v2/commands/start"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// DefaultRootCommandSetup sets up the default commands for the root command.
func DefaultRootCommandSetup[
	NodeT Node,
	ConsensusParamsT ConsensusParams[ConsensusParamsT],
	GenesisStateT GenesisState[GenesisStateT],
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
](
	root *Root[NodeT],
	chainSpec common.ChainSpec,
) {
	// // Setup the custom start command options.
	// startCmdOptions := server.StartCmdOptions[T]{
	// 	AddFlags: flags.AddBeaconKitFlags,
	// }

	// Add all the commands to the root command.
	root.Command.AddCommand(
		// `comet`
		// cometbft.Commands(appCreator),
		// `client`
		// client.Commands(),
		// `config`
		// confixcmd.ConfigCommand(),
		// `init`
		initcli.Command[ConsensusParamsT, GenesisStateT](),
		// `genesis`
		genesis.Commands(chainSpec),
		// `deposit`
		deposit.Commands[ExecutionPayloadT](chainSpec),
		// `jwt`
		jwt.Commands(),
		// // `keys`
		// keys.Commands(),
		// // `prune`
		// pruning.Cmd(appCreator),
		// // `rollback`
		// server.NewRollbackCmd(appCreator),
		// // `snapshots`
		// snapshot.Cmd(appCreator),
		// `start`
		start.Command(root.Node),
		// // `status`
		// server.StatusCommand(),
		// // `version`
		// version.NewVersionCommand(),
	)
}
