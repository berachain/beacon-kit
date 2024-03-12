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

package bls12381

import "github.com/itsdevbear/comet-bls12-381/bls/blst"

// VerifySignature checks if a given signature is valid for a message and public
// key.
// It returns true if the signature is valid, otherwise it panics if an error
// occurs during the verification process.
func VerifySignature(
	pubKey [PubKeyLength]byte,
	msg []byte,
	signature [SignatureLength]byte,
) bool {
	pubkey, err := blst.PublicKeyFromBytes(pubKey[:])
	if err != nil {
		panic(err)
	}
	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		panic(err)
	}

	return sig.Verify(pubkey, msg)
}
