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

package params

import (
	flags "github.com/berachain/beacon-kit/config/flags"
	"github.com/berachain/beacon-kit/io/cli/parser"
	"github.com/berachain/beacon-kit/primitives"
)

//nolint:lll // struct tags may create long lines.
type BeaconChainConfig struct {
	// Gwei value constants.
	//
	// MinDepositAmount is the minimum deposit amount per deposit
	// transaction.
	MinDepositAmount uint64 `mapstructure:"min-deposit-amount"`
	// MaxEffectiveBalance is the maximum effective balance allowed for a
	// validator.
	MaxEffectiveBalance uint64 `mapstructure:"max-effective-balance"`
	// EffectiveBalanceIncrement is the effective balance increment.
	EffectiveBalanceIncrement uint64 `mapstructure:"effective-balance-increment"`

	// Time parameters constants.
	//
	// SlotsPerEpoch is the number of slots per epoch.
	SlotsPerEpoch uint64 `mapstructure:"slots-per-epoch"`
	// SlotsPerHistoricalRoot is the number of slots per historical root.
	SlotsPerHistoricalRoot uint64 `mapstructure:"slots-per-historical-root"`

	// Eth1-related values.
	//
	// DepositContractAddress is the address of the deposit contract.
	DepositContractAddress primitives.ExecutionAddress `mapstructure:"deposit-contract-address"`

	// Fork-related values.
	//
	// ElectraForkEpoch is the epoch at which the Electra fork is activated.
	ElectraForkEpoch primitives.Epoch `mapstructure:"electra-fork-epoch"`

	// State list lengths
	//
	// EpochsPerHistoricalVector is the number of epochs in the historical
	// vector.
	EpochsPerHistoricalVector uint64 `mapstructure:"epochs-per-historical-vector"`

	// Max operations per block constants.
	//
	// MaxDepositsPerBlock specifies the maximum number of deposit operations
	// allowed per block.
	MaxDepositsPerBlock uint64 `mapstructure:"max-deposits-per-block"`
	// MaxWithdrawalsPerPayload indicates the maximum number of withdrawal
	// operations allowed in a single payload.
	MaxWithdrawalsPerPayload uint64 `mapstructure:"max-withdrawals-per-payload"`
}

// Template returns the configuration template.
func (c BeaconChainConfig) Template() string {
	//nolint:lll // long lines are okay.
	return `
[beacon-kit.beacon-chain]

########### Gwei Values ###########
# MinDepositAmount is the minimum deposit amount per deposit transaction.
min-deposit-amount = {{.BeaconKit.Beacon.MinDepositAmount}}

# MaxEffectiveBalance is the maximum effective balance allowed for a validator.
max-effective-balance = {{.BeaconKit.Beacon.MaxEffectiveBalance}}

# EffectiveBalanceIncrement is the effective balance increment.
effective-balance-increment = {{.BeaconKit.Beacon.EffectiveBalanceIncrement}}

########### Time Parameters ##########
# SlotsPerEpoch is the number of slots per epoch.
slots-per-epoch = {{.BeaconKit.Beacon.SlotsPerEpoch}}

# SlotsPerHistoricalRoot is the number of slots per historical root.
slots-per-historical-root = {{.BeaconKit.Beacon.SlotsPerHistoricalRoot}}

########### Eth1 Data ###########
# DepositContractAddress is the address of the deposit contract.
deposit-contract-address = "{{.BeaconKit.Beacon.DepositContractAddress}}"

########### Forks ###########
# Electra fork epoch
electra-fork-epoch = {{.BeaconKit.Beacon.ElectraForkEpoch}}

########### State List Lengths ###########
# EpochsPerHistoricalVector is the number of epochs in the historical vector.
epochs-per-historical-vector = {{.BeaconKit.Beacon.EpochsPerHistoricalVector}}

########### Max Operations ###########
# MaxDepositsPerBlock specifies the maximum number of deposit operations allowed per block.
max-deposits-per-block = {{.BeaconKit.Beacon.MaxDepositsPerBlock}}

# MaxWithdrawalsPerPayload indicates the maximum number of withdrawal operations allowed in a single payload.
max-withdrawals-per-payload = {{.BeaconKit.Beacon.MaxWithdrawalsPerPayload}}
`
}

func (c BeaconChainConfig) Parse(
	parser parser.AppOptionsParser,
) (*BeaconChainConfig, error) {
	var err error

	if c.MinDepositAmount, err = parser.GetUint64(
		flags.MinDepositAmount,
	); err != nil {
		return nil, err
	}

	if c.MaxEffectiveBalance, err = parser.GetUint64(
		flags.MaxEffectiveBalance,
	); err != nil {
		return nil, err
	}

	if c.EffectiveBalanceIncrement, err = parser.GetUint64(
		flags.EffectiveBalanceIncrement,
	); err != nil {
		return nil, err
	}

	if c.SlotsPerEpoch, err = parser.GetUint64(
		flags.SlotsPerEpoch,
	); err != nil {
		return nil, err
	}

	if c.SlotsPerHistoricalRoot, err = parser.GetUint64(
		flags.SlotsPerHistoricalRoot,
	); err != nil {
		return nil, err
	}

	if c.DepositContractAddress, err = parser.GetExecutionAddress(
		flags.DepositContractAddress,
	); err != nil {
		return nil, err
	}

	if c.ElectraForkEpoch, err = parser.GetEpoch(
		flags.ElectraForkEpoch,
	); err != nil {
		return nil, err
	}

	if c.EpochsPerHistoricalVector, err = parser.GetUint64(
		flags.EpochsPerHistoricalVector,
	); err != nil {
		return nil, err
	}

	if c.MaxDepositsPerBlock, err = parser.GetUint64(
		flags.MaxDepositsPerBlock,
	); err != nil {
		return nil, err
	}

	if c.MaxWithdrawalsPerPayload, err = parser.GetUint64(
		flags.MaxWithdrawalsPerPayload,
	); err != nil {
		return nil, err
	}

	return &c, nil
}
