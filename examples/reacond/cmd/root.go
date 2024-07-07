package main

import (
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/client"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/cometbft"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/deposit"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/genesis"
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands/jwt"
	"github.com/berachain/beacon-kit/mod/cli/pkg/flags"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	rollserv "github.com/rollkit/cosmos-sdk-starter/server"
	rollconf "github.com/rollkit/rollkit/config"
	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command) {
	flags.AddBeaconKitFlags(cmd)
	rollconf.AddFlags(cmd)
}

// RollKitRootCommandSetup sets up the default commands for the root command.
func RollKitRootCommandSetup[T types.Node](
	root *commands.Root,
	mm *module.Manager,
	appCreator servertypes.AppCreator[T],
	chainSpec common.ChainSpec,
) {
	// Setup the custom start command options.
	startCmdOptions := server.StartCmdOptions[T]{
		AddFlags:            addFlags,
		StartCommandHandler: rollserv.StartHandler[T],
	}

	// Add all the commands to the root command.
	root.Command().AddCommand(
		// `comet`
		cometbft.Commands(appCreator),
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
		pruning.Cmd(appCreator),
		// `rollback`
		server.NewRollbackCmd(appCreator),
		// `snapshots`
		snapshot.Cmd(appCreator),
		// `start`
		server.StartCmdWithOptions(appCreator, startCmdOptions),
		// `status`
		server.StatusCommand(),
		// `version`
		version.NewVersionCommand(),
	)
}
