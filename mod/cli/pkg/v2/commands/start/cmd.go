package start

import (
	"github.com/spf13/cobra"
)

// TODO: this function needs the memory address of where the
// real node is going to be assigned.
func Command[NodeT Node](node NodeT) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the node",
		Long:  "Start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: wire real context
			return node.Start(cmd.Context())
		},
	}
}
