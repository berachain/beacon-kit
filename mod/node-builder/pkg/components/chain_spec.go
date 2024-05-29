package components

import (
	"path/filepath"

	"cosmossdk.io/depinject"
	viperlib "github.com/berachain/beacon-kit/mod/node-builder/pkg/config/viper"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/chain"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const ChainSpecFileName = "chainspec.toml"

type ChainSpecInput struct {
	depinject.In
	ClientCtx client.Context
}

func ProvideChainSpec(in ChainSpecInput) (primitives.ChainSpec, error) {
	configDir := filepath.Join(in.ClientCtx.HomeDir, "config")
	chainSpecFile := filepath.Join(configDir, ChainSpecFileName)

	v := viper.New()
	v.AddConfigPath(chainSpecFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// type specUnmarshaller struct {
	// 	Data primitives.ChainSpecData `mapstructure:"specData"`
	// }

	var specData primitives.ChainSpecData
	if err := v.Unmarshal(
		&specData,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			viperlib.StringToExecutionAddressFunc(),
		)),
	); err != nil {
		return nil, err
	}

	return chain.NewChainSpec(specData), nil
}
