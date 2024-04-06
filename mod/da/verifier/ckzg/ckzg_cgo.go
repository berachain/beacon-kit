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

//go:build ckzg

package ckzg

import (
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
)

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v Verifier) VerifyKZGProof(
	commitment kzg.Commitment,
	point kzg.Point,
	claim kzg.Claim,
	proof kzg.Proof,
) error {
	valid, err := ckzg4844.VerifyKZGProof(
		(ckzg4844.Bytes48)(commitment),
		(ckzg4844.Bytes32)(point),
		(ckzg4844.Bytes32)(claim),
		(ckzg4844.Bytes48)(proof),
	)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidProof
	}
	return nil
}

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v Verifier) VerifyBlobProof(
	blob *kzg.Blob,
	commitment kzg.Commitment,
	proof kzg.Proof,
) error {
	valid, err := ckzg4844.VerifyBlobKZGProof(
		(*ckzg4844.Blob)(blob),
		(ckzg4844.Bytes48)(commitment),
		(ckzg4844.Bytes48)(proof),
	)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidProof
	}
	return nil
}
