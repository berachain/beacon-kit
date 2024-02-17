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

package logs

import (
	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
	eth "github.com/itsdevbear/bolaris/execution/engine/ethclient"
)

// Option is a function that applies a specific configuration to the Processor.
type Option func(*Processor) error

// WithEthClient is an Option that sets the Ethereum client for the Processor.
func WithEthClient(eth1Client *eth.Eth1Client) Option {
	return func(p *Processor) error {
		p.eth1Client = eth1Client
		return nil
	}
}

// WithHandlers is an Option that sets the contract address for the Processor.
func WithHandlers(handlers map[common.Address]callback.LogHandler) Option {
	return func(p *Processor) error {
		p.handlers = handlers
		return nil
	}
}

// WithLogger is an Option that sets the logger for the Processor.
func WithLogger(logger log.Logger) Option {
	return func(p *Processor) error {
		p.logger = logger.With("module", "beacon-kit-log-processor")
		return nil
	}
}
