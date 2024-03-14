package types

// SideCars is a slice of blob side cars to be included in the block
type BlobSidecars struct {
	BlobSidecars []*BlobTxSidecar `ssz-max:"16"`
}

// BlobTxSidecar is a struct that contains blobs and their associated information.
type BlobTxSidecar struct {
	Blob           []byte   `ssz-size:"131072"`
	KzgCommitment  []byte   `ssz-size:"48"`
	KzgProof       []byte   `ssz-size:"48"`
	InclusionProof [][]byte `ssz-size:"17,32"`
}
