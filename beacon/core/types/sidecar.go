package types

// SideCars is a slice of blob side cars to be included in the block
type BlobSidecars struct {
	Sidecars []*BlobSidecar `ssz-max:"6"`
}

// BlobSidecar is a struct that contains blobs and their associated information.
type BlobSidecar struct {
	Index          uint64
	Blob           []byte   `ssz-size:"131072"`
	KzgCommitment  []byte   `ssz-size:"48"`
	KzgProof       []byte   `ssz-size:"48"`
	InclusionProof [][]byte `ssz-size:"17,32"`
}
