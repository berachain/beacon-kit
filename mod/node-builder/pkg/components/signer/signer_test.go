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

package signer_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/node-builder/pkg/components/signer"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	blst "github.com/itsdevbear/comet-bls12-381/bls/blst"
	"github.com/stretchr/testify/require"
)

func TestSigner_SignAndVerify(t *testing.T) {
	// Generate a random secret key for testing
	key, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	signer, err := signer.NewBLSSigner([32]byte(key.Marshal()))
	require.NoError(t, err, "Failed to create signer")

	// Test message
	message := []byte("test message")

	// Sign the message
	signature, err := signer.Sign(message)
	require.NoError(t, err, "Failed to sign message")

	// Extract the public key from the signer
	pubKeyBytes := signer.SecretKey.PublicKey().Marshal()
	pubkey, err := blst.PublicKeyFromBytes(pubKeyBytes)
	require.NoError(t, err, "Failed to create public key")
	sig, err := blst.SignatureFromBytes(signature[:])
	require.NoError(t, err, "Failed to create signature")

	// Verify the signature
	require.True(t, sig.Verify(pubkey, message), "Signature should be valid")
}

func TestSigner_FailOnInvalidSignature(t *testing.T) {
	// Generate a random secret key for testing
	key1, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	_signer, err := signer.NewBLSSigner(
		[constants.BLSSecretKeyLength]byte(key1.Marshal()),
	)
	require.NoError(t, err, "Failed to create signer")

	// Generate a second random key.
	key2, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	otherSigner, err := signer.NewBLSSigner([32]byte(key2.Marshal()))
	require.NoError(t, err, "Failed to create signer")

	// Test message
	message := []byte("test message")

	// Incorrect signature
	signedByWrongKey, err := otherSigner.Sign(message)
	require.NoError(t, err, "Failed to sign message")

	pubKeyBytes := _signer.SecretKey.PublicKey().Marshal()

	// Attempt to verify the incorrect signature
	pubkey, err := blst.PublicKeyFromBytes(pubKeyBytes)
	require.NoError(t, err, "Failed to create public key")
	sig, err := blst.SignatureFromBytes(signedByWrongKey[:])
	require.NoError(t, err, "Failed to create signature")
	require.False(t, sig.Verify(pubkey, message), "Signature should be invalid")
}
