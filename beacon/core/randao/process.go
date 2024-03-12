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

	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	"github.com/berachain/beacon-kit/beacon/core/state"
	"github.com/berachain/beacon-kit/beacon/signing"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/primitives"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

type BeaconStateProvider interface {
	// BeaconState returns the current beacon state.
	BeaconState(context.Context) state.BeaconState
}

// Processor is the randao processor.
type Processor struct {
	stateProvider BeaconStateProvider
	signer        bls12381.Signer
}

func NewProcessor(
	beaconStateProvider BeaconStateProvider,
	signer bls12381.Signer,
) *Processor {
	return &Processor{
		stateProvider: beaconStateProvider,
		signer:        signer,
	}
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
func (rs *Processor) BuildReveal(
	_ context.Context,
	epoch primitives.Epoch,
) (types.Reveal, error) {
	signingRoot := rs.GetSigningRoot(epoch)

	return rs.signer.Sign(signingRoot), nil
}

func (rs *Processor) GetSigningRoot(
	epoch primitives.Epoch,
) []byte {
	return rs.computeSigningRoot(epoch, rs.getDomain(epoch))
}

func (rs *Processor) computeSigningRoot(
	epoch primitives.Epoch,
	_ signing.Domain,
) []byte {
	return sdktypes.Uint64ToBigEndian(epoch)
}

func (rs *Processor) getDomain(
	_ primitives.Epoch,
) signing.Domain {
	return signing.Domain{}
}

// VerifyReveal verifies the reveal of the proposer.
func (rs *Processor) VerifyReveal(
	proposerPubkey [bls12381.PubKeyLength]byte,
	msg []byte,
	reveal types.Reveal,
) bool {
	return bls12381.VerifySignature(proposerPubkey, msg, reveal)
}
