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

//go:build !ckzg

package ckzg_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/mod/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestVerifyBlobKZGProof(t *testing.T) {
	validBlob, validProof, validCommitment := setupTestData(
		t, "test_data.json")

	testCases := []struct {
		name        string
		blob        *eip4844.Blob
		proof       eip4844.KZGProof
		commitment  eip4844.KZGCommitment
		expectError bool
	}{
		{
			name:        "Valid Proof",
			blob:        validBlob,
			proof:       validProof,
			commitment:  validCommitment,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := globalVerifier.VerifyBlobProof(
				tc.blob,
				tc.proof,
				tc.commitment,
			)
			if tc.expectError {
				require.Error(t, err, "cgo is not enabled")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestVerifyBlobProofBatch tests the valid proofs in batch.
func TestVerifyBlobProofBatch(t *testing.T) {
	if globalVerifier == nil {
		t.Fatal("globalVerifier is not initialized")
	}
	// Load the test data
	fs := afero.NewOsFs()
	fullPath := filepath.Join(baseDir, "test_data_batch.json")
	file, err := afero.ReadFile(fs, fullPath)

	require.NoError(t, err)

	// Unmarshal the JSON data
	var data struct {
		Blobs       []string `json:"blobs"`
		Proofs      []string `json:"proofs"`
		Commitments []string `json:"commitments"`
	}
	err = json.Unmarshal(file, &data)
	require.NoError(t, err)

	// Convert the data to the types expected by VerifyBlobProofBatch
	args := &types.BlobProofArgs{
		Blobs:       make([]*eip4844.Blob, len(data.Blobs)),
		Proofs:      make([]eip4844.KZGProof, len(data.Proofs)),
		Commitments: make([]eip4844.KZGCommitment, len(data.Commitments)),
	}
	for i := range data.Blobs {
		var blob eip4844.Blob
		err = blob.UnmarshalJSON([]byte(`"` + data.Blobs[i] + `"`))
		require.NoError(t, err)
		args.Blobs[i] = &blob

		var proof eip4844.KZGProof
		err = proof.UnmarshalJSON([]byte(`"` + data.Proofs[i] + `"`))
		require.NoError(t, err)
		args.Proofs[i] = proof

		var commitment eip4844.KZGCommitment
		err = commitment.UnmarshalJSON([]byte(`"` + data.Commitments[i] + `"`))
		require.NoError(t, err)
		args.Commitments[i] = commitment
	}

	err = globalVerifier.VerifyBlobProofBatch(args)
	require.Error(t, err, "cgo is not enabled")
}
