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
	"fmt"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/state"
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/forks"
	"github.com/berachain/beacon-kit/mod/forks/version"
	"github.com/berachain/beacon-kit/mod/node-builder/config"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/go-faster/xor"
	blst "github.com/itsdevbear/comet-bls12-381/bls/blst"
	sha256 "github.com/minio/sha256-simd"
)

// Processor is the randao processor.
type Processor struct {
	cfg    *config.Config
	signer core.BLSSigner
	logger log.Logger
}

// NewProcessor creates a new randao processor.
func NewProcessor(
	opts ...Option,
) *Processor {
	p := &Processor{}
	for _, opt := range opts {
		if err := opt(p); err != nil {
			panic(err)
		}
	}
	return p
}

// ProcessRandao processes the randao reveal.
// process_randao in the Ethereum 2.0 specification.
func (p *Processor) ProcessRandao(
	st state.BeaconState,
	blk beacontypes.BeaconBlock,
) error {
	// proposer := blk.
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

	epoch := p.cfg.Beacon.SlotToEpoch(slot)
	signingRoot, err := p.computeSigningRoot(epoch, root)
	if err != nil {
		return err
	}

	reveal := blk.GetBody().GetRandaoReveal()
	if !blst.VerifySignaturePubkeyBytes(
		proposer.Pubkey[:],
		signingRoot[:],
		reveal[:],
	) {
		return ErrInvalidSignature
	}

	prevMix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % p.cfg.Beacon.EpochsPerHistoricalVector,
	)
	if err != nil {
		return err
	}

	mix := p.buildMix(prevMix, blk.GetBody().GetRandaoReveal())
	p.logger.Info("randao mix updated 🎲", "new_mix", mix)
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch)%p.cfg.Beacon.EpochsPerHistoricalVector,
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
func (p *Processor) BuildReveal(
	st state.BeaconState,
) (primitives.BLSSignature, error) {
	genesisValidatorsRoot, err := st.GetGenesisValidatorsRoot()
	if err != nil {
		return primitives.BLSSignature{}, err
	}

	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return primitives.BLSSignature{}, err
	}

	return p.buildReveal(
		genesisValidatorsRoot,
		p.cfg.Beacon.SlotToEpoch(slot),
	)
}

// ProcessRandaoMixesReset resets the randao mixes.
// process_randao_mixes_reset in the Ethereum 2.0 specification.
func (p *Processor) ProcessRandaoMixesReset(st state.BeaconState) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	epoch := p.cfg.Beacon.SlotToEpoch(slot)
	mix, err := st.GetRandaoMixAtIndex(
		uint64(epoch) % p.cfg.Beacon.EpochsPerHistoricalVector,
	)
	if err != nil {
		return err
	}
	return st.UpdateRandaoMixAtIndex(
		uint64(epoch+1)%p.cfg.Beacon.EpochsPerHistoricalVector,
		mix,
	)
}

// buildReveal creates a reveal for the proposer.
func (p *Processor) buildReveal(
	genesisValidatorsRoot primitives.Root,
	epoch primitives.Epoch,
) (primitives.BLSSignature, error) {
	signingRoot, err := p.computeSigningRoot(epoch, genesisValidatorsRoot)
	if err != nil {
		return primitives.BLSSignature{}, err
	}
	return p.signer.Sign(signingRoot[:])
}

// buildMix builds a new mix from a given mix and reveal.
func (p *Processor) buildMix(
	mix primitives.Bytes32,
	reveal primitives.BLSSignature,
) primitives.Bytes32 {
	newMix := make([]byte, constants.RootLength)
	revealHash := sha256.Sum256(reveal[:])
	// Apparently this library giga fast? Good project? lmeow.
	_ = xor.Bytes(newMix, mix[:], revealHash[:])
	return primitives.Bytes32(newMix)
}

// computeSigningRoot computes the signing root for the epoch.
func (p *Processor) computeSigningRoot(
	epoch primitives.Epoch,
	genesisValidatorsRoot primitives.Root,
) (primitives.Root, error) {
	fd := forks.NewForkData(
		version.FromUint32(
			p.cfg.Beacon.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	signingDomain, err := fd.ComputeDomain(primitives.DomainTypeRandao)
	if err != nil {
		return primitives.Root{}, err
	}

	signingRoot, err := primitives.ComputeSigningRootUInt64(
		uint64(epoch),
		signingDomain,
	)

	if err != nil {
		return primitives.Root{},
			fmt.Errorf("failed to compute signing root: %w", err)
	}
	return signingRoot, nil
}
