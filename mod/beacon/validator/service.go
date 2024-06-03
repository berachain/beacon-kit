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
	"sync"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
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
	DepositStoreT DepositStore[*types.Deposit],
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
	bsb StorageBackend[BeaconStateT, *types.Deposit, DepositStoreT]
	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		*transition.Context,
	]
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
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once
}

// NewService creates a new validator service.
func NewService[
	BeaconBlockT BeaconBlock[BeaconBlockT, BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody[
		*types.Deposit, *types.Eth1Data, *types.ExecutionPayload],
	BeaconStateT BeaconState[BeaconStateT],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore[*types.Deposit],
](
	cfg *Config,
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	bsb StorageBackend[BeaconStateT, *types.Deposit, DepositStoreT],
	stateProcessor StateProcessor[BeaconBlockT, BeaconStateT, *transition.Context],
	signer crypto.BLSSigner,
	blobFactory BlobFactory[
		BeaconBlockT, BeaconBlockBodyT, BlobSidecarsT,
	],
	localPayloadBuilder PayloadBuilder[BeaconStateT],
	remotePayloadBuilders []PayloadBuilder[BeaconStateT],
	ts TelemetrySink,
) *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
] {
	return &Service[
		BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositStoreT,
	]{
		cfg:                   cfg,
		logger:                logger,
		bsb:                   bsb,
		chainSpec:             chainSpec,
		signer:                signer,
		stateProcessor:        stateProcessor,
		blobFactory:           blobFactory,
		localPayloadBuilder:   localPayloadBuilder,
		remotePayloadBuilders: remotePayloadBuilders,
		metrics:               newValidatorMetrics(ts),
		forceStartupSyncOnce:  new(sync.Once),
	}
}

// Name returns the name of the service.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) Name() string {
	return "validator"
}

// Start starts the service.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) Start(
	context.Context,
) error {
	s.logger.Info(
		"starting validator service 🛜 ",
		"optimistic_payload_builds", s.cfg.EnableOptimisticPayloadBuilds,
	)
	return nil
}

// Status returns the status of the service.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) Status() error {
	return nil
}

// WaitForHealthy waits for the service to become healthy.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) WaitForHealthy(
	context.Context,
) {
}

// RequestBestBlock builds a new beacon block.
//
//nolint:funlen // todo:fix.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
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

	// Force a sync of the startup head if we haven't done so already.
	//
	// TODO: This is a super hacky. It should be handled better elsewhere,
	// ideally via some broader sync service.
	s.forceStartupSyncOnce.Do(func() { s.forceStartupHead(ctx, st) })

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
	deposits, err := s.bsb.DepositStore(ctx).GetDepositsByIndex(
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
//
//nolint:gocognit // todo fix.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) VerifyIncomingBlock(
	ctx context.Context,
	blk BeaconBlockT,
) error {
	// Grab a copy of the state to verify the incoming block.
	st := s.bsb.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	//
	// TODO: This is a super hacky. It should be handled better elsewhere,
	// ideally via some broader sync service.
	s.forceStartupSyncOnce.Do(func() { s.forceStartupHead(ctx, st) })

	// If the block is nil or a nil pointer, exit early.
	if blk.IsNil() {
		s.logger.Error(
			"aborting block verification on nil block ⛔️ ",
		)

		if s.localPayloadBuilder.Enabled() &&
			s.cfg.EnableOptimisticPayloadBuilds {
			go func() {
				if pErr := s.rebuildPayloadForRejectedBlock(
					ctx, st,
				); pErr != nil {
					s.logger.Error(
						"failed to rebuild payload for nil block",
						"error", pErr,
					)
				}
			}()
		}

		return ErrNilBlk
	}

	s.logger.Info(
		"received incoming beacon block 📫 ",
		"state_root", blk.GetStateRoot(),
	)

	// We purposefully make a copy of the BeaconState in orer
	// to avoid modifying the underlying state, for the event in which
	// we have to rebuild a payload for this slot again, if we do not agree
	// with the incoming block.
	stCopy := st.Copy()

	// Verify the state root of the incoming block.
	if err := s.verifyStateRoot(
		ctx, stCopy, blk,
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
			primitives.Root(localStateRoot),
			"error",
			err,
		)

		if s.localPayloadBuilder.Enabled() &&
			s.cfg.EnableOptimisticPayloadBuilds {
			go func() {
				if pErr := s.rebuildPayloadForRejectedBlock(
					ctx, st,
				); pErr != nil {
					s.logger.Error(
						"failed to rebuild payload for rejected block",
						"for_slot", blk.GetSlot(),
						"error", pErr,
					)
				}
			}()
		}

		return err
	}

	s.logger.Info(
		"state root verification succeeded - accepting incoming block 🏎️ ",
		"state_root", blk.GetStateRoot(),
	)

	if s.localPayloadBuilder.Enabled() && s.cfg.EnableOptimisticPayloadBuilds {
		go func() {
			if err := s.optimisticPayloadBuild(ctx, stCopy, blk); err != nil {
				s.logger.Error(
					"failed to build optimistic payload",
					"for_slot", blk.GetSlot()+1,
					"error", err,
				)
			}
		}()
	}

	return nil
}

