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

	"github.com/berachain/beacon-kit/beacon/deposits"
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
		if finalizeErr != nil {
			finalizeErr = fmt.Errorf("force sync upon finalize: failed retrieving parent proposer pubkey: %w", finalizeErr)
		} else {
			finalizeErr = s.forceSyncUponFinalize(ctx, blk, parentProposerPubkey)
		}
	})
	if finalizeErr != nil {
		return nil, finalizeErr
	}

	// STEP 2: Finalize sidecars first (block will check for sidecar availability).
	if err = s.FinalizeSidecars(ctx, req.SyncingToHeight, signedBlk, blobs); err != nil {
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

// FinalizeSidecars makes sure the blobs a finalized block commits to are persisted in the availability store,
// keeping "finalized within the DA window" equivalent to "data available". Outside the DA window there is
// nothing to enforce and the call is a no-op.
//
// Where the sidecars come from depends on how the block reached us. Below the blob consensus enable height
// they ride in the request txs, and for a block verified at the tip they were already fetched during
// ProcessProposal. In both cases the caller passes them in, and they are persisted under a strict availability
// check. With blob consensus enabled and no sidecars in hand (block sync, or a tip block this node never voted
// on), the store is checked first, then the same tiered retrieval ProcessProposal uses if the block is at the
// tip, and failing that a background fetch is queued. Queueing never fails the block: it is already committed
// by 2/3+ of the voting power, which implies the data was available to the network, and the fetcher keeps
// retrying until the slot leaves the DA window while the node reports itself as still syncing.
func (s *Service) FinalizeSidecars(
	ctx sdk.Context,
	syncingToHeight int64,
	signedBlk *ctypes.SignedBeaconBlock,
	blobs datypes.BlobSidecars,
) error {
	blk := signedBlk.GetBeaconBlock()

	// SyncingToHeight is always the tip of the chain both during sync and when
	// caught up. We don't need to process sidecars unless they are within DA period.
	//
	//#nosec: G115 // SyncingToHeight will never be negative.
	if !s.chainSpec.WithinDAPeriod(blk.GetSlot(), math.Slot(syncingToHeight)) {
		// Here outside Data Availability window. Just log if needed
		if len(blobs) > 0 {
			s.logger.Info("Skipping blob processing outside of Data Availability Period", "slot", blk.GetSlot().Base10(), "head", syncingToHeight)
		}
		return nil
	}

	processAndCheck := func(sidecars datypes.BlobSidecars) error {
		err := s.blobProcessor.ProcessSidecars(
			s.storageBackend.AvailabilityStore(),
			sidecars,
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

	//#nosec: G115 // slots are never negative.
	blobConsensusEnabled := s.chainSpec.IsBlobConsensusEnabled(int64(blk.GetSlot().Unwrap()))

	// Legacy layout (sidecars rode as a consensus tx) or sidecars handed to us by the caller (verified during ProcessProposal):
	// persist and enforce availability strictly.
	if !blobConsensusEnabled || len(blobs) > 0 {
		return processAndCheck(blobs)
	}

	// Blob consensus is enabled and no sidecars were provided with the block.
	if len(blk.GetBody().GetBlobKzgCommitments()) == 0 {
		return nil
	}

	// Already stored and bound to this exact block (e.g. restart replay after a crash mid-finalization, or a previous fetch already
	// persisted them). Presence alone would not be enough: sidecars from a different proposal at this slot with identical commitments
	// would pass an availability check while carrying the wrong header and signature.
	if sidecarsAlreadyStored(s.storageBackend.AvailabilityStore(), blk.GetHeader(), blk.GetBody().GetBlobKzgCommitments()) {
		return nil
	}

	// At the tip of the chain a node can reach FinalizeBlock without a successful ProcessProposal (it voted nil while 2/3+
	// committed, or it just joined). Use the same tiered retrieval as ProcessProposal.
	//#nosec: G115 // slots are never negative.
	if int64(blk.GetSlot().Unwrap()) == syncingToHeight {
		sidecars, err := s.fetchProposalSidecars(ctx, signedBlk)
		if err == nil {
			return processAndCheck(sidecars)
		}
		s.logger.Warn("Failed to fetch blob sidecars at tip, queueing background fetch",
			"slot", blk.GetSlot().Unwrap(), "error", err)
	}

	// Syncing (or the tip fetch failed): queue an asynchronous fetch. The background fetcher retries until the slot leaves the DA
	// window, and the node does not report itself as synced while in-window fetches are pending. Failing the block here is not an
	// option: the block is already committed by 2/3+ of the voting power, which also implies they held and verified the data.
	if err := s.blobFetcher.QueueBlobRequest(signedBlk); err != nil {
		return fmt.Errorf("failed to queue blob fetch request for slot %d: %w",
			blk.GetSlot().Unwrap(), err)
	}
	return nil
}

func (s *Service) PostFinalizeBlockOps(ctx sdk.Context, blk *ctypes.BeaconBlock) error {
	// TODO: consider extracting LatestExecutionPayloadHeader instead of using state here
	st := s.storageBackend.StateFromContext(ctx)

	// Before Fulu, deposits must be fetched from the EL (at the eth1 follow distance).
	deposits.FetchPreviousDepositsPreFulu(
		ctx, s.depositContract, blk, s.eth1FollowDistance, s.storageBackend.DepositStore(), s.logger,
	)

	// Store the finalized block in the KVStore.
	slot := blk.GetSlot()
	if err := s.storageBackend.BlockStore().Set(blk); err != nil {
		s.logger.Error(
			"failed to store block", "slot", slot, "error", err,
		)
		return err
	}

	// Update the head slot for blob distribution (peers are served based on it, and pending fetches outside the DA window are
	// dropped against it).
	s.blobFetcher.SetHeadSlot(slot)

	// Prune the availability and deposit store.
	if err := s.processPruning(ctx, blk); err != nil {
		s.logger.Error("failed to processPruning", "error", err)
	}

	if err := s.sendPostBlockFCU(ctx, st); err != nil {
		return fmt.Errorf("sendPostBlockFCU failed: %w", err)
	}

	// reset latest verified payload in block builder to signal
	// that no payload is available to reuse for blk.Slot
	s.localBuilder.CacheLatestVerifiedPayload(blk.Slot, nil)

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

	// If on the first block of Fulu, catchup the previous block's deposits.
	if err := deposits.CatchupFuluDeposits(
		ctx, s.depositContract, st, beaconBlk, s.chainSpec, s.storageBackend.DepositStore(), s.logger,
	); err != nil {
		return nil, err
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
