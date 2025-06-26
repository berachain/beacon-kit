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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/cache"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
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

	// Check whether currently block hash is already available. If so
	// we may speed up block finalization.
	hash := string(req.Hash)
	switch cached, err := s.cachedStates.GetCached(hash); {
	case err == nil:
		// Block with height equal to initial height is special and we can't rely on its cache.
		// This is because Genesis state is cached but not committed (and purged from s.cachedStates)
		// We handle the case below.
		if req.Height > s.initialHeight {
			if err = s.cachedStates.MarkAsFinal(hash); err != nil {
				return nil, fmt.Errorf("failed marking state as final, hash %s, height %d: %w", hash, req.Height, err)
			}

			finalState := cached.State
			var (
				signedBlk *ctypes.SignedBeaconBlock
				sidecars  datypes.BlobSidecars
			)
			signedBlk, sidecars, err = s.Blockchain.ParseBeaconBlock(req)
			if err != nil {
				return nil, fmt.Errorf("finalize block: failed parsing block: %w", err)
			}
			blk := signedBlk.GetBeaconBlock()

			if err = s.Blockchain.FinalizeSidecars(
				finalState.Context(),
				req.SyncingToHeight,
				blk,
				sidecars,
			); err != nil {
				return nil, fmt.Errorf("failed finalizing sidecars: %w", err)
			}
			if err = s.Blockchain.PostFinalizeBlockOps(
				finalState.Context(),
				blk,
			); err != nil {
				return nil, fmt.Errorf("finalize block: failed post finalize block ops: %w", err)
			}
			return formatFinalizeBlockResponse(cached.ValUpdates, len(req.Txs), s.cmtConsensusParams)
		}

	case errors.Is(err, cache.ErrStateNotFound):
		// this is a benign error, it just signal it's the first time we see the block being finalized
		// Keep processing below

	default:
		return nil, fmt.Errorf("failed checking cached state, hash %s, height %d: %w", hash, req.Height, err)
	}

	// Block has not been cached already, as it happens when we block sync.
	// Normally we would directly process it but first block post initial height
	// needs special handling.
	var finalState *cache.State
	cachedFinalHash, cachedFinalState, err := s.cachedStates.GetFinal()
	switch {
	case errors.Is(err, cache.ErrNoFinalState):
		hash = string(req.Hash) // not necessary but clarifies intent
		finalState = s.resetState(ctx)

	case err == nil:
		hash = cachedFinalHash
		finalState = cachedFinalState

		// Preserve the CosmosSDK context while using the correct base ctx.
		finalState.SetContext(finalState.Context().WithContext(ctx))

	default:
		return nil, fmt.Errorf("failed checking cached final state, hash %s, height %d: %w", hash, req.Height, err)
	}

	valUpdates, err := s.Blockchain.FinalizeBlock(finalState.Context(), req)
	if err != nil {
		return nil, err
	}

	// Cache and mark state as final
	s.cachedStates.Cache(hash, &cache.Element{
		State:      finalState,
		ValUpdates: valUpdates,
	})
	if err = s.cachedStates.MarkAsFinal(hash); err != nil {
		return nil, fmt.Errorf("failed marking state as final, hash %s, height %d: %w", hash, req.Height, err)
	}

	return formatFinalizeBlockResponse(valUpdates, len(req.Txs), s.cmtConsensusParams)
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
	_, finalState, err := s.cachedStates.GetFinal()
	if err != nil {
		// This is unexpected since CometBFT should call Commit only
		// after FinalizeBlock has been called. Panic appeases nilaway.
		panic(fmt.Errorf("workingHash: %w", err))
	}
	finalState.Write()

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

func formatFinalizeBlockResponse(
	valUpdates transition.ValidatorUpdates,
	txsLen int,
	cmtConsensusParams *cmttypes.ConsensusParams,
) (*cmtabci.FinalizeBlockResponse, error) {
	// This result format is expected by Comet. That actual execution will happen as part of the state transition.
	txResults := make([]*cmtabci.ExecTxResult, txsLen)
	for i := range txsLen {
		//nolint:mnd // its okay for now.
		txResults[i] = &cmtabci.ExecTxResult{
			Codespace: "sdk",
			Code:      2,
			Log:       "skip decoding",
			GasWanted: 0,
			GasUsed:   0,
		}
	}

	formattedValUpdates, err := iter.MapErr(
		valUpdates,
		convertValidatorUpdate[cmtabci.ValidatorUpdate],
	)
	if err != nil {
		return nil, err
	}

	cp := cmtConsensusParams.ToProto()
	return &cmtabci.FinalizeBlockResponse{
		TxResults:             txResults,
		ValidatorUpdates:      formattedValUpdates,
		ConsensusParamUpdates: &cp,
	}, nil
}
