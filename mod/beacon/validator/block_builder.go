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

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// buildBlockAndSidecars builds a new beacon block.
func (s *Service[
	BeaconBlockT, _, _, BlobSidecarsT, _, _, _, _, _, _,
]) buildBlockAndSidecars(
	ctx context.Context,
	requestedSlot math.Slot,
) (BeaconBlockT, BlobSidecarsT, error) {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		startTime = time.Now()
		g, _      = errgroup.WithContext(ctx)
	)

	defer s.metrics.measureRequestBlockForProposalTime(startTime)

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propagating the finalized
	// and safe block hashes to the execution client.
	st := s.bsb.BeaconState()
	fmt.Println("BEFORE PROCESSING SLOT!!")
	// Prepare the state such that it is ready to build a block for
	// the requested slot
	if _, err := s.stateProcessor.ProcessSlots(st, requestedSlot); err != nil {
		return blk, sidecars, err
	}
	fmt.Println("AFTER PROCESSING SLOT!!")

	// Build the reveal for the current slot.
	// TODO: We can optimize to pre-compute this in parallel?
	reveal, err := s.buildRandaoReveal(st, requestedSlot)
	if err != nil {
		return blk, sidecars, err
	}

	// Create a new empty block from the current state.
	blk, err = s.getEmptyBeaconBlockForSlot(
		st, requestedSlot,
	)
	if err != nil {
		return blk, sidecars, err
	}

	// Get the payload for the block.
	envelope, err := s.retrieveExecutionPayload(ctx, st, blk)
	if err != nil {
		return blk, sidecars, err
	} else if envelope == nil {
		return blk, sidecars, ErrNilPayload
	}

	// We have to assemble the block body prior to producing the sidecars
	// since we need to generate the inclusion proofs.
	if err = s.buildBlockBody(st, blk, reveal, envelope); err != nil {
		return blk, sidecars, err
	}

	// Produce blob sidecars, we produce them in parallel to computing the state
	// root as an optimization.
	//
	// TODO: Figure out a clean way to break "BlockAndSidecars" into two
	// functions
	// without giving up the parallelization benefits.
	g.Go(func() error {
		sidecars, err = s.blobFactory.BuildSidecars(
			blk, envelope.GetBlobsBundle(),
		)
		return err
	})

	// Compute the state root for the block.
	g.Go(func() error {
		return s.computeAndSetStateRoot(ctx, st, blk)
	})

	// Wait for all the goroutines to finish.
	if err = g.Wait(); err != nil {
		return blk, sidecars, err
	}

	s.logger.Info(
		"Beacon block successfully built",
		"slot", requestedSlot.Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	return blk, sidecars, nil
}

// getEmptyBeaconBlockForSlot creates a new empty block.
func (s *Service[
	BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _,
]) getEmptyBeaconBlockForSlot(
	st BeaconStateT, requestedSlot math.Slot,
) (BeaconBlockT, error) {
	var blk BeaconBlockT
	// Create a new block.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		(requestedSlot.Unwrap() - 1) % s.chainSpec.SlotsPerHistoricalRoot(),
	)

	if err != nil {
		return blk, errors.Newf(
			"failed to get block root at index: %w",
			err,
		)
	}

	// Get the proposer index for the slot.
	proposerIndex, err := st.ValidatorIndexByPubkey(
		s.signer.PublicKey(),
	)
	if err != nil {
		return blk, errors.Newf(
			"failed to get validator by pubkey: %w",
			err,
		)
	}

	return blk.NewWithVersion(
		requestedSlot,
		proposerIndex,
		parentBlockRoot,
		s.chainSpec.ActiveForkVersionForSlot(requestedSlot),
	)
}

// buildRandaoReveal builds a randao reveal for the given slot.
func (s *Service[
	_, _, BeaconStateT, _, _, _, _, _, _, ForkDataT,
]) buildRandaoReveal(
	st BeaconStateT,
	slot math.Slot,
) (crypto.BLSSignature, error) {
	var (
		forkData ForkDataT
		epoch    = s.chainSpec.SlotToEpoch(slot)
	)

	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	signingRoot, err := forkData.New(
		version.FromUint32[common.Version](
			s.chainSpec.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	).ComputeRandaoSigningRoot(
		s.chainSpec.DomainTypeRandao(),
		epoch,
	)

	if err != nil {
		return crypto.BLSSignature{}, err
	}
	return s.signer.Sign(signingRoot[:])
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	BeaconBlockT, _, BeaconStateT, _, _, _, _,
	ExecutionPayloadT, ExecutionPayloadHeaderT, _,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
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
	if err != nil {
		s.metrics.failedToRetrievePayload(
			blk.GetSlot(),
			err,
		)

		// The latest execution payload header will be from the previous block
		// during the block building phase.
		var lph ExecutionPayloadHeaderT
		lph, err = st.GetLatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}

		// If we failed to retrieve the payload, request a synchronous payload.
		//
		// NOTE: The state here is properly configured by the
		// prepareStateForBuilding
		//
		// call that needs to be called before requesting the Payload.
		// TODO: We should decouple the PayloadBuilder from BeaconState to make
		// this less confusing.
		return s.localPayloadBuilder.RequestPayloadSync(
			ctx,
			st,
			blk.GetSlot(),
			// TODO: this is hood.
			max(
				//#nosec:G701
				uint64(time.Now().Unix()+1),
				lph.GetTimestamp().Unwrap()+1,
			),
			blk.GetParentBlockRoot(),
			lph.GetBlockHash(),
			lph.GetParentHash(),
		)
	}
	return envelope, nil
}

// BuildBlockBody assembles the block body with necessary components.
func (s *Service[
	BeaconBlockT, _, BeaconStateT, _,
	_, _, Eth1DataT, ExecutionPayloadT, _, _,
]) buildBlockBody(
	st BeaconStateT,
	blk BeaconBlockT,
	reveal crypto.BLSSignature,
	envelope engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT],
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

	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return ErrNilDepositIndexStart
	}

	// Dequeue deposits from the state.
	deposits, err := s.bsb.DepositStore().GetDepositsByIndex(
		depositIndex,
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return err
	}

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	var eth1Data Eth1DataT
	// TODO: assemble real eth1data.
	body.SetEth1Data(eth1Data.New(
		common.Bytes32{},
		0,
		gethprimitives.ZeroHash,
	))

	// Set the graffiti on the block body.
	body.SetGraffiti(bytes.ToBytes32([]byte(s.cfg.Graffiti)))

	return body.SetExecutionData(envelope.GetExecutionPayload())
}

// computeAndSetStateRoot computes the state root of an outgoing block
// and sets it in the block.
func (s *Service[
	BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _,
]) computeAndSetStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	stateRoot, err := s.computeStateRoot(ctx, st, blk)
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
	BeaconBlockT, _, BeaconStateT, _, _, _, _, _, _, _,
]) computeStateRoot(
	ctx context.Context,
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
		},
		st, blk,
	); err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot()
}
