// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package validator

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// GetEmptyBlock creates a new empty block.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) GetEmptyBeaconBlock(
	st BeaconStateT, slot math.Slot,
) (BeaconBlockT, error) {
	var blk BeaconBlockT
	// Create a new block.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		uint64(slot) % s.chainSpec.SlotsPerHistoricalRoot(),
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
		slot,
		proposerIndex,
		parentBlockRoot,
		s.chainSpec.ActiveForkVersionForSlot(slot),
	)
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error) {
	// Get the payload for the block.
	envelope, err := s.localPayloadBuilder.
		RetrievePayload(
			ctx,
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
		)
	if err != nil {
		s.metrics.failedToRetrieveOptimisticPayload(
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
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

// prepareStateForBuilding ensures that the state is at the requested slot
// before building a block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) prepareStateForBuilding(
	st BeaconStateT,
	requestedSlot math.Slot,
) error {
	// Get the current state slot.
	stateSlot, err := st.GetSlot()
	if err != nil {
		return err
	}

	slotDifference := requestedSlot - stateSlot
	switch {
	case slotDifference == 1:
		// If our BeaconState is not up to date, we need to process
		// a slot to get it up to date.
		if _, err = s.stateProcessor.ProcessSlot(st); err != nil {
			return err
		}

		// Request the slot again, it should've been incremented by 1.
		stateSlot, err = st.GetSlot()
		if err != nil {
			return err
		}

		// If after doing so, we aren't exactly at the requested slot,
		// we should return an error.
		if requestedSlot != stateSlot {
			return errors.Newf(
				"requested slot could not be processed, requested: %d, state: %d",
				requestedSlot,
				stateSlot,
			)
		}
	case slotDifference > 1:
		return errors.Newf(
			"requested slot is too far ahead, requested: %d, state: %d",
			requestedSlot,
			stateSlot,
		)
	case slotDifference < 1:
		return errors.Newf(
			"requested slot is in the past, requested: %d, state: %d",
			requestedSlot,
			stateSlot,
		)
	}

	return nil
}

// buildRandaoReveal builds a randao reveal for the given slot.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT, BlobSidecarsT,
]) buildRandaoReveal(
	st BeaconStateT,
	slot math.Slot,
) (crypto.BLSSignature, error) {
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	epoch := s.chainSpec.SlotToEpoch(slot)
	signingRoot, err := types.NewForkData(
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
