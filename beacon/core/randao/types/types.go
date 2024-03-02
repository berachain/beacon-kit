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

package types

import (
	"github.com/berachain/comet-bls12-381/bls/blst"
	bls "github.com/itsdevbear/bolaris/crypto/bls12_381"
	"github.com/itsdevbear/bolaris/crypto/sha256"
)

type Domain [32]byte

func BuildDomain() Domain {
	return Domain{}
}

// Reveal represents the signature of the RANDAO reveal.
type Reveal [bls.SignatureLength]byte

// BuildReveal creates a new Reveal from the given signature.
//
// https://github.com/ethereum/consensus-specs/blob/v1.3.0/specs/phase0/validator.md#randao-reveal
//
// nolint:lll
func BuildReveal(signData []byte, privKey [bls.SecretKeyLength]byte) Reveal {
	pk, err := blst.SecretKeyFromBytes(privKey[:])
	if err != nil {
		panic(err)
	}

	sig := pk.Sign(signData)

	var reveal Reveal
	copy(reveal[:], sig.Marshal())

	return reveal
}

// Mix
// We can think of a RANDAO as being like a deck of cards that's passed round
// the table, each person shuffling it in turn: the deck gets repeatedly
// re-randomised. Even if one contributor's randomness is weak, the cumulative
// result has a high level of entropy.
//
// In order to accomplish this, we need to keep track of the "current mix",
// which represents the current state of this shuffled deck of cards, as it
// is passed around the table.
type Mix [32]byte

// MixinNewReveal
// the new Hash into the current Mix. We first take the existing reveal and
// hash itt hen we XOR the result into the current mix.
func (m Mix) MixinNewReveal(reveal Reveal) Mix {
	for idx, b := range sha256.Hash(reveal[:]) {
		m[idx] ^= b
	}
	return m
}
