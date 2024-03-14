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

package config

import (
	"github.com/berachain/beacon-kit/config/flags"
	"github.com/berachain/beacon-kit/io/cli/parser"
)

// Chain conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Chain] = &Chain{}

const (
	defaultChainID = 1
)

// DefaultChainConfig returns the default chain general configuration.
func DefaultChainConfig() Chain {
	return Chain{
		ChainID: defaultChainID,
	}
}

// Chain represents the configuration struct for
// the general configuration of the beacon chain.
type Chain struct {
	ChainID uint64
}

// Parse parses the configuration.
func (c Chain) Parse(parser parser.AppOptionsParser) (*Chain, error) {
	var err error
	if c.ChainID, err = parser.GetUint64(
		flags.ChainID,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Chain) Template() string {
	
	return `
[beacon-kit.beacon-config.chain]
# ChainID is the unique identifier for the beacon chain.
chain-id = {{.BeaconKit.Beacon.Chain.ChainID}}
`
}
