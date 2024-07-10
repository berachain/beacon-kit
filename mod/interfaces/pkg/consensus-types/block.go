package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock is the interface for a beacon block.
type BeaconBlock[
	T any,
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT any,
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	RawBeaconBlock[
		BeaconBlockBodyT, BeaconBlockHeaderT, DepositT, Eth1DataT, ExecutionPayloadT,
	]
	constraints.EmptyWithVersion[T]
	constraints.NewFromSSZable[T]
	// NewWithVersion creates a new beacon block with a given version.
	NewWithVersion(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		forkVersion uint32,
	) (T, error)
}

// RawBeaconBlock is the interface for a beacon block.
type RawBeaconBlock[
	BeaconBlockBodyT RawBeaconBlockBody[DepositT, Eth1DataT, ExecutionPayloadT],
	BeaconBlockHeaderT any,
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] interface {
	constraints.SSZMarshallable
	constraints.Nillable
	constraints.Versionable
	// SetStateRoot sets the state root of the block.
	SetStateRoot(common.Root)
	// GetStateRoot returns the state root of the block.
	GetStateRoot() common.Root
	// GetSlot returns the slot of the block.
	GetSlot() math.Slot
	// GetProposerIndex returns the proposer index of the block.
	GetProposerIndex() math.ValidatorIndex
	// GetParentBlockRoot returns the parent block root of the block.
	GetParentBlockRoot() common.Root
	// GetBody returns the body of the block.
	GetBody() BeaconBlockBodyT
	// GetHeader returns the header of the block.
	GetHeader() BeaconBlockHeaderT
}
