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

package core

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"golang.org/x/sync/errgroup"
)

// processExecutionPayload processes the execution payload and ensures it
// matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT,
	_, _, _, ExecutionPayloadHeaderT, _, _, _, _, _, _,
]) processExecutionPayload(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	var (
		body    = blk.GetBody()
		payload = body.GetExecutionPayload()
		header  ExecutionPayloadHeaderT
		g, gCtx = errgroup.WithContext(context.Background())
	)

	// Skip payload verification if the context is configured as such.
	if !ctx.GetSkipPayloadVerification() {
		g.Go(func() error {
			return sp.validateExecutionPayload(
				gCtx, st, blk, ctx.GetOptimisticEngine(),
			)
		})
	}

	// Get the execution payload header.
	g.Go(func() error {
		var err error
		header, err = payload.ToHeader(
			sp.txsMerkleizer, sp.cs.MaxWithdrawalsPerPayload(),
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// Set the latest execution payload header.
	return st.SetLatestExecutionPayloadHeader(header)
}

// validateExecutionPayload validates the execution payload against both local
// state
// and the execution engine.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT,
	_, _, _, _, _, _, _, _, _, _, _,
]) validateExecutionPayload(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	optimisticEngine bool,
) error {
	body := blk.GetBody()
	payload := body.GetExecutionPayload()

	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	// We want to check to ensure the chain is canonical with respect to the
	// parent hash before we let the execution client know about the
	// payload,
	// this is to prevent Polygon style re-orgs from being triggered by a
	// malicious actor who tries to force clients to accept a non-canonical
	// block that passes block validity checks.
	if safeHash := lph.GetBlockHash(); safeHash != payload.GetParentHash() {
		return errors.Wrapf(
			ErrParentPayloadHashMismatch,
			"parent block with hash %x is not finalized, expected finalized hash %x",
			payload.GetParentHash(),
			safeHash,
		)
	}

	parentBeaconBlockRoot := blk.GetParentBlockRoot()
	if err = sp.executionEngine.VerifyAndNotifyNewPayload(
		ctx, engineprimitives.BuildNewPayloadRequest(
			payload,
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			optimisticEngine,
		),
	); err != nil {
		return err
	}

	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// When we are verifying a payload we expect that it was produced by
	// the proposer for the slot that it is for.
	expectedMix, err := st.GetRandaoMixAtIndex(
		uint64(sp.cs.SlotToEpoch(slot)) % sp.cs.EpochsPerHistoricalVector())
	if err != nil {
		return err
	}

	// Ensure the prev randao matches the local state.
	if payload.GetPrevRandao() != expectedMix {
		return errors.Wrapf(
			ErrRandaoMixMismatch,
			"prev randao does not match, expected: %x, got: %x",
			expectedMix, payload.GetPrevRandao(),
		)
	}

	// TODO: Verify timestamp data once Clock is done.
	// if expectedTime, err := spec.TimeAtSlot(slot, genesisTime); err !=
	// nil { 	return errors.Newf("slot or genesis time in state is corrupt,
	// cannot
	// compute time: %v", err)
	// } else if payload.Timestamp != expectedTime {
	// 	return errors.Newf("state at slot %d, genesis time %d, expected
	// execution
	// payload time %d, but got %d",
	// 		slot, genesisTime, expectedTime, payload.Timestamp)
	// }

	// Verify the number of blobs.
	blobKzgCommitments := body.GetBlobKzgCommitments()
	if uint64(len(blobKzgCommitments)) > sp.cs.MaxBlobsPerBlock() {
		return errors.Wrapf(
			ErrExceedsBlockBlobLimit,
			"expected: %d, got: %d",
			sp.cs.MaxBlobsPerBlock(), len(blobKzgCommitments),
		)
	}

	// Verify the number of withdrawals.
	// TODO: This is in the wrong spot I think.
	if withdrawals := payload.GetWithdrawals(); uint64(
		len(payload.GetWithdrawals()),
	) > sp.cs.MaxWithdrawalsPerPayload() {
		return errors.Newf(
			"too many withdrawals, expected: %d, got: %d",
			sp.cs.MaxWithdrawalsPerPayload(), len(withdrawals),
		)
	}
	return nil
}
