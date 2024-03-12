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

package types_test

import (
	"crypto/rand"
	"testing"

	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	"github.com/berachain/beacon-kit/crypto/sha256"
	"github.com/stretchr/testify/require"
)

func TestMix(t *testing.T) {
	var initMix types.Mix = sha256.Hash([]byte("init"))

	var randomReveal1 [bls12381.SignatureLength]byte
	_, err := rand.Read(randomReveal1[:])
	require.NoError(t, err)

	// Reveal 1, first signer
	mix1 := initMix.MixinNewReveal(randomReveal1)
	require.Equal(
		t,
		mix1,
		xor(sha256.Hash([]byte("init")), sha256.Hash(randomReveal1[:])),
	)

	// Reveal 2, second signer
	var randomReveal2 [bls12381.SignatureLength]byte
	_, err = rand.Read(randomReveal2[:])
	require.NoError(t, err)

	mix2 := mix1.MixinNewReveal(randomReveal2)
	require.Equal(t, mix2, xor(mix1, sha256.Hash(randomReveal2[:])))

	// XOR is commutative
	altMix1 := initMix.MixinNewReveal(randomReveal2)
	require.Equal(
		t,
		altMix1,
		xor(sha256.Hash([]byte("init")), sha256.Hash(randomReveal2[:])),
	)

	altMix2 := altMix1.MixinNewReveal(randomReveal1)
	require.Equal(t, altMix2, xor(altMix1, sha256.Hash(randomReveal1[:])))

	require.NotEqual(t, mix1, altMix1) // first mix is different
	require.Equal(t, mix2, altMix2)    // but the second mix is the same
}

func xor(a, b [32]byte) types.Mix {
	dst := make([]byte, len(a))
	for i := range a {
		dst[i] = a[i] ^ b[i]
	}

	return types.Mix(dst)
}
