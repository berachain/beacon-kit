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

package primitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Validator as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#validator
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen --path validator.go -objs Validator -include ./pkg/crypto,./pkg/bytes,./primitives.go,./withdrawal_credentials.go,./pkg/math -output validator.ssz.go
type Validator struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey crypto.BLSPubkey `json:"pubkey"                     ssz-size:"48"`
	// WithdrawalCredentials are an address that controls the validator.
	WithdrawalCredentials WithdrawalCredentials `json:"withdrawalCredentials"      ssz-size:"32"`
	// EffectiveBalance is the validator's current effective balance in gwei.
	EffectiveBalance math.Gwei `json:"effectiveBalance"`
	// Slashed indicates whether the validator has been slashed.
	Slashed bool `json:"slashed"`
	// ActivationEligibilityEpoch is the epoch in which the validator became
	// eligible for activation.
	ActivationEligibilityEpoch math.Epoch `json:"activationEligibilityEpoch"`
	// ActivationEpoch is the epoch in which the validator activated.
	ActivationEpoch math.Epoch `json:"activationEpoch"`
	// ExitEpoch is the epoch in which the validator exited.
	ExitEpoch math.Epoch `json:"exitEpoch"`
	// WithdrawableEpoch is the epoch in which the validator can withdraw.
	WithdrawableEpoch math.Epoch `json:"withdrawableEpoch"`
}

// NewValidatorFromDeposit creates a new Validator from the
// given public key, withdrawal credentials, and amount.
//
// As defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#deposits
//
//nolint:lll
func NewValidatorFromDeposit(
	pubkey crypto.BLSPubkey,
	withdrawalCredentials WithdrawalCredentials,
	amount math.Gwei,
	effectiveBalanceIncrement math.Gwei,
	maxEffectiveBalance math.Gwei,
) *Validator {
	return &Validator{
		Pubkey:                pubkey,
		WithdrawalCredentials: withdrawalCredentials,
		EffectiveBalance: min(
			amount-amount%effectiveBalanceIncrement,
			maxEffectiveBalance,
		),
		Slashed:                    false,
		ActivationEligibilityEpoch: math.Epoch(constants.FarFutureEpoch),
		ActivationEpoch:            math.Epoch(constants.FarFutureEpoch),
		ExitEpoch:                  math.Epoch(constants.FarFutureEpoch),
		WithdrawableEpoch:          math.Epoch(constants.FarFutureEpoch),
	}
}

// GetPubkey returns the public key of the validator.
func (v *Validator) GetPubkey() crypto.BLSPubkey {
	return v.Pubkey
}

// GetEffectiveBalance returns the effective balance of the validator.
func (v *Validator) GetEffectiveBalance() math.Gwei {
	return v.EffectiveBalance
}

// IsActive as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_active_validator
//
//nolint:lll
func (v Validator) IsActive(epoch math.Epoch) bool {
	return v.ActivationEpoch <= epoch && epoch < v.ExitEpoch
}

// IsEligibleForActivation as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_eligible_for_activation_queue
//
//nolint:lll
func (v Validator) IsEligibleForActivation(
	finalizedEpoch math.Epoch,
) bool {
	return v.ActivationEligibilityEpoch <= finalizedEpoch &&
		v.ActivationEpoch == math.Epoch(constants.FarFutureEpoch)
}

// IsEligibleForActivationQueue as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_eligible_for_activation_queue
//
//nolint:lll
func (v Validator) IsEligibleForActivationQueue(
	maxEffectiveBalance math.Gwei,
) bool {
	return v.ActivationEligibilityEpoch == math.Epoch(
		constants.FarFutureEpoch,
	) &&
		v.EffectiveBalance == maxEffectiveBalance
}

// IsSlashed as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_slashable_validator
//
//nolint:lll
func (v Validator) IsSlashable(epoch math.Epoch) bool {
	return !v.Slashed && v.ActivationEpoch <= epoch &&
		epoch < v.WithdrawableEpoch
}

// IsFullyWithdrawable as defined in the Ethereum 2.0 specfication:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#is_fully_withdrawable_validator
//
//nolint:lll
func (v Validator) IsFullyWithdrawable(
	balance math.Gwei,
	epoch math.Epoch,
) bool {
	return v.HasEth1WithdrawalCredentials() && v.WithdrawableEpoch <= epoch &&
		balance > 0
}

// IsPartiallyWithdrawable as defined in the Ethereum 2.0 specfication:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#is_partially_withdrawable_validator
//
//nolint:lll
func (v Validator) IsPartiallyWithdrawable(
	balance, maxEffectiveBalance math.Gwei,
) bool {
	hasExcessBalance := balance > maxEffectiveBalance
	return v.HasEth1WithdrawalCredentials() &&
		v.HasMaxEffectiveBalance(maxEffectiveBalance) && hasExcessBalance
}

// IsWithdrawable as defined in the Ethereum 2.0 specfication:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#has_eth1_withdrawal_credential
//
//nolint:lll
func (v Validator) HasEth1WithdrawalCredentials() bool {
	return v.WithdrawalCredentials[0] == EthSecp256k1CredentialPrefix
}

// HasMaxEffectiveBalance determines if the validator has the maximum effective
// balance.
func (v Validator) HasMaxEffectiveBalance(
	maxEffectiveBalance math.Gwei,
) bool {
	return v.EffectiveBalance == maxEffectiveBalance
}
