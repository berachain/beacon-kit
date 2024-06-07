package spec

import (
	"github.com/berachain/beacon-kit/mod/errors"
	viperlib "github.com/berachain/beacon-kit/mod/node-core/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	cmttypes "github.com/cometbft/cometbft/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// MustReadFromAppOpts reads the configuration options from the given
// application options.
func MustReadFromAppOpts(
	opts servertypes.AppOptions,
) primitives.ChainSpec {
	spec, err := ReadFromAppOpts(opts)
	if err != nil {
		panic(err)
	}
	return spec
}

// ReadFromAppOpts reads the configuration options from the given
// application options.
func ReadFromAppOpts(
	opts servertypes.AppOptions,
) (primitives.ChainSpec, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil,
			errors.Newf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		ChainSpec primitives.ChainSpecData `mapstructure:"chain-spec"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDomainTypeFunc(),
			viperlib.StringToCometConsensusParamsFunc[*cmttypes.ConsensusParams](),
		)),
	); err != nil {
		return nil, errors.Newf(
			"failed to decode chain-spec configuration: %w",
			err,
		)
	}

	return chain.NewChainSpec(cfg.ChainSpec), nil
}
