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

package ethclient

import (
	"net/url"
	"time"

	"cosmossdk.io/log"
	"github.com/itsdevbear/bolaris/io/jwt"
)

// Option is a functional option for the Eth1Client.
type Option func(s *Eth1Client) error

// WithEndpointDialURL for authenticating the execution node JSON-RPC endpoint.
func WithEndpointDialURL(dialURL string) Option {
	return func(s *Eth1Client) error {
		u, err := url.Parse(dialURL)
		if err != nil {
			return err
		}
		s.dialURL = u
		return nil
	}
}

// WithJWTSecret sets the JWT secret for the Eth1Client.
func WithJWTSecret(secret *jwt.Secret) Option {
	return func(s *Eth1Client) error {
		if secret == nil {
			return ErrNilJWTSecret
		}
		s.jwtSecret = secret
		return nil
	}
}

// WithStartupRetryInterval sets the startup retry interval for the Eth1Client.
func WithStartupRetryInterval(interval time.Duration) Option {
	return func(s *Eth1Client) error {
		s.startupRetryInterval = interval
		return nil
	}
}

// WithJWTRefreshInterval sets the JWT refresh interval for the Eth1Client.
func WithJWTRefreshInterval(interval time.Duration) Option {
	return func(s *Eth1Client) error {
		s.jwtRefreshInterval = interval
		return nil
	}
}

// WithHealthCheckInterval sets the health check interval for the Eth1Client.
func WithHealthCheckInterval(interval time.Duration) Option {
	return func(s *Eth1Client) error {
		s.healthCheckInterval = interval
		return nil
	}
}

// WithLogger sets the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *Eth1Client) error {
		s.logger = logger.With("module", "beacon-kit-execution")
		return nil
	}
}

// WithRequiredChainID sets the required
// chain ID for the Eth1Client.
func WithRequiredChainID(chainID uint64) Option {
	return func(s *Eth1Client) error {
		s.chainID = chainID
		return nil
	}
}
