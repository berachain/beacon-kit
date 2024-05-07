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

package gokzg

import (
	"unsafe"

	prooftypes "github.com/berachain/beacon-kit/mod/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// Verifier is a KZG verifier that uses the Go implementation of KZG.
type Verifier struct {
	*gokzg4844.Context
}

// NewVerifier creates a new GoKZGVerifier.
func NewVerifier(ts *gokzg4844.JSONTrustedSetup) (*Verifier, error) {
	ctx, err := gokzg4844.NewContext4096(ts)
	if err != nil {
		return nil, err
	}
	return &Verifier{ctx}, nil
}

// VerifyProof verifies the KZG proof that the polynomial represented by the
// blob evaluated at the given point is the claimed value.
func (v Verifier) VerifyBlobProof(
	blob *eip4844.Blob,
	proof eip4844.KZGProof,
	commitment eip4844.KZGCommitment,
) error {
	return v.Context.
		VerifyBlobKZGProof(
			(*gokzg4844.Blob)(blob),
			(gokzg4844.KZGCommitment)(commitment),
			(gokzg4844.KZGProof)(proof))
}

// VerifyBlobProofBatch verifies the KZG proof that the polynomial represented
// by the blob evaluated at the given point is the claimed value.
// It is more efficient than VerifyBlobProof when verifying multiple proofs.
func (v Verifier) VerifyBlobProofBatch(
	args *prooftypes.BlobProofArgs,
) error {
	blobs := make([]gokzg4844.Blob, len(args.Blobs))
	for i := range args.Blobs {
		blobs[i] = *(*gokzg4844.Blob)(args.Blobs[i])
	}

	//#nosec:G103 // "use of unsafe calls should be audited" lmeow.
	return v.Context.
		VerifyBlobKZGProofBatch(
			blobs,
			*(*[]gokzg4844.KZGCommitment)(
				unsafe.Pointer(&args.Commitments)),
			*(*[]gokzg4844.KZGProof)(unsafe.Pointer(&args.Proofs)),
		)
}
