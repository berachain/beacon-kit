// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package builder

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/config"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/berachain/beacon-kit/runtime/service"
)

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
	GetBestPayload(
		ctx context.Context,
		slot primitives.Slot,
		parentBlockRoot [32]byte,
		parentEth1Hash primitives.ExecutionHash,
	) (enginetypes.ExecutionPayload, *enginetypes.BlobsBundleV1, bool, error)
}

type RandaoProcessor interface {
	BuildReveal(
		ctx context.Context,
		epoch primitives.Epoch,
	) (types.Reveal, error)
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

	randaoProcessor RandaoProcessor
}

// LocalBuilder returns the local builder.
func (s *Service) LocalBuilder() PayloadBuilder {
	return s.localBuilder
}

// RequestBestBlock builds a new beacon block.
func (s *Service) RequestBestBlock(
	ctx context.Context, slot primitives.Slot,
) (beacontypes.BeaconBlock, error) {
	s.Logger().Info("our turn to propose a block 🙈", "slot", slot)
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.

	reveal, err := s.randaoProcessor.BuildReveal(
		ctx, s.BeaconCfg().SlotToEpoch(slot),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build reveal: %w", err)
	}

	parentBlockRoot := s.BeaconState(ctx).GetParentBlockRoot()

	// Create a new empty block from the current state.
	beaconBlock, err := beacontypes.EmptyBeaconBlock(
		slot, parentBlockRoot, s.ActiveForkVersionForSlot(slot), reveal,
	)
	if err != nil {
		return nil, err
	} else if beaconBlock == nil {
		return nil, beacontypes.ErrNilBlk
	}

	// Get the payload for the block.
	payload, blobsBundle, overrideBuilder, err := s.localBuilder.GetBestPayload(
		ctx,
		slot,
		parentBlockRoot,
		s.ForkchoiceStore(ctx).JustifiedPayloadBlockHash(),
	)
	if err != nil {
		return beaconBlock, err
	}

	// TODO: Dencun
	_ = blobsBundle

	// TODO: allow external block builders to override the payload.
	_ = overrideBuilder

	// Assemble a new block with the payload.
	body := beaconBlock.GetBody()
	if body.IsNil() {
		return nil, beacontypes.ErrNilBlkBody
	}

	// Dequeue deposits from the state.
	deposits, err := s.BeaconState(ctx).ExpectedDeposits(
		s.BeaconCfg().Limits.MaxDepositsPerBlock,
	)
	if err != nil {
		return nil, err
	}

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	// if err = b
	if err = body.SetExecutionData(payload); err != nil {
		return nil, err
	}

	s.Logger().Info("finished assembling beacon block 🛟",
		"slot", slot,
		"deposits", len(deposits),
	)

	// Return the block.
	return beaconBlock, nil
}
