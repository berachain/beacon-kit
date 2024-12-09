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

package gokzg

import (
	"unsafe"

	"github.com/berachain/beacon-kit/da/kzg/types"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

const Implementation = "crate-crypto/go-kzg-4844"

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

// GetImplementation returns the implementation of the verifier.
func (v Verifier) GetImplementation() string {
	return Implementation
}

// VerifyBlobProof verifies the KZG proof that the polynomial represented by the
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
	args *types.BlobProofArgs,
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
