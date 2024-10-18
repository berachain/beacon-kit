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

package consensustypes

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// Validator is the custom Validator struct used
// in Berachain.
//
//nolint:lll // long withdrawal credentials
type Validator struct {
	// Pubkey is the validator's 48-byte BLS public key.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// WithdrawalCredentials are an address that controls the validator.
	WithdrawalCredentials types.WithdrawalCredentials `json:"withdrawal_credentials"`
	// EffectiveBalance is the validator's current effective balance in gwei.
	EffectiveBalance math.Gwei `json:"effective_balance"`
	// ActivationEpoch is the epoch in which the validator activated.
	ActivationEpoch math.Epoch `json:"activation_epoch"`
	// ExitEpoch is the epoch in which the validator exited.
	ExitEpoch math.Epoch `json:"exit_epoch"`
	// WithdrawableEpoch is the epoch in which the validator can withdraw.
	WithdrawableEpoch math.Epoch `json:"withdrawable_epoch"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

func NewValidatorFromDeposit(
	pubkey crypto.BLSPubkey,
	withdrawalCredentials types.WithdrawalCredentials,
	effectiveBalance math.Gwei,
	activationEpoch math.Epoch,
) *Validator {
	return &Validator{
		Pubkey:                pubkey,
		WithdrawalCredentials: withdrawalCredentials,
		EffectiveBalance:      effectiveBalance,
		ActivationEpoch:       activationEpoch,
	}
}

// Empty creates an empty Validator.
func (*Validator) Empty() *Validator {
	return &Validator{}
}

// New creates a new Validator with the given public key, withdrawal
// credentials,.
func (v *Validator) New(
	pubkey crypto.BLSPubkey,
	withdrawalCredentials types.WithdrawalCredentials,
	amount math.Gwei,
	effectiveBalanceIncrement math.Gwei,
	_ math.Gwei,
) *Validator {
	return NewValidatorFromDeposit(
		pubkey,
		withdrawalCredentials,
		amount,
		effectiveBalanceIncrement,
	)
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// ValidatorSize is the size of the Validator struct in bytes.
const ValidatorSize = 112

// SizeSSZ returns the size of the Validator object in SSZ encoding.
func (*Validator) SizeSSZ() uint32 {
	return ValidatorSize
}

// DefineSSZ defines the SSZ encoding for the Validator object.
func (v *Validator) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &v.Pubkey)
	ssz.DefineStaticBytes(codec, &v.WithdrawalCredentials)
	ssz.DefineUint64(codec, &v.EffectiveBalance)
	ssz.DefineUint64(codec, &v.ActivationEpoch)
	ssz.DefineUint64(codec, &v.ExitEpoch)
	ssz.DefineUint64(codec, &v.WithdrawableEpoch)
}

// HashTreeRoot computes the SSZ hash tree root of the Validator object.
func (v *Validator) HashTreeRoot() common.Root {
	return ssz.HashSequential(v)
}

// MarshalSSZ marshals the Validator object to SSZ format.
func (v *Validator) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, v.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, v)
}

// UnmarshalSSZ unmarshals the Validator object from SSZ format.
func (v *Validator) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, v)
}

// HashTreeRootWith ssz hashes the Validator object with a hasher.
func (v *Validator) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Pubkey'
	hh.PutBytes(v.Pubkey[:])

	// Field (1) 'WithdrawalCredentials'
	hh.PutBytes(v.WithdrawalCredentials[:])

	// Field (2) 'EffectiveBalance'
	hh.PutUint64(uint64(v.EffectiveBalance))

	// Field (3) 'ActivationEpoch'
	hh.PutUint64(uint64(v.ActivationEpoch))

	// Field (4) 'ExitEpoch'
	hh.PutUint64(uint64(v.ExitEpoch))

	// Field (5) 'WithdrawableEpoch'
	hh.PutUint64(uint64(v.WithdrawableEpoch))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Validator object.
func (v *Validator) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(v)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

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
	return v.WithdrawalCredentials[0] == types.EthSecp256k1CredentialPrefix
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
func (v Validator) GetWithdrawalCredentials() types.WithdrawalCredentials {
	return v.WithdrawalCredentials
}

// GetActivationEpoch returns the activation epoch of the validator.
func (v *Validator) GetActivationEpoch() math.Epoch {
	return v.ActivationEpoch
}

// GetExitEpoch returns the exit epoch of the validator.
func (v *Validator) GetExitEpoch() math.Epoch {
	return v.ExitEpoch
}
