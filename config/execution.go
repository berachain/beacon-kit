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
	"github.com/berachain/beacon-kit/primitives"
	"github.com/ethereum/go-ethereum/common"
)

// Limits conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Execution] = &Execution{}

// DefaultValidatorConfig returns the default validator configuration.
func DefaultExecutionConfig() Execution {
	return Execution{
		DepositContractAddress: common.HexToAddress(
			"0x18Df82C7E422A42D47345Ed86B0E935E9718eBda",
		),
	}
}

// Execution represents the configuration struct for the
// execution layer on the beacon chain.
type Execution struct {
	DepositContractAddress primitives.ExecutionAddress
}

// Parse parses the configuration.
func (c Execution) Parse(parser parser.AppOptionsParser) (*Execution, error) {
	var err error
	if c.DepositContractAddress, err = parser.GetExecutionAddress(
		flags.DepositContractAddress,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

// Template returns the configuration template.
func (c Execution) Template() string {
	//nolint:lll
	return `
[beacon-kit.beacon-config.execution]
# DepositContractAddress is the address of the deposit contract.
deposit-contract-address = "{{.BeaconKit.Beacon.Execution.DepositContractAddress}}"
`
}
