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

	"github.com/berachain/beacon-kit/chain-spec/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/crypto"
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
	// chainSpec contains the chain specification
	chainSpec chain.ChainSpec
}

// newVerifier creates a new Verifier with the given proof verifier.
func newVerifier(
	proofVerifier kzg.BlobProofVerifier,
	telemetrySink TelemetrySink,
	chainSpec chain.ChainSpec,
) *verifier {
	return &verifier{
		proofVerifier: proofVerifier,
		metrics:       newVerifierMetrics(telemetrySink),
		chainSpec:     chainSpec,
	}
}

// verifySidecars verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *verifier) verifySidecars(
	sidecars datypes.BlobSidecars,
	blkHeader *ctypes.BeaconBlockHeader,
	verifierFn func(
		blkHeader *ctypes.BeaconBlockHeader,
		signature crypto.BLSSignature,
	) error,
) error {
	defer bv.metrics.measureVerifySidecarsDuration(
		time.Now(), math.U64(len(sidecars)),
		bv.proofVerifier.GetImplementation(),
	)

	g, _ := errgroup.WithContext(context.Background())

	// create lookup table to check for duplicate commitments
	duplicateCommitment := make(map[eip4844.KZGCommitment]struct{})

	// Validate sidecar fields against data from the BeaconBlock.
	for i, s := range sidecars {
		// Check if sidecar's kzgCommitment is duplicate. Along with the
		// length check and the inclusion proof, this fully verifies that
		// the KzgCommitments in the BlobSidecar are the exact same as the
		// ones in the BeaconBlockBody without having to explicitly compare.
		if _, exists := duplicateCommitment[s.GetKzgCommitment()]; exists {
			return errors.New(
				"found duplicate KzgCommitments in BlobSidecars",
			)
		}
		duplicateCommitment[s.GetKzgCommitment()] = struct{}{}

		// This check happens outside the goroutines so that we do not
		// process the inclusion proofs before validating the index.
		if s.GetIndex() >= bv.chainSpec.MaxBlobsPerBlock() {
			return fmt.Errorf("invalid sidecar Index: %d", i)
		}
		g.Go(func() error {
			// Verify the signature.
			var sigHeader = s.GetSignedBeaconBlockHeader()

			// Check BlobSidecar.Header equality with BeaconBlockHeader
			if !sigHeader.GetHeader().Equals(blkHeader) {
				return fmt.Errorf("unequal block header: idx: %d", i)
			}

			// Verify BeaconBlockHeader with signature
			if err := verifierFn(
				blkHeader,
				sigHeader.GetSignature(),
			); err != nil {
				return err
			}
			return nil
		})
	}

	// Verify the inclusion proofs on the blobs concurrently.
	g.Go(func() error {
		return bv.verifyInclusionProofs(sidecars, blkHeader.GetSlot())
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
	slot math.Slot,
) error {
	startTime := time.Now()
	defer bv.metrics.measureVerifyInclusionProofsDuration(
		startTime, math.U64(len(scs)),
	)

	// Grab the KZG offset for the fork version.
	kzgOffset, err := ctypes.BlockBodyKZGOffset(
		slot, bv.chainSpec,
	)
	if err != nil {
		return err
	}

	// Grab the inclusion proof depth for the fork version.
	inclusionProofDepth, err := ctypes.KZGCommitmentInclusionProofDepth(
		slot, bv.chainSpec,
	)
	if err != nil {
		return err
	}

	return scs.VerifyInclusionProofs(kzgOffset, inclusionProofDepth)
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
