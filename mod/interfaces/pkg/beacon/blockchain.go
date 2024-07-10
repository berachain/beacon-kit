package beacon

import (
	"context"

	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BlockchainService defines the interface for interacting with the blockchain
// state and processing blocks.
type BlockchainService[
	BeaconBlockT any,
	BlobSidecarsT constraints.SSZMarshallable,
	DepositT any,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	GenesisT types.Genesis[DepositT, ExecutionPayloadHeaderT],
] interface {
	// ProcessGenesisData processes the genesis data and initializes the beacon
	// state.
	ProcessGenesisData(
		context.Context,
		GenesisT,
	) (transition.ValidatorUpdates, error)
	// ProcessBeaconBlock processes the given beacon block and associated
	// blobs sidecars.
	ProcessBeaconBlock(
		context.Context,
		BeaconBlockT,
	) (transition.ValidatorUpdates, error)
	// ReceiveBlock receives a beacon block and
	// associated blobs sidecars for processing.
	ReceiveBlock(ctx context.Context, blk BeaconBlockT) error
	// VerifyIncomingBlock verifies the state root of an incoming block
	// and logs the process.
	VerifyIncomingBlock(ctx context.Context, blk BeaconBlockT) error
}
