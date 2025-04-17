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

package validator

import (
	"context"
	"fmt"
	"time"

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// BuildBlockAndSidecars builds a new beacon block.
//
//nolint:funlen // comments are pretty verbose
func (s *Service) BuildBlockAndSidecars(
	ctx context.Context,
	slotData *types.SlotData,
) ([]byte, []byte, error) {
	startTime := time.Now()
	defer s.metrics.measureRequestBlockForProposalTime(startTime)

	if !s.localPayloadBuilder.Enabled() {
		// node is not supposed to build blocks
		return nil, nil, builder.ErrPayloadBuilderDisabled
	}

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propagating the finalized
	// and safe block hashes to the execution client.
	st := s.sb.StateFromContext(ctx)

	// blkSlot is the height for the next block, which consensus is requesting BeaconKit to build.
	blkSlot := slotData.GetSlot()

	// Prepare the state such that it is ready to build a block for the requested slot.
	if _, err := s.stateProcessor.ProcessSlots(st, blkSlot); err != nil {
		return nil, nil, err
	}

	// Grab parent block root for payload request.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		(blkSlot.Unwrap() - 1) % s.chainSpec.SlotsPerHistoricalRoot(),
	)
	if err != nil {
		return nil, nil, err
	}

	// Get the payload for the block.
	envelope, err := s.retrieveExecutionPayload(ctx, st, parentBlockRoot, slotData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed retrieving execution payload: %w", err)
	}

	// We introduce hard forks with the expectation that the first block proposed after the
	// hard fork timestamp is when new rules apply. When building blocks, we provide the Execution
	// Layer client with a timestamp, and it will create its payload based on that timestamp. We
	// must use this same timestamp from the payload to build the beacon block. This ensures that
	// we are building on the same fork version as the Execution Layer.
	timestamp := envelope.GetExecutionPayload().GetTimestamp()

	// Build forkdata used for the signing root of the reveal and the sidecars.
	forkData, err := s.buildForkData(st, timestamp)
	if err != nil {
		return nil, nil, err
	}

	// Create a new empty block from the current state.
	blk, err := s.getEmptyBeaconBlockForSlot(st, blkSlot, forkData.CurrentVersion, parentBlockRoot)
	if err != nil {
		return nil, nil, err
	}

	// Build the reveal for the current slot.
	// TODO: We can optimize to pre-compute this in parallel?
	reveal, err := s.buildRandaoReveal(forkData, blkSlot)
	if err != nil {
		return nil, nil, err
	}

	// We have to assemble the block body prior to producing the sidecars
	// since we need to generate the inclusion proofs.
	if err = s.buildBlockBody(ctx, st, blk, reveal, envelope); err != nil {
		return nil, nil, fmt.Errorf("failed build block body: %w", err)
	}

	// Compute the state root for the block.
	if err = s.computeAndSetStateRoot(
		ctx,
		slotData.GetProposerAddress(),
		slotData.GetConsensusTime(),
		st,
		blk,
	); err != nil {
		return nil, nil, err
	}

	// Craft the signature and signed beacon block.
	signedBlk, err := ctypes.NewSignedBeaconBlock(blk, forkData, s.chainSpec, s.signer)
	if err != nil {
		return nil, nil, err
	}

	// Produce blob sidecars with new StateRoot
	sidecars, err := s.blobFactory.BuildSidecars(signedBlk, envelope.GetBlobsBundle())
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info(
		"Beacon block successfully built",
		"slot", blkSlot.Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	signedBlkBytes, bbErr := signedBlk.MarshalSSZ()
	if bbErr != nil {
		return nil, nil, bbErr
	}
	sidecarsBytes, scErr := sidecars.MarshalSSZ()
	if scErr != nil {
		return nil, nil, scErr
	}

	return signedBlkBytes, sidecarsBytes, nil
}

// getEmptyBeaconBlockForSlot creates a new empty block.
func (s *Service) getEmptyBeaconBlockForSlot(
	st *statedb.StateDB, requestedSlot math.Slot,
	forkVersion common.Version, parentBlockRoot common.Root,
) (*ctypes.BeaconBlock, error) {
	// Get the proposer index for the slot.
	proposerIndex, err := st.ValidatorIndexByPubkey(
		s.signer.PublicKey(),
	)
	if err != nil {
		return nil, err
	}

	// Create a new block.
	return ctypes.NewBeaconBlockWithVersion(
		requestedSlot,
		proposerIndex,
		parentBlockRoot,
		forkVersion,
	)
}

func (s *Service) buildForkData(st *statedb.StateDB, timestamp math.U64) (*ctypes.ForkData, error) {
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	return ctypes.NewForkData(
		s.chainSpec.ActiveForkVersionForTimestamp(timestamp),
		genesisValidatorsRoot,
	), nil
}

// buildRandaoReveal builds a randao reveal for the given slot.
func (s *Service) buildRandaoReveal(
	forkData *ctypes.ForkData, slot math.Slot,
) (crypto.BLSSignature, error) {
	signingRoot := forkData.ComputeRandaoSigningRoot(
		s.chainSpec.DomainTypeRandao(),
		s.chainSpec.SlotToEpoch(slot),
	)
	signature, err := s.signer.Sign(signingRoot[:])
	if err != nil {
		return signature, fmt.Errorf("block building failed randao checks: %w", err)
	}
	return signature, nil
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service) retrieveExecutionPayload(
	ctx context.Context,
	st *statedb.StateDB,
	parentBlockRoot common.Root,
	slotData *types.SlotData,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	// TODO: Add external block builders to this flow.
	//
	// Get the payload for the block.
	slot := slotData.GetSlot()
	envelope, err := s.localPayloadBuilder.RetrievePayload(ctx, slot, parentBlockRoot)
	if err == nil {
		return envelope, nil
	}

	// If we failed to retrieve the payload, request a synchronous payload.
	//
	// NOTE: The state here is properly configured by the
	// prepareStateForBuilding
	//
	// call that needs to be called before requesting the Payload.
	// TODO: We should decouple the PayloadBuilder from BeaconState to make
	// this less confusing.
	s.metrics.failedToRetrievePayload(slot, err)

	// The latest execution payload header will be from the previous block
	// during the block building phase.
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	// We must prepare the state for the fork version of the new block being built to handle
	// the case where the new block is on a new fork version. Although we do not have the
	// confirmed timestamp by the EL, we will assume it to be `nextPayloadTimestamp` to decide
	// the new block's fork version.
	nextPayloadTimestamp := payloadtime.Next(
		slotData.GetConsensusTime(),
		lph.GetTimestamp(),
		false, // buildOptimistically
	)
	err = s.stateProcessor.ProcessFork(st, nextPayloadTimestamp, false)
	if err != nil {
		return nil, err
	}

	return s.localPayloadBuilder.RequestPayloadSync(
		ctx,
		st,
		slot,
		nextPayloadTimestamp,
		parentBlockRoot,
		lph.GetBlockHash(),
		lph.GetParentHash(),
	)
}

// BuildBlockBody assembles the block body with necessary components.
func (s *Service) buildBlockBody(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	reveal crypto.BLSSignature,
	envelope ctypes.BuiltExecutionPayloadEnv,
) error {
	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body == nil {
		return ErrNilBlkBody
	}

	// Set the reveal on the block body.
	body.SetRandaoReveal(reveal)

	// If we get returned a nil blobs bundle, we should return an error.
	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle == nil {
		return ErrNilBlobsBundle
	}

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// Dequeue deposits from the state.
	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return fmt.Errorf("failed loading eth1 deposit index: %w", err)
	}

	// Grab all previous deposits from genesis up to the current index + max deposits per block.
	deposits, err := s.sb.DepositStore().GetDepositsByIndex(
		ctx,
		constants.FirstDepositIndex,
		depositIndex+s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}
	if uint64(len(deposits)) < depositIndex {
		return errors.Wrapf(ErrDepositStoreIncomplete,
			"all historical deposits not available, expected: %d, got: %d",
			depositIndex, len(deposits),
		)
	}

	eth1Data := ctypes.NewEth1Data(deposits.HashTreeRoot())
	body.SetEth1Data(eth1Data)

	s.logger.Info(
		"Building block body with local deposits",
		"start_index", depositIndex, "num_deposits", uint64(len(deposits))-depositIndex,
	)
	body.SetDeposits(deposits[depositIndex:])

	// Set the graffiti on the block body.
	sizedGraffiti := bytes.ExtendToSize([]byte(s.cfg.Graffiti), bytes.B32Size)
	graffiti, err := bytes.ToBytes32(sizedGraffiti)
	if err != nil {
		return fmt.Errorf("failed processing graffiti: %w", err)
	}
	body.SetGraffiti(graffiti)

	// Fill in unused field with non-nil value
	body.SetSyncAggregate(&ctypes.SyncAggregate{})

	// Set the execution payload on the block body.
	body.SetExecutionPayload(envelope.GetExecutionPayload())

	if version.EqualsOrIsAfter(body.GetForkVersion(), version.Electra()) {
		// TODO(pectra): Remove the conversion once DecodeExecutionRequests constructor changed.
		encodedReqs := envelope.GetEncodedExecutionRequests()
		result := make([][]byte, len(encodedReqs))
		for i, req := range encodedReqs {
			result[i] = req // conversion from ExecutionRequest to []byte
		}

		var requests *ctypes.ExecutionRequests
		if requests, err = ctypes.DecodeExecutionRequests(result); err != nil {
			return err
		}
		if err = body.SetExecutionRequests(requests); err != nil {
			return err
		}
	}

	return nil
}

// computeAndSetStateRoot computes the state root of an outgoing block
// and sets it in the block.
func (s *Service) computeAndSetStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	stateRoot, err := s.computeStateRoot(
		ctx,
		proposerAddress,
		consensusTime,
		st,
		blk,
	)
	if err != nil {
		s.logger.Error(
			"failed to compute state root while building block ❗️ ",
			"slot", blk.GetSlot().Base10(),
			"error", err,
		)
		return err
	}
	blk.SetStateRoot(stateRoot)
	return nil
}

// computeStateRoot computes the state root of an outgoing block.
func (s *Service) computeStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
) (common.Root, error) {
	startTime := time.Now()
	defer s.metrics.measureStateRootComputationTime(startTime)

	// TODO: Think about how this would affect the proposer when
	// the payload in their block has come from a remote builder.
	txCtx := transition.NewTransitionCtx(
		ctx,
		consensusTime,
		proposerAddress,
	).
		WithVerifyPayload(false).
		WithVerifyRandao(false).
		WithVerifyResult(false).
		WithMeterGas(false)

	if _, err := s.stateProcessor.Transition(txCtx, st, blk); err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot(), nil
}
