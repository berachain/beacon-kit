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

	"github.com/berachain/comet-bls12-381/bls/blst"
	"github.com/itsdevbear/bolaris/beacon/core/randao/types"
	"github.com/itsdevbear/bolaris/beacon/core/state"
	bls12_381 "github.com/itsdevbear/bolaris/crypto/bls12_381"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

type beaconStateProvider interface {
	// GetBeaconState returns the current beacon state.
	BeaconState(context.Context) state.BeaconState
}

type blsSigner interface {
	Sign([]byte) [bls12_381.SignatureLength]byte
	Verify([32]byte, []byte, [bls12_381.SignatureLength]byte) bool
}

// Processor is the randao processor.
type Processor struct {
	beaconStateProvider
	signer blsSigner
	cfg    *Config
}

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
	domain := rs.getDomain(epoch, nil)
	signingRoot := rs.computeSigningRoot(epoch, domain)
	return rs.signer.Sign(signingRoot), nil
}

// def process_randao(state: BeaconState, body: BeaconBlockBody) -> None:
//
//	epoch = get_current_epoch(state)
//	# Verify RANDAO reveal
//	proposer = state.validators[get_beacon_proposer_index(state)]
//	signing_root = compute_signing_root(epoch, get_domain(state, DOMAIN_RANDAO))
//	assert bls.Verify(proposer.pubkey, signing_root, body.randao_reveal)
//	# Mix in RANDAO reveal
//	mix = xor(get_randao_mix(state, epoch), hash(body.randao_reveal))
//	state.randao_mixes[epoch % EPOCHS_PER_HISTORICAL_VECTOR] = mix
func (rs *Processor) ProcessRandao(
	ctx context.Context,
	epoch primitives.Epoch,
	proposerPubkey [32]byte,
	prevReveal types.Reveal,
) error {
	st := rs.BeaconState(ctx)
	signingRoot := rs.computeSigningRoot(epoch, rs.getDomain(epoch, nil))
	rs.signer.Verify(proposerPubkey, signingRoot, prevReveal)
	mix, err := st.RandaoMix()
	if err != nil {
		return err
	}

	return st.SetRandaoMix(
		epoch%rs.cfg.EpochsPerHistoricalVector,
		mix.MixinNewReveal(prevReveal),
	)
}

func (rs *Processor) computeSigningRoot(
	_ primitives.Epoch,
	_ types.Domain,
) []byte {
	return []byte{}
}

func (rs *Processor) getDomain(
	_ primitives.Epoch,
	_ []byte,
) types.Domain {
	return types.BuildDomain()
}

// VerifyReveal verifies the reveal of the proposer.
func (Processor) VerifyReveal(
	proposerPubkey []byte,
	signature [bls12_381.SignatureLength]byte,
	reveal types.Reveal,
) bool {
	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		panic(err)
	}

	p, err := blst.PublicKeyFromBytes(proposerPubkey)
	if err != nil {
		panic(err)
	}

	return sig.Verify(p, reveal[:])
}
