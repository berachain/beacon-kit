package runner

import (
	"github.com/berachain/beacon-kit/mod/cli/pkg/commands"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// Runner is a type that runs the root command.
type Runner struct {
	// nodeHome is the home directory of the node.
	nodeHome string
	// rootCmd is the root command.
	rootCmd *commands.Root
}

// New returns a new Runner with the given root command.
func New[NodeT servertypes.Application](
	rootCmd *commands.Root,
	appCreator servertypes.AppCreator[NodeT],
) *Runner {
	commands.SetupRootCmdWithNode[NodeT](
		rootCmd,
		appCreator,
	)
	return &Runner{rootCmd: rootCmd}
}

// Run runs the root command.
func (runner *Runner) Run() error {
	return runner.rootCmd.Run(runner.nodeHome)
}
