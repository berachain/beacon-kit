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

package blob

import (
	"context"
	"fmt"
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

// VerifySidecars verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *Verifier) VerifySidecars(
	sidecars *types.BlobSidecars, kzgOffset uint64,
) error {
	var (
		g, _      = errgroup.WithContext(context.Background())
		startTime = time.Now()
	)

	defer bv.metrics.measureVerifySidecarsDuration(
		startTime, math.U64(len(sidecars.Sidecars)),
		bv.proofVerifier.GetImplementation(),
	)

	fmt.Println("Starting sidecar verification",
		"num_sidecars", len(sidecars.Sidecars),
		"kzg_offset", kzgOffset)

	// Verify the inclusion proofs on the blobs concurrently.
	g.Go(func() error {
		// TODO: KZGOffset needs to be configurable and not
		// passed in.
		err := bv.VerifyInclusionProofs(sidecars, kzgOffset)
		if err != nil {
			fmt.Println("Inclusion proof verification failed", "error", err)
		} else {
			fmt.Println("Inclusion proof verification succeeded")
		}
		return err
	})

	// Verify the KZG proofs on the blobs concurrently.
	g.Go(func() error {
		err := bv.VerifyKZGProofs(sidecars)
		if err != nil {
			fmt.Println("KZG proof verification failed", "error", err)
		} else {
			fmt.Println("KZG proof verification succeeded")
		}
		return err
	})

	g.Go(func() error {
		err := sidecars.ValidateBlockRoots()
		if err != nil {
			fmt.Println("Block root validation failed", "error", err)
		} else {
			fmt.Println("Block root validation succeeded")
		}
		return err
	})

	// Wait for all goroutines to finish and return the result.
	err := g.Wait()
	if err != nil {
		fmt.Println("Sidecar verification failed", "error", err)
	} else {
		fmt.Println("Sidecar verification succeeded")
	}
	return err
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
