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

	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/services/builder/config"
)

// Service is responsible for building beacon blocks.
type Service struct {
	service.BaseService
	cfg *config.Config

	// signer is used to retrieve the public key of this node.
	signer core.BLSSigner

	// localBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks is done by submitting forkchoice updates through.
	// The local Builder.
	localBuilder PayloadBuilder

	// remoteBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remoteBuilders []PayloadBuilder

	// randaoProcessor is responsible for building the reveal for the
	// current slot.
	randaoProcessor RandaoProcessor
}

// LocalBuilder returns the local builder.
func (s *Service) LocalBuilder() PayloadBuilder {
	return s.localBuilder
}

// RequestBestBlock builds a new beacon block.
//
//nolint:funlen // todo:fix.
func (s *Service) RequestBestBlock(
	ctx context.Context,
	st state.BeaconState,
	slot primitives.Slot,
) (beacontypes.BeaconBlock, *datypes.BlobSidecars, error) {
	s.Logger().Info("our turn to propose a block 🙈", "slot", slot)
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	reveal, err := s.randaoProcessor.BuildReveal(st)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build reveal: %w", err)
	}

	parentBlockRoot, err := st.GetBlockRootAtIndex(
		uint64(slot) % s.BeaconCfg().SlotsPerHistoricalRoot,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get block root at index: %w",
			err,
		)
	}
	// Get the proposer index for the slot.
	proposerIndex, err := st.ValidatorIndexByPubkey(
		s.signer.PublicKey(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get validator by pubkey: %w",
			err,
		)
	}

	// Compute the state root for the block.
	stateRoot, err := s.computeStateRoot(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to compute state root: %w",
			err,
		)
	}

	// Create a new empty block from the current state.
	blk, err := beacontypes.EmptyBeaconBlock(
		slot,
		proposerIndex,
		parentBlockRoot,
		stateRoot,
		s.BeaconCfg().ActiveForkVersionForSlot(slot),
		reveal,
	)
	if err != nil {
		return nil, nil, err
	} else if blk == nil {
		return nil, nil, beacontypes.ErrNilBlk
	}

	parentEth1BlockHash, err := st.GetEth1BlockHash()
	if err != nil {
		return nil, nil, err
	}

	// Get the payload for the block.
	payload, blobsBundle, overrideBuilder, err := s.localBuilder.GetBestPayload(
		ctx,
		st,
		slot,
		parentBlockRoot,
		parentEth1BlockHash,
	)
	if err != nil {
		return blk, nil, fmt.Errorf(
			"failed to get block root at index: %w",
			err,
		)
	}

	// TODO: allow external block builders to override the payload.
	_ = overrideBuilder

	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body.IsNil() {
		return nil, nil, beacontypes.ErrNilBlkBody
	}

	// If we get returned a nil blobs bundle, we should return an error.
	if blobsBundle == nil {
		return nil, nil, beacontypes.ErrNilBlobsBundle
	}

	commitments := make([][48]byte, len(blobsBundle.Commitments))
	for i, c := range blobsBundle.Commitments {
		commitments[i] = [48]byte(c)
	}
	body.SetBlobKzgCommitments(commitments)

	// Dequeue deposits from the state.
	deposits, err := st.ExpectedDeposits(
		s.BeaconCfg().MaxDepositsPerBlock,
	)
	if err != nil {
		return nil, nil, err
	}

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	// if err = b
	if err = body.SetExecutionData(payload); err != nil {
		return nil, nil, err
	}

	// Build the blob sidecars.
	blobSidecars, err := beacontypes.BuildBlobSidecar(blk, blobsBundle)
	if err != nil {
		return nil, nil, err
	}

	s.Logger().Info("finished assembling beacon block 🛟",
		"slot", slot, "deposits", len(deposits))

	return blk, blobSidecars, nil
}
