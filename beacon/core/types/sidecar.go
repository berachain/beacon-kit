package types

// BlobTxSidecar is a struct that contains blobs and their associated information.
type BlobTxSidecar struct {
	Blob           []byte   `ssz-size:"131072"`
	KzgCommitment  []byte   `ssz-size:"48"`
	KzgProof       []byte   `ssz-size:"48"`
	InclusionProof [][]byte `ssz-size:"17,32"`
}
