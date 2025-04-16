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

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
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

//nolint:funlen // not an issue
func (s *Service) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) error {
	if countTx := len(req.Txs); countTx > MaxConsensusTxsCount {
		return fmt.Errorf("max expected %d, got %d: %w",
			MaxConsensusTxsCount, countTx,
			ErrTooManyConsensusTxs,
		)
	}

	forkVersion := s.chainSpec.ActiveForkVersionForTimestamp(math.U64(req.GetTime().Unix())) //#nosec: G115
	// Decode signed block and sidecars.
	signedBlk, sidecars, err := encoding.ExtractBlobsAndBlockFromRequest(
		req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		forkVersion,
	)
	if err != nil {
		return err
	}
	if signedBlk == nil {
		s.logger.Warn(
			"Aborting block verification - beacon block not found in proposal",
		)
		return ErrNilBlk
	}
	if sidecars == nil {
		s.logger.Warn(
			"Aborting block verification - blob sidecars not found in proposal",
		)
		return ErrNilBlob
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
	blkVersion := s.chainSpec.ActiveForkVersionForTimestamp(blk.GetTimestamp())
	if !version.Equals(blkVersion, forkVersion) {
		return fmt.Errorf("CometBFT version %v, BeaconBlock version %v: %w",
			forkVersion, blkVersion,
			ErrVersionMismatch,
		)
	}

	// Make sure we have the right number of BlobSidecars
	blobKzgCommitments := blk.GetBody().GetBlobKzgCommitments()
	numCommitments := len(blobKzgCommitments)
	if numCommitments != len(sidecars) {
		return fmt.Errorf("expected %d sidecars, got %d: %w",
			numCommitments, len(sidecars),
			ErrSidecarCommitmentMismatch,
		)
	}
	if uint64(numCommitments) > s.chainSpec.MaxBlobsPerBlock() {
		return fmt.Errorf("expected less than %d sidecars, got %d: %w",
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
			return fmt.Errorf("%w, idx: %d", ErrSidecarSignatureMismatch, i)
		}
	}
	err = s.VerifyIncomingBlockSignature(ctx, blk, signedBlk.GetSignature())
	if err != nil {
		return err
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
			return err
		}
	}

	// Process the block.
	consensusBlk := types.NewConsensusBlock(
		blk,
		req.GetProposerAddress(),
		req.GetTime(),
	)
	err = s.VerifyIncomingBlock(
		ctx,
		consensusBlk.GetBeaconBlock(),
		consensusBlk.GetConsensusTime(),
		consensusBlk.GetProposerAddress(),
	)
	if err != nil {
		s.logger.Error("failed to verify incoming block", "error", err)
		return err
	}

	return nil
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
//nolint:funlen // not an issue
func (s *Service) VerifyIncomingBlock(
	ctx context.Context,
	beaconBlk *ctypes.BeaconBlock,
	consensusTime math.U64,
	proposerAddress []byte,
) error {
	// Grab a copy of the state to verify the incoming block.
	preState := s.storageBackend.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	// TODO: Address the need for calling forceStartupSyncOnce in ProcessProposal. On a running
	// network (such as mainnet), it should be theoretically impossible to hit the case where
	// ProcessProposal is called before FinalizeBlock. It may be the case that new networks run
	// into this case during the first block after genesis.
	// TODO: Consider panicing here if this fails. If our node cannot successfully run
	// forceStartupSync, then we should shut down the node and fix the problem.
	s.forceStartupSyncOnce.Do(func() { s.forceSyncUponProcess(ctx, preState) })

	s.logger.Info(
		"Received incoming beacon block",
		"state_root", beaconBlk.GetStateRoot(),
		"slot", beaconBlk.GetSlot(),
	)

	// We purposefully make a copy of the BeaconState in order
	// to avoid modifying the underlying state, for the event in which
	// we have to rebuild a payload for this slot again, if we do not agree
	// with the incoming block.
	postState := preState.Copy(ctx)

	// verify block slot
	stateSlot, err := postState.GetSlot()
	if err != nil {
		s.logger.Error(
			"failed loading state slot to verify block slot",
			"reason", err,
		)
		return err
	}

	blkSlot := beaconBlk.GetSlot()
	if blkSlot != stateSlot+1 {
		s.logger.Error(
			"Rejecting incoming beacon block ❌ ",
			"state slot", stateSlot.Base10(),
			"block slot", blkSlot.Base10(),
			"reason", ErrUnexpectedBlockSlot.Error(),
		)
		return ErrUnexpectedBlockSlot
	}

	// Verify the state root of the incoming block.
	err = s.verifyStateRoot(
		ctx,
		postState,
		beaconBlk,
		consensusTime,
		proposerAddress)
	if err != nil {
		s.logger.Error(
			"Rejecting incoming beacon block ❌ ",
			"state_root", beaconBlk.GetStateRoot(),
			"reason", err,
		)

		if s.shouldBuildOptimisticPayloads() {
			lph, lphErr := preState.GetLatestExecutionPayloadHeader()
			if lphErr != nil {
				return errors.Join(
					err,
					fmt.Errorf("failed getting LatestExecutionPayloadHeader: %w", lphErr),
				)
			}

			// If we are rejecting the incoming block, let's optimistically build a candidate
			// payload for this slot (in case we are selected as the next proposer for this slot).
			go s.handleRebuildPayloadForRejectedBlock(
				ctx,
				preState,
				payloadtime.Next(
					consensusTime,
					lph.GetTimestamp(),
					true, // buildOptimistically
				),
			)
		}

		return err
	}

	s.logger.Info(
		"State root verification succeeded - accepting incoming beacon block",
		"state_root",
		beaconBlk.GetStateRoot(),
	)

	if s.shouldBuildOptimisticPayloads() {
		lph, lphErr := postState.GetLatestExecutionPayloadHeader()
		if lphErr != nil {
			return fmt.Errorf("failed loading LatestExecutionPayloadHeader: %w", lphErr)
		}

		go s.handleOptimisticPayloadBuild(
			ctx,
			postState,
			beaconBlk,
			payloadtime.Next(
				consensusTime,
				lph.GetTimestamp(),
				true, // buildOptimistically
			),
		)
	}

	return nil
}

// verifyStateRoot verifies the state root of an incoming block.
func (s *Service) verifyStateRoot(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	consensusTime math.U64,
	proposerAddress []byte,
) error {
	startTime := time.Now()
	defer s.metrics.measureStateRootVerificationTime(startTime)

	txCtx := transition.NewTransitionCtx(
		ctx,
		consensusTime,
		proposerAddress,
	).
		WithVerifyPayload(true).
		WithVerifyRandao(true).
		WithVerifyResult(true).
		WithMeterGas(false)

	_, err := s.stateProcessor.Transition(txCtx, st, blk)
	return err
}

// shouldBuildOptimisticPayloads returns true if optimistic
// payload builds are enabled.
func (s *Service) shouldBuildOptimisticPayloads() bool {
	return s.optimisticPayloadBuilds && s.localBuilder.Enabled()
}
