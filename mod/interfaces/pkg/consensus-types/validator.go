package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Validator[
	T any,
	WithdrawalCredentialsT any,
] interface {
	constraints.SSZMarshallable
	// New returns a new validator.
	New(
		pubkey crypto.BLSPubkey,
		withdrawalCredentials WithdrawalCredentialsT,
		amount math.Gwei,
		effectiveBalanceIncrement math.Gwei,
		maxEffectiveBalance math.Gwei,
	) T
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator.
	GetEffectiveBalance() math.Gwei
	// GetWithdrawable returns the epoch when the validator can withdraw.
	GetWithdrawableEpoch() math.Epoch
	// GetWithdrawalCredentials returns the withdrawal credentials of the validator.
	GetWithdrawalCredentials() WithdrawalCredentialsT
	// SetEffectiveBalance sets the effective balance of the validator.
	SetEffectiveBalance(effectiveBalance math.Gwei)
	// IsActive as defined in the Ethereum 2.0 Spec
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_active_validator
	IsActive(epoch math.Epoch) bool
	// IsEligibleForActivation as defined in the Ethereum 2.0 Spec
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_eligible_for_activation_queue
	IsEligibleForActivation(finalizedEpoch math.Epoch) bool
	// IsEligibleForActivationQueue as defined in the Ethereum 2.0 Spec
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_eligible_for_activation_queue
	IsEligibleForActivationQueue(maxEffectiveBalance math.Gwei) bool
	// IsSlashable as defined in the Ethereum 2.0 Spec
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#is_slashable_validator
	IsSlashable(epoch math.Epoch) bool
	// IsSlashed returns whether the validator has been slashed.
	IsSlashed() bool
	// IsFullyWithdrawable as defined in the Ethereum 2.0 specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#is_fully_withdrawable_validator
	IsFullyWithdrawable(balance math.Gwei, epoch math.Epoch) bool
	// IsPartiallyWithdrawable as defined in the Ethereum 2.0 specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#is_partially_withdrawable_validator
	IsPartiallyWithdrawable(balance, maxEffectiveBalance math.Gwei) bool
	// HasEth1WithdrawalCredentials as defined in the Ethereum 2.0 specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#has_eth1_withdrawal_credential
	HasEth1WithdrawalCredentials() bool
	// HasMaxEffectiveBalance determines if the validator has the maximum effective
	// balance.
	HasMaxEffectiveBalance(maxEffectiveBalance math.Gwei) bool
}
