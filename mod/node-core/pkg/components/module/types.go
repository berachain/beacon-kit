package beacon

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"go.starlark.net/lib/proto"
)

// AttestationData is an interface for accessing the attestation data.
type AttestationData[AttestationDataT any] interface {
	// GetIndex returns the index of the attestation data.
	GetIndex() math.U64
	// New creates a new attestation data instance.
	New(math.U64, math.U64, common.Root) AttestationDataT
}

// BeaconState is an interface for accessing the beacon state.
type BeaconState interface {
	// GetValidatorIndexByCometBFTAddress returns the validator index by the
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
	// HashTreeRoot returns the hash tree root of the beacon state.
	HashTreeRoot() ([32]byte, error)
}

// Middleware is the interface for the CometBFT middleware.
type Middleware[
	AttestationDataT,
	SlashingInfoT any,
	SlotDataT SlotData[AttestationDataT, SlashingInfoT, SlotDataT],
] interface {
	InitGenesis(
		ctx context.Context, bz []byte,
	) (transition.ValidatorUpdates, error)
	PrepareProposal(context.Context, SlotDataT) ([]byte, []byte, error)
	ProcessProposal(
		ctx context.Context, req proto.Message,
	) error
	EndBlock(ctx context.Context) (transition.ValidatorUpdates, error)
}

// SlashingInfo is an interface for accessing the slashing info.
type SlashingInfo[SlashingInfoT any] interface {
	// New creates a new slashing info instance.
	New(math.U64, math.U64) SlashingInfoT
}

// SlotData is an interface for accessing the slot data.
type SlotData[AttestationDataT, SlashingInfoT, SlotDataT any] interface {
	// New creates a new slot data instance.
	New(math.Slot, []AttestationDataT, []SlashingInfoT) SlotDataT
}

// StorageBackend defines an interface for accessing various storage components
// required by the beacon node.
type StorageBackend[BeaconStateT BeaconState] interface {
	// StateFromContext retrieves the beacon state from the given context.
	StateFromContext(context.Context) BeaconStateT
}
