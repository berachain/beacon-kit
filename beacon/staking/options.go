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

package staking

import (
	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/itsdevbear/bolaris/execution/engine/ethclient"
)

// Option is a function type that takes a pointer to an Service and returns an error.
type Option func(*Service) error

// WithLogger is an Option that sets the logger for the Service.
func WithLogger(logger log.Logger) Option {
	return func(s *Service) error {
		s.logger = logger.With("staking", "log-handler")
		return nil
	}
}

// WithEth1Client is an Option that sets the Ethereum client for the Service.
func WithEth1Client(eth1Client *eth.Eth1Client) Option {
	return func(s *Service) error {
		s.eth1Client = eth1Client
		return nil
	}
}

// WithDepositContractAddress is an Option that sets the deposit contract address for the Service.
func WithDepositContractAddress(depositContractAddress common.Address) Option {
	return func(s *Service) error {
		s.depositContractAddress = depositContractAddress
		return nil
	}
}
