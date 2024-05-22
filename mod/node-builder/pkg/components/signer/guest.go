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

package signer

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/itsdevbear/comet-bls12-381/bls"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

type GuestSigner struct {
	bls.SecretKey
}

// NewGuestSigner creates a new Signer instance given a secret key.
func NewGuestSigner(
	keyBz [constants.BLSSecretKeyLength]byte,
) (*GuestSigner, error) {
	secretKey, err := blst.SecretKeyFromBytes(keyBz[:])
	if err != nil {
		return nil, err
	}
	return &GuestSigner{
		SecretKey: secretKey,
	}, nil
}

// PublicKey returns the public key of the signer.
func (b *GuestSigner) PublicKey() crypto.BLSPubkey {
	return crypto.BLSPubkey(b.SecretKey.PublicKey().Marshal())
}

// Sign generates a signature for a given message using the signer's secret key.
// It returns the signature and any error encountered during the signing
// process.
func (b *GuestSigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	return crypto.BLSSignature(b.SecretKey.Sign(msg).Marshal()), nil
}

// VerifySignature verifies a signature for a given message using the signer's
// public key.
func (GuestSigner) VerifySignature(
	pubKey crypto.BLSPubkey,
	msg []byte,
	signature crypto.BLSSignature,
) error {
	pubkey, err := blst.PublicKeyFromBytes(pubKey[:])
	if err != nil {
		return err
	}

	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		return err
	}

	if !sig.Verify(pubkey, msg) {
		return ErrInvalidSignature
	}
	return nil
}
