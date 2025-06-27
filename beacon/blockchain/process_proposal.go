// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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

package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// BeaconBlockTxIndex represents the index of the beacon block transaction.
	// It is the first transaction in the tx list.
	BeaconBlockTxIndex uint = iota
	// BlobSidecarsTxIndex represents the index of the blob sidecar transaction.
	// It follows the beacon block transaction in the tx list.
	BlobSidecarsTxIndex

	// A Consensus block has at most two transactions (block and blob).
	MaxConsensusTxsCount = 2
)

func (s *Service) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (transition.ValidatorUpdates, error) {
	signedBlk, sidecars, err := s.ParseBeaconBlock(req)
	if err != nil {
		s.logger.Error("Failed to decode block and blobs", "error", err)
		return nil, fmt.Errorf("failed to decode block and blobs: %w", err)
	}
	blk := signedBlk.GetBeaconBlock()

	// There are two different timestamps:
	//     - The "consensus time" is determined by CometBFT consensus and can be retrieved with `req.GetTime()`
	//     - The "block time" is determined by beacon-kit consensus and can be retrieved with `blk.GetTimestamp()`
	// The "consensus time" is what the network agrees the current time is based on CometBFT PBTS.
	// This "consensus time" is used to constrain the timestamp set as the "block time" by the
	// beacon-kit app, but they are not always equal in value. The "block time" is used by the
	// beacon-kit consensus and execution layers to determine the active fork version.
	//
	// When unmarshaling the BeaconBlock, we do not yet have access to the "block time", so we
	// must rely on the "consensus time" as our best estimation of the "block time" needed to
	// determine the current fork version. Since the two timestamps could be different, we need to
	// ensure that the fork version for these timestamps are the same. This may result in a failed
	// proposal or two at the start of the fork.
	forkVersion := s.chainSpec.ActiveForkVersionForTimestamp(math.U64(req.GetTime().Unix())) //#nosec: G115
	blkVersion := s.chainSpec.ActiveForkVersionForTimestamp(blk.GetTimestamp())
	if !version.Equals(blkVersion, forkVersion) {
		return nil, fmt.Errorf("CometBFT version %v, BeaconBlock version %v: %w",
			forkVersion, blkVersion,
			ErrVersionMismatch,
		)
	}

	// Make sure we have the right number of BlobSidecars
	blobKzgCommitments := blk.GetBody().GetBlobKzgCommitments()
	numCommitments := len(blobKzgCommitments)
	if numCommitments != len(sidecars) {
		return nil, fmt.Errorf("expected %d sidecars, got %d: %w",
			numCommitments, len(sidecars),
			ErrSidecarCommitmentMismatch,
		)
	}
	if uint64(numCommitments) > s.chainSpec.MaxBlobsPerBlock() {
		return nil, fmt.Errorf("expected less than %d sidecars, got %d: %w",
			s.chainSpec.MaxBlobsPerBlock(), numCommitments,
			core.ErrExceedsBlockBlobLimit,
		)
	}

	// Verify the block and sidecar signatures. We can simply verify the block
	// signature and then make sure the sidecar signatures match the block.
	blkSignature := signedBlk.GetSignature()
	for i, sidecar := range sidecars {
		sidecarSignature := sidecar.GetSignature()
		if !bytes.Equal(blkSignature[:], sidecarSignature[:]) {
			return nil, fmt.Errorf("%w, idx: %d", ErrSidecarSignatureMismatch, i)
		}
	}
	err = s.VerifyIncomingBlockSignature(ctx, blk, signedBlk.GetSignature())
	if err != nil {
		return nil, err
	}

	if numCommitments > 0 {
		// Process the blob sidecars
		//
		// In theory, swapping the order of verification between the sidecars
		// and the incoming block should not introduce any inconsistencies
		// in the state on which the sidecar verification depends on (notably
		// the currently active fork). ProcessProposal should only
		// keep the state changes as candidates (which is what we do in
		// VerifyIncomingBlock).
		err = s.VerifyIncomingBlobSidecars(ctx, sidecars, blk.GetHeader(), blobKzgCommitments)
		if err != nil {
			s.logger.Error("failed to verify incoming blob sidecars", "error", err)
			return nil, err
		}
	}

	// Process the block.
	s.logger.Debug( // needed for some checks in sim tests
		"Processing block with fork version",
		"block", req.Height,
		"fork", blkVersion.String(),
	)
	consensusBlk := types.NewConsensusBlock(
		blk,
		req.GetProposerAddress(),
		req.GetTime(),
	)

	var valUpdates transition.ValidatorUpdates
	valUpdates, err = s.VerifyIncomingBlock(ctx, consensusBlk)
	if err != nil {
		s.logger.Error("failed to verify incoming block", "error", err)
		return nil, err
	}

	return valUpdates.CanonicalSort(), nil
}

func (s *Service) VerifyIncomingBlockSignature(
	ctx context.Context,
	beaconBlk *ctypes.BeaconBlock,
	signature crypto.BLSSignature,
) error {
	// Get the sidecar verification function from the state processor
	signatureVerifierFn, err := s.stateProcessor.GetSignatureVerifierFn(
		s.storageBackend.StateFromContext(ctx),
	)
	if err != nil {
		return errors.New("failed to create block signature verifier")
	}
	err = signatureVerifierFn(beaconBlk, signature)
	if err != nil {
		return fmt.Errorf("failed verifying incoming block signature: %w", err)
	}
	return err
}

