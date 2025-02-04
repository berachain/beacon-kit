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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"golang.org/x/sync/errgroup"
)

// processExecutionPayload processes the execution payload and ensures it
// matches the local state.
func (sp *StateProcessor) processExecutionPayload(
	ctx transition.ReadOnlyContext,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	var (
		body    = blk.GetBody()
		payload = body.GetExecutionPayload()
		header  = &ctypes.ExecutionPayloadHeader{} // appeases nilaway
		g, gCtx = errgroup.WithContext(ctx.ConsensusCtx())
	)

	payloadTimestamp := payload.GetTimestamp().Unwrap()
	consensusTimestamp := ctx.ConsensusTime().Unwrap()

	sp.metrics.gaugeTimestamps(payloadTimestamp, consensusTimestamp)

	sp.logger.Info("processExecutionPayload",
		"consensus height", blk.GetSlot().Unwrap(),
		"payload height", payload.GetNumber().Unwrap(),
		"payload timestamp", payloadTimestamp,
		"consensus timestamp", consensusTimestamp,
		"verify payload", ctx.VerifyPayload(),
	)

	// Perform payload verification only if the context is configured as such.
	if ctx.VerifyPayload() {
		g.Go(func() error {
			return sp.validateExecutionPayload(gCtx, st, blk, ctx.OptimisticEngine())
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

	if ctx.MeterGas() {
		sp.metrics.gaugeBlockGasUsed(
			payload.GetNumber(), payload.GetGasUsed(), payload.GetBlobGasUsed(),
		)
	}

	// Set the latest execution payload header.
	return st.SetLatestExecutionPayloadHeader(header)
}

// validateExecutionPayload validates the execution payload against both local
// state and the execution engine.
func (sp *StateProcessor) validateExecutionPayload(
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
func (sp *StateProcessor) validateStatelessPayload(blk *ctypes.BeaconBlock) error {
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

	// No need to verify bounded number of commitments here, since it is
	// verified early on in ProcessProposal.
	return nil
}

// validateStatefulPayload performs stateful checks on the execution payload.
func (sp *StateProcessor) validateStatefulPayload(
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
