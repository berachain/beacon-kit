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
	"os"

	"cosmossdk.io/depinject"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/itsdevbear/bolaris/primitives"
	"github.com/spf13/cast"

	"github.com/itsdevbear/bolaris/beacon/core/randao/types"
	"github.com/itsdevbear/bolaris/beacon/core/state"
	bls12381 "github.com/itsdevbear/bolaris/crypto/bls12_381"
)

type BeaconStateProvider interface {
	// BeaconState returns the current beacon state.
	BeaconState(context.Context) state.BeaconState
}

// Processor is the randao processor.
type Processor struct {
	stateProvider BeaconStateProvider
	signer        bls12381.BlsSigner
	cfg           *Config
}

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In

	BeaconState BeaconStateProvider
	AppOpts     servertypes.AppOptions
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out

	RandaoProcessor *Processor
}

func ProvideRandaoProcessor(in DepInjectInput) DepInjectOutput {
	homeDir := cast.ToString(in.AppOpts.Get(flags.FlagHome))
	fmt.Println("HomeDir: ", homeDir)
	key, err := p2p.LoadNodeKey(fmt.Sprintf("%s/config/priv_validator_key.json", homeDir))
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Key: ", key.PrivKey)

	var pk [32]byte
	copy(pk[:], key.PrivKey.Bytes())

	signer := bls12381.NewBlsSigner(pk)
	processor := NewProcessor(in.BeaconState, signer, &Config{
		EpochsPerHistoricalVector: 0,
		ConfiguredPubKeyLength:    0,
	})

	return DepInjectOutput{
		RandaoProcessor: processor,
	}
}

func NewProcessor(beaconStateProvider BeaconStateProvider, signer bls12381.BlsSigner, cfg *Config) *Processor {
	return &Processor{stateProvider: beaconStateProvider, signer: signer, cfg: cfg}
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
	ctx context.Context,
	epoch primitives.Epoch,
) (types.Reveal, error) {
	st := rs.stateProvider.BeaconState(ctx)
	root := st.GetParentBlockRoot()
	domain := rs.getDomain(epoch, root[:])
	signingRoot := rs.computeSigningRoot(epoch, domain)

	return rs.signer.Sign(signingRoot)
}

// ProcessRandao
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
	proposerPubkey [bls12381.PubKeyLength]byte,
	prevReveal types.Reveal,
) error {
	st := rs.stateProvider.BeaconState(ctx)
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
	epoch primitives.Epoch,
	d types.Domain,
) []byte {
	epochSSZUInt64 := primitives.SSZUint64(epoch)
	sszBz, err := epochSSZUInt64.MarshalSSZ()
	if err != nil {
		// don't actually panic
		panic(err)
	}

	return sszBz
}

func (rs *Processor) getDomain(
	epoch primitives.Epoch,
	_ []byte,
) types.Domain {
	epochSSZUInt64 := primitives.SSZUint64(epoch)
	sszBz, err := epochSSZUInt64.MarshalSSZ()
	if err != nil {
		// don't actually panic
		panic(err)
	}

	_ = sszBz

	// We can also get the has tree root (trivial because this is one item but
	// yeah)

	htr, err := epochSSZUInt64.HashTreeRoot()
	if err != nil {
		panic(err)
	}

	_ = htr
	return types.BuildDomain()
}

// VerifyReveal verifies the reveal of the proposer.
func (rs Processor) VerifyReveal(
	proposerPubkey [bls12381.PubKeyLength]byte,
	signature [bls12381.SignatureLength]byte,
	reveal types.Reveal,
) bool {
	return rs.signer.Verify(proposerPubkey, reveal[:], signature)
}
