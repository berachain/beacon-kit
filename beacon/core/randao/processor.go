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
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	"github.com/berachain/beacon-kit/beacon/core/signing"
	"github.com/berachain/beacon-kit/beacon/core/state"
	"github.com/berachain/beacon-kit/config"
	crypto "github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
)

// Processor is the randao processor.
type Processor struct {
	cfg    *config.Config
	signer crypto.Signer[[bls12381.SignatureLength]byte]
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
) (types.Reveal, error) {
	return p.buildReveal(
		st.GetChainID(),
		p.cfg.Beacon.SlotToEpoch(st.GetSlot()),
	)
}

// buildReveal creates a reveal for the proposer.
func (p *Processor) buildReveal(
	chainID string,
	epoch primitives.Epoch,
) (types.Reveal, error) {
	signingRoot, err := p.computeSigningRoot(chainID, epoch)
	if err != nil {
		return types.Reveal{}, err
	}
	return p.signer.Sign(signingRoot[:]), nil
}

// VerifyReveal verifies the reveal of the proposer.
func (p *Processor) VerifyReveal(
	st state.BeaconState,
	proposerPubkey [bls12381.PubKeyLength]byte,
	reveal types.Reveal,
) error {
	return p.verifyReveal(
		proposerPubkey,
		st.GetChainID(),
		p.cfg.Beacon.SlotToEpoch(st.GetSlot()),
		reveal,
	)
}

// VerifyReveal verifies the reveal of the proposer.
func (p *Processor) verifyReveal(
	proposerPubkey [bls12381.PubKeyLength]byte,
	chainID string,
	epoch primitives.Epoch,
	reveal types.Reveal,
) error {
	signingRoot, err := p.computeSigningRoot(chainID, epoch)
	if err != nil {
		return err
	}
	if ok := bls12381.VerifySignature(
		proposerPubkey,
		signingRoot[:],
		reveal,
	); !ok {
		return ErrInvalidSignature
	}

	p.logger.Info("randao reveal successfully verified 🐸",
		"reveal", reveal,
	)
	return nil
}

// MixinNewReveal mixes in a new reveal.
func (p *Processor) MixinNewReveal(
	st state.BeaconState,
	reveal types.Reveal,
) error {
	// Get last slots randao mix.
	mix, err := st.RandaoMixAtIndex(
		st.GetSlot() % p.cfg.Beacon.Limits.EpochsPerHistoricalVector,
	)
	if err != nil {
		return fmt.Errorf("failed to get randao mix: %w", err)
	}

	// Mix in the reveal with the previous slotsmix.
	newMix := mix.MixinNewReveal(reveal)

	// Set this slots mix to the new mix.
	if err = st.UpdateRandaoMixAtIndex(
		st.GetSlot()%p.cfg.Beacon.Limits.EpochsPerHistoricalVector,
		newMix,
	); err != nil {
		return fmt.Errorf("failed to set new randao mix: %w", err)
	}
	p.logger.Info("randao mix updated 🎲", "new_mix", newMix)
	return nil
}

func (p *Processor) computeSigningRoot(
	chainID string,
	epoch primitives.Epoch,
) (primitives.HashRoot, error) {
	signingDomain, err := signing.GetDomain(
		p.cfg,
		chainID,
		signing.DomainRandao,
		epoch,
	)
	if err != nil {
		return primitives.HashRoot{}, fmt.Errorf(
			"failed to get domain: %w",
			err,
		)
	}
	signingRoot, err := signing.ComputeSigningRoot(
		primitives.SSZEpoch(epoch),
		signingDomain,
	)
	if err != nil {
		return primitives.HashRoot{},
			fmt.Errorf("failed to compute signing root: %w", err)
	}
	return signingRoot, nil
}
