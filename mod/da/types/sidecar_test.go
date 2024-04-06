// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/da/types"
	primitives "github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSidecarMarshalling(t *testing.T) {
	// Create a sample BlobSidecar
	blob := kzg.Blob{}
	for i := range blob {
		blob[i] = byte(i % 256)
	}
	sidecar := types.BlobSidecar{
		Index:             1,
		Blob:              blob[:],
		KzgCommitment:     [48]byte{},
		KzgProof:          [48]byte{},
		BeaconBlockHeader: &primitives.BeaconBlockHeader{},
		InclusionProof: [][]byte{
			[]byte("00000000000000000000000000000001"),
			[]byte("00000000000000000000000000000002"),
			[]byte("00000000000000000000000000000003"),
			[]byte("00000000000000000000000000000004"),
			[]byte("00000000000000000000000000000005"),
			[]byte("00000000000000000000000000000006"),
			[]byte("00000000000000000000000000000007"),
			[]byte("00000000000000000000000000000008"),
		},
	}

	// Marshal the sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(t, err, "Marshalling should not produce an error")
	require.NotNil(t, marshalled, "Marshalling should produce a result")

	// Unmarshal the sidecar
	unmarshalled := types.BlobSidecar{}
	err = unmarshalled.UnmarshalSSZ(marshalled)
	require.NoError(t, err, "Unmarshalling should not produce an error")

	// Compare the original and unmarshalled sidecars
	assert.Equal(
		t,
		sidecar,
		unmarshalled,
		"The original and unmarshalled sidecars should be equal",
	)
}
