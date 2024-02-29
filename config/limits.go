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

// Limits conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Limits] = &Limits{}

const (
	defaultMaxDepositsPerBlock      = 16
	defaultMaxWithdrawalsPerPayload = 16
)

// DefaultValidatorConfig returns the default validator configuration.
func DefaultLimitsConfig() Limits {
	return Limits{
		MaxDepositsPerBlock:      defaultMaxDepositsPerBlock,
		MaxWithdrawalsPerPayload: defaultMaxWithdrawalsPerPayload,
	}
}

// Limits represents the configuration struct for the limits on the beacon
// chain.
type Limits struct {
	MaxDepositsPerBlock      uint64 `mapstructure:"max-deposits-per-block"`
	MaxWithdrawalsPerPayload uint64 `mapstructure:"max-withdrawals-per-payload"`
}

// Template returns the configuration template.
func (c Limits) Template() string {
	//nolint:lll
	return `
[beacon-kit.beacon-config.limits]
# MaxDepositsPerBlock is the maximum number of Deposits allowed in a block.
max-deposits-per-block = {{.BeaconKit.Beacon.Limits.MaxDepositsPerBlock}}
# MaxWithdrawalsPerPayload is the maximum number of Withdrawals allowed in a payload.
max-withdrawals-per-payload = {{.BeaconKit.Beacon.Limits.MaxWithdrawalsPerPayload}}
`
}
