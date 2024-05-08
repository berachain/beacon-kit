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
	"github.com/cometbft/cometbft/p2p"
	bls "github.com/itsdevbear/comet-bls12-381/bls"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// BLSSigner holds the secret key used for signing operations.
type BLSSigner struct {
	bls.SecretKey
}

// NewBLSSigner creates a new Signer instance given a secret key.
func NewBLSSigner(
	keyBz [constants.BLSSecretKeyLength]byte,
) (*BLSSigner, error) {
	secretKey, err := blst.SecretKeyFromBytes(keyBz[:])
	if err != nil {
		return nil, err
	}
	return &BLSSigner{
		SecretKey: secretKey,
	}, nil
}

// NewSignerFromFile creates a new Signer instance given a
// file path to a CometBFT node key.
//
// TODO: In order to make RANDAO signing compatible with the BLS12-381
// we need to extend the PrivVal interface to allow signing arbitrary things.
// @tac0turtle.
func NewFromCometBFTNodeKey(filePath string) (*BLSSigner, error) {
	key, err := p2p.LoadNodeKey(filePath)
	if err != nil {
		return nil, err
	}

	return NewBLSSigner(
		[constants.BLSSecretKeyLength]byte(key.PrivKey.Bytes()))
}

// PublicKey returns the public key of the signer.
func (b *BLSSigner) PublicKey() crypto.BLSPubkey {
	return crypto.BLSPubkey(b.SecretKey.PublicKey().Marshal())
}

// Sign generates a signature for a given message using the signer's secret key.
// It returns the signature and any error encountered during the signing
// process.
func (b *BLSSigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	return crypto.BLSSignature(b.SecretKey.Sign(msg).Marshal()), nil
}

// VerifySignature verifies a signature for a given message using the signer's
// public key.
func (BLSSigner) VerifySignature(
	pubKey []byte,
	msg []byte,
	signature []byte,
) error {
	pubkey, err := blst.PublicKeyFromBytes(pubKey)
	if err != nil {
		return err
	}

	sig, err := blst.SignatureFromBytes(signature)
	if err != nil {
		return err
	}

	if !sig.Verify(pubkey, msg) {
		return ErrInvalidSignature
	}
	return nil
}
