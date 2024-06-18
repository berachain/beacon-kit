package components

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	flags "github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

type LoggerInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// CreateSDKLogger creates a the default SDK logger.
// It reads the log level and format from the server context.
func ProvideLogger(
	in LoggerInput,
) (log.Logger, error) {
	logLvlStr := cast.ToString(in.AppOpts.Get(flags.FlagLogLevel))

	return phuslu.NewLogger[any, log.Logger](logLvlStr), nil
}
