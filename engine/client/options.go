// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package client

import (
	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/config"
	eth "github.com/itsdevbear/bolaris/engine/client/ethclient"
	"github.com/itsdevbear/bolaris/io/jwt"
)

// Option is a function type that takes a pointer to an engineClient and returns
// an error.
type Option func(*EngineClient) error

// WithBeaconConfig is an option to set the beacon configuration.
func WithBeaconConfig(beaconCfg *config.Beacon) Option {
	return func(s *EngineClient) error {
		s.beaconCfg = beaconCfg
		return nil
	}
}

// WithEngineConfig is a function that returns an Option.
func WithEngineConfig(cfg *config.Engine) Option {
	return func(s *EngineClient) error {
		var err error

		// Load the JWT secret from the config if it's not already set.
		// Get JWT Secret for eth1 connection.
		s.jwtSecret, err = jwt.NewFromFile(cfg.JWTSecretPath)
		if err != nil {
			return err
		}

		s.cfg = cfg
		return nil
	}
}

// WithEth1Client is a function that returns an Option.
func WithEth1Client(eth1Client *eth.Eth1Client) Option {
	return func(s *EngineClient) error {
		s.Eth1Client = eth1Client
		return nil
	}
}

// WithLogger is an option to set the logger for the EngineClient.
func WithLogger(logger log.Logger) Option {
	return func(s *EngineClient) error {
		s.logger = logger.With("module", "beacon-kit.engine.client")
		return nil
	}
}
