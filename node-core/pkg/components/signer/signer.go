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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/primitives/pkg/crypto"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/types"
)

// BLSSigner utilize an underlying PrivValidator signer using data persisted to
// disk to prevent double signing.
type BLSSigner struct {
	types.PrivValidator
}

// NewBLSSigner creates a new BLSSigner instance using the provided key and
// state
// file paths.
// If the key file does not exist, the program will exit.
func NewBLSSigner(keyFilePath string, stateFilePath string) *BLSSigner {
	filePV := privval.LoadFilePV(keyFilePath, stateFilePath)
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
	} else if len(sig) != constants.BLSSignatureLength {
		return crypto.BLSSignature{}, errors.Wrapf(
			ErrInvalidSignature, "expected signature length %d, got %d",
			constants.BLSSignatureLength, len(sig),
		)
	}
	return crypto.BLSSignature(sig), nil
}

// VerifySignature verifies a signature against a message and a public key.
func (f BLSSigner) VerifySignature(
	pubKey crypto.BLSPubkey,
	msg []byte,
	signature crypto.BLSSignature,
) error {
	if ok := bls12381.PubKey(pubKey[:]).
		VerifySignature(msg, signature[:]); !ok {
		return ErrInvalidSignature
	}
	return nil
}
