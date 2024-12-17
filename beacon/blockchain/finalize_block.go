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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blockchain

import (
    "context"
    "strings"
    "time"

    "github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
    "github.com/berachain/beacon-kit/consensus/types"
    "github.com/berachain/beacon-kit/primitives/math"
    "github.com/berachain/beacon-kit/primitives/transition"
    cmtabci "github.com/cometbft/cometbft/abci/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *Service[
    _, _, ConsensusBlockT, BeaconBlockT, _, _,
    _, _, _, GenesisT, ConsensusSidecarsT, BlobSidecarsT, _,
]) FinalizeBlock(
    ctx sdk.Context,
    req *cmtabci.FinalizeBlockRequest,
) (transition.ValidatorUpdates, error) {
    var (
        valUpdates  transition.ValidatorUpdates
        finalizeErr error
    )

    // STEP 1: Decode blok and blobs
    blk, blobs, err := encoding.
        ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](
        req,
        BeaconBlockTxIndex,
        BlobSidecarsTxIndex,
        s.chainSpec.ActiveForkVersionForSlot(
            math.Slot(req.Height),
        ))
    if err != nil {
        //nolint:nilerr // If we don't have a block, we can't do anything.
        return nil, nil
    }

    // STEP 2: Finalize sidecars first (block will check for
    // sidecar availability)
    err = s.blobProcessor.ProcessSidecars(
        s.storageBackend.AvailabilityStore(),
        blobs,
    )
    if err != nil {
        s.logger.Error("Failed to process blob sidecars", "error", err)
    }

    // STEP 3: finalize the block
    var consensusBlk *types.ConsensusBlock[BeaconBlockT]
    consensusBlk = consensusBlk.New(
        blk,
        req.GetProposerAddress(),
        req.GetTime(),
    )

    cBlk, ok := any(consensusBlk).(ConsensusBlockT)
    if !ok {
        panic("failed to convert consensusBlk to ConsensusBlockT")
    }

    st := s.storageBackend.StateFromContext(ctx)
    valUpdates, finalizeErr = s.finalizeBeaconBlock(ctx, st, cBlk)
    if finalizeErr != nil {
        s.logger.Error("Failed to process verified beacon block",
            "error", finalizeErr,
        )
    }

    // STEP 4: Post Finalizations cleanups

    // fetch and store the deposit for the block
    blockNum := blk.GetBody().GetExecutionPayload().GetNumber()
    s.depositFetcher(ctx, blockNum)

    // store the finalized block in the KVStore.
    slot := blk.GetSlot()
    if err = s.blockStore.Set(blk); err != nil {
        s.logger.Error(
            "failed to store block", "slot", slot, "error", err,
        )
    }

    // prune the availability and deposit store
    err = s.processPruning(blk)
    if err != nil {
        s.logger.Error("failed to processPruning", "error", err)
    }

    // New implementation with retries and backoff
    go func() {
        backoff := time.Second
        maxBackoff := time.Minute * 5
        maxAttempts := 5

        for attempt := 0; attempt < maxAttempts; attempt++ {
            err := s.sendPostBlockFCU(ctx, st, cBlk)
            if err == nil {
                return
            }

            // Log the error
            s.logger.Error("FCU failed", "attempt", attempt+1, "error", err)

            if !isTemporaryError(err) {
                s.pauseNodeParticipation()
                return
            }

            // Wait before next attempt
            time.Sleep(backoff)
            backoff = min(backoff*2, maxBackoff)
        }

        // If all attempts failed
        s.pauseNodeParticipation()
    }()

    return valUpdates, nil
}

// finalizeBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service[
    _, _, ConsensusBlockT, _, _, BeaconStateT, _, _, _, _, _, _, _,
]) finalizeBeaconBlock(
    ctx context.Context,
    st BeaconStateT,
    blk ConsensusBlockT,
) (transition.ValidatorUpdates, error) {
    beaconBlk := blk.GetBeaconBlock()

    // If the block is nil, exit early.
    if beaconBlk.IsNil() {
        return nil, ErrNilBlk
    }

    valUpdates, err := s.executeStateTransition(ctx, st, blk)
    if err != nil {
        return nil, err
    }

    // If the blobs needed to process the block are not available, we
    // return an error. It is safe to use the slot off of the beacon block
    // since it has been verified as correct already.
    if !s.storageBackend.AvailabilityStore().IsDataAvailable(
        ctx, beaconBlk.GetSlot(), beaconBlk.GetBody(),
    ) {
        return nil, ErrDataNotAvailable
    }
    return valUpdates.CanonicalSort(), nil
}

// executeStateTransition runs the stf.
func (s *Service[
    _, _, ConsensusBlockT, _, _, BeaconStateT, _, _, _, _, _, _, _,
]) executeStateTransition(
    ctx context.Context,
    st BeaconStateT,
    blk ConsensusBlockT,
) (transition.ValidatorUpdates, error) {
    startTime := time.Now()
    defer s.metrics.measureStateTransitionDuration(startTime)
    valUpdates, err := s.stateProcessor.Transition(
        &transition.Context{
            Context: ctx,

            // We set `OptimisticEngine` to true since this is called during
            // FinalizeBlock. We want to assume the payload is valid. If it
            // ends up not being valid later, the node will simply AppHash,
            // which is completely fine. This means we were syncing from a
            // bad peer, and we would likely AppHash anyways.
            OptimisticEngine: true,

            // When we are NOT synced to the tip, process proposal
            // does NOT get called and thus we must ensure that
            // NewPayload is called to get the execution
            // client the payload.
            //
            // When we are synced to the tip, we can skip the
            // NewPayload call since we already gave our execution client
            // the payload in process proposal.
            //
            // In both cases the payload was already accepted by a majority
            // of validators in their process proposal call and thus
            // the "verification aspect" of this NewPayload call is
            // actually irrelevant at this point.
            SkipPayloadVerification: false,

            ProposerAddress: blk.GetProposerAddress(),
            ConsensusTime:   blk.GetConsensusTime(),
        },
        st,
        blk.GetBeaconBlock(),
    )
    return valUpdates, err
}

// Helper functions
func isTemporaryError(err error) bool {
    return strings.Contains(err.Error(), "timeout") || 
           strings.Contains(err.Error(), "connection refused")
}

func (s *Service) pauseNodeParticipation() {
    s.logger.Error("Node participation paused due to FCU failures")
    // Implementation of node participation pause
    // Add actual implementation based on your node's architecture
}

func min(a, b time.Duration) time.Duration {
    if a < b {
        return a
    }
    return b
}
