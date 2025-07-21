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
	"context"
	"fmt"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *Service) FinalizeBlock(
	ctx sdk.Context,
	req *cmtabci.FinalizeBlockRequest,
) (transition.ValidatorUpdates, error) {
	// STEP 1: Decode block and blobs.
	signedBlk, blobs, err := s.ParseBeaconBlock(req)
	if err != nil {
		s.logger.Error("Failed to decode block and blobs", "error", err)
		return nil, fmt.Errorf("failed to decode block and blobs: %w", err)
	}
	blk := signedBlk.GetBeaconBlock()
	st := s.storageBackend.StateFromContext(ctx)

	// Send an FCU to force the HEAD of the chain on the EL on startup.
	var finalizeErr error
	s.forceStartupSyncOnce.Do(func() {
		var parentProposerPubkey *crypto.BLSPubkey
		parentProposerPubkey, finalizeErr = st.ParentProposerPubkey(blk.GetTimestamp())
		if err != nil {
			finalizeErr = fmt.Errorf("force sync upon finalize: failed retrieving parent proposer pubkey: %w", finalizeErr)
		} else {
			finalizeErr = s.forceSyncUponFinalize(ctx, blk, parentProposerPubkey)
		}
	})
	if finalizeErr != nil {
		return nil, finalizeErr
	}

	// STEP 2: Finalize sidecars first (block will check for sidecar availability).
	if err = s.FinalizeSidecars(ctx, req.SyncingToHeight, blk, blobs); err != nil {
		return nil, fmt.Errorf("failed finalizing sidecars: %w", err)
	}

	// STEP 3: Finalize the block.
	consensusBlk := types.NewConsensusBlock(blk, req.GetProposerAddress(), req.GetTime())
	valUpdates, err := s.finalizeBeaconBlock(ctx, st, consensusBlk)
	if err != nil {
		s.logger.Error("Failed to process verified beacon block",
			"error", err,
		)
		return nil, err
	}

	// STEP 4: Post Finalizations cleanups.
	return valUpdates, s.PostFinalizeBlockOps(ctx, blk)
}

func (s *Service) FinalizeSidecars(
	ctx sdk.Context,
	syncingToHeight int64,
	blk *ctypes.BeaconBlock,
	blobs datypes.BlobSidecars,
) error {
	// SyncingToHeight is always the tip of the chain both during sync and when
	// caught up. We don't need to process sidecars unless they are within DA period.
	//
	//#nosec: G115 // SyncingToHeight will never be negative.
	if s.chainSpec.WithinDAPeriod(blk.GetSlot(), math.Slot(syncingToHeight)) {
		err := s.blobProcessor.ProcessSidecars(
			s.storageBackend.AvailabilityStore(),
			blobs,
		)
		if err != nil {
			s.logger.Error("Failed to process blob sidecars", "error", err)
			return fmt.Errorf("failed to process blob sidecars: %w", err)
		}

		// Ensure we can access the data using the commitments from the block.
		if !s.storageBackend.AvailabilityStore().IsDataAvailable(
			ctx, blk.GetSlot(), blk.GetBody(),
		) {
			return ErrDataNotAvailable
		}
		return nil
	}

	// Here outside Data Availability window. Just log if needed
	if len(blobs) > 0 {
		s.logger.Info(
			"Skipping blob processing outside of Data Availability Period",
			"slot", blk.GetSlot().Base10(), "head", syncingToHeight,
		)
	}
	return nil
}

func (s *Service) PostFinalizeBlockOps(ctx sdk.Context, blk *ctypes.BeaconBlock) error {
	// TODO: consider extracting LatestExecutionPayloadHeader instead of using state here
	st := s.storageBackend.StateFromContext(ctx)

	// Fetch and store the deposit for the block.
	blockNum := blk.GetBody().GetExecutionPayload().GetNumber()
	s.depositFetcher(ctx, blockNum)

	// Store the finalized block in the KVStore.
	//
	// TODO: Store full SignedBeaconBlock with all data in storage
	slot := blk.GetSlot()
	if err := s.storageBackend.BlockStore().Set(blk); err != nil {
		s.logger.Error(
			"failed to store block", "slot", slot, "error", err,
		)
		return err
	}

	// Prune the availability and deposit store.
	if err := s.processPruning(ctx, blk); err != nil {
		s.logger.Error("failed to processPruning", "error", err)
	}

	if err := s.sendPostBlockFCU(ctx, st); err != nil {
		return fmt.Errorf("sendPostBlockFCU failed: %w", err)
	}

	return nil
}

// finalizeBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) finalizeBeaconBlock(
	ctx context.Context,
	st *statedb.StateDB,
	blk *types.ConsensusBlock,
) (transition.ValidatorUpdates, error) {
	beaconBlk := blk.GetBeaconBlock()

	// If the block is nil, exit early.
	if beaconBlk == nil {
		return nil, ErrNilBlk
	}

	valUpdates, err := s.executeStateTransition(ctx, st, blk)
	if err != nil {
		return nil, err
	}
	return valUpdates.CanonicalSort(), nil
}

// executeStateTransition runs the stf.
func (s *Service) executeStateTransition(
	ctx context.Context,
	st *statedb.StateDB,
	blk *types.ConsensusBlock,
) (transition.ValidatorUpdates, error) {
	startTime := time.Now()
	defer s.metrics.measureStateTransitionDuration(startTime)

	// Notes about context attributes:
	// - VerifyPayload: set to true. When we are NOT synced to the tip,
	// process proposal does NOT get called and thus we must ensure that
	// NewPayload is called to get the execution client the payload.
	// When we are synced to the tip, we can skip the
	// NewPayload call since we already gave our execution client
	// the payload in process proposal.
	// In both cases the payload was already accepted by a majority
	// of validators in their process proposal call and thus
	// the "verification aspect" of this NewPayload call is
	// actually irrelevant at this point.
	// - VerifyRandao: set to false. We skip randao validation in FinalizeBlock
	// since either
	//   1. we validated it during ProcessProposal at the head of the chain OR
	//   2. we are bootstrapping and implicitly trust that the randao was validated by
	//    the super majority during ProcessProposal of the given block height.
	txCtx := transition.NewTransitionCtx(
		ctx,
		blk.GetConsensusTime(),
		blk.GetProposerAddress(),
	).
		WithVerifyPayload(true).
		WithVerifyRandao(false).
		WithVerifyResult(false).
		WithMeterGas(true)

	return s.stateProcessor.Transition(
		txCtx,
		st,
		blk.GetBeaconBlock(),
	)
}
