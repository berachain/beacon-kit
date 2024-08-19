package context

import (
	"os"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetServerContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetServerContextFromCmd(cmd *cobra.Command) *Context {
	if v := cmd.Context().Value(server.ServerContextKey); v != nil {
		serverCtxPtr := v.(*Context)
		return serverCtxPtr
	}

	return NewDefaultContext()
}

func NewDefaultContext() *Context {
	return NewContext(
		viper.New(),
		cmtcfg.DefaultConfig(),
		log.NewLogger(os.Stdout),
	)
}