// VerifyIncomingBlobSidecars verifies the BlobSidecars of an incoming
// proposal and logs the process.
func (s *Service) VerifyIncomingBlobSidecars(
	ctx context.Context,
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
) error {
	// Verify the blobs and ensure they match the local state.
	err := s.blobProcessor.VerifySidecars(ctx, sidecars, blkHeader, kzgCommitments)
	if err != nil {
		s.logger.Error(
			"Blob sidecars verification failed - rejecting incoming blob sidecars",
			"reason", err, "slot", blkHeader.GetSlot(),
		)
		return err
	}

	s.logger.Info(
		"Blob sidecars verification succeeded - accepting incoming blob sidecars",
		"num_blobs", len(sidecars), "slot", blkHeader.GetSlot(),
	)
	return nil
}

// VerifyIncomingBlock verifies the state root of an incoming block
// and logs the process.
//
//nolint:funlen // abundantly commented
func (s *Service) VerifyIncomingBlock(
	ctx context.Context,
	blk *types.ConsensusBlock,
) (transition.ValidatorUpdates, error) {
	beaconBlk := blk.GetBeaconBlock()
	state := s.storageBackend.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	// TODO: Address the need for calling forceStartupSyncOnce in ProcessProposal. On a running
	// network (such as mainnet), it should be theoretically impossible to hit the case where
	// ProcessProposal is called before FinalizeBlock. It may be the case that new networks run
	// into this case during the first block after genesis.
	// TODO: Consider panicing here if this fails. If our node cannot successfully run
	// forceStartupSync, then we should shut down the node and fix the problem.
	s.forceStartupSyncOnce.Do(func() { s.forceSyncUponProcess(ctx, state) })

	s.logger.Info(
		"Received incoming beacon block",
		"state_root", beaconBlk.GetStateRoot(),
		"slot", beaconBlk.GetSlot(),
	)

	// verify block slot
	stateSlot, err := state.GetSlot()
	if err != nil {
		s.logger.Error(
			"failed loading state slot to verify block slot",
			"reason", err,
		)
		return nil, err
	}

	blkSlot := beaconBlk.GetSlot()
	if blkSlot != stateSlot+1 {
		s.logger.Error(
			"Rejecting incoming beacon block ❌ ",
			"state slot", stateSlot.Base10(),
			"block slot", blkSlot.Base10(),
			"reason", ErrUnexpectedBlockSlot.Error(),
		)
		return nil, ErrUnexpectedBlockSlot
	}

	var (
		nextBlockData *builder.RequestPayloadData
		errFetch      error
	)

	if s.shouldBuildOptimisticPayloads() {
		// state copy makes sure that preFetchBuildData does not affect state
		copiedState := state.Copy(ctx)
		nextBlockData, errFetch = s.preFetchBuildData(copiedState, blk.GetConsensusTime())
		if errFetch != nil {
			// We don't return with err if pre-fetch fails. Instead we log the issue
			// and still move to process the current block. Next block can always be
			// built right after current height is finalized.
			s.logger.Warn(
				"Failed pre fetching data for optimistic block building",
				"case", "block rejectiong",
				"err", errFetch,
			)
		}
	}

	// Verify the state root of the incoming block.
	var valUpdates transition.ValidatorUpdates
	valUpdates, err = s.verifyStateRoot(ctx, state, blk)
	if err != nil {
		s.logger.Error(
			"Rejecting incoming beacon block ❌ ",
			"state_root", beaconBlk.GetStateRoot(),
			"reason", err,
		)

		if s.shouldBuildOptimisticPayloads() {
			if nextBlockData == nil {
				// Failed fetching data to build next block. Just return block error
				return nil, err
			}
			go s.handleRebuildPayloadForRejectedBlock(ctx, nextBlockData)
		}

		return nil, err
	}

	s.logger.Info(
		"State root verification succeeded - accepting incoming beacon block",
		"state_root", beaconBlk.GetStateRoot(),
	)

	if s.shouldBuildOptimisticPayloads() {
		// state copy makes sure that preFetchBuildDataForSuccess does not affect state
		copiedState := state.Copy(ctx)
		nextBlockData, errFetch = s.preFetchBuildData(copiedState, blk.GetConsensusTime())
		if errFetch != nil {
			// We don't mark the block as rejected if it is valid but pre-fetch fails.
			// Instead we log the issue and move to process the current block.
			// Next block can always be built right after current height is finalized.
			s.logger.Warn(
				"Failed pre fetching data for optimistic block building",
				"case", "block success",
				"err", errFetch,
			)
			return valUpdates, nil
		}
		go s.handleOptimisticPayloadBuild(ctx, nextBlockData)
	}

	return valUpdates, nil
}

// verifyStateRoot verifies the state root of an incoming block.
func (s *Service) verifyStateRoot(
	ctx context.Context,
	st *statedb.StateDB,
	blk *types.ConsensusBlock,
) (transition.ValidatorUpdates, error) {
	startTime := time.Now()
	defer s.metrics.measureStateTransitionDuration(startTime)

	txCtx := transition.NewTransitionCtx(
		ctx,
		blk.GetConsensusTime(),
		blk.GetProposerAddress(),
	).
		WithVerifyPayload(true).
		WithVerifyRandao(true).
		WithVerifyResult(true).
		WithMeterGas(true)

	valUpdates, err := s.stateProcessor.Transition(txCtx, st, blk.GetBeaconBlock())
	return valUpdates, err
}

// shouldBuildOptimisticPayloads returns true if optimistic
// payload builds are enabled.
func (s *Service) shouldBuildOptimisticPayloads() bool {
	return s.optimisticPayloadBuilds && s.localBuilder.Enabled()
}
