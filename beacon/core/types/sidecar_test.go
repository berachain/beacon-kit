package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/stretchr/testify/assert"
)

func TestBlobSidecarMarshalUnmarshal(t *testing.T) {
	original := &types.BlobSidecar{
		Index:          1,
		Blob:           make([]byte, 131072),
		KzgCommitment:  make([]byte, 48),
		KzgProof:       make([]byte, 48),
		InclusionProof: make([][]byte, 17),
	}
	for i := range original.InclusionProof {
		original.InclusionProof[i] = make([]byte, 32)
	}

	marshaled, err := original.MarshalSSZ()
	assert.NoError(t, err, "marshaling should not produce an error")

	unmarshaled := &types.BlobSidecar{}
	err = unmarshaled.UnmarshalSSZ(marshaled)
	assert.NoError(t, err, "unmarshaling should not produce an error")

	assert.Equal(t, original, unmarshaled, "unmarshaled object should equal the original")
}

func TestBlobSidecarsMarshalUnmarshal(t *testing.T) {
	original := &types.BlobSidecars{
		Sidecars: make([]*types.BlobSidecar, 0),
	}

	marshaled, err := original.MarshalSSZ()
	assert.NoError(t, err, "marshaling should not produce an error")

	unmarshaled := &types.BlobSidecars{}
	err = unmarshaled.UnmarshalSSZ(marshaled)
	assert.NoError(t, err, "unmarshaling should not produce an error")

	assert.Equal(t, original, unmarshaled, "unmarshaled object should equal the original")
}
