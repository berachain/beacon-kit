package core

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// StateProcessor defines the interface for processing various state transitions
// in the beacon chain.
type StateProcessor[
	BeaconBlockT,
	BeaconStateT,
	BlobSidecarsT,
	ContextT,
	DepositT,
	ExecutionPayloadHeaderT any,
] interface {
	// InitializePreminedBeaconStateFromEth1 initializes the premined beacon
	// state
	// from the eth1 deposits.
	InitializePreminedBeaconStateFromEth1(
		st BeaconStateT,
		deposits []DepositT,
		payloadHeader ExecutionPayloadHeaderT,
		version common.Version,
	) (transition.ValidatorUpdates, error)
	// ProcessBlock processes the state transition for a given block.
	ProcessBlock(
		ctx ContextT,
		st BeaconStateT,
		blk BeaconBlockT,
	) error
	// ProcessSlots processes the state transition for a range of slots.
	ProcessSlots(
		st BeaconStateT,
		slot math.Slot,
	) (transition.ValidatorUpdates, error)
	// Transition processes the state transition for a given block.
	Transition(
		ctx ContextT,
		st BeaconStateT,
		blk BeaconBlockT,
	) (transition.ValidatorUpdates, error)
}
