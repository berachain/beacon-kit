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

package validator

import (
	"context"
	"fmt"
	"time"

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/types"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
)

// BuildBlockAndSidecars builds a new beacon block.
func (s *Service[
	BeaconBlockT, _, _, BlobSidecarsT,
	_, SlashingInfoT, SlotDataT,
]) BuildBlockAndSidecars(
	ctx context.Context,
	slotData types.SlotData[ctypes.SlashingInfo],
) ([]byte, []byte, error) {
	var (
		blk      BeaconBlockT
		sidecars BlobSidecarsT
		forkData *ctypes.ForkData
	)

	startTime := time.Now()
	defer s.metrics.measureRequestBlockForProposalTime(startTime)

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propagating the finalized
	// and safe block hashes to the execution client.
	st := s.sb.StateFromContext(ctx)

	// Prepare the state such that it is ready to build a block for
	// the requested slot
	if _, err := s.stateProcessor.ProcessSlots(
		st,
		slotData.GetSlot(),
	); err != nil {
		return nil, nil, err
	}

	// Build forkdata used for the signing root of the reveal and the sidecars
	forkData, err := s.buildForkData(st, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	// Build the reveal for the current slot.
	// TODO: We can optimize to pre-compute this in parallel?
	reveal, err := s.buildRandaoReveal(forkData, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	// Create a new empty block from the current state.
	blk, err = s.getEmptyBeaconBlockForSlot(st, slotData.GetSlot())
	if err != nil {
		return nil, nil, err
	}

	// Get the payload for the block.
	envelope, err := s.retrieveExecutionPayload(ctx, st, blk, slotData)
	if err != nil {
		return nil, nil, err
	}
	if envelope == nil {
		return nil, nil, ErrNilPayload
	}

	// We have to assemble the block body prior to producing the sidecars
	// since we need to generate the inclusion proofs.
	if err = s.buildBlockBody(
		ctx, st, blk, reveal, envelope, slotData,
	); err != nil {
		return nil, nil, err
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

	// Produce blob sidecars with new StateRoot
	sidecars, err = s.blobFactory.BuildSidecars(
		blk,
		envelope.GetBlobsBundle(),
		s.signer,
		forkData,
	)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info(
		"Beacon block successfully built",
		"slot", slotData.GetSlot().Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	blkBytes, bbErr := blk.MarshalSSZ()
	if bbErr != nil {
		return nil, nil, bbErr
	}
	sidecarsBytes, scErr := sidecars.MarshalSSZ()
	if scErr != nil {
		return nil, nil, scErr
	}

	return blkBytes, sidecarsBytes, nil
}

// getEmptyBeaconBlockForSlot creates a new empty block.
func (s *Service[
	BeaconBlockT, BeaconStateT, _, _, _, _, _,
]) getEmptyBeaconBlockForSlot(
	st BeaconStateT, requestedSlot math.Slot,
) (BeaconBlockT, error) {
	var blk BeaconBlockT
	// Create a new block.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		(requestedSlot.Unwrap() - 1) % s.chainSpec.SlotsPerHistoricalRoot(),
	)

	if err != nil {
		return blk, err
	}

	// Get the proposer index for the slot.
	proposerIndex, err := st.ValidatorIndexByPubkey(
		s.signer.PublicKey(),
	)
	if err != nil {
		return blk, err
	}

	return blk.NewWithVersion(
		requestedSlot,
		proposerIndex,
		parentBlockRoot,
		s.chainSpec.ActiveForkVersionForSlot(requestedSlot),
	)
}

func (s *Service[
	_, BeaconStateT, _, _, _, _, _,
]) buildForkData(
	st BeaconStateT,
	slot math.Slot,
) (*ctypes.ForkData, error) {
	var (
		epoch = s.chainSpec.SlotToEpoch(slot)
	)

	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	return ctypes.NewForkData(
		version.FromUint32[common.Version](
			s.chainSpec.ActiveForkVersionForEpoch(epoch),
		),
		genesisValidatorsRoot,
	), nil
}

// buildRandaoReveal builds a randao reveal for the given slot.
func (s *Service[
	_, BeaconStateT, _, _, _, _, _,
]) buildRandaoReveal(
	forkData *ctypes.ForkData,
	slot math.Slot,
) (crypto.BLSSignature, error) {
	var epoch = s.chainSpec.SlotToEpoch(slot)
	signingRoot := forkData.ComputeRandaoSigningRoot(
		s.chainSpec.DomainTypeRandao(),
		epoch,
	)
	return s.signer.Sign(signingRoot[:])
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	BeaconBlockT, BeaconStateT, _, _, _,
	SlashingInfoT, SlotDataT,
]) retrieveExecutionPayload(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	slotData types.SlotData[ctypes.SlashingInfo],
) (ctypes.BuiltExecutionPayloadEnv, error) {
	//
	// TODO: Add external block builders to this flow.
	//
	// Get the payload for the block.
	envelope, err := s.localPayloadBuilder.
		RetrievePayload(
			ctx,
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
		)
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

	s.metrics.failedToRetrievePayload(
		blk.GetSlot(),
		err,
	)

	// The latest execution payload header will be from the previous block
	// during the block building phase.
	var lph *ctypes.ExecutionPayloadHeader
	lph, err = st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	return s.localPayloadBuilder.RequestPayloadSync(
		ctx,
		st,
		blk.GetSlot(),
		payloadtime.Next(
			slotData.GetConsensusTime(),
			lph.GetTimestamp(),
			false, // buildOptimistically
		).Unwrap(),
		blk.GetParentBlockRoot(),
		lph.GetBlockHash(),
		lph.GetParentHash(),
	)
}

// BuildBlockBody assembles the block body with necessary components.
func (s *Service[
	BeaconBlockT, BeaconStateT, _, _, _,
	SlashingInfoT, SlotDataT,
]) buildBlockBody(
	_ context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	reveal crypto.BLSSignature,
	envelope ctypes.BuiltExecutionPayloadEnv,
	slotData types.SlotData[ctypes.SlashingInfo],
) error {
	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body.IsNil() {
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
		return ErrNilDepositIndexStart
	}

	// Bartio and Boonet pre Fork2 have deposit broken and undervalidated
	// Any other network should build deposits the right way
	if !(s.chainSpec.DepositEth1ChainID() == spec.BartioChainID ||
		(s.chainSpec.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
			blk.GetSlot() < math.U64(spec.BoonetFork2Height))) {
		depositIndex++
	}
	deposits, err := s.sb.DepositStore().GetDepositsByIndex(
		depositIndex,
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}

	// Set the deposits on the block body.
	s.logger.Info(
		"Building block body with local deposits",
		"start_index", depositIndex, "num_deposits", len(deposits),
	)
	body.SetDeposits(deposits)

	var eth1Data *ctypes.Eth1Data
	// TODO: assemble real eth1data.
	body.SetEth1Data(eth1Data.New(
		common.Root{},
		0,
		common.ExecutionHash{},
	))

	// Set the graffiti on the block body.
	sizedGraffiti := bytes.ExtendToSize([]byte(s.cfg.Graffiti), bytes.B32Size)
	graffiti, err := bytes.ToBytes32(sizedGraffiti)
	if err != nil {
		return fmt.Errorf("failed processing graffiti: %w", err)
	}
	body.SetGraffiti(graffiti)

	// Get the epoch to find the active fork version.
	epoch := s.chainSpec.SlotToEpoch(blk.GetSlot())
	activeForkVersion := s.chainSpec.ActiveForkVersionForEpoch(
		epoch,
	)
	if activeForkVersion >= version.DenebPlus {
		body.SetAttestations(slotData.GetAttestationData())

		// Set the slashing info on the block body.
		// TODO: Remove conversion once generics have been replaced with
		// concrete types.
		slashingInfo := slotData.GetSlashingInfo()
		body.SetSlashingInfo(convertSlashingInfo[SlashingInfoT](
			slashingInfo,
		))
	}

	body.SetExecutionPayload(envelope.GetExecutionPayload())
	return nil
}

// computeAndSetStateRoot computes the state root of an outgoing block
// and sets it in the block.
func (s *Service[
	BeaconBlockT, BeaconStateT, _, _, _, _, _,
]) computeAndSetStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st BeaconStateT,
	blk BeaconBlockT,
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
func (s *Service[
	BeaconBlockT, BeaconStateT, _, _, _, _, _,
]) computeStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	st BeaconStateT,
	blk BeaconBlockT,
) (common.Root, error) {
	startTime := time.Now()
	defer s.metrics.measureStateRootComputationTime(startTime)
	if _, err := s.stateProcessor.Transition(
		// TODO: We should think about how having optimistic
		// engine enabled here would affect the proposer when
		// the payload in their block has come from a remote builder.
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
	); err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot(), nil
}

func convertSlashingInfo[
	SlashingInfoT any,
](
	data []ctypes.SlashingInfo,
) []*ctypes.SlashingInfo {
	converted := make([]*ctypes.SlashingInfo, len(data))
	for i, d := range data {
		// #nosec G601 // TODO remove once we remove the SlashingInfoT generic type
		val, ok := any(&d).(*ctypes.SlashingInfo)
		if !ok {
			panic(fmt.Sprintf("failed to convert slashing info at index %d", i))
		}
		converted[i] = val
	}
	return converted
}
