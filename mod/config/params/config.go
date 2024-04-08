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
	"github.com/berachain/beacon-kit/mod/primitives"
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
	// EpochsPerSlashingsVector is the number of epochs in the slashings vector.
	EpochsPerSlashingsVector uint64 `mapstructure:"epochs-per-slashings-vector"`

	// Max operations per block constants.
	//
	// MaxDepositsPerBlock specifies the maximum number of deposit operations
	// allowed per block.
	MaxDepositsPerBlock uint64 `mapstructure:"max-deposits-per-block"`
	// MaxWithdrawalsPerPayload indicates the maximum number of withdrawal
	// operations allowed in a single payload.
	MaxWithdrawalsPerPayload uint64 `mapstructure:"max-withdrawals-per-payload"`
	// MaxBlobsPerBlock specifies the maximum number of blobs allowed per block.
	MaxBlobsPerBlock uint64 `mapstructure:"max-blobs-per-block"`

	// Rewards and penalties constants.
	//
	// ProportionalSlashingMultiplier is the slashing multiplier relative to the
	// base penalty.
	ProportionalSlashingMultiplier uint64 `mapstructure:"proportional-slashing-multiplier"`

	// Deneb Values
	//
	// MinEpochsForBlobsSidecarsRequest is the minimum number of epochs the node
	// will keep the blobs for.
	MinEpochsForBlobsSidecarsRequest uint64 `mapstructure:"min-epochs-for-blobs-sidecars-request"`
}
