package components

import (
	"io"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	flags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/rs/zerolog"
)

// CreateSDKLogger creates a the default SDK logger.
// It reads the log level and format from the server context.
func ProvideLogger(ctx *server.Context, out io.Writer) (log.Logger, error) {
	var opts []log.Option
	if ctx.Viper.GetString(flags.FlagLogFormat) == flags.OutputFormatJSON {
		opts = append(opts, log.OutputJSONOption())
	}
	opts = append(opts,
		log.ColorOption(!ctx.Viper.GetBool(flags.FlagLogNoColor)),
		// We use CometBFT flag (cmtcli.TraceFlag) for trace logging.
		log.TraceOption(ctx.Viper.GetBool(server.FlagTrace)))

	// check and set filter level or keys for the logger if any
	logLvlStr := ctx.Viper.GetString(flags.FlagLogLevel)
	if logLvlStr == "" {
		return log.NewLogger(out, opts...), nil
	}

	logLvl, err := zerolog.ParseLevel(logLvlStr)
	switch {
	case err != nil:
		// If the log level is not a valid zerolog level, then we try to parse it as a key filter.
		filterFunc, err := log.ParseLogLevel(logLvlStr)
		if err != nil {
			return nil, err
		}

		opts = append(opts, log.FilterOption(filterFunc))
	default:
		opts = append(opts, log.LevelOption(logLvl))
	}

	return phuslu.NewLogger[any, log.Logger](logLvlStr), nil
}
