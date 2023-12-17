// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
