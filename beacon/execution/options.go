package execution

import (
	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/cosmos/config"
)

type Option func(*engineCaller) error

// WithLogger is an option to set the logger for the Eth1Client.
func WithBeaconConfig(beaconCfg *config.Beacon) Option {
	return func(s *engineCaller) error {
		s.beaconCfg = beaconCfg
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *engineCaller) error {
		s.logger = logger
		return nil
	}
}
