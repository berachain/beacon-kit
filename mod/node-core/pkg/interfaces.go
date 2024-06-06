package node

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
)

type NodeI interface {
	servertypes.Application

	Run() error

	SetAppName(name string)
	SetAppDescription(description string)
	SetRootCmd(cmd *cobra.Command)
	SetApplication(app servertypes.Application)
}
