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

package core

import (
	"crypto/sha256"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/go-faster/xor"
)

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, ContextT,
	DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	ForkT, ForkDataT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
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
		version.FromUint32[primitives.Version](
			sp.cs.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	if !skipVerification {
		var signingRoot primitives.Root
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

	return st.UpdateRandaoMixAtIndex(
		uint64(epoch)%sp.cs.EpochsPerHistoricalVector(),
		sp.buildRandaoMix(prevMix, body.GetRandaoReveal()),
	)
}

// processRandaoMixesReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#randao-mixes-updates
//
//nolint:lll
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, ContextT,
	DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	ForkT, ForkDataT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
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
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, ContextT,
	DepositT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	ForkT, ForkDataT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
]) buildRandaoMix(
	mix primitives.Bytes32,
	reveal crypto.BLSSignature,
) primitives.Bytes32 {
	newMix := make([]byte, constants.RootLength)
	revealHash := sha256.Sum256(reveal[:])
	// Apparently this library giga fast? Good project? lmeow.
	_ = xor.Bytes(newMix, mix[:], revealHash[:])
	return primitives.Bytes32(newMix)
}
