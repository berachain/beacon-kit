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
	"errors"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/types"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// BLSSigner utilize an underlying PrivValidator signer using data persisted to disk to prevent
// double signing.
type BLSSigner struct {
	types.PrivValidator
}

func NewBLSSigner(keyFilePath string, stateFilePath string) *BLSSigner {
	filePV := privval.LoadOrGenFilePV(keyFilePath, stateFilePath)
	return &BLSSigner{PrivValidator: filePV}
}

// ========================== Implements BLS Signer ==========================

// PublicKey returns the public key of the signer.
func (f BLSSigner) PublicKey() crypto.BLSPubkey {
	key, err := f.PrivValidator.GetPubKey()
	if err != nil {
		return crypto.BLSPubkey{}
	}
	return crypto.BLSPubkey(key.Bytes())
}

// Sign generates a signature for a given message using the signer's secret key.
func (f BLSSigner) Sign(msg []byte) (crypto.BLSSignature, error) {
	sig, err := f.PrivValidator.SignBytes(msg)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	return crypto.BLSSignature(sig), nil
}

// VerifySignature verifies a signature against a message and a public key.
func (f BLSSigner) VerifySignature(
	blsPk crypto.BLSPubkey,
	msg []byte,
	signature crypto.BLSSignature) error {
	pk, err := blst.PublicKeyFromBytes(blsPk[:])
	if err != nil {
		return err
	}

	sig, err := blst.SignatureFromBytes(signature[:])
	if err != nil {
		return err
	}

	if !sig.Verify(pk, msg) {
		return errors.New("signature verification failed")
	}
	return nil
}
