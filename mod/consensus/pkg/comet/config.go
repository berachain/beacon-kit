package cometbft

import (
	"github.com/cometbft/cometbft/config"
)

type Config struct {
	cmtConfig config.Config

	NodeKeyFile            string
	PrivValidatorKeyFile   string
	PrivValidatorStateFile string
}
