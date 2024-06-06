// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// prepareStateForBuilding ensures that the state is at the requested slot
// before building a block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
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

// computeStateRoot computes the state root of an outgoing block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) computeStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) (primitives.Root, error) {
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
		return primitives.Root{}, err
	}

	return st.HashTreeRoot()
}

// verifyStateRoot verifies the state root of an incoming block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) verifyStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	startTime := time.Now()
	defer s.metrics.measureStateRootVerificationTime(startTime)
	if _, err := s.stateProcessor.Transition(
		// We run with a non-optimistic engine here to ensure
		// that the proposer does not try to push through a bad block.
		&transition.Context{
			Context:                 ctx,
			OptimisticEngine:        false,
			SkipPayloadVerification: false,
			SkipValidateResult:      false,
			SkipValidateRandao:      false,
		},
		st, blk,
	); errors.Is(err, engineerrors.ErrAcceptedPayloadStatus) {
		// It is safe for the validator to ignore this error since
		// the state transition will enforce that the block is part
		// of the canonical chain.
		//
		// TODO: this is only true because we are assuming SSF.
		return nil
	} else if err != nil {
		return err
	}

	return nil
}
