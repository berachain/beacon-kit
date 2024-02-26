package builder

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/config"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
	GetBestPayload(
		ctx context.Context,
		slot primitives.Slot,
		parentBlockRoot [32]byte,
		parentEth1Hash common.Hash,
	) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error)
}

// Service is responsible for building beacon blocks.
type Service struct {
	service.BaseService
	cfg *config.Builder

	// localBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks is done by submitting forkchoice updates through.
	// The local Builder.
	localBuilder   PayloadBuilder
	remoteBuilders []PayloadBuilder
}

// LocalBuilder returns the local builder.
func (s *Service) LocalBuilder() PayloadBuilder {
	return s.localBuilder
}

// RequestBestBlock builds a new beacon block.
func (s *Service) RequestBestBlock(
	ctx context.Context, slot primitives.Slot,
) (consensus.BeaconKitBlock, error) {
	s.Logger().Info("our turn to propose a block 🙈", "slot", slot)
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.

	// // // TODO: SIGN UR RANDAO THINGY HERE OR SOMETHING.
	// _ = s.beaconKitValKey
	// // _, err := s.beaconKitValKey.Key.PrivKey.Sign([]byte("hello world"))
	// // if err != nil {
	// // 	return nil, err
	// // }

	parentBlockRoot := s.BeaconState(ctx).GetParentBlockRoot()

	// Create a new empty block from the current state.
	beaconBlock, err := consensus.EmptyBeaconKitBlock(
		slot, parentBlockRoot, s.ActiveForkVersionForSlot(slot),
	)
	if err != nil {
		return nil, err
	}

	// TODO: right now parent is ALWAYS the previously finalized. But
	// maybe not forever?
	parentEth1Hash := s.BeaconState(ctx).GetFinalizedEth1BlockHash()

	// Get the payload for the block.
	payload, blobsBundle, overrideBuilder, err := s.localBuilder.GetBestPayload(
		ctx, slot, parentBlockRoot, parentEth1Hash,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Dencun
	_ = blobsBundle

	// TODO: allow external block builders to override the payload.
	_ = overrideBuilder

	// Assemble a new block with the payload.
	if err = beaconBlock.AttachExecution(payload); err != nil {
		return nil, err
	}

	// Return the block.
	return beaconBlock, nil
}
