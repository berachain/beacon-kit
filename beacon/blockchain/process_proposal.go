// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
)

func (s *Service) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) error {
	// Decode signed block and sidecars.
	signedBlk, sidecars, err := encoding.
		ExtractBlobsAndBlockFromRequest(
			req,
			BeaconBlockTxIndex,
			BlobSidecarsTxIndex,
			s.chainSpec.ActiveForkVersionForSlot(math.Slot(req.Height))) // #nosec G115
	if err != nil {
		return err
	}
	if signedBlk.IsNil() {
		s.logger.Warn(
			"Aborting block verification - beacon block not found in proposal",
		)
		return ErrNilBlk
	}
	if sidecars.IsNil() {
		s.logger.Warn(
			"Aborting block verification - blob sidecars not found in proposal",
		)
		return ErrNilBlob
	}

	blk := signedBlk.GetMessage()
	// Make sure we have the right number of BlobSidecars
	numCommitments := len(blk.GetBody().GetBlobKzgCommitments())
	if numCommitments != len(sidecars) {
		err = fmt.Errorf("expected %d sidecars, got %d: %w",
			numCommitments, len(sidecars),
			ErrSidecarCommittmentMismatch,
		)
		s.logger.Warn(err.Error())
		return err
	}

	// Verify the block and sidecar signatures. We can simply verify the block
	// signature and then make sure the sidecar signatures match the block.
	blkSignature := signedBlk.GetSignature()
	for i, sidecar := range sidecars {
		sidecarSignature := sidecar.GetSignedBeaconBlockHeader().GetSignature()
		if !bytes.Equal(blkSignature[:], sidecarSignature[:]) {
			return fmt.Errorf("%w, idx: %d", ErrSidecarSignatureMismatch, i)
		}
	}
	err = s.VerifyIncomingBlockSignature(ctx, signedBlk.GetMessage(), signedBlk.GetSignature())
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
		err = s.VerifyIncomingBlobSidecars(sidecars, blk.GetHeader(), blk.GetBody().GetBlobKzgCommitments())
		if err != nil {
			s.logger.Error("failed to verify incoming blob sidecars", "error", err)
			return err
		}
	}

	// Process the block
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
	return signatureVerifierFn(beaconBlk, signature)
}

// VerifyIncomingBlobSidecars verifies the BlobSidecars of an incoming
// proposal and logs the process.
func (s *Service) VerifyIncomingBlobSidecars(
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
) error {
	s.logger.Info("Received incoming blob sidecars")

	// Verify the blobs and ensure they match the local state.
	err := s.blobProcessor.VerifySidecars(sidecars, blkHeader, kzgCommitments)
	if err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars",
			"reason", err,
		)
		return err
	}

	s.logger.Info(
		"Blob sidecars verification succeeded - accepting incoming blob sidecars",
		"num_blobs",
		len(sidecars),
	)
	return nil
}

// VerifyIncomingBlock verifies the state root of an incoming block
// and logs the process.
func (s *Service) VerifyIncomingBlock(
	ctx context.Context,
	beaconBlk *ctypes.BeaconBlock,
	consensusTime math.U64,
	proposerAddress []byte,
) error {
	// Grab a copy of the state to verify the incoming block.
	preState := s.storageBackend.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	//
	// TODO: This is a super hacky. It should be handled better elsewhere,
	// ideally via some broader sync service.
	s.forceStartupSyncOnce.Do(func() { s.forceStartupHead(ctx, preState) })

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

	// Verify the state root of the incoming block.
	err := s.verifyStateRoot(
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
	_, err := s.stateProcessor.Transition(
		// We run with a non-optimistic engine here to ensure
		// that the proposer does not try to push through a bad block.
		&transition.Context{
			Context:                 ctx,
			MeterGas:                false,
			OptimisticEngine:        false,
			SkipPayloadVerification: false,
			SkipValidateResult:      false,
			SkipValidateRandao:      false,
			ProposerAddress:         proposerAddress,
			ConsensusTime:           consensusTime,
		},
		st, blk,
	)

	return err
}

// shouldBuildOptimisticPayloads returns true if optimistic
// payload builds are enabled.
func (s *Service) shouldBuildOptimisticPayloads() bool {
	return s.optimisticPayloadBuilds && s.localBuilder.Enabled()
}
