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
	"fmt"

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"golang.org/x/sync/errgroup"
)

// processExecutionPayload processes the execution payload and ensures it
// matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT,
	_, _, _, ExecutionPayloadHeaderT, _, _, _, _, _, _, _, _,
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

	sp.logger.Info("processExecutionPayload",
		fmt.Sprintf("consensus height %d", ctx.GetConsensusBlockHeight()),
		fmt.Sprintf("payload height %d", payload.GetNumber()),
		fmt.Sprintf("payload timestamp %d", payload.GetTimestamp()),
		fmt.Sprintf("bound timestamp %d", payload.GetTimestamp()),
		fmt.Sprintf("skip verification %t", ctx.GetSkipPayloadVerification()),
	)

	// Skip payload verification if the context is configured as such.
	if !ctx.GetSkipPayloadVerification() {
		g.Go(func() error {
			return sp.validateExecutionPayload(
				gCtx, st, blk,
				ctx.GetNextPayloadTimestamp(),
				ctx.GetOptimisticEngine(),
			)
		})
	}

	// Get the execution payload header. TODO: This is live on bArtio with a bug
	// and needs to be hardforked off of. We check for version and convert to
	// header based on that version as a temporary solution to avoid breaking
	// changes.
	g.Go(func() error {
		var err error
		header, err = payload.ToHeader(
			sp.cs.MaxWithdrawalsPerPayload(),
			sp.cs.DepositEth1ChainID(),
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
// state and the execution engine.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT,
	_, _, _, _, _, _, _, _, _, _, _, _, _,
]) validateExecutionPayload(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	nextPayloadTimestamp math.U64,
	optimisticEngine bool,
) error {
	if err := sp.validateStatelessPayload(
		blk,
		nextPayloadTimestamp,
	); err != nil {
		return err
	}
	return sp.validateStatefulPayload(ctx, st, blk, optimisticEngine)
}

// validateStatelessPayload performs stateless checks on the execution payload.
func (sp *StateProcessor[
	BeaconBlockT, _, _, _,
	_, _, _, _, _, _, _, _, _, _, _, _, _,
]) validateStatelessPayload(
	blk BeaconBlockT,
	nextPayloadTimestamp math.U64,
) error {
	body := blk.GetBody()
	payload := body.GetExecutionPayload()

	// We skip timestamp check on Bartio for backward compatibility reasons
	// TODO: enforce the check when we drop other Bartio special cases.
	if sp.cs.DepositEth1ChainID() != spec.BartioChainID {
		if pt := payload.GetTimestamp(); pt >= nextPayloadTimestamp {
			return errors.Wrapf(
				ErrTooFarInTheFuture,
				"payload timestamp, max: %d, got: %d",
				nextPayloadTimestamp, pt,
			)
		}
	}

	// Verify the number of withdrawals.
	withdrawals := payload.GetWithdrawals()
	if uint64(len(withdrawals)) > sp.cs.MaxWithdrawalsPerPayload() {
		return errors.Wrapf(
			ErrExceedMaximumWithdrawals,
			"too many withdrawals, expected: %d, got: %d",
			sp.cs.MaxWithdrawalsPerPayload(), len(withdrawals),
		)
	}

	// Verify the number of blobs.
	blobKzgCommitments := body.GetBlobKzgCommitments()
	if uint64(len(blobKzgCommitments)) > sp.cs.MaxBlobsPerBlock() {
		return errors.Wrapf(
			ErrExceedsBlockBlobLimit,
			"expected: %d, got: %d",
			sp.cs.MaxBlobsPerBlock(), len(blobKzgCommitments),
		)
	}

	return nil
}

// validateStatefulPayload performs stateful checks on the execution payload.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT,
	_, _, _, _, _, _, _, _, _, _, _, _, _,
]) validateStatefulPayload(
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

	// Check chain canonicity
	safeHash := lph.GetBlockHash()
	if safeHash != payload.GetParentHash() {
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

	// Verify RANDAO
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	expectedMix, err := st.GetRandaoMixAtIndex(
		sp.cs.SlotToEpoch(slot).Unwrap() % sp.cs.EpochsPerHistoricalVector())
	if err != nil {
		return err
	}

	if payload.GetPrevRandao() != expectedMix {
		return errors.Wrapf(
			ErrRandaoMixMismatch,
			"prev randao does not match, expected: %x, got: %x",
			expectedMix, payload.GetPrevRandao(),
		)
	}

	return nil
}
