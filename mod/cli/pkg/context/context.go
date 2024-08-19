package context

import (
	"github.com/berachain/beacon-kit/mod/log"

	cmtcfg "github.com/cometbft/cometbft/config"
	"github.com/spf13/viper"
)

// Context is a struct that contains the context for the CLI.
// TODO: rename this to Server or something more appropriate.
type Context struct {
	Viper  *viper.Viper
	Config *cmtcfg.Config
	Logger log.Logger[any]
}

// NewContext creates a new Context.
func NewContext(v *viper.Viper, config *cmtcfg.Config, logger log.Logger[any]) *Context {
	return &Context{v, config, logger}
}
