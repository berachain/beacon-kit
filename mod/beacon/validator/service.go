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

package validator

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Service is responsible for building beacon blocks.
type Service[
	BeaconStateT BeaconState,
	BlobSidecarsT BlobSidecars,
] struct {
	// cfg is the validator config.
	cfg *Config

	// logger is a logger.
	logger log.Logger[any]

	// chainSpec is the chain spec.
	chainSpec primitives.ChainSpec

	// signer is used to retrieve the public key of this node.
	signer crypto.BLSSigner

	// blobFactory is used to create blob sidecars for blocks.
	blobFactory BlobFactory[BlobSidecarsT, types.BeaconBlockBody]

	// randaoProcessor is responsible for building the reveal for the
	// current slot.
	randaoProcessor RandaoProcessor[BeaconStateT]

	// ds is used to retrieve deposits that have been
	// queued up for inclusion in the next block.
	ds DepositStore

	// localBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks is done by submitting forkchoice updates through.
	// The local Builder.
	localBuilder PayloadBuilder[BeaconStateT]

	// remoteBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remoteBuilders []PayloadBuilder[BeaconStateT]
}

// NewService creates a new validator service.
func NewService[
	BeaconStateT BeaconState,
	BlobSidecarsT BlobSidecars,
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	blobFactory BlobFactory[BlobSidecarsT, types.BeaconBlockBody],
	randaoProcessor RandaoProcessor[BeaconStateT],
	ds DepositStore,
	localBuilder PayloadBuilder[BeaconStateT],
	remoteBuilders []PayloadBuilder[BeaconStateT],
) *Service[BeaconStateT, BlobSidecarsT] {
	return &Service[BeaconStateT, BlobSidecarsT]{
		cfg:             cfg,
		logger:          logger,
		chainSpec:       chainSpec,
		signer:          signer,
		blobFactory:     blobFactory,
		randaoProcessor: randaoProcessor,
		ds:              ds,
		localBuilder:    localBuilder,
		remoteBuilders:  remoteBuilders,
	}
}

// Name returns the name of the service.
func (s *Service[BeaconStateT, BlobSidecarsT]) Name() string {
	return "validator"
}

// Start starts the service.
func (s *Service[BeaconStateT, BlobSidecarsT]) Start(context.Context) {}

// Status returns the status of the service.
func (s *Service[BeaconStateT, BlobSidecarsT]) Status() error { return nil }

// WaitForHealthy waits for the service to become healthy.
func (s *Service[BeaconStateT, BlobSidecarsT]) WaitForHealthy(
	context.Context,
) {
}

// RequestBestBlock builds a new beacon block.
//
//nolint:funlen // todo:fix.
func (s *Service[BeaconStateT, BlobSidecarsT]) RequestBestBlock(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
) (types.BeaconBlock, BlobSidecarsT, error) {
	var sidecars BlobSidecarsT
	s.logger.Info("our turn to propose a block 🙈", "slot", slot)
	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	reveal, err := s.randaoProcessor.BuildReveal(st)
	if err != nil {
		return nil, sidecars, errors.Newf("failed to build reveal: %w", err)
	}

	parentBlockRoot, err := st.GetBlockRootAtIndex(
		uint64(slot) % s.chainSpec.SlotsPerHistoricalRoot(),
	)
	if err != nil {
		return nil, sidecars, errors.Newf(
			"failed to get block root at index: %w",
			err,
		)
	}
	// Get the proposer index for the slot.
	proposerIndex, err := st.ValidatorIndexByPubkey(
		s.signer.PublicKey(),
	)
	if err != nil {
		return nil, sidecars, errors.Newf(
			"failed to get validator by pubkey: %w",
			err,
		)
	}

	// Compute the state root for the block.
	// TODO: IMPLEMENT RN THIS DOES NOTHING.
	stateRoot, err := s.computeStateRoot(ctx)
	if err != nil {
		return nil, sidecars, errors.Newf(
			"failed to compute state root: %w",
			err,
		)
	}

	// Create a new empty block from the current state.
	blk, err := types.EmptyBeaconBlock(
		slot,
		proposerIndex,
		parentBlockRoot,
		stateRoot,
		s.chainSpec.ActiveForkVersionForSlot(slot),
	)
	if err != nil {
		return nil, sidecars, err
	}

	// The latest execution payload header, will be from the previous block
	// during the block building phase.
	parentExecutionPayload, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, sidecars, err
	}

	// Get the payload for the block.
	envelope, err := s.localBuilder.RetrieveOrBuildPayload(
		ctx,
		st,
		slot,
		parentBlockRoot,
		parentExecutionPayload.GetBlockHash(),
	)
	if err != nil {
		return blk, sidecars, errors.Newf(
			"failed to get block root at index: %w",
			err,
		)
	} else if envelope == nil {
		return nil, sidecars, ErrNilPayload
	}

	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body.IsNil() {
		return nil, sidecars, ErrNilBlkBody
	}

	// If we get returned a nil blobs bundle, we should return an error.
	// TODO: allow external block builders to override the payload.
	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle == nil {
		return nil, sidecars, ErrNilBlobsBundle
	}

	// Dequeue deposits from the state.
	deposits, err := s.ds.ExpectedDeposits(
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return nil, sidecars, err
	}

	payload := envelope.GetExecutionPayload()
	if payload == nil || payload.IsNil() {
		return nil, sidecars, ErrNilPayload
	}

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	// TODO: assemble real eth1data.
	body.SetEth1Data(&types.Eth1Data{
		DepositRoot:  primitives.Bytes32{},
		DepositCount: 0,
		BlockHash:    common.ExecutionHash{},
	})

	// Set the reveal on the block body.
	body.SetRandaoReveal(reveal)

	// Set the execution data.
	if err = body.SetExecutionData(payload); err != nil {
		return nil, sidecars, err
	}

	// Build the sidecars for the block.
	blobSidecars, err := s.blobFactory.BuildSidecars(blk, blobsBundle)
	if err != nil {
		return nil, sidecars, err
	}

	s.logger.Info("finished assembling beacon block 🛟",
		"slot", slot, "deposits", len(deposits))

	// Set the execution payload on the block body.
	return blk, blobSidecars, nil
}
