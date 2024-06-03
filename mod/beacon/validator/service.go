// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload,
	],
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
	// metrics is a metrics collector.
	metrics *validatorMetrics
}

// NewService creates a new validator service.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload],
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
	ds DepositStore[*types.Deposit],
	localPayloadBuilder PayloadBuilder[BeaconStateT],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT],
	ts TelemetrySink,
) *Service[BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT] {
	return &Service[BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT]{
		cfg:                   cfg,
		logger:                logger,
		bsb:                   bsb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		ds:                    ds,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
		metrics:               newValidatorMetrics(ts),
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
	defer s.metrics.measureRequestBestBlockTime(startTime)
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
	reveal, err := s.buildRandaoReveal(st, requestedSlot)
	if err != nil {
		return blk, sidecars, err
	}

	// Create a new empty block from the current state.
	blk, err = s.getEmptyBeaconBlock(
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
	envelope, err := s.retrieveExecutionPayload(ctx, st, blk)
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

	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return blk, sidecars, ErrNilDepositIndexStart
	}

	// Dequeue deposits from the state.
	deposits, err := s.ds.GetDepositsByIndex(
		depositIndex,
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

	// Compute and set the state root for the block.
	g.Go(func() error {
		s.logger.Info(
			"computing state root for block 🌲",
			"slot", blk.GetSlot(),
		)

		var stateRoot primitives.Root
		stateRoot, err = s.computeStateRoot(ctx, st, blk)
		if err != nil {
			s.logger.Error(
				"failed to compute state root while building block ❗️ ",
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

// verifyIncomingBlockStateRoot verifies the state root of an incoming block
// and logs the process.
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
		"state root verification succeeded - accepting incoming block 🏎️ ",
		"state_root", blk.GetStateRoot(),
	)
	return nil
}
