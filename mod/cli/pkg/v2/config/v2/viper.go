package config

import (
	viperlib "github.com/berachain/beacon-kit/mod/cli/pkg/v2/config/v2/viper"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// AppOptions is from the SDK, we should look to remove its usage.
type AppOptions interface {
	Get(string) interface{}
}

// MustReadConfigFromAppOpts reads the configuration options from the given
// application options.
func MustReadConfigFromAppOpts[ConfigT any](opts AppOptions) *ConfigT {
	cfg, err := ReadConfigFromAppOpts[ConfigT](opts)
	if err != nil {
		panic(err)
	}
	return cfg
}

// ReadConfigFromAppOpts reads the configuration options from the given
// application options.
func ReadConfigFromAppOpts[ConfigT any](opts AppOptions) (*ConfigT, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.Newf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		BeaconKit ConfigT `mapstructure:"beacon-kit"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDialURLFunc(),
			viperlib.StringToConnectionURLFunc(),
		))); err != nil {
		return nil, errors.Newf(
			"failed to decode beacon-kit configuration: %w",
			err,
		)
	}

	return &cfg.BeaconKit, nil
}
