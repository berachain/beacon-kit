package kzg

import "github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"

// BlobProofVerifier is a verifier for blobs.
type BlobProofVerifier[
	BlobProofArgsT any,
] interface {
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
	VerifyBlobProofBatch(BlobProofArgsT) error
}
