// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

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

// LegacyKeyFromString returns a LegacyKey from a hex-encoded string.
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
