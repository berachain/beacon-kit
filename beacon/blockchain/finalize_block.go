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

	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/consensus/types"
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
	var (
		valUpdates  transition.ValidatorUpdates
		finalizeErr error
	)

	// STEP 1: Decode block and blobs
	signedBlk, blobs, err := encoding.ExtractBlobsAndBlockFromRequest(
		req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		s.chainSpec.ActiveForkVersionForSlot(math.Slot(req.Height))) // #nosec G115
	if err != nil {
		s.logger.Error("Failed to decode block and blobs", "error", err)
		return nil, fmt.Errorf("failed to decode block and blobs: %w", err)
	}

	// STEP 2: Finalize sidecars first (block will check for
	// sidecar availability)
	err = s.blobProcessor.ProcessSidecars(
		s.storageBackend.AvailabilityStore(),
		blobs,
	)
	if err != nil {
		s.logger.Error("Failed to process blob sidecars", "error", err)
		return nil, fmt.Errorf("failed to process blob sidecars: %w", err)
	}

	// STEP 3: finalize the block
	blk := signedBlk.GetMessage()
	if blk == nil {
		s.logger.Error("SignedBeaconBlock contains nil BeaconBlock during FinalizeBlock")
		return nil, ErrNilBlk
	}
	consensusBlk := types.NewConsensusBlock(
		blk,
		req.GetProposerAddress(),
		req.GetTime(),
	)

	st := s.storageBackend.StateFromContext(ctx)
	valUpdates, finalizeErr = s.finalizeBeaconBlock(ctx, st, consensusBlk)
	if finalizeErr != nil {
		s.logger.Error("Failed to process verified beacon block",
			"error", finalizeErr,
		)
		return nil, finalizeErr
	}

	// STEP 4: Post Finalizations cleanups

	// fetch and store the deposit for the block
	blockNum := blk.GetBody().GetExecutionPayload().GetNumber()
	s.depositFetcher(ctx, blockNum)

	// store the finalized block in the KVStore.
	// TODO: Store full SignedBeaconBlock with all data in storage
	slot := blk.GetSlot()
	if err = s.storageBackend.BlockStore().Set(blk); err != nil {
		s.logger.Error(
			"failed to store block", "slot", slot, "error", err,
		)
		return nil, err
	}

	// prune the availability and deposit store
	err = s.processPruning(ctx, blk)
	if err != nil {
		s.logger.Error("failed to processPruning", "error", err)
	}

	go s.sendPostBlockFCU(ctx, st, consensusBlk)

	return valUpdates, nil
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
func (s *Service) executeStateTransition(
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
		WithVerifyResult(false).
		WithMeterGas(true).
		WithOptimisticEngine(true)

	return s.stateProcessor.Transition(
		txCtx,
		st,
		blk.GetBeaconBlock(),
	)
}
