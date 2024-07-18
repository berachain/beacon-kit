package storage

import (
	"context"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// BlockStore is the interface for block storage.
type BlockStore[BeaconBlockT any] interface {
	// Get retrieves the block at the given slot.
	Get(slot uint64) (BeaconBlockT, error)
	// Set sets the block at the given slot.
	Set(slot uint64, block BeaconBlockT) error
	// Prune prunes the block store of [start, end).
	Prune(start, end uint64) error
}

// Deposit is a struct that represents a deposit.
type Deposit interface {
	constraints.SSZMarshallable
	GetIndex() uint64
}

// DepositStore defines the interface for deposit storage.
type DepositStore[DepositT any] interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]DepositT, error)
	// Prune prunes the deposit store of [start, end)
	Prune(start, end uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

// Validator represents an interface for a validator with generic withdrawal
// credentials. WithdrawalCredentialsT is a type parameter that must implement
// the WithdrawalCredentials interface.
type Validator[WithdrawalCredentialsT WithdrawalCredentials] interface {
	// GetWithdrawalCredentials returns the withdrawal credentials of the
	// validator.
	GetWithdrawalCredentials() WithdrawalCredentialsT
	// IsFullyWithdrawable checks if the validator is fully withdrawable given a
	// certain Gwei amount and epoch.
	IsFullyWithdrawable(amount math.Gwei, epoch math.Epoch) bool
	// IsPartiallyWithdrawable checks if the validator is partially withdrawable
	// given two Gwei amounts.
	IsPartiallyWithdrawable(amount1 math.Gwei, amount2 math.Gwei) bool
}

// Withdrawal represents an interface for a withdrawal.
type Withdrawal[T any] interface {
	New(
		index math.U64,
		validator math.ValidatorIndex,
		address gethprimitives.ExecutionAddress,
		amount math.Gwei,
	) T
}

// WithdrawalCredentials represents an interface for withdrawal credentials.
type WithdrawalCredentials interface {
	// ToExecutionAddress converts the withdrawal credentials to an execution
	// address.
	ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
}
