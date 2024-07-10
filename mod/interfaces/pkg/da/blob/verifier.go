package blob

// Verifier is responsible for verifying blobs, including their
// inclusion and KZG proofs.
type Verifier[BlobSidecarsT any] interface {
	// VerifySidecars verifies the blobs for both inclusion as well
	// as the KZG proofs.
	VerifySidecars(
		sidecars BlobSidecarsT, kzgOffset uint64,
	) error
	// VerifyInclusionProofs verifies the inclusion proofs on the sidecars.
	VerifyInclusionProofs(
		scs BlobSidecarsT,
		kzgOffset uint64,
	) error
	// VerifyKZGProofs verifies the sidecars.
	VerifyKZGProofs(
		scs BlobSidecarsT,
	) error
}
