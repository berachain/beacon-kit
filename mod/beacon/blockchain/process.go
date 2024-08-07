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
	"context"
	"time"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// ProcessGenesisData processes the genesis state and initializes the beacon
// state.
func (s *Service[
	_, _, _, _, _, _, _, _, GenesisT, _, _,
]) ProcessGenesisData(
	ctx context.Context,
	genesisData GenesisT,
) (transition.ValidatorUpdates, error) {
	return s.sp.InitializePreminedBeaconStateFromEth1(
		s.sb.StateFromContext(ctx),
		genesisData.GetDeposits(),
		genesisData.GetExecutionPayloadHeader(),
		genesisData.GetForkVersion(),
	)
}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _, _,
]) ProcessBeaconBlock(
	ctx context.Context,
	blk BeaconBlockT,
) (transition.ValidatorUpdates, error) {
	// If the block is nil, exit early.
	if blk.IsNil() {
		return nil, ErrNilBlk
	}

	// We set `OptimisticEngine` to true since this is called during
	// FinalizeBlock. We want to assume the payload is valid. If it
	// ends up not being valid later, the node will simply AppHash,
	// which is completely fine. This means we were syncing from a
	// bad peer, and we would likely AppHash anyways.
	st := s.sb.StateFromContext(ctx)
	valUpdates, err := s.executeStateTransition(ctx, st, blk)
	if err != nil {
		return nil, err
	}

	// If the blobs needed to process the block are not available, we
	// return an error. It is safe to use the slot off of the beacon block
	// since it has been verified as correct already.
	if !s.sb.AvailabilityStore().IsDataAvailable(
		ctx, blk.GetSlot(), blk.GetBody(),
	) {
		return nil, ErrDataNotAvailable
	}

	// If required, we want to forkchoice at the end of post
	// block processing.
	// TODO: this is hood as fuck.
	// We won't send a fcu if the block is bad, should be addressed
	// via ticker later.
	if err = s.blkBroker.Publish(ctx,
		asynctypes.NewEvent(
			ctx, messages.BeaconBlockFinalized, blk,
		),
	); err != nil {
		return nil, err
	}

	go s.sendPostBlockFCU(ctx, st, blk)

	return valUpdates.RemoveDuplicates().Sort(), nil
}

// executeStateTransition runs the stf.
func (s *Service[
	_, BeaconBlockT, _, _, BeaconStateT, _, _, _, _, _, _,
]) executeStateTransition(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) (transition.ValidatorUpdates, error) {
	startTime := time.Now()
	defer s.metrics.measureStateTransitionDuration(startTime)
	valUpdates, err := s.sp.Transition(
		&transition.Context{
			Context:          ctx,
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
		},
		st,
		blk,
	)
	return valUpdates, err
}
