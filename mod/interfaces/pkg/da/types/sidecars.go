package types

import "github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"

type BlobSidecars[T any, BlobSidecarT any] interface {
	constraints.Nillable
	NewFromSidecars(sidecars []BlobSidecarT) T
	// Len returns the length of the sidecars.
	Len() int
	// Get returns the sidecar at the given index.
	Get(idx uint32) (BlobSidecarT, error)
	// GetSidecars returns all sidecars.
	GetSidecars() []BlobSidecarT
	// ValidateBlockRoots checks to make sure that all blobs
	// in the sidecar are from the same block.
	ValidateBlockRoots() error
	// VerifyInclusionProofs verifies the inclusion proofs for all sidecars.
	VerifyInclusionProofs(kzgOffset uint64) error
}
