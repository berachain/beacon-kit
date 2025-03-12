// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"golang.org/x/sync/errgroup"
)

// verifier is responsible for verifying blobs, including their
// inclusion and KZG proofs.
type verifier struct {
	// proofVerifier is used to verify the KZG proofs of the blobs.
	proofVerifier kzg.BlobProofVerifier
	// metrics collects and reports metrics related to the verification process.
	metrics *verifierMetrics
}

// newVerifier creates a new Verifier with the given proof verifier.
func newVerifier(
	proofVerifier kzg.BlobProofVerifier,
	telemetrySink TelemetrySink,
) *verifier {
	return &verifier{
		proofVerifier: proofVerifier,
		metrics:       newVerifierMetrics(telemetrySink),
	}
}

// verifySidecars verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *verifier) verifySidecars(
	ctx context.Context,
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	kzgCommitments eip4844.KZGCommitments[common.ExecutionHash],
) error {
	numSidecars := uint64(len(sidecars))
	defer bv.metrics.measureVerifySidecarsDuration(
		time.Now(), math.U64(numSidecars),
		bv.proofVerifier.GetImplementation(),
	)

	g, _ := errgroup.WithContext(ctx)

	// Create lookup table for each blob sidecar commitment and indicies.
	blobSidecarCommitments := make(map[eip4844.KZGCommitment]struct{})
	blobSidecarIndicies := make(map[uint64]struct{})

	// Validate sidecar fields against data from the BeaconBlock.
	for i, s := range sidecars {
		// Fill lookup table with commitments from the blob sidecars.
		blobSidecarCommitments[s.GetKzgCommitment()] = struct{}{}

		// We should only have unique indexes.
		if _, exists := blobSidecarIndicies[s.GetIndex()]; exists {
			return fmt.Errorf("duplicate sidecar Index: %d", i)
		}
		blobSidecarIndicies[s.GetIndex()] = struct{}{}

		// This check happens outside the goroutines so that we do not
		// process the inclusion proofs before validating the index.
		if s.GetIndex() >= numSidecars {
			return fmt.Errorf("invalid sidecar Index: %d", i)
		}

		// Verify the signature.
		var sigHeader = s.GetSignedBeaconBlockHeader()

		// Check BlobSidecar.Header equality with BeaconBlockHeader
		if !sigHeader.GetHeader().Equals(blkHeader) {
			return fmt.Errorf("unequal block header: idx: %d", i)
		}
	}

	// Ensure each commitment from the BeaconBlock has a corresponding sidecar commitment.
	for _, kzgCommitment := range kzgCommitments {
		if _, exists := blobSidecarCommitments[kzgCommitment]; !exists {
			return fmt.Errorf("missing kzg commitment: %s", kzgCommitment)
		}
	}

	// Verify the inclusion proofs on the blobs concurrently.
	g.Go(func() error {
		return bv.verifyInclusionProofs(sidecars)
	})

	// Verify the KZG proofs on the blobs concurrently.
	g.Go(func() error {
		return bv.verifyKZGProofs(sidecars)
	})

	// Wait for all goroutines to finish and return the result.
	return g.Wait()
}

func (bv *verifier) verifyInclusionProofs(
	scs datypes.BlobSidecars,
) error {
	startTime := time.Now()
	defer bv.metrics.measureVerifyInclusionProofsDuration(
		startTime, math.U64(len(scs)),
	)

	return scs.VerifyInclusionProofs()
}

// verifyKZGProofs verifies the sidecars.
func (bv *verifier) verifyKZGProofs(
	scs datypes.BlobSidecars,
) error {
	start := time.Now()
	defer bv.metrics.measureVerifyKZGProofsDuration(
		start, math.U64(len(scs)),
		bv.proofVerifier.GetImplementation(),
	)

	switch len(scs) {
	case 0:
		return nil
	case 1:
		blob := scs[0].GetBlob()
		// This method is fastest for a single blob.
		return bv.proofVerifier.VerifyBlobProof(
			&blob,
			scs[0].GetKzgProof(),
			scs[0].GetKzgCommitment(),
		)
	default:
		// For multiple blobs batch verification is more performant
		// than verifying each blob individually (even when done in parallel).
		return bv.proofVerifier.VerifyBlobProofBatch(kzg.ArgsFromSidecars(scs))
	}
}
