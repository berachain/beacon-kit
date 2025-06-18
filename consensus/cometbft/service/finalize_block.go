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
	"fmt"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/sourcegraph/conc/iter"
)

func (s *Service) finalizeBlock(
	ctx context.Context,
	req *cmtabci.FinalizeBlockRequest,
) (*cmtabci.FinalizeBlockResponse, error) {
	res, err := s.finalizeBlockInternal(ctx, req)
	if res != nil {
		res.AppHash = s.workingHash()
	}

	return res, err
}

func (s *Service) finalizeBlockInternal(
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
	var finalizeBlockState *state
	if s.finalStateHash == nil {
		finalizeBlockState = s.resetState(ctx)
	} else {
		// Preserve the CosmosSDK context while using the correct base ctx.

		st, found := s.candidateStates[*s.finalStateHash]
		if !found {
			panic(fmt.Errorf("missing candidate state for hash %s", *s.finalStateHash))
		}
		st.state.SetContext(st.state.Context().WithContext(ctx))
		stateHash := string(req.Hash)
		s.finalStateHash = &stateHash
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

	valUpdates, err := s.Blockchain.FinalizeBlock(finalizeBlockState.Context(), req)
	if err != nil {
		return nil, err
	}
	s.candidateStates[string(req.Hash)] = &CacheElement{
		state:      finalizeBlockState,
		valUpdates: valUpdates,
	}

	formattedValUpdates, err := iter.MapErr(
		valUpdates,
		convertValidatorUpdate[cmtabci.ValidatorUpdate],
	)
	if err != nil {
		return nil, err
	}

	cp := s.cmtConsensusParams.ToProto()
	return &cmtabci.FinalizeBlockResponse{
		TxResults:             txResults,
		ValidatorUpdates:      formattedValUpdates,
		ConsensusParamUpdates: &cp,
	}, nil
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
	if s.finalStateHash == nil {
		// this is unexpected since workingHash is called only after
		// internalFinalizeBlock. Panic appeases nilaway.
		panic(fmt.Errorf("workingHash: %w", errNilFinalizeBlockState))
	}
	s.candidateStates[*s.finalStateHash].state.ms.Write()

	// Get the hash of all writes in order to return the apphash to the comet in
	// finalizeBlock.
	commitHash := s.sm.GetCommitMultiStore().WorkingHash()
	s.logger.Debug(
		"hash of all writes",
		"workingHash",
		fmt.Sprintf("%X", commitHash),
	)

	return commitHash
}

func (s *Service) validateFinalizeBlockHeight(
	req *cmtabci.FinalizeBlockRequest,
) error {
	if req.Height < 1 {
		return fmt.Errorf(
			"finalizeBlock at height %v: %w",
			req.Height,
			errInvalidHeight,
		)
	}

	lastBlockHeight := s.LastBlockHeight()

	// expectedHeight holds the expected height to validate
	var expectedHeight int64
	if lastBlockHeight == 0 && s.initialHeight > 1 {
		// In this case, we're validating the first block of the chain, i.e no
		// previous commit. The height we're expecting is the initial height.
		expectedHeight = s.initialHeight
	} else {
		// This case can mean two things:
		//
		// - Either there was already a previous commit in the store, in which
		// case we increment the version from there.
		// - Or there was no previous commit, in which case we start at version
		// 1.
		expectedHeight = lastBlockHeight + 1
	}

	if req.Height != expectedHeight {
		return fmt.Errorf(
			"invalid height: %d; expected: %d",
			req.Height,
			expectedHeight,
		)
	}

	return nil
}
