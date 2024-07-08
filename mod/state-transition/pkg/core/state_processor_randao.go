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

package core

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto/sha256"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/go-faster/xor"
)

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT,
	_, _, _, _, _, _, ForkDataT, _, _, _, _,
]) processRandaoReveal(
	st BeaconStateT,
	blk BeaconBlockT,
	skipVerification bool,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Ensure the proposer index is valid.
	proposer, err := st.ValidatorByIndex(blk.GetProposerIndex())
	if err != nil {
		return err
	}

	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}

	epoch := sp.cs.SlotToEpoch(slot)
	body := blk.GetBody()

	var fd ForkDataT
	fd = fd.New(
		version.FromUint32[common.Version](
			sp.cs.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	if !skipVerification {
		var signingRoot common.Root
		signingRoot, err = fd.ComputeRandaoSigningRoot(
			sp.cs.DomainTypeRandao(), epoch)
		if err != nil {
			return err
		}

		reveal := body.GetRandaoReveal()
		if err = sp.signer.VerifySignature(
			proposer.GetPubkey(),
			signingRoot[:],
			reveal,
		); err != nil {
			return err
		}
	}

	prevMix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % sp.cs.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}

	mix, err := sp.buildRandaoMix(prevMix, body.GetRandaoReveal())
	if err != nil {
		return err
	}

	return st.UpdateRandaoMixAtIndex(
		uint64(epoch)%sp.cs.EpochsPerHistoricalVector(), mix,
	)
}

// processRandaoMixesReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#randao-mixes-updates
//
//nolint:lll
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processRandaoMixesReset(
	st BeaconStateT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	epoch := sp.cs.SlotToEpoch(slot)
	mix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % sp.cs.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch+1)%sp.cs.EpochsPerHistoricalVector(),
		mix,
	)
}

// buildRandaoMix as defined in the Ethereum 2.0 specification.
func (sp *StateProcessor[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) buildRandaoMix(
	mix common.Bytes32,
	reveal crypto.BLSSignature,
) (common.Bytes32, error) {
	newMix := make([]byte, constants.RootLength)
	revealHash := sha256.Hash(reveal[:])
	// Apparently this library giga fast? Good project? lmeow.
	if numXor := xor.Bytes(
		newMix, mix[:], revealHash[:],
	); numXor != constants.RootLength {
		return common.Bytes32{}, ErrXorInvalid
	}
	return common.Bytes32(newMix), nil
}
