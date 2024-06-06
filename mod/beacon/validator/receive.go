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

	"github.com/berachain/beacon-kit/mod/primitives"
)

// VerifyIncomingBlock verifies the state root of an incoming block
// and logs the process.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) VerifyIncomingBlock(
	ctx context.Context,
	blk BeaconBlockT,
) error {
	// Grab a copy of the state to verify the incoming block.
	preState := s.bsb.StateFromContext(ctx)

	// Force a sync of the startup head if we haven't done so already.
	//
	// TODO: This is a super hacky. It should be handled better elsewhere,
	// ideally via some broader sync service.
	s.forceStartupSyncOnce.Do(func() { s.forceStartupHead(ctx, preState) })

	// If the block is nil or a nil pointer, exit early.
	if blk.IsNil() {
		s.logger.Error(
			"aborting block verification - beacon block not found in proposal 🚫 ",
		)

		if s.shouldBuildOptimisticPayloads() {
			go s.handleRebuildPayloadForRejectedBlock(ctx, preState)
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
	postState := preState.Copy()

	// Verify the state root of the incoming block.
	if err := s.verifyStateRoot(
		ctx, postState, blk,
	); err != nil {
		// TODO: this is expensive because we are not caching the
		// previous result of HashTreeRoot().
		localStateRoot, htrErr := preState.HashTreeRoot()
		if htrErr != nil {
			return htrErr
		}

		s.logger.Error(
			"rejecting incoming beacon block ❌ ",
			"block_state_root",
			blk.GetStateRoot(),
			"local_state_root",
			primitives.Root(localStateRoot),
			"error",
			err,
		)

		if s.shouldBuildOptimisticPayloads() {
			go s.handleRebuildPayloadForRejectedBlock(ctx, preState)
		}

		return err
	}

	s.logger.Info(
		"state root verification succeeded - accepting incoming beacon block 🏎️ ",
		"state_root",
		blk.GetStateRoot(),
	)

	if s.shouldBuildOptimisticPayloads() {
		go s.handleOptimisticPayloadBuild(ctx, postState, blk)
	}

	return nil
}

// VerifyIncomingBlobs receives blobs from the network and processes them.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) VerifyIncomingBlobs(
	_ context.Context,
	blk BeaconBlockT,
	sidecars BlobSidecarsT,
) error {
	if blk.IsNil() {
		s.logger.Error(
			"aborting blob verification - beacon block not found in proposal 🚫 ",
		)
		return ErrNilBlk
	}

	// If there are no blobs to verify, return early.
	if sidecars.Len() == 0 {
		s.logger.Info(
			"no blob sidecars to verify, skipping verifier 🧢 ",
			"slot",
			blk.GetSlot(),
		)
		return nil
	}

	s.logger.Info(
		"received incoming blob sidecars 🚔 ",
		"state_root", blk.GetStateRoot(),
	)

	// Verify the blobs and ensure they match the local state.
	if err := s.blobProcessor.VerifyBlobs(blk.GetSlot(), sidecars); err != nil {
		s.logger.Error(
			"rejecting incoming blob sidecars ❌ ",
			"error", err,
		)
		return err
	}

	s.logger.Info(
		"blob sidecars verification succeeded - accepting incoming blob sidecars 💦 ",
		"num_blobs",
		sidecars.Len(),
	)

	return nil
}
