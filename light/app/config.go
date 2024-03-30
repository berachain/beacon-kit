package app

import (
	"github.com/berachain/beacon-kit/light/client/provider"
	"github.com/berachain/beacon-kit/light/client/provider/comet"
)

type Config struct {
	Comet    *comet.Config
	Provider *provider.Config
}

func NewConfig(comet *comet.Config, provider *provider.Config) *Config {
	return &Config{
		Comet:    comet,
		Provider: provider,
	}
}
