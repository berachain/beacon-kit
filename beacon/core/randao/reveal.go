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
	"encoding/binary"

	"github.com/berachain/comet-bls12-381/bls"
	"github.com/berachain/comet-bls12-381/bls/blst"
)

// Reveal represents the reveal of the validator for the current epoch.
type Reveal [BLSSignatureLength]byte

func (r Reveal) Verify(pubKey bls.PubKey, signingData SigningData) bool {
	bytes, err := blst.SignatureFromBytes(r[:])
	if err != nil {
		panic(err)
	}

	return bytes.Verify(pubKey, signingData.Marshall())
}

type Epoch uint64

type SigningData struct {
	Epoch   Epoch
	ChainID string
}

// Marshall converts the signing data into a byte slice to be signed.
// this includes the chain-id and the epoch.
func (s SigningData) Marshall() []byte {
	var buf []byte

	// TODO maybe caching?
	binary.LittleEndian.AppendUint64(buf, uint64(s.Epoch))
	buf = append(buf, []byte(s.ChainID)...)

	return buf
}

// NewRandaoReveal creates the randao reveal for a given signing data and private key.
func NewRandaoReveal(signingData SigningData, privKey bls.SecretKey) (Reveal, error) {
	sig := privKey.Sign(signingData.Marshall())

	var reveal Reveal
	copy(reveal[:], sig.Marshal())

	return reveal, nil
}
