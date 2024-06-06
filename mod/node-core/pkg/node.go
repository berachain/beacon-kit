package node

import (
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

type Node struct {
	*app.BeaconApp

	name        string
	description string

	rootCmd *cobra.Command
}

func New[NodeT NodeI]() NodeT {
	return NodeI(&Node{}).(NodeT)
}

func (n *Node) Run() error {
	return svrcmd.Execute(
		n.rootCmd, "", components.DefaultNodeHome,
	)
}

func (n *Node) SetAppName(name string) {
	n.name = name
}

func (n *Node) SetAppDescription(description string) {
	n.description = description
}

func (n *Node) SetApplication(a servertypes.Application) {
	n.BeaconApp = a.(*app.BeaconApp)
}

func (n *Node) SetRootCmd(cmd *cobra.Command) {
	n.rootCmd = cmd
}
