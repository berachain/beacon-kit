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
	"encoding/hex"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/itsdevbear/comet-bls12-381/bls"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// LegacySigner is a BLS12-381 signer that uses a bls.SecretKey for signing.
type LegacySigner struct {
	bls.SecretKey
}

// NewLegacySigner creates a new Signer instance given a secret key.
func NewLegacySigner(
	keyBz LegacyKey,
) (*LegacySigner, error) {
	secretKey, err := blst.SecretKeyFromBytes(keyBz[:])
	if err != nil {
		return nil, err
	}
	return &LegacySigner{
		SecretKey: secretKey,
	}, nil
}

// PublicKey returns the public key of the signer.
func (b *LegacySigner) PublicKey() crypto.BLSPubkey {
	return crypto.BLSPubkey(b.SecretKey.PublicKey().Marshal())
}

// Sign generates a signature for a given message using the signer's secret key.
// It returns the signature and any error encountered during the signing
// process.
func (b *LegacySigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	return crypto.BLSSignature(b.SecretKey.Sign(msg).Marshal()), nil
}

// VerifySignature verifies a signature against a message and public key.
func (LegacySigner) VerifySignature(
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

// LegacyKey is a byte array that represents a BLS12-381 secret key.
type LegacyKey [constants.BLSSecretKeyLength]byte

// GetLegacyInput returns a LegacyInput from a hex-encoded string.
func LegacyKeyFromString(privKey string) (LegacyKey, error) {
	privKeyBz, err := hex.DecodeString(privKey)
	if err != nil {
		return LegacyKey{}, err
	}
	if len(privKeyBz) != constants.BLSSecretKeyLength {
		return LegacyKey{}, ErrInvalidValidatorPrivateKeyLength
	}
	return LegacyKey(privKeyBz), nil
}
