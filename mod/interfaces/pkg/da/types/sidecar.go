package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type BlobSidecar[
	T any,
	BeaconBlockHeaderT any,
] interface {
	constraints.Nillable
	constraints.SSZMarshallable
	// New creates a new sidecar.
	New(
		index math.U64,
		header BeaconBlockHeaderT,
		blob *eip4844.Blob,
		commitment eip4844.KZGCommitment,
		proof eip4844.KZGProof,
		inclusionProof [][32]byte,
	) T
	// HasValidInclusionProof checks if the sidecar has a valid inclusion proof.
	HasValidInclusionProof(kzgOffset uint64) bool
	// GetSlot returns the slot for which the sidecar belongs to.
	GetSlot() math.Slot
	// GetBlob returns the blob for the sidecar.
	GetBlob() *eip4844.Blob
	// GetCommitment returns the KZG commitment for the sidecar.
	GetCommitment() eip4844.KZGCommitment
	// GetProof returns the KZG proof for the sidecar.
	GetProof() eip4844.KZGProof
}
