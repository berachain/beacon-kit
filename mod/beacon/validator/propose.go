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
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"golang.org/x/sync/errgroup"
)

// RequestBlockForProposal builds a new beacon block.
//
//nolint:funlen // todo:fix.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) RequestBlockForProposal(
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
	s.logger.Info(
		"requesting beacon block assembly 🙈",
		"slot", requestedSlot.Base10(),
	)

	// The goal here is to acquire a payload whose parent is the previously
	// finalized block, such that, if this payload is accepted, it will be
	// the next finalized block in the chain. A byproduct of this design
	// is that we get the nice property of lazily propagating the finalized
	// and safe block hashes to the execution client.
	st := s.bsb.StateFromContext(ctx)

	// Prepare the state such that it is ready to build a block for
	// the requested slot
	if _, err := s.stateProcessor.ProcessSlots(st, requestedSlot); err != nil {
		return blk, sidecars, err
	}

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

	// Assemble a new block with the payload.
	body := blk.GetBody()
	if body.IsNil() {
		return blk, sidecars, ErrNilBlkBody
	}

	// Set the reveal on the block body.
	body.SetRandaoReveal(reveal)

	// Get the payload for the block.
	envelope, err := s.retrieveExecutionPayload(ctx, st, blk)
	if err != nil {
		return blk, sidecars, err
	} else if envelope == nil {
		return blk, sidecars, ErrNilPayload
	}

	// If we get returned a nil blobs bundle, we should return an error.
	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle == nil {
		return blk, sidecars, ErrNilBlobsBundle
	}

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	depositIndex, err := st.GetEth1DepositIndex()
	if err != nil {
		return blk, sidecars, ErrNilDepositIndexStart
	}

	// Dequeue deposits from the state.
	deposits, err := s.bsb.DepositStore(ctx).GetDepositsByIndex(
		depositIndex,
		s.chainSpec.MaxDepositsPerBlock(),
	)
	if err != nil {
		return blk, sidecars, err
	}

	// Set the deposits on the block body.
	body.SetDeposits(deposits)

	// Set the KZG commitments on the block body.
	body.SetBlobKzgCommitments(blobsBundle.GetCommitments())

	// TODO: assemble real eth1data.
	body.SetEth1Data(&types.Eth1Data{
		DepositRoot:  primitives.Bytes32{},
		DepositCount: 0,
		BlockHash:    common.ZeroHash,
	})

	// Set the execution data.
	if err = body.SetExecutionData(
		envelope.GetExecutionPayload(),
	); err != nil {
		return blk, sidecars, err
	}

	// Produce block sidecars.
	g.Go(func() error {
		var sidecarErr error
		sidecars, sidecarErr = s.blobFactory.BuildSidecars(
			blk,
			envelope.GetBlobsBundle(),
		)
		return sidecarErr
	})

	g.Go(func() error {
		return s.computeAndSetStateRoot(ctx, st, blk)
	})

	if err = g.Wait(); err != nil {
		return blk, sidecars, err
	}

	s.logger.Info(
		"beacon block successfully built 🛠️ ",
		"slot", requestedSlot.Base10(),
		"state_root", blk.GetStateRoot(),
		"duration", time.Since(startTime).String(),
	)

	return blk, sidecars, nil
}

// GetEmptyBlock creates a new empty block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) getEmptyBeaconBlockForSlot(
	st BeaconStateT, requestedSlot math.Slot,
) (BeaconBlockT, error) {
	var blk BeaconBlockT
	// Create a new block.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		uint64(requestedSlot-1) % s.chainSpec.SlotsPerHistoricalRoot(),
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
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) buildRandaoReveal(
	st BeaconStateT,
	slot math.Slot,
) (crypto.BLSSignature, error) {
	var forkData ForkDataT
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	epoch := s.chainSpec.SlotToEpoch(slot)
	signingRoot, err := forkData.New(
		version.FromUint32[primitives.Version](
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
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error) {
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
		var lph *types.ExecutionPayloadHeader
		lph, err = st.GetLatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}

		// If we failed to retrieve the payload, request a synchrnous payload.
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
				uint64((lph.GetTimestamp()+1)),
			),
			blk.GetParentBlockRoot(),
			lph.GetBlockHash(),
			lph.GetParentHash(),
		)
	}
	return envelope, nil
}
