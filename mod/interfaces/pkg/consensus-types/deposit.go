package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Deposit[
	T any,
	ForkDataT any,
	WithdrawalCredentialsT any,
] interface {
	constraints.SSZMarshallable
	// New creates a new deposit.
	New(
		pubkey crypto.BLSPubkey,
		credentials WithdrawalCredentialsT,
		amount math.Gwei,
		signature crypto.BLSSignature,
		index uint64,
	) T
	// VerifySignature verifies the signature of the deposit.
	VerifySignature(
		forkData ForkDataT,
		domainType common.DomainType,
		signatureVerificationFn func(
			pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
		) error,
	) error
	// GetAmount returns the amount of the deposit.
	GetAmount() math.Gwei
	// GetIndex returns the index of the deposit.
	GetIndex() uint64
	// GetPubkey returns the public key of the deposit.
	GetPubkey() crypto.BLSPubkey
	// GetSignature returns the signature of the deposit.
	GetSignature() crypto.BLSSignature
	// GetWithdrawalCredentials returns the withdrawal credentials of the deposit.
	GetWithdrawalCredentials() WithdrawalCredentialsT
}
