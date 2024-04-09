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

package da

import (
	"github.com/berachain/beacon-kit/mod/da/proof"
	"github.com/berachain/beacon-kit/mod/da/types"
)

// BlobProofVerifier is a verifier for blobs.
type BlobVerifier struct {
	proofVerifier proof.BlobProofVerifier
}

// NewBlobVerifier creates a new BlobVerifier with the given proof verifier.
func NewBlobVerifier(
	proofVerifier proof.BlobProofVerifier,
) *BlobVerifier {
	return &BlobVerifier{
		proofVerifier: proofVerifier,
	}
}

// VerifyKZGProofs verifies the sidecars.
func (bv *BlobVerifier) VerifyKZGProofs(
	scs *types.BlobSidecars,
) error {
	switch len(scs.Sidecars) {
	case 0:
		return nil
	case 1:
		// This method is fastest for a single blob.
		return bv.proofVerifier.VerifyBlobProof(
			&scs.Sidecars[0].Blob,
			scs.Sidecars[0].KzgProof,
			scs.Sidecars[0].KzgCommitment,
		)
	default:
		// For multiple blobs batch verification is more performant
		// than verifying each blob individually (even when done in parallel).
		return bv.proofVerifier.VerifyBlobProofBatch(
			proof.ArgsFromSidecars(scs))
	}
}
