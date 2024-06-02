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
	"time"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"golang.org/x/sync/errgroup"
)

// Verifier is responsible for verifying blobs, including their
// inclusion and KZG proofs.
type Verifier struct {
	// proofVerifier is used to verify the KZG proofs of the blobs.
	proofVerifier kzg.BlobProofVerifier
	// metrics collects and reports metrics related to the verification process.
	metrics *verifierMetrics
}

// NewVerifier creates a new Verifier with the given proof verifier.
func NewVerifier(
	proofVerifier kzg.BlobProofVerifier,
	telemetrySink TelemetrySink,
) *Verifier {
	return &Verifier{
		proofVerifier: proofVerifier,
		metrics:       newVerifierMetrics(telemetrySink),
	}
}

// VerifyBlobs verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *Verifier) VerifyBlobs(
	sidecars *types.BlobSidecars, kzgOffset uint64,
) error {
	var (
		g, _      = errgroup.WithContext(context.Background())
		startTime = time.Now()
	)

	defer bv.metrics.measureVerifyBlobsDuration(
		startTime, math.U64(len(sidecars.Sidecars)),
		bv.proofVerifier.GetImplementation(),
	)

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

	g.Go(func() error {
		return sidecars.ValidateBlockRoots()
	})

	// Wait for all goroutines to finish and return the result.
	return g.Wait()
}

func (bv *Verifier) VerifyInclusionProofs(
	scs *types.BlobSidecars,
	kzgOffset uint64,
) error {
	startTime := time.Now()
	defer bv.metrics.measureVerifyInclusionProofsDuration(
		startTime, math.U64(len(scs.Sidecars)),
	)
	return scs.VerifyInclusionProofs(kzgOffset)
}

// VerifyKZGProofs verifies the sidecars.
func (bv *Verifier) VerifyKZGProofs(
	scs *types.BlobSidecars,
) error {
	start := time.Now()
	defer bv.metrics.measureVerifyKZGProofsDuration(
		start, math.U64(len(scs.Sidecars)),
		bv.proofVerifier.GetImplementation(),
	)

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
