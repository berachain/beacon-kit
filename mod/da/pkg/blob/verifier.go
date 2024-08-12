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
	"time"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"golang.org/x/sync/errgroup"
)

// Verifier is responsible for verifying blobs, including their
// inclusion and KZG proofs.
type Verifier[
	BeaconBlockHeaderT BeaconBlockHeader,
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT, BeaconBlockHeaderT],
] struct {
	// proofVerifier is used to verify the KZG proofs of the blobs.
	proofVerifier kzg.BlobProofVerifier
	// metrics collects and reports metrics related to the verification process.
	metrics *verifierMetrics
}

// NewVerifier creates a new Verifier with the given proof verifier.
func NewVerifier[
	BeaconBlockHeaderT BeaconBlockHeader,
	BlobSidecarT Sidecar[BeaconBlockHeaderT],
	BlobSidecarsT Sidecars[BlobSidecarT, BeaconBlockHeaderT],
](
	proofVerifier kzg.BlobProofVerifier,
	telemetrySink TelemetrySink,
) *Verifier[BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT] {
	return &Verifier[BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT]{
		proofVerifier: proofVerifier,
		metrics:       newVerifierMetrics(telemetrySink),
	}
}

// VerifySidecars verifies the blobs for both inclusion as well
// as the KZG proofs.
func (bv *Verifier[_, _, BlobSidecarsT]) VerifySidecars(
	sidecars BlobSidecarsT, kzgOffset uint64,
) error {
	var (
		g, _      = errgroup.WithContext(context.Background())
		startTime = time.Now()
	)

	defer bv.metrics.measureVerifySidecarsDuration(
		startTime, math.U64(sidecars.Len()),
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

func (bv *Verifier[_, _, BlobSidecarsT]) VerifyInclusionProofs(
	scs BlobSidecarsT,
	kzgOffset uint64,
) error {
	startTime := time.Now()
	defer bv.metrics.measureVerifyInclusionProofsDuration(
		startTime, math.U64(scs.Len()),
	)
	return scs.VerifyInclusionProofs(kzgOffset)
}

// VerifyKZGProofs verifies the sidecars.
func (bv *Verifier[_, _, BlobSidecarsT]) VerifyKZGProofs(
	scs BlobSidecarsT,
) error {
	start := time.Now()
	defer bv.metrics.measureVerifyKZGProofsDuration(
		start, math.U64(scs.Len()),
		bv.proofVerifier.GetImplementation(),
	)

	switch scs.Len() {
	case 0:
		return nil
	case 1:
		blob := scs.Get(0).GetBlob()
		// This method is fastest for a single blob.
		return bv.proofVerifier.VerifyBlobProof(
			&blob,
			scs.Get(0).GetKzgProof(),
			scs.Get(0).GetKzgCommitment(),
		)
	default:
		// For multiple blobs batch verification is more performant
		// than verifying each blob individually (even when done in parallel).
		return bv.proofVerifier.VerifyBlobProofBatch(kzg.ArgsFromSidecars(scs))
	}
}
