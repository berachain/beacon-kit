package blob

import "github.com/berachain/beacon-kit/mod/primitives/pkg/math"

// SidecarFactory is the interface for a factory for sidecars.
type SidecarFactory[
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BlobsBundleT any,
	BlobSidecarsT any,
] interface {
	// BuildSidecars builds a sidecar.
	BuildSidecars(
		blk BeaconBlockT,
		bundle BlobsBundleT,
	) (BlobSidecarsT, error)
	// BuildKZGInclusionProof builds a KZG inclusion proof.
	BuildKZGInclusionProof(
		body BeaconBlockBodyT,
		index math.U64,
	) ([][32]byte, error)
	// BuildBlockBodyProof builds a block body proof.
	BuildBlockBodyProof(
		body BeaconBlockBodyT,
	) ([][32]byte, error)
	// BuildCommitmentProof builds a commitment proof.
	BuildCommitmentProof(
		body BeaconBlockBodyT,
		index math.U64,
	) ([][32]byte, error)
}
