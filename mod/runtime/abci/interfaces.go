package abci

import (
	"context"

	"github.com/berachain/beacon-kit/mod/core/state"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
)

type BuilderService interface {
	RequestBestBlock(
		context.Context,
		state.BeaconState,
		math.Slot,
	) (beacontypes.BeaconBlock, *datypes.BlobSidecars, error)
}

type BlockchainService interface {
	ProcessSlot(state.BeaconState) error
	BeaconState(context.Context) state.BeaconState
	ProcessBeaconBlock(
		context.Context,
		state.BeaconState,
		beacontypes.ReadOnlyBeaconBlock,
		*datypes.BlobSidecars,
	) error
	PostBlockProcess(
		context.Context,
		state.BeaconState,
		beacontypes.ReadOnlyBeaconBlock,
	) error
	ChainSpec() primitives.ChainSpec
	ValidatePayloadOnBlk(context.Context, beacontypes.ReadOnlyBeaconBlock) error
}
