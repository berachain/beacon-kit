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
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	crypto "github.com/berachain/beacon-kit/crypto"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/tree"
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
//	signing_root = compute_signing_root(
//						compute_epoch_at_slot(block.slot),
//						domain)
//
//	return bls.Sign(privkey, signing_root)
func (p *Processor) BuildReveal(
	ctx context.Context,
	epoch primitives.Epoch,
) (types.Reveal, error) {
	return p.signer.Sign(p.computeSigningRoot(epoch, p.getDomain(ctx, epoch))), nil
}

// VerifyReveal verifies the reveal of the proposer.
func (p *Processor) VerifyReveal(
	ctx context.Context,
	proposerPubkey [bls12381.PubKeyLength]byte,
	epoch primitives.Epoch,
	reveal types.Reveal,
) bool {
	ok := bls12381.VerifySignature(
		proposerPubkey,
		p.computeSigningRoot(epoch, p.getDomain(ctx, epoch)),
		reveal,
	)
	if ok {
		p.logger.Info("randao reveal successfully verified 🤫 ",
			"reveal", reveal,
		)
	} else {
		p.logger.Info("reveal verification failed")
	}
	return ok
}

// MixinNewReveal mixes in a new reveal.
func (p *Processor) MixinNewReveal(
	ctx context.Context,
	blk beacontypes.BeaconBlock,
) error {
	st := p.BeaconState(ctx)
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
func (p *Processor) computeSigningRoot(
	epoch primitives.Epoch,
	domain common.BLSDomain,
) []byte {
	// Convert epoch to little endian bytes.
	epochBytes := make([]byte, 32)
	binary.LittleEndian.PutUint64(epochBytes, uint64(epoch))
	var epoch32 [32]byte
	copy(epoch32[:], epochBytes)

	// Compute the signing root.
	signingRoot := common.ComputeSigningRoot(epoch32, domain)

	// Convert the signing root to bytes.
	signingRootBytes := make([]byte, 32)
	copy(signingRootBytes[:], signingRoot[:])
	return signingRootBytes
}

// getDomain returns the domain.
func (p *Processor) getDomain(
	ctx context.Context,
	_ primitives.Epoch,
) common.BLSDomain {
	st := p.BeaconStateProvider.BeaconState(ctx)
	return common.ComputeDomain(
		common.DOMAIN_RANDAO,
		common.Version{0, 0, 0}, // TODO: this out with fork version
		tree.Root(st.GenesisValidatorsRoot()),
	)
}
