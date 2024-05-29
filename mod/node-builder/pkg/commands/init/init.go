package initutil

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

// TODO: unhood this way of wrapping cosmos' init command
func Command(mm *module.Manager) *cobra.Command {
	cmd := genutilcli.InitCmd(mm)

	// TODO: Temp holder of existing cosmos init runE function,
	// not sure if we need this or if calling it directly
	// inside of the new assignment will work.
	runFn := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := runFn(cmd, args); err != nil {
			return err
		}
		// if err := ; err != nil {
		// 	return err
		// }

		return createChainSpecCmd(cmd, args)
	}

	return cmd
}
