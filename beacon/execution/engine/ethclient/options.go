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

package eth

import (
	"net/url"

	"cosmossdk.io/log"
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

// WithJWTSecret is an option to set the JWT secret for the Eth1Client.
func WithJWTSecret(secret [jwtLength]byte) Option {
	return func(s *Eth1Client) error {
		if len(secret) != jwtLength {
			return ErrInvalidJWTSecretLength
		}
		s.jwtSecret = secret
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *Eth1Client) error {
		s.logger = logger.With("module", "beacon-kit-execution")
		return nil
	}
}

// WithRequiredChainID is an option to set the required
// chain ID for the Eth1Client.
func WithRequiredChainID(chainID uint64) Option {
	return func(s *Eth1Client) error {
		s.chainID = chainID
		return nil
	}
}
