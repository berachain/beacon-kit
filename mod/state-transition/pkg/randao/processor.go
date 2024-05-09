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

package randao

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/go-faster/xor"
	sha256 "github.com/minio/sha256-simd"
)

// Processor is the randao processor.
type Processor[
	BeaconBlockBodyT BeaconBlockBody,
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconStateT BeaconState,
] struct {
	chainSpec primitives.ChainSpec
	signer    crypto.BLSSigner
	logger    log.Logger[any]
}

// NewProcessor creates a new randao processor.
func NewProcessor[
	BeaconBlockBodyT BeaconBlockBody,
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconStateT BeaconState,
](
	chainSpec primitives.ChainSpec,
	signer crypto.BLSSigner,
	logger log.Logger[any],
) *Processor[BeaconBlockBodyT, BeaconBlockT, BeaconStateT] {
	return &Processor[BeaconBlockBodyT, BeaconBlockT, BeaconStateT]{
		chainSpec: chainSpec,
		signer:    signer,
		logger:    logger,
	}
}

// ProcessRandao processes the randao reveal.
// process_randao in the Ethereum 2.0 specification.
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) ProcessRandao(
	st BeaconStateT,
	blk BeaconBlockT,
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

	root, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}

	epoch := p.chainSpec.SlotToEpoch(slot)
	signingRoot, err := p.computeSigningRoot(epoch, root)
	if err != nil {
		return err
	}

	body := blk.GetBody()
	reveal := body.GetRandaoReveal()
	if err = p.signer.VerifySignature(
		proposer.Pubkey[:],
		signingRoot[:],
		reveal[:],
	); err != nil {
		return err
	}

	prevMix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % p.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}

	mix := p.buildMix(prevMix, body.GetRandaoReveal())
	p.logger.Info("randao mix updated 🎲", "new_mix", mix)
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch)%p.chainSpec.EpochsPerHistoricalVector(),
		mix,
	)
}

// BuildReveal creates a reveal for the proposer.
// def get_epoch_signature(state: BeaconState, block: BeaconBlock, privkey: int)
// -> BLSSignature:
//
//	domain = get_domain(state, DOMAIN_RANDAO, compute_epoch_at_slot(block.slot))
//	signing_root = compute_signing_root(
//						compute_epoch_at_slot(block.slot),
//						domain)
//
//	return bls.Sign(privkey, signing_root)
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) BuildReveal(
	st BeaconStateT,
) (crypto.BLSSignature, error) {
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return crypto.BLSSignature{}, err
	}

	return p.buildReveal(
		genesisValidatorsRoot,
		p.chainSpec.SlotToEpoch(slot),
	)
}

// ProcessRandaoMixesReset resets the randao mixes.
// process_randao_mixes_reset in the Ethereum 2.0 specification.
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) ProcessRandaoMixesReset(
	st BeaconStateT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	epoch := p.chainSpec.SlotToEpoch(slot)
	mix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % p.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return err
	}
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch+1)%p.chainSpec.EpochsPerHistoricalVector(),
		mix,
	)
}

// buildReveal creates a reveal for the proposer.
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) buildReveal(
	genesisValidatorsRoot primitives.Root,
	epoch math.Epoch,
) (crypto.BLSSignature, error) {
	signingRoot, err := p.computeSigningRoot(epoch, genesisValidatorsRoot)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	return p.signer.Sign(signingRoot[:])
}

// buildMix builds a new mix from a given mix and reveal.
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) buildMix(
	mix primitives.Bytes32,
	reveal crypto.BLSSignature,
) primitives.Bytes32 {
	newMix := make([]byte, constants.RootLength)
	revealHash := sha256.Sum256(reveal[:])
	// Apparently this library giga fast? Good project? lmeow.
	_ = xor.Bytes(newMix, mix[:], revealHash[:])
	return primitives.Bytes32(newMix)
}

// computeSigningRoot computes the signing root for the epoch.
func (p *Processor[
	BeaconBlockBodyT, BeaconBlockT, BeaconStateT,
]) computeSigningRoot(
	epoch math.Epoch,
	genesisValidatorsRoot primitives.Root,
) (primitives.Root, error) {
	fd := types.NewForkData(
		version.FromUint32[primitives.Version](
			p.chainSpec.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	signingDomain, err := fd.ComputeDomain(p.chainSpec.DomainTypeRandao())
	if err != nil {
		return primitives.Root{}, err
	}

	signingRoot, err := ssz.ComputeSigningRootUInt64(
		uint64(epoch),
		signingDomain,
	)

	if err != nil {
		return primitives.Root{},
			errors.Newf("failed to compute signing root: %w", err)
	}
	return signingRoot, nil
}
