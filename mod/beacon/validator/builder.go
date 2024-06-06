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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// GetEmptyBlock creates a new empty block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) getEmptyBeaconBlock(
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
