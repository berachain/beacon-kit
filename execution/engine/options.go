// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package engine

import (
	"time"

	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/config"
	eth "github.com/itsdevbear/bolaris/execution/engine/ethclient"
)

// Option is a function type that takes a pointer to an engineClient and returns an error.
type Option func(*engineClient) error

// WithEth1Client is a function that returns an Option.
func WithEth1Client(eth1Client *eth.Eth1Client) Option {
	return func(s *engineClient) error {
		s.Eth1Client = eth1Client
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithBeaconConfig(beaconCfg *config.Beacon) Option {
	return func(s *engineClient) error {
		s.beaconCfg = beaconCfg
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *engineClient) error {
		s.logger = logger.With("module", "beacon-kit.engine")
		return nil
	}
}

// WithEngineTimeout is an option to set the timeout for the engine.
func WithEngineTimeout(engineTimeout time.Duration) Option {
	return func(s *engineClient) error {
		s.engineTimeout = engineTimeout
		return nil
	}
}
