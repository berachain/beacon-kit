// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// Licensed work is provided "as is" without warranties or conditions of any kind.

package validator

import (
	"context"
	"fmt"
	"time"

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// BuildBlockAndSidecars builds a new beacon block and its associated sidecars.
func (s *Service[_]) BuildBlockAndSidecars(
	ctx context.Context,
	slotData types.SlotData,
) ([]byte, []byte, error) {
	startTime := time.Now()
	defer s.metrics.measureRequestBlockForProposalTime(startTime)

	// Initialize variables for block, sidecars, and fork data.
	var (
		blk      *ctypes.BeaconBlock
		sidecars datypes.BlobSidecars
		forkData *ctypes.ForkData
	)

	st := s.sb.StateFromContext(ctx)

	// Prepare state for the requested slot.
	if err := s.prepareStateForSlot(st, slotData.GetSlot()); err != nil {
		return nil, nil, err
	}

	// Generate fork data and Randao reveal.
	forkData, err := s.buildForkData(st, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	reveal, err := s.buildRandaoReveal(forkData, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	// Create an empty beacon block.
	blk, err = s.getEmptyBeaconBlockForSlot(st, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	// Retrieve execution payload.
	envelope, err := s.retrieveExecutionPayload(ctx, st, blk, slotData)
	if err != nil || envelope == nil {
		return nil, nil, fmt.Errorf("failed to retrieve execution payload: %w", err)
	}

	// Build the block body.
	if err := s.buildBlockBody(ctx, st, blk, reveal, envelope, slotData); err != nil {
		return nil, nil, err
	}

	// Compute and set the state root.
	if err := s.computeAndSetStateRoot(ctx, slotData.GetProposerAddress(), slotData.GetConsensusTime(), st, blk); err != nil {
		return nil, nil, err
	}

	// Build blob sidecars.
	sidecars, err = s.blobFactory.BuildSidecars(blk, envelope.GetBlobsBundle(), s.signer, forkData)
	if err != nil {
		return nil, nil, err
	}

	// Marshal the block and sidecars into SSZ format.
	blkBytes, blkErr := blk.MarshalSSZ()
	sidecarsBytes, sidecarsErr := sidecars.MarshalSSZ()

	if blkErr != nil || sidecarsErr != nil {
		return nil, nil, fmt.Errorf("failed to marshal block or sidecars: %v, %v", blkErr, sidecarsErr)
	}

	s.logger.Info(
		"Beacon block successfully built",
		"slot", slotData.GetSlot().Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	return blkBytes, sidecarsBytes, nil
}

// prepareStateForSlot prepares the state for the given slot.
func (s *Service[_]) prepareStateForSlot(st *statedb.StateDB, slot math.Slot) error {
	_, err := s.stateProcessor.ProcessSlots(st, slot)
	return err
}

// buildForkData constructs fork data for a given slot.
func (s *Service[_]) buildForkData(st *statedb.StateDB, slot math.Slot) (*ctypes.ForkData, error) {
	epoch := s.chainSpec.SlotToEpoch(slot)
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	return ctypes.NewForkData(
		version.FromUint32[common.Version](s.chainSpec.ActiveForkVersionForEpoch(epoch)),
		genesisValidatorsRoot,
	), nil
}

// buildRandaoReveal generates the Randao reveal for a slot.
func (s *Service[_]) buildRandaoReveal(forkData *ctypes.ForkData, slot math.Slot) (crypto.BLSSignature, error) {
	epoch := s.chainSpec.SlotToEpoch(slot)
	signingRoot := forkData.ComputeRandaoSigningRoot(s.chainSpec.DomainTypeRandao(), epoch)
	return s.signer.Sign(signingRoot[:])
}

// retrieveExecutionPayload fetches the execution payload for a block.
func (s *Service[_]) retrieveExecutionPayload(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	slotData types.SlotData,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	envelope, err := s.localPayloadBuilder.RetrievePayload(ctx, blk.GetSlot(), blk.GetParentBlockRoot())
	if err == nil {
		return envelope, nil
	}

	// Handle fallback payload retrieval.
	s.metrics.failedToRetrievePayload(blk.GetSlot(), err)

	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	return s.localPayloadBuilder.RequestPayloadSync(
		ctx,
		st,
		blk.GetSlot(),
		payloadtime.Next(slotData.GetConsensusTime(), lph.GetTimestamp(), false).Unwrap(),
		blk.GetParentBlockRoot(),
		lph.GetBlockHash(),
		lph.GetParentHash(),
	)
}

// computeAndSetStateRoot computes and sets the state root for a block.
func (s *Service[_]) computeAndSetStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	stateRoot, err := s.computeStateRoot(ctx, proposerAddress, consensusTime, st, blk)
	if err != nil {
		s.logger.Error("Failed to compute state root", "slot", blk.GetSlot().Base10(), "error", err)
		return err
	}
	blk.SetStateRoot(stateRoot)
	return nil
}

// computeStateRoot computes the state root for a block.
func (s *Service[_]) computeStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
) (common.Root, error) {
	startTime := time.Now()
	defer s.metrics.measureStateRootComputationTime(startTime)

	_, err := s.stateProcessor.Transition(
		&transition.Context{
			Context:                 ctx,
			OptimisticEngine:        true,
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			SkipValidateRandao:      true,
			ProposerAddress:         proposerAddress,
			ConsensusTime:           consensusTime,
		},
		st, blk,
	)
	if err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot(), nil
}
