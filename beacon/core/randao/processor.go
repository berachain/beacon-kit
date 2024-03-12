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
	"context"
	"fmt"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/beacon/signing"
	crypto "github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// Processor is the randao processor.
type Processor struct {
	BeaconStateProvider
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
//	signing_root = compute_signing_root(compute_epoch_at_slot(block.slot),
//
// domain)
//
//	return bls.Sign(privkey, signing_root)
func (p *Processor) BuildReveal(
	_ context.Context,
	epoch primitives.Epoch,
) (types.Reveal, error) {
	signingRoot := p.GetSigningRoot(epoch)

	return p.signer.Sign(signingRoot), nil
}

// VerifyReveal verifies the reveal of the proposer.
func (p *Processor) VerifyReveal(
	proposerPubkey [bls12381.PubKeyLength]byte,
	msg []byte,
	reveal types.Reveal,
) bool {
	return bls12381.VerifySignature(proposerPubkey, msg, reveal)
}

// MixinNewReveal mixes in a new reveal.
func (p *Processor) MixinNewReveal(
	ctx context.Context,
	blk beacontypes.BeaconBlock,
) error {
	st := p.BeaconState(ctx)
	randaoMix, err := st.RandaoMix()
	if err != nil {
		return fmt.Errorf("failed to get randao mix: %w", err)
	}
	reveal := blk.GetRandaoReveal()

	newMix := randaoMix.MixinNewReveal(reveal)
	if err = st.SetRandaoMix(newMix); err != nil {
		return fmt.Errorf("failed to set new randao mix: %w", err)
	}
	p.logger.Info("updated randao mix ", "new_mix", newMix)
	return nil
}

// GetSigningRoot returns the signing root.
// // TODO: COMPLETE
func (p *Processor) GetSigningRoot(
	epoch primitives.Epoch,
) []byte {
	return p.computeSigningRoot(epoch, p.getDomain(epoch))
}

// computeSigningRoot computes the signing root.
// // TODO: COMPLETE
func (p *Processor) computeSigningRoot(
	epoch primitives.Epoch,
	_ signing.Domain,
) []byte {
	return sdktypes.Uint64ToBigEndian(epoch)
}

// getDomain returns the domain.
// TODO: COMPLETE
func (p *Processor) getDomain(
	_ primitives.Epoch,
) signing.Domain {
	return signing.Domain{}
}
