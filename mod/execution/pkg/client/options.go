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
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	eth "github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
)

// Option is a function type that takes a pointer to an engineClient and returns
// an error.
type Option func(*EngineClient) error

// WithEngineConfig is a function that returns an Option.
func WithEngineConfig(cfg *Config) Option {
	return func(s *EngineClient) error {
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

// WithJWTSecret is a function that returns an Option.
func WithJWTSecret(secret *jwt.Secret) Option {
	return func(s *EngineClient) error {
		s.jwtSecret = secret
		return nil
	}
}

// WithLogger is an option to set the logger for the EngineClient.
func WithLogger(logger log.Logger[any]) Option {
	return func(s *EngineClient) error {
		s.logger = logger
		return nil
	}
}

// WithCacheConfig is an option to create a new cache
// with the given config for the EngineClient.
func WithCacheConfig(config cache.Config) Option {
	return func(s *EngineClient) error {
		s.engineCache = cache.NewEngineCache(config)
		return nil
	}
}
