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

package cometbft

import (
	"context"
	"errors"
	"fmt"

	statem "github.com/berachain/beacon-kit/consensus/cometbft/service/state"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/sourcegraph/conc/iter"
)

func (s *Service) finalizeBlock(
	ctx context.Context,
	req *cmtabci.FinalizeBlockRequest,
) (*cmtabci.FinalizeBlockResponse, error) {
	if err := s.validateFinalizeBlockHeight(req); err != nil {
		return nil, err
	}

	// finalizeBlockState should be set on InitChain or ProcessProposal. If it
	// is nil, it means we are replaying this block and we need to set the state
	// here given that during block replay ProcessProposal is not executed by
	// CometBFT.
	stateCtx, err := s.stateHandler.GetFinalizeStateContext()
	switch {
	case err == nil:
		// Preserve the CosmosSDK context while using the correct base ctx.
		stateCtx = stateCtx.WithContext(ctx)
	case errors.Is(err, statem.ErrNilFinalizeBlockState):
		stateCtx, err = s.stateHandler.NewStateCtx(ctx, req.Height)
		if err != nil {
			panic(fmt.Errorf("finalize block: failed retrieving state context after reset: %w", err))
		}
	default:
		panic(fmt.Errorf("finalize block: failed retrieving state context: %w", err))
	}

	// This result format is expected by Comet. That actual execution will happen as part of the state transition.
	txResults := make([]*cmtabci.ExecTxResult, len(req.Txs))
	for i := range req.Txs {
		//nolint:mnd // its okay for now.
		txResults[i] = &cmtabci.ExecTxResult{
			Codespace: "sdk",
			Code:      2,
			Log:       "skip decoding",
			GasWanted: 0,
			GasUsed:   0,
		}
	}

	finalizeBlock, err := s.Blockchain.FinalizeBlock(stateCtx, req)
	if err != nil {
		return nil, err
	}

	valUpdates, err := iter.MapErr(
		finalizeBlock,
		convertValidatorUpdate[cmtabci.ValidatorUpdate],
	)
	if err != nil {
		return nil, err
	}

	cp := s.cmtConsensusParams.ToProto()
	appHash := s.workingHash() // to be done only if no errors above
	res := &cmtabci.FinalizeBlockResponse{
		TxResults:             txResults,
		ValidatorUpdates:      valUpdates,
		ConsensusParamUpdates: &cp,
		AppHash:               appHash,
	}
	return res, nil
}

// workingHash gets the apphash that will be finalized in commit.
// These writes will be persisted to the root multi-store
// (s.smGet.CommitMultiStore()) and flushed to disk  in the
// Commit phase. This means when the ABCI client requests
// Commit(), the application state transitions will be flushed
// to disk and as a result, but we already have an application
// Merkle root.
func (s *Service) workingHash() []byte {
	// Write the FinalizeBlock state into branched storage and commit the
	// MultiStore. The write to the FinalizeBlock state writes all state
	// transitions to the root MultiStore (s.sm.GetCommitMultiStore())
	// so when Commit() is called it persists those values.
	commitHash, err := s.stateHandler.WriteFinalizeState()
	if err != nil {
		panic(fmt.Errorf("workingHash: %w", err))
	}

	// Get the hash of all writes in order to return the apphash to the comet in
	// finalizeBlock.
	s.logger.Debug(
		"hash of all writes",
		"workingHash", fmt.Sprintf("%X", commitHash),
	)

	return commitHash
}

func (s *Service) validateFinalizeBlockHeight(req *cmtabci.FinalizeBlockRequest) error {
	if req.Height < statem.InitialHeight {
		return fmt.Errorf(
			"finalizeBlock at height %v: %w",
			req.Height,
			errInvalidHeight,
		)
	}

	// expectedHeight holds the expected height to validate
	lastBlockHeight := s.LastBlockHeight()
	expectedHeight := lastBlockHeight + 1

	if req.Height != expectedHeight {
		return fmt.Errorf(
			"invalid height: %d; expected: %d",
			req.Height,
			expectedHeight,
		)
	}

	return nil
}
