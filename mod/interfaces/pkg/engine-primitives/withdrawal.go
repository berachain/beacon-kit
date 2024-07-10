package engineprimitives

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Withdrawal represents a validator withdrawal from the consensus layer.
type Withdrawal[T any] interface {
	// New creates a new Withdrawal instance.
	New(
		index math.U64,
		validator math.ValidatorIndex,
		address gethprimitives.ExecutionAddress,
		amount math.Gwei,
	) T
	// Equals returns true if the Withdrawal is equal to the other.
	Equals(other T) bool
	// GetIndex returns the unique identifier for the withdrawal.
	GetIndex() math.U64
	// GetValidatorIndex returns the index of the validator initiating the
	// withdrawal.
	GetValidatorIndex() math.ValidatorIndex
	// GetAddress returns the execution address where the withdrawal will be sent.
	GetAddress() gethprimitives.ExecutionAddress
	// GetAmount returns the amount of Gwei to be withdrawn.
	GetAmount() math.Gwei
	// IsFixed returns true if the Withdrawal has a fixed size.
	IsFixed() bool
	// Type returns the type of the Withdrawal.
	Type() types.Type
}
