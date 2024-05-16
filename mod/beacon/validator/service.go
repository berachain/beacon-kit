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
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"golang.org/x/sync/errgroup"
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

	// bsb is the beacon state backend.
	bsb BeaconStorageBackend[BeaconStateT]

	// randaoProcessor is responsible for building the reveal for the
	// current slot.
	randaoProcessor RandaoProcessor[BeaconStateT]

	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor[BeaconStateT]

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
	bsb BeaconStorageBackend[BeaconStateT],
	stateProcessor StateProcessor[BeaconStateT],
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
		bsb:             bsb,
		chainSpec:       chainSpec,
		signer:          signer,
		stateProcessor:  stateProcessor,
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
	requestedSlot math.Slot,
) (types.BeaconBlock, BlobSidecarsT, error) {
	var (
		sidecars  BlobSidecarsT
		startTime = time.Now()
		g, gCtx   = errgroup.WithContext(ctx)
	)

	s.logger.Info("requesting beacon block assembly 🙈", "slot", requestedSlot)
	defer func() {
		s.logger.Info(
			"finished beacon block assembly 🛟 ",
			"slot",
			requestedSlot,
			"duration",
			time.Since(startTime).String(),
		)
	}()

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	st := s.bsb.StateFromContext(ctx)

	// Get the current state slot.
	stateSlot, err := st.GetSlot()
	if err != nil {
		return nil, sidecars, err
	}

	slotDifference := requestedSlot - stateSlot
	switch {
	case slotDifference == 1:
		// If our BeaconState is not up to date, we need to process
		// a slot to get it up to date.
		if err = s.stateProcessor.ProcessSlot(st); err != nil {
			return nil, sidecars, err
		}

		// Request the slot again, it should've been incremented by 1.
		stateSlot, err = st.GetSlot()
		if err != nil {
			return nil, sidecars, err
		}

		// If after doing so, we aren't exactly at the requested slot,
		// we should return an error.
		if requestedSlot != stateSlot {
			return nil, sidecars, errors.Newf(
				"requested slot could not be processed, requested: %d, state: %d",
				requestedSlot,
				stateSlot,
			)
		}
	case slotDifference > 1:
		return nil, sidecars, errors.Newf(
			"requested slot is too far ahead, requested: %d, state: %d",
			requestedSlot,
			stateSlot,
		)
	case slotDifference < 1:
		return nil, sidecars, errors.Newf(
			"requested slot is in the past, requested: %d, state: %d",
			requestedSlot,
			stateSlot,
		)
	}

	// Create a new empty block from the current state.
	blk, err := s.GetEmptyBeaconBlock(
		st,
		requestedSlot,
	)
	if err != nil {
		return nil, sidecars, err
	}

	// Get the payload for the block.
	envelope, err := s.RetrievePayload(ctx, st, blk)
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

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// Build the reveal for the current slot.
	// TODO: We can optimize to pre-compute this in parallel.
	reveal, err := s.randaoProcessor.BuildReveal(st)
	if err != nil {
		return nil, sidecars, errors.Newf("failed to build reveal: %w", err)
	}

	// Set the reveal on the block body.
	body.SetRandaoReveal(reveal)

	// Dequeue deposits from the state.
	deposits, err := s.ds.ExpectedDeposits(
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return nil, sidecars, err
	}

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// TODO: assemble real eth1data.
	body.SetEth1Data(&types.Eth1Data{
		DepositRoot:  primitives.Bytes32{},
		DepositCount: 0,
		BlockHash:    common.ZeroHash,
	})

	// Set the execution data.
	if err = body.SetExecutionData(
		envelope.GetExecutionPayload(),
	); err != nil {
		return nil, sidecars, err
	}

	// Produce block sidecars.
	g.Go(func() error {
		var sidecarErr error
		sidecars, sidecarErr = s.BuildSidecars(
			blk,
			envelope.GetBlobsBundle(),
		)
		return sidecarErr
	})

	// Set the state root on the BeaconBlock.
	g.Go(func() error {
		return s.SetBlockStateRoot(gCtx, st, blk)
	})

	// Set the execution payload on the block body.
	return blk, sidecars, g.Wait()
}
