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

	payloadtime "github.com/berachain/beacon-kit/mod/beacon/payload-time"
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

	payloadTimestamp := payload.GetTimestamp().Unwrap()
	consensusTimestamp := ctx.GetConsensusTime().Unwrap()

	// add telemetry for the diff between payloadTimestamp and
	// consensusTimestamp
	err := sp.metrics.gaugePayloadConsensusTimestampDiff(
		payloadTimestamp, consensusTimestamp)
	if err != nil {
		sp.logger.Error("failed to gauge timestamp diff", "error", err)
	}

	sp.logger.Info("processExecutionPayload",
		"consensus height", blk.GetSlot().Unwrap(),
		"payload height", payload.GetNumber().Unwrap(),
		"payload timestamp", payloadTimestamp,
		"consensus timestamp", consensusTimestamp,
		"skip payload verification", ctx.GetSkipPayloadVerification(),
	)

	// Skip payload verification if the context is configured as such.
	if !ctx.GetSkipPayloadVerification() {
		g.Go(func() error {
			return sp.validateExecutionPayload(
				gCtx, st, blk,
				ctx.GetConsensusTime(),
				ctx.GetOptimisticEngine(),
			)
		})
	}

	// Get the execution payload header. TODO: This is live on bArtio with a bug
	// and needs to be hardforked off of. We check for version and convert to
	// header based on that version as a temporary solution to avoid breaking
	// changes.
	g.Go(func() error {
		header, err = payload.ToHeader()
		return err
	})

	if err = g.Wait(); err != nil {
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
	consensusTime math.U64,
	optimisticEngine bool,
) error {
	if err := sp.validateStatelessPayload(blk); err != nil {
		return err
	}
	return sp.validateStatefulPayload(
		ctx,
		st,
		blk,
		consensusTime,
		optimisticEngine,
	)
}

// validateStatelessPayload performs stateless checks on the execution payload.
func (sp *StateProcessor[
	BeaconBlockT, _, _, _,
	_, _, _, _, _, _, _, _, _, _, _, _, _,
]) validateStatelessPayload(
	blk BeaconBlockT,
) error {
	body := blk.GetBody()
	payload := body.GetExecutionPayload()

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
	consensusTime math.U64,
	optimisticEngine bool,
) error {
	body := blk.GetBody()
	payload := body.GetExecutionPayload()

	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	// We skip timestamp check on Bartio for backward compatibility reasons
	// TODO: enforce the check when we drop other Bartio special cases.
	if sp.cs.DepositEth1ChainID() != spec.BartioChainID {
		if err = payloadtime.Verify(
			consensusTime,
			lph.GetTimestamp(),
			payload.GetTimestamp(),
		); err != nil {
			return err
		}
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
