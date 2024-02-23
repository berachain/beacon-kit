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

package runtime

import (
	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/config"
	"github.com/itsdevbear/bolaris/runtime/service"
)

// Option is a function that modifies the BeaconKitRuntime.
type Option func(*BeaconKitRuntime) error

// WithConfig sets the configuration of the BeaconKitRuntime.
func WithConfig(cfg *config.Config) Option {
	return func(r *BeaconKitRuntime) error {
		r.cfg = cfg
		return nil
	}
}

// WithServiceRegistry sets the service registry of the BeaconKitRuntime.
func WithServiceRegistry(reg *service.Registry) Option {
	return func(r *BeaconKitRuntime) error {
		r.services = reg
		return nil
	}
}

// WithLogger sets the logger of the BeaconKitRuntime.
func WithLogger(logger log.Logger) Option {
	return func(r *BeaconKitRuntime) error {
		r.logger = logger.With("module", "beacon-kit-runtime")
		return nil
	}
}

// WithBeaconStateProvider sets the BeaconStateProvider
// of the BeaconKitRuntime.
func WithBeaconStateProvider(fscp BeaconStateProvider) Option {
	return func(r *BeaconKitRuntime) error {
		r.fscp = fscp
		return nil
	}
}
