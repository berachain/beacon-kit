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

package cometbft

import (
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	cmtcmd "github.com/cometbft/cometbft/cmd/cometbft/commands"
	"github.com/spf13/cobra"
)

// Commands creates a new command for managing CometBFT
// related commands.
func Commands[T types.Node](
	appCreator servertypes.AppCreator[T],
) *cobra.Command {
	cometCmd := &cobra.Command{
		Use:     "comet",
		Aliases: []string{"cometbft", "tendermint"},
		Short:   "CometBFT subcommands",
	}

	cometCmd.AddCommand(
		server.ShowNodeIDCmd(),
		server.ShowValidatorCmd(),
		server.ShowAddressCmd(),
		server.VersionCmd(),
		cmtcmd.ResetAllCmd,
		cmtcmd.ResetStateCmd,
		server.BootstrapStateCmd(appCreator),
	)

	return cometCmd
}
