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

package kzg

import (
	"github.com/berachain/beacon-kit/mod/da/kzg/ckzg"
	"github.com/berachain/beacon-kit/mod/da/kzg/gokzg"
	prooftypes "github.com/berachain/beacon-kit/mod/da/kzg/types"
	"github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/cockroachdb/errors"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// BlobProofVerifier is a verifier for blobs.
type BlobProofVerifier interface {
	// VerifyBlobProof verifies that the blob data corresponds to the provided
	// commitment.
	VerifyBlobProof(
		blob *primitives.Blob,
		proof primitives.Proof,
		commitment primitives.Commitment,
	) error

	// VerifyBlobProofBatch verifies the KZG proof that the polynomial
	// represented
	// by the blob evaluated at the given point is the claimed value.
	// For most implementations it is more efficient than VerifyBlobProof when
	// verifying multiple proofs.
	VerifyBlobProofBatch(
		*prooftypes.BlobProofArgs,
	) error
}

const (
	// crateCryptoGoKzg4844 is the crate-crypto/go-kzg-4844 implementation.
	crateCryptoGoKzg4844 = "crate-crypto/go-kzg-4844"
	// ethereumCKzg4844 is the ethereum/c-kzg-4844 implementation.
	ethereumCKzg4844 = "ethereum/c-kzg-4844"
)

// NewBlobProofVerifier creates a new BlobVerifier with the given
// implementation.
func NewBlobProofVerifier(
	impl string,
	ts *gokzg4844.JSONTrustedSetup,
) (BlobProofVerifier, error) {
	switch impl {
	case crateCryptoGoKzg4844:
		return gokzg.NewVerifier(ts)
	case ethereumCKzg4844:
		return ckzg.NewVerifier(ts)
	default:
		return nil, errors.Wrapf(
			ErrUnsupportedKzgImplementation,
			"supplied: %s, supported: %s, %s",
			impl, crateCryptoGoKzg4844, ethereumCKzg4844,
		)
	}
}

// ArgsFromSidecars converts a BlobSidecars to a slice of BlobProofArgs.
func ArgsFromSidecars(
	scs *types.BlobSidecars,
) *prooftypes.BlobProofArgs {
	proofArgs := &prooftypes.BlobProofArgs{
		Blobs:       make([]*primitives.Blob, len(scs.Sidecars)),
		Proofs:      make([]primitives.Proof, len(scs.Sidecars)),
		Commitments: make([]primitives.Commitment, len(scs.Sidecars)),
	}
	for i, sidecar := range scs.Sidecars {
		proofArgs.Blobs[i] = &sidecar.Blob
		proofArgs.Proofs[i] = sidecar.KzgProof
		proofArgs.Commitments[i] = sidecar.KzgCommitment
	}
	return proofArgs
}
