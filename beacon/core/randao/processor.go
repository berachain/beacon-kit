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
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	crypto "github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// Processor is the randao processor.
type Processor struct {
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
	epoch primitives.Epoch,
) (types.Reveal, error) {
	return p.signer.Sign(
		p.computeSigningRoot(epoch, p.getDomain(epoch)),
	), nil
}

// VerifyReveal verifies the reveal of the proposer.
func (p *Processor) VerifyReveal(
	proposerPubkey [bls12381.PubKeyLength]byte,
	epoch primitives.Epoch,
	reveal types.Reveal,
) error {
	if ok := bls12381.VerifySignature(
		proposerPubkey,
		p.computeSigningRoot(epoch, p.getDomain(epoch)),
		reveal,
	); !ok {
		return ErrInvalidSignature
	}

	p.logger.Info("randao reveal successfully verified 🤫 ",
		"reveal", reveal,
	)
	return nil
}

// MixinNewReveal mixes in a new reveal.
func (p *Processor) MixinNewReveal(
	st state.BeaconState,
	blk beacontypes.BeaconBlock,
) error {
	mix, err := st.RandaoMix()
	if err != nil {
		return fmt.Errorf("failed to get randao mix: %w", err)
	}

	newMix := mix.MixinNewReveal(blk.GetRandaoReveal())
	if err = st.SetRandaoMix(newMix); err != nil {
		return fmt.Errorf("failed to set new randao mix: %w", err)
	}
	p.logger.Info("randao mix updated 🎲", "new_mix", newMix)
	return nil
}

// computeSigningRoot computes the signing root.
// // TODO: COMPLETE THIS IS CURRENTLY WRONG.
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
