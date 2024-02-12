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

package config

import (
	"github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/config/parser"
)

// Limits conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Limits] = &Limits{}

// DefaultValidatorConfig returns the default validator configuration.
func DefaultLimitsConfig() Limits {
	return Limits{
		MaxDeposits:    16, //nolint:gomnd
		MaxWithdrawals: 16, //nolint:gomnd
	}
}

// Limits represents the configuration struct for the limits on the beacon chain.
type Limits struct {
	MaxDeposits    uint64
	MaxWithdrawals uint64
}

// Parse parses the configuration.
func (c Limits) Parse(parser parser.AppOptionsParser) (*Limits, error) {
	var err error
	if c.MaxDeposits, err = parser.GetUint64(flags.MaxDeposits); err != nil {
		return nil, err
	}
	if c.MaxWithdrawals, err = parser.GetUint64(flags.MaxWithdrawals); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Limits) Template() string {
	return `
[beacon-kit.beacon-config.limits]
# MaxDeposits is the maximum number of Deposits allowed in a block.
max-deposits = {{.BeaconKit.Beacon.Limits.MaxDeposits}}
# MaxWithdrawals is the maximum number of Withdrawals allowed in a block.
max-withdrawals = {{.BeaconKit.Beacon.Limits.MaxWithdrawals}}
`
}