// rebuildPayloadForRejectedBlock rebuilds a payload for the current
// slot, if the incoming block was rejected.
//
// NOTE: We cannot use any data off the incoming block and must recompute
// any required information from our local state. We do this since we have
// rejected the incoming block and it would be unsafe to use any
// information from it.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositT,
]) rebuildPayloadForRejectedBlock(
	ctx context.Context,
	st BeaconStateT,
) error {
	var (
		previousBlockRoot primitives.Root
		lph               engineprimitives.ExecutionPayloadHeader
		slot              math.Slot
	)

	s.logger.Info("rebuilding payload for rejected block ⏳ ")

	// In order to rebuild a payload for the current slot, we need to know the
	// previous block root, since we know that this is unmodified state.
	// We can safely get the latest block header and then rebuild the
	// previous block and it's root.
	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}

	stateRoot, err := st.HashTreeRoot()
	if err != nil {
		return err
	}

	latestHeader.StateRoot = stateRoot
	previousBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	// We need to get the *last* finalized execution payload, thus
	// the BeaconState that was passed in must be `unmodified`.
	lph, err = st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	slot, err = st.GetSlot()
	if err != nil {
		return err
	}

	// Submit a request for a new payload.
	if _, err = s.localPayloadBuilder.RequestPayloadAsync(
		ctx,
		st,
		// We are rebuilding for the current slot.
		slot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+1),
			uint64((lph.GetTimestamp()+1)),
		),
		// We set the parent root to the previous block root.
		previousBlockRoot,
		// We set the head of our chain to previous finalized block.
		lph.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// TODO: This is making an assumption about the consensus rules
		// and possibly should be made more explicit later on.
		lph.GetBlockHash(),
	); err != nil {
		s.metrics.markRebuildPayloadForRejectedBlockFailure(slot, err)
		return err
	}
	s.metrics.markRebuildPayloadForRejectedBlockSuccess(slot)
	return nil
}

// optimisticPayloadBuild builds a payload for the next slot.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositT,
]) optimisticPayloadBuild(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	s.logger.Info("optimistically triggering payload build for next slot 🛩️ ")

	// We know that this block was properly formed so we can
	// calculate the block hash easily.
	blkRoot, err := blk.HashTreeRoot()
	if err != nil {
		return err
	}

	// We process the slot to update any RANDAO values.
	if _, err = s.stateProcessor.ProcessSlot(
		st,
	); err != nil {
		return err
	}

	// We are building for the next slot, so we increment the slot relative
	// to the block we just processed.
	slot := blk.GetSlot() + 1

	// We then trigger a request for the next payload.
	payload := blk.GetBody().GetExecutionPayload()
	if _, err = s.localPayloadBuilder.RequestPayloadAsync(
		ctx, st,
		slot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+int64(s.chainSpec.TargetSecondsPerEth1Block())),
			uint64((payload.GetTimestamp()+1)),
		),
		// The previous block root is simply the root of the block we just
		// processed.
		blkRoot,
		// We set the head of our chain to the block we just processed.
		payload.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// This is safe to do since this block was accepted and the thus the
		// parent hash was deemed valid by the state transition function we
		// just processed.
		payload.GetParentHash(),
	); err != nil {
		s.metrics.markOptimisticPayloadBuildFailure(slot, err)
		return err
	}
	s.metrics.markOptimisticPayloadBuildSuccess(slot)
	return nil
}

// verifyIncomingBlockStateRoot verifies the state root of an incoming block
// and logs the process.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) forceStartupHead(
	ctx context.Context,
	st BeaconStateT,
) {
	slot, err := st.GetSlot()
	if err != nil {
		s.logger.Error(
			"failed to get slot for force startup head",
			"error", err,
		)
		return
	}

	// TODO: Verify if the slot number is correct here, I believe in current
	// form
	// it should be +1'd. Not a big deal until hardforks are in play though.
	if err = s.localPayloadBuilder.SendForceHeadFCU(ctx, st, slot+1); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}
