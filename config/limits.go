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
)

// Limits conforms to the BeaconKitConfig interface.
var _ BeaconKitConfig[Limits] = &Limits{}

// Default configuration limits for the beacon chain.
const (
	// defaultEpochsPerHistoricalVector defines the default number of epochs.
	defaultEpochsPerHistoricalVector = 1
	// defaultSlotsPerHistoricalRoot defines the default limit for historical
	// roots.
	defaultSlotsPerHistoricalRoot = 32
	// defaultMaxDepositsPerBlock specifies the default maximum number of
	// deposits allowed per block.
	defaultMaxDepositsPerBlock = 16
	// defaultMaxRedirectsPerBlock specifies the default maximum number of
	// redirects allowed per block.
	defaultMaxRedirectsPerBlock = 16
	// defaultMaxWithdrawalsPerPayload indicates the default maximum number of
	// withdrawals allowed per payload.
	defaultMaxWithdrawalsPerPayload = 16
)

// DefaultValidatorConfig returns the default validator configuration.
func DefaultLimitsConfig() Limits {
	return Limits{
		EpochsPerHistoricalVector: defaultEpochsPerHistoricalVector,
		SlotsPerHistoricalRoot:    defaultSlotsPerHistoricalRoot,
		MaxDepositsPerBlock:       defaultMaxDepositsPerBlock,
		MaxRedirectsPerBlock:      defaultMaxRedirectsPerBlock,
		MaxWithdrawalsPerPayload:  defaultMaxWithdrawalsPerPayload,
	}
}

// Limits represents the configuration struct for the limits on the beacon
// chain.
type Limits struct {
	// EpochsPerHistoricalVector defines the number of epochs per historical
	// vector. (i.e RANDAO)
	EpochsPerHistoricalVector primitives.Epoch
	// SlotsPerHistoricalRoot defines the maximum number of historical root
	// entries that can be stored.
	SlotsPerHistoricalRoot uint64
	// MaxDepositsPerBlock specifies the maximum number of deposit operations
	// allowed per block.
	MaxDepositsPerBlock uint64
	// MaxRedirectsPerBlock specifies the maximum number of redirect operations
	// allowed per block.
	MaxRedirectsPerBlock uint64
	// MaxWithdrawalsPerPayload indicates the maximum number of withdrawal
	// operations allowed in a single payload.
	MaxWithdrawalsPerPayload uint64
}

// Parse parses the configuration.
func (c Limits) Parse(parser parser.AppOptionsParser) (*Limits, error) {
	var err error

	if c.EpochsPerHistoricalVector, err = parser.GetEpoch(
		flags.EpochsPerHistoricalVector,
	); err != nil {
		return nil, err
	}

	if c.SlotsPerHistoricalRoot, err = parser.GetUint64(
		flags.SlotsPerHistoricalRoot,
	); err != nil {
		return nil, err
	}

	if c.MaxDepositsPerBlock, err = parser.GetUint64(
		flags.MaxDeposits,
	); err != nil {
		return nil, err
	}

	if c.MaxRedirectsPerBlock, err = parser.GetUint64(
		flags.MaxRedirects,
	); err != nil {
		return nil, err
	}

	if c.MaxWithdrawalsPerPayload, err = parser.GetUint64(
		flags.MaxWithdrawals,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Template returns the configuration template.
func (c Limits) Template() string {
	//nolint:lll
	return `
[beacon-kit.beacon-config.limits]
# EpochsPerHistoricalVector is the number of epochs per historical vector.
epochs-per-historical-vector = {{.BeaconKit.Beacon.Limits.EpochsPerHistoricalVector}}

# SlotsPerHistoricalRoot is the maximum number of historical roots that will be stored.
slots-per-historical-root = {{.BeaconKit.Beacon.Limits.SlotsPerHistoricalRoot}}

# MaxDepositsPerBlock is the maximum number of Deposits allowed in a block.
max-deposits-per-block = {{.BeaconKit.Beacon.Limits.MaxDepositsPerBlock}}

# MaxRedirectsPerBlock is the maximum number of Redirects allowed in a block.
max-redirects-per-block = {{.BeaconKit.Beacon.Limits.MaxRedirectsPerBlock}}

# MaxWithdrawalsPerPayload is the maximum number of Withdrawals allowed in a payload.
max-withdrawals-per-payload = {{.BeaconKit.Beacon.Limits.MaxWithdrawalsPerPayload}}
`
}
