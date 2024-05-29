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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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

	// Create a new empty block from the current state.
	return types.EmptyBeaconBlock[BeaconBlockT](
		slot,
		proposerIndex,
		parentBlockRoot,
		s.chainSpec.ActiveForkVersionForSlot(slot),
	)
}

func (s *Service[
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) retrievePayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// The latest execution payload header, will be from the previous block
	// during the block building phase.
	parentExecutionPayload, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	// Get the payload for the block.
	envelope, err := s.localPayloadBuilder.
		RetrieveOrBuildPayload(
			ctx,
			st,
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
			parentExecutionPayload.GetBlockHash(),
		)
	if err != nil {
		return nil, err
	} else if envelope == nil {
		return nil, ErrNilPayload
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
