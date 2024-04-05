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

package verifier

import (
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// GoKZGVerifier is a KZG verifier that uses the Go implementation of KZG.
type GoKZGVerifier struct {
	*gokzg4844.Context
}

// NewGoKZGVerifier creates a new GoKZGVerifier.
func NewGoKZGVerifier(ts *gokzg4844.JSONTrustedSetup) (*GoKZGVerifier, error) {
	ctx, err := gokzg4844.NewContext4096(ts)
	if err != nil {
		return nil, err
	}
	return &GoKZGVerifier{ctx}, nil
}

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v GoKZGVerifier) VerifyKZGProof(
	commitment kzg.Commitment,
	point kzg.Point,
	claim kzg.Claim,
	proof kzg.Proof,
) error {
	return v.Context.
		VerifyKZGProof((gokzg4844.KZGCommitment)(commitment),
			(gokzg4844.Scalar)(point),
			(gokzg4844.Scalar)(claim),
			(gokzg4844.KZGProof)(proof),
		)
}

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v GoKZGVerifier) VerifyBlobProof(
	blob *kzg.Blob,
	commitment kzg.Commitment,
	proof kzg.Proof,
) error {
	return v.Context.
		VerifyBlobKZGProof(
			(*gokzg4844.Blob)(blob),
			(gokzg4844.KZGCommitment)(commitment),
			(gokzg4844.KZGProof)(proof))
}
