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

package blob

import (
	"context"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"golang.org/x/sync/errgroup"
)

// Verifier is a verifier for blobs.
type Verifier struct {
	proofVerifier kzg.BlobProofVerifier
}

// NewVerifier creates a new Verifier with the given proof verifier.
func NewVerifier(
	proofVerifier kzg.BlobProofVerifier,
) *Verifier {
	return &Verifier{
		proofVerifier: proofVerifier,
	}
}

// VerifyBlobs verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *Verifier) VerifyBlobs(
	sidecars *types.BlobSidecars, kzgOffset uint64,
) error {
	g, _ := errgroup.WithContext(context.Background())

	// Verify the inclusion proofs on the blobs concurrently.
	g.Go(func() error {
		// TODO: KZGOffset needs to be configurable and not
		// passed in.
		return bv.VerifyInclusionProofs(
			sidecars, kzgOffset,
		)
	})

	// Verify the KZG proofs on the blobs concurrently.
	g.Go(func() error {
		return bv.VerifyKZGProofs(sidecars)
	})

	// Wait for all goroutines to finish and return the result.
	return g.Wait()
}

func (bv *Verifier) VerifyInclusionProofs(
	scs *types.BlobSidecars,
	kzgOffset uint64,
) error {
	return scs.VerifyInclusionProofs(kzgOffset)
}

// VerifyKZGProofs verifies the sidecars.
func (bv *Verifier) VerifyKZGProofs(
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
		return bv.proofVerifier.VerifyBlobProofBatch(kzg.ArgsFromSidecars(scs))
	}
}
