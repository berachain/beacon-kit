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

package kzg

import (
	"github.com/berachain/beacon-kit/da/pkg/kzg/ckzg"
	"github.com/berachain/beacon-kit/da/pkg/kzg/gokzg"
	kzgtypes "github.com/berachain/beacon-kit/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// BlobProofVerifier is a verifier for blobs.
type BlobProofVerifier interface {
	// GetImplementation returns the implementation of the verifier.
	GetImplementation() string
	// VerifyBlobProof verifies that the blob data corresponds to the provided
	// commitment.
	VerifyBlobProof(
		blob *eip4844.Blob,
		proof eip4844.KZGProof,
		commitment eip4844.KZGCommitment,
	) error
	// VerifyBlobProofBatch verifies the KZG proof that the polynomial
	// represented
	// by the blob evaluated at the given point is the claimed value.
	// For most implementations it is more efficient than VerifyBlobProof when
	// verifying multiple proofs.
	VerifyBlobProofBatch(*kzgtypes.BlobProofArgs) error
}

// NewBlobProofVerifier creates a new BlobVerifier with the given
// implementation.
func NewBlobProofVerifier(
	impl string,
	ts *gokzg4844.JSONTrustedSetup,
) (BlobProofVerifier, error) {
	switch impl {
	case gokzg.Implementation:
		return gokzg.NewVerifier(ts)
	case ckzg.Implementation:
		return ckzg.NewVerifier(ts)
	default:
		return nil, errors.Wrapf(
			ErrUnsupportedKzgImplementation,
			"supplied: %s, supported: %s, %s",
			impl, gokzg.Implementation, ckzg.Implementation,
		)
	}
}

// ArgsFromSidecars converts a BlobSidecars to a slice of BlobProofArgs.
func ArgsFromSidecars[
	BlobSidecarT kzgtypes.BlobSidecar,
	BlobSidecarsT kzgtypes.BlobSidecars[BlobSidecarT],
](
	scs BlobSidecarsT,
) *kzgtypes.BlobProofArgs {
	proofArgs := &kzgtypes.BlobProofArgs{
		Blobs:       make([]*eip4844.Blob, scs.Len()),
		Proofs:      make([]eip4844.KZGProof, scs.Len()),
		Commitments: make([]eip4844.KZGCommitment, scs.Len()),
	}
	for i, sidecar := range scs.GetSidecars() {
		blob := sidecar.GetBlob()
		proofArgs.Blobs[i] = &blob
		proofArgs.Proofs[i] = sidecar.GetKzgProof()
		proofArgs.Commitments[i] = sidecar.GetKzgCommitment()
	}
	return proofArgs
}
