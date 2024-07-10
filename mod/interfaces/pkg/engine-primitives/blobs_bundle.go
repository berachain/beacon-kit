package engineprimitives

import "github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"

// BlobsBundle is an interface for the blobs bundle.
type BlobsBundle[C, P ~[48]byte, B ~[131072]byte] interface {
	constraints.Nillable
	// GetCommitments returns the commitments in the blobs bundle.
	GetCommitments() []C
	// GetProofs returns the proofs in the blobs bundle.
	GetProofs() []P
	// GetBlobs returns the blobs in the blobs bundle.
	GetBlobs() []*B
}
