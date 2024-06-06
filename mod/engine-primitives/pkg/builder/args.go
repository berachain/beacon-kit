package builder

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// PayloadRequest represents the request payload for a block.
type PayloadRequest struct {
	// Attribute Args
	//
	// Slot represents the slot for which the block is proposed.
	Slot math.Slot
	// Timestamp represents the time at which the block is proposed.
	Timestamp uint64
	// ParentBlockRoot represents the root of the parent block.
	ParentBlockRoot common.Root
	// ExepectedWithdrawals represents the expected withdrawals for the block.
	ExpectedWithdrawals []common.Root
	// RandaoMix represents the Randao mix for the block.
	RandaoMix common.Root

	// State Args
	//
	// HeadEth1BlockHash represents the hash of the head Ethereum 1 block.
	HeadEth1BlockHash common.ExecutionHash
	// SafeEth1BlockHash represents the hash of the final Ethereum 1 block.
	SafeEth1BlockHash common.ExecutionHash
	// FinalEth1BlockHash represents the hash of the final Ethereum 1 block.
	FinalEth1BlockHash common.ExecutionHash

	// BuilderArgs
	//
	// ForceUpdate indicates whether to force the underlying forkchoice update
	// to attempt to build a block. This will force a payload build and will ignore
	// any caching checks, forcing payload builds can result in nil payloadIDs being
	// returned from the
	ForceUpdate bool
}
