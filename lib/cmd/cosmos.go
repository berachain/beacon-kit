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

package cmd

import (
	authcmd "cosmossdk.io/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

// GenesisCommands builds genesis-related `simd genesis` command. Users may
// provide application specific commands as a parameter.
func GenesisCommands(
	txConfig client.TxConfig,
	mm *module.Manager,
	appExport servertypes.AppExporter,
	cmds ...*cobra.Command,
) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, mm, appExport)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

// QueryCommands constructs a new cobra.Command to interact with querying
// subcommands.
func QueryCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:gomnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	// Adding subcommands for querying blockchain data.
	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		server.QueryBlockResultsCmd(),
		authcmd.QueryTxCmd(),
	)

	return cmd
}

// TxCommands creates a new cobra.Command for transaction-related subcommands.
func TxCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:gomnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)

	return cmd
}
