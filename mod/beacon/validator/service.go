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
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"golang.org/x/sync/errgroup"
)

// Service is responsible for building beacon blocks.
type Service[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[*types.Deposit, *types.Eth1Data],
	BeaconStateT BeaconState[BeaconStateT],
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
	blobFactory BlobFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
	]

	// bsb is the beacon state backend.
	bsb StorageBackend[BeaconStateT]

	// randaoProcessor is responsible for building the reveal for the
	// current slot.
	randaoProcessor RandaoProcessor[BeaconStateT]

	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
	]

	// ds is used to retrieve deposits that have been
	// queued up for inclusion in the next block.
	ds DepositStore[*types.Deposit]

	// localPayloadBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks is done by submitting forkchoice updates through.
	// The local Builder.
	localPayloadBuilder PayloadBuilder[BeaconStateT]

	// remotePayloadBuilders represents a list of remote block builders, these
	// builders are connected to other execution clients via the EngineAPI.
	remotePayloadBuilders []PayloadBuilder[BeaconStateT]
}

// NewService creates a new validator service.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[*types.Deposit, *types.Eth1Data],
	BeaconStateT BeaconState[BeaconStateT],
	BlobSidecarsT BlobSidecars,
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	bsb StorageBackend[BeaconStateT],
	stateProcessor StateProcessor[BeaconBlockT, BeaconStateT, *transition.Context],
	signer crypto.BLSSigner,
	blobFactory BlobFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
	],
	randaoProcessor RandaoProcessor[BeaconStateT],
	ds DepositStore[*types.Deposit],
	localPayloadBuilder PayloadBuilder[BeaconStateT],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT],
) *Service[BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT] {
	return &Service[BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT]{
		cfg:                   cfg,
		logger:                logger,
		bsb:                   bsb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		randaoProcessor:       randaoProcessor,
		ds:                    ds,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
	}
}

// Name returns the name of the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) Name() string {
	return "validator"
}

// Start starts the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) Start(
	context.Context,
) error {
	return nil
}

// Status returns the status of the service.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) Status() error {
	return nil
}

// WaitForHealthy waits for the service to become healthy.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) WaitForHealthy(
	context.Context,
) {
}

// RequestBestBlock builds a new beacon block.
//
//nolint:funlen // todo:fix.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) RequestBestBlock(
	ctx context.Context,
	requestedSlot math.Slot,
) (BeaconBlockT, BlobSidecarsT, error) {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		startTime = time.Now()
		g, _      = errgroup.WithContext(ctx)
	)

	s.logger.Info("requesting beacon block assembly 🙈", "slot", requestedSlot)

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propogating the finalized
	// and safe block hashes to the execution client.
	st := s.bsb.StateFromContext(ctx)

	// Prepare the state such that it is ready to build a block for
	// the request slot
	if err := s.prepareStateForBuilding(st, requestedSlot); err != nil {
		return blk, sidecars, err
	}

	// Build the reveal for the current slot.
	// TODO: We can optimize to pre-compute this in parallel.
	reveal, err := s.randaoProcessor.BuildReveal(st)
	if err != nil {
		return blk, sidecars, err
	}

	// Create a new empty block from the current state.
	blk, err = s.GetEmptyBeaconBlock(
		st, requestedSlot,
	)
	if err != nil {
		return blk, sidecars, err
	}

	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body.IsNil() {
		return blk, sidecars, ErrNilBlkBody
	}

	// Set the reveal on the block body.
	body.SetRandaoReveal(reveal)

	// Get the payload for the block.
	envelope, err := s.retrievePayload(ctx, st, blk)
	if err != nil {
		return blk, sidecars, err
	} else if envelope == nil {
		return blk, sidecars, ErrNilPayload
	}

	// If we get returned a nil blobs bundle, we should return an error.
	// TODO: allow external block builders to override the payload.
	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle == nil {
		return blk, sidecars, ErrNilBlobsBundle
	}

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// Dequeue deposits from the state.
	deposits, err := s.ds.ExpectedDeposits(
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return blk, sidecars, err
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
		return blk, sidecars, err
	}

	// Produce block sidecars.
	g.Go(func() error {
		var sidecarErr error
		sidecars, sidecarErr = s.blobFactory.BuildSidecars(
			blk,
			envelope.GetBlobsBundle(),
		)
		return sidecarErr
	})

	// Set the state root on the BeaconBlock.
	g.Go(func() error {
		// Compute the state root for the block.
		var stateRoot primitives.Root

		s.logger.Info(
			"computing state root for block 🌲",
			"slot", blk.GetSlot(),
		)

		stateRoot, err = s.computeStateRoot(ctx, st, blk)
		if err != nil {
			s.logger.Error(
				"failed to compute state root for block ❌ ",
				"slot", requestedSlot,
				"error", err,
			)
			return err
		}

		s.logger.Info("state root computed for block 💻 ",
			"slot", requestedSlot,
			"state_root", stateRoot,
		)
		blk.SetStateRoot(stateRoot)
		return nil
	})

	if err = g.Wait(); err != nil {
		return blk, sidecars, err
	}

	s.logger.Info(
		"beacon block successfully built 🛠️ ",
		"slot", requestedSlot,
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	return blk, sidecars, nil
}

// verifyIncomingBlockStateRoot verifies the state root of an incoming block and
// logs the process.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) VerifyIncomingBlock(
	ctx context.Context,
	blk BeaconBlockT,
) error {
	s.logger.Info(
		"received incoming beacon block 📫 ",
		"state_root", blk.GetStateRoot(),
	)

	// Grab a copy of the state to verify the incoming block.
	st := s.bsb.StateFromContext(ctx)

	// Verify the state root of the incoming block.
	if err := s.verifyStateRoot(
		ctx, st, blk,
	); err != nil {
		// TODO: this is expensive because we are not caching the
		// previous result of HashTreeRoot().
		localStateRoot, htrErr := st.HashTreeRoot()
		if htrErr != nil {
			return htrErr
		}

		s.logger.Error(
			"rejecting incoming block ❌ ",
			"block_state_root",
			blk.GetStateRoot(),
			"local_state_root",
			localStateRoot,
			"error",
			err,
		)
		return err
	}

	s.logger.Info(
		"state root verification succeeded - accepting incoming block ✅ ",
		"state_root", blk.GetStateRoot(),
	)
	return nil
}
