package spec

import (
	"bytes"
	"os"
	"text/template"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	viperlib "github.com/berachain/beacon-kit/mod/node-builder/pkg/config/viper"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// WriteToFile writes the chain spec data to a file.
func WriteToFile(filepath string, spec primitives.ChainSpecData) error {
	tmpl, err := template.New("specTemplate").Parse(Template)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, spec); err != nil {
		return err
	}

	return os.WriteFile(filepath, buffer.Bytes(), 0o600)
}

// MustReadFromAppOpts reads the configuration options from the given
// application options.
func MustReadFromAppOpts(opts servertypes.AppOptions) *primitives.ChainSpecData {
	specData, err := ReadFromAppOpts(opts)
	if err != nil {
		panic(err)
	}
	return specData
}

// ReadFromAppOpts reads the configuration options from the given
// application options.
func ReadFromAppOpts(opts servertypes.AppOptions) (*primitives.ChainSpecData, error) {
	v, ok := opts.(*viper.Viper)
	if !ok {
		return nil, errors.Newf("invalid application options type: %T", opts)
	}

	type cfgUnmarshaller struct {
		ChainSpec primitives.ChainSpecData `mapstructure:"chain-spec"`
	}
	cfg := cfgUnmarshaller{}
	if err := v.Unmarshal(&cfg,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			viperlib.StringToExecutionAddressFunc(),
			viperlib.StringToDomainTypeFunc(),
		))); err != nil {
		return nil, errors.Newf(
			"failed to decode chain-spec configuration: %w",
			err,
		)
	}

	return &cfg.ChainSpec, nil
}
