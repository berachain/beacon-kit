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
	"github.com/prysmaticlabs/prysm/v4/network"
	"github.com/prysmaticlabs/prysm/v4/network/authorization"

	"cosmossdk.io/log"
)

type Option func(s *Eth1Client) error

// WithHTTPEndpointAndJWTSecret for authenticating the execution node JSON-RPC endpoint.
func WithHTTPEndpointAndJWTSecret(endpointString string, secret []byte) Option {
	return func(s *Eth1Client) error {
		if len(secret) == 0 {
			return nil
		}
		// Overwrite authorization type for all endpoints to be of a bearer type.
		hEndpoint := network.HttpEndpoint(endpointString)
		hEndpoint.Auth.Method = authorization.Bearer
		hEndpoint.Auth.Value = string(secret)

		s.cfg.currHTTPEndpoint = hEndpoint
		return nil
	}
}

// WithLogger is an option to set the logger for the Eth1Client.
func WithLogger(logger log.Logger) Option {
	return func(s *Eth1Client) error {
		s.logger = logger
		return nil
	}
}

// WithHeaders is an option to set the headers for the Eth1Client.
func WithHeaders(headers []string) Option {
	return func(s *Eth1Client) error {
		s.cfg.headers = headers
		return nil
	}
}

// WithRequiredChainID is an option to set the required
// chain ID for the Eth1Client.
func WithRequiredChainID(chainID uint64) Option {
	return func(s *Eth1Client) error {
		s.cfg.chainID = chainID
		return nil
	}
}
