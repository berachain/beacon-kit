package engineprimitives

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// PayloadAttributes represents the attributes of a beacon block payload.
type PayloadAttributes[
	T any,
	WithdrawalT any,
] interface {
	constraints.Versionable
	constraints.Nillable
	// New creates a new payload attributes instance.
	New(
		forkVersion uint32,
		timestamp uint64,
		prevRandao common.Bytes32,
		suggestedFeeRecipient gethprimitives.ExecutionAddress,
		withdrawals []WithdrawalT,
		parentBeaconBlockRoot common.Root,
	) (T, error)
	// Validate validates the payload attributes.
	Validate() error
	// GetSuggestedFeeRecipient returns the suggested fee recipient for the
	// block.
	GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress
}
