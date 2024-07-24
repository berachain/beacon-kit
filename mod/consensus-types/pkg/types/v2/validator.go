// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

// Validator as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#validator
//
//nolint:lll
type Validator struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// WithdrawalCredentials are an address that controls the validator.
	WithdrawalCredentials WithdrawalCredentials `json:"withdrawalCredentials"`
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

// New creates a new Validator with the given public key, withdrawal
// credentials,.
func (v *Validator) New(
	pubkey crypto.BLSPubkey,
	withdrawalCredentials WithdrawalCredentials,
	amount math.Gwei,
	effectiveBalanceIncrement math.Gwei,
	maxEffectiveBalance math.Gwei,
) *Validator {
	return NewValidatorFromDeposit(
		pubkey,
		withdrawalCredentials,
		amount,
		effectiveBalanceIncrement,
		maxEffectiveBalance,
	)
}

// SizeSSZ returns the size of the Validator object for SSZ encoding.
func (v *Validator) SizeSSZ() (size uint32) {
	size = 121
	// PublicKey (48),
	// WithdrawalCredentials (32),
	// EffectiveBalance (8),
	// Slashed (1),
	// ActivationEligibilityEpoch (8),
	// ActivationEpoch (8),
	// ExitEpoch (8),
	// WithdrawableEpoch (8)
	return
}

// DefineSSZ defines the SSZ encoding for the Validator object.
func (v *Validator) DefineSSZ(c *ssz.Codec) {
	ssz.DefineStaticBytes(c, &v.Pubkey)
	ssz.DefineStaticBytes(c, &v.WithdrawalCredentials)
	ssz.DefineUint64(c, &v.EffectiveBalance)
	ssz.DefineBool(c, &v.Slashed)
	ssz.DefineUint64(c, &v.ActivationEligibilityEpoch)
	ssz.DefineUint64(c, &v.ActivationEpoch)
	ssz.DefineUint64(c, &v.ExitEpoch)
	ssz.DefineUint64(c, &v.WithdrawableEpoch)
}

// MarshalSSZTo serializes the Validator object into a writer.
func (v *Validator) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, v)
}

// MarshalSSZ serializes the Validator object into a writer.
func (v *Validator) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, v.SizeSSZ())
	return v.MarshalSSZTo(buf)
}

// UnmarshalSSZ deserializes the Validator object from SSZ encoding.
func (v *Validator) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, v)
}

// HashTreeRoot returns the hash tree root of the Validator object.
func (v *Validator) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(v), nil
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

// IsSlashable as defined in the Ethereum 2.0 Spec
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_slashable_validator
//
//nolint:lll
func (v Validator) IsSlashable(epoch math.Epoch) bool {
	return !v.Slashed && v.ActivationEpoch <= epoch &&
		epoch < v.WithdrawableEpoch
}

// IsSlashed returns whether the validator has been slashed.
func (v Validator) IsSlashed() bool {
	return v.Slashed
}

// IsFullyWithdrawable as defined in the Ethereum 2.0 specification:
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

// IsPartiallyWithdrawable as defined in the Ethereum 2.0 specification:
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

// HasEth1WithdrawalCredentials as defined in the Ethereum 2.0 specification:
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

// SetEffectiveBalance sets the effective balance of the validator.
func (v *Validator) SetEffectiveBalance(balance math.Gwei) {
	v.EffectiveBalance = balance
}

// GetWithdrawableEpoch returns the epoch when the validator can withdraw.
func (v Validator) GetWithdrawableEpoch() math.Epoch {
	return v.WithdrawableEpoch
}

// GetWithdrawalCredentials returns the withdrawal credentials of the validator.
func (v Validator) GetWithdrawalCredentials() WithdrawalCredentials {
	return v.WithdrawalCredentials
}
