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

package blockchain

import (
	"cosmossdk.io/log"

	"github.com/itsdevbear/bolaris/config"
)

// Option is a function type that takes a pointer to a Service and returns an error.
type Option func(*Service) error

// WithExecutionService is a function that returns an Option.
// It sets the Service of the Service to the provided Service.
func WithExecutionService(en ExecutionService) Option {
	return func(s *Service) error {
		s.en = en
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithBeaconConfig(beaconCfg *config.Beacon) Option {
	return func(s *Service) error {
		s.beaconCfg = beaconCfg
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *Service) error {
		s.logger = logger.With("module", "beacon-kit-blockchain")
		return nil
	}
}

// WithForkChoiceStoreProvider is an option to set the ForkChoiceStoreProvider for the Service.
func WithForkChoiceStoreProvider(fcsp ForkChoiceStoreProvider) Option {
	return func(s *Service) error {
		s.fcsp = fcsp
		return nil
	}
}
