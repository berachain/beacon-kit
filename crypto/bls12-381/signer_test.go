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

package bls12381_test

import (
	"testing"

	bls12_381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	blst "github.com/itsdevbear/comet-bls12-381/bls/blst"
	"github.com/stretchr/testify/require"
)

func TestSigner_SignAndVerify(t *testing.T) {
	// Generate a random secret key for testing
	key, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	signer, err := bls12_381.NewSigner([32]byte(key.Marshal()))
	require.NoError(t, err, "Failed to create signer")

	// Test message
	message := []byte("test message")

	// Sign the message
	signature := signer.Sign(message)

	// Extract the public key from the signer
	pubKeyBytes := signer.SecretKey.PublicKey().Marshal()

	// Verify the signature
	valid := bls12_381.VerifySignature(
		[48]byte(pubKeyBytes),
		message,
		signature,
	)
	require.True(t, valid, "Signature should be valid")
}

func TestSigner_FailOnInvalidSignature(t *testing.T) {
	// Generate a random secret key for testing
	key1, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	signer, err := bls12_381.NewSigner(
		[bls12_381.SecretKeyLength]byte(key1.Marshal()),
	)
	require.NoError(t, err, "Failed to create signer")

	// Generate a second random key.
	key2, err := blst.RandKey()
	require.NoError(t, err, "Failed to generate random key")
	otherSigner, err := bls12_381.NewSigner([32]byte(key2.Marshal()))
	require.NoError(t, err, "Failed to create signer")

	// Test message
	message := []byte("test message")

	// Incorrect signature
	signedByWrongKey := otherSigner.Sign(message)

	pubKeyBytes := signer.SecretKey.PublicKey().Marshal()

	// Attempt to verify the incorrect signature
	valid := bls12_381.VerifySignature(
		[bls12_381.PubKeyLength]byte(pubKeyBytes),
		message,
		signedByWrongKey,
	)
	require.False(t, valid, "Signature should be invalid")
}
