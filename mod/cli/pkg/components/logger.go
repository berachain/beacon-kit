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
	// out := os.Stdout
	// var opts []log.Option
	// if in.AppOpts.Get(flags.FlagLogFormat) == flags.OutputFormatJSON {
	// 	opts = append(opts, log.OutputJSONOption())
	// }
	// opts = append(opts,
	// 	log.ColorOption(!cast.ToBool(in.AppOpts.Get(flags.FlagLogNoColor))),
	// 	// We use CometBFT flag (cmtcli.TraceFlag) for trace logging.
	// 	log.TraceOption(cast.ToBool(in.AppOpts.Get(server.FlagTrace))))

	// check and set filter level or keys for the logger if any
	logLvlStr := cast.ToString(in.AppOpts.Get(flags.FlagLogLevel))
	// if logLvlStr == "" {
	// 	return log.NewLogger(out, opts...), nil
	// }

	// logLvl, err := zerolog.ParseLevel(logLvlStr)
	// switch {
	// case err != nil:
	// 	// If the log level is not a valid zerolog level, then we try to parse it as a key filter.
	// 	filterFunc, err := log.ParseLogLevel(logLvlStr)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	opts = append(opts, log.FilterOption(filterFunc))
	// default:
	// 	opts = append(opts, log.LevelOption(logLvl))
	// }

	return phuslu.NewLogger[any, log.Logger](logLvlStr), nil
}
