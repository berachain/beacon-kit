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
	"github.com/berachain/comet-bls12-381/bls/blst"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
)

func TestMix(t *testing.T) {
	var initMix Mix = sha256.Hash([]byte("init"))

	pvKey1, err := blst.RandKey()
	require.NoError(t, err)

	pvKey2, err := blst.RandKey()
	require.NoError(t, err)

	signingData := SigningData{
		Epoch:   12,
		ChainID: "berachain-1",
	}

	// Reveal 1, first signer
	reveal1, err := NewRandaoReveal(signingData, pvKey1)
	require.NoError(t, err)

	mix1 := initMix.MixWithReveal(reveal1)
	require.NoError(t, err)
	require.Equal(t, mix1, xor(sha256.Hash([]byte("init")), sha256.Hash(reveal1.Marshal())))

	// Reveal 2, second signer
	reveal2, err := NewRandaoReveal(signingData, pvKey2)
	require.NoError(t, err)
	require.Equal(t, mix1, xor(sha256.Hash([]byte("init")), sha256.Hash(reveal1.Marshal())))

	mix2 := mix1.MixWithReveal(reveal2)
	require.NoError(t, err)
	require.Equal(t, mix2, xor(mix1, sha256.Hash(reveal2.Marshal())))

	// XOR is commutative
	altMix1 := initMix.MixWithReveal(reveal2)
	altMix2 := altMix1.MixWithReveal(reveal1)
	require.NotEqual(t, mix1, altMix1) // first mix is different
	require.Equal(t, mix2, altMix2)    // but the second mix is the same
}

func xor(a, b [32]byte) Mix {
	dst := make([]byte, len(a))
	for i := range a {
		dst[i] = a[i] ^ b[i]
	}

	return Mix(dst)
}
