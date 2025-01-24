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

package core

import (
	"context"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"golang.org/x/sync/errgroup"
)

// processExecutionPayload processes the execution payload and ensures it
// matches the local state.
func (sp *StateProcessor[ContextT]) processExecutionPayload(
	ctx ContextT, st *statedb.StateDB, blk *ctypes.BeaconBlock,
) error {
	var (
		body    = blk.GetBody()
		payload = body.GetExecutionPayload()
		header  *ctypes.ExecutionPayloadHeader
		g, gCtx = errgroup.WithContext(context.Background())
	)

	payloadTimestamp := payload.GetTimestamp().Unwrap()
	consensusTimestamp := ctx.GetConsensusTime().Unwrap()

	sp.metrics.gaugeTimestamps(payloadTimestamp, consensusTimestamp)

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
			return sp.validateExecutionPayload(gCtx, st, blk, ctx.GetOptimisticEngine())
		})
	}

	// Get the execution payload header.
	g.Go(func() error {
		var err error
		header, err = payload.ToHeader()
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if ctx.GetMeterGas() {
		sp.metrics.gaugeBlockGasUsed(
			payload.GetNumber(), payload.GetGasUsed(), payload.GetBlobGasUsed(),
		)
	}

	// Set the latest execution payload header.
	return st.SetLatestExecutionPayloadHeader(header)
}

// validateExecutionPayload validates the execution payload against both local
// state and the execution engine.
func (sp *StateProcessor[_]) validateExecutionPayload(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	optimisticEngine bool,
) error {
	if err := sp.validateStatelessPayload(blk); err != nil {
		return err
	}
	return sp.validateStatefulPayload(ctx, st, blk, optimisticEngine)
}

// validateStatelessPayload performs stateless checks on the execution payload.
func (sp *StateProcessor[_]) validateStatelessPayload(blk *ctypes.BeaconBlock) error {
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
func (sp *StateProcessor[_]) validateStatefulPayload(
	ctx context.Context, st *statedb.StateDB, blk *ctypes.BeaconBlock, optimisticEngine bool,
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
		ctx, ctypes.BuildNewPayloadRequest(
			payload,
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			optimisticEngine,
		),
	); err != nil {
		switch {
		// Skip ErrAcceptedPayloadStatus here. This status will be resolved
		// by the following forkchoice update, turning the payload into VALID
		// or INVALID.
		case errors.Is(err, engineerrors.ErrAcceptedPayloadStatus):
		default:
			return err
		}
	}

	// Since we have single slot finality, the previous block is already
	// considered final from cometBFT perspective. The newPayload being
	// submitted from this proposal should be set as the new head of the
	// chain, and we can update the finalized block with the EL. We send the
	// FCU now which fully verifies that the block is VALID or INVALID.
	_, _, err = sp.executionEngine.NotifyForkchoiceUpdate(
		ctx, ctypes.BuildForkchoiceUpdateRequestNoAttrs(
			&engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      payload.GetBlockHash(),
				SafeBlockHash:      safeHash,
				FinalizedBlockHash: safeHash,
			},
			sp.cs.ActiveForkVersionForSlot(blk.GetSlot()),
		),
	)
	// If we are unable to set the block as the head for any reason, it is
	// considered INVALID.
	if err != nil {
		return err
	}

	// Verify RANDAO
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	expectedMix, err := st.GetRandaoMixAtIndex(
		sp.cs.SlotToEpoch(slot).Unwrap() % sp.cs.EpochsPerHistoricalVector(),
	)
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
