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

//go:build ckzg

package ckzg_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	ckzg "github.com/berachain/beacon-kit/mod/da/pkg/kzg/ckzg"
	prooftypes "github.com/berachain/beacon-kit/mod/da/pkg/kzg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// Mock data for testing
var (
	testBlob        = &eip4844.Blob{}
	testProof       = eip4844.KZGProof{}
	testCommitment  = eip4844.KZGCommitment{}
	validBlob       = &eip4844.Blob{}
	validProof      = eip4844.KZGProof{}
	validCommitment = eip4844.KZGCommitment{}
)

var verifier *ckzg.Verifier

// TestMain sets up the trusted setup before running the tests
func TestMain(m *testing.M) {
	fs := afero.NewOsFs()
	file, err := afero.ReadFile(fs, "./files/kzg-trusted-setup.json")

	dummyT := &testing.T{}
	require.NoError(dummyT, err)
	var ts gokzg4844.JSONTrustedSetup
	err = json.Unmarshal(file, &ts)
	require.NoError(dummyT, err)

	verifier, err = ckzg.NewVerifier(&ts)

	require.NoError(dummyT, err)
	require.NotNil(dummyT, verifier)
	// Run the tests
	os.Exit(m.Run())
}

func TestVerifyBlobKZGProofCgoEnabled(t *testing.T) {
	setup(t, "./files/test_data.json")

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
			expectError: false,
		},
		{
			name:        "Bad arguments",
			blob:        testBlob,
			proof:       testProof,
			commitment:  testCommitment,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifier.VerifyBlobProof(tc.blob, tc.proof, tc.commitment)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func setup(t *testing.T, filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}
	require.NoError(t, err)
	type Test struct {
		Input struct {
			Blob       string `json:"blob"`
			Commitment string `json:"commitment"`
			Proof      string `json:"proof"`
		}
		Output *bool `json:"output"`
	}
	var test Test

	err = json.Unmarshal(data, &test)
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	errBlob := validBlob.UnmarshalJSON([]byte(`"` + test.Input.Blob + `"`))
	fmt.Println("errBlob", errBlob)
	require.NoError(t, errBlob)

	if errBlob != nil {
		require.Nil(t, test.Output)
		return
	}

	err = validCommitment.UnmarshalJSON([]byte(`"` + test.Input.Commitment + `"`))
	if err != nil {
		require.Nil(t, test.Output)
		return
	}

	err = validProof.UnmarshalJSON([]byte(`"` + test.Input.Proof + `"`))
	if err != nil {
		require.Nil(t, test.Output)
		return
	}

}

// TestVerifyBlobProofBatch tests the VerifyBlobProofBatch function for valid proofs
func TestVerifyBlobProofBatch(t *testing.T) {
	// Load the test data
	fs := afero.NewOsFs()
	file, err := afero.ReadFile(fs, "./files/test_data_batch.json")
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
	args := &prooftypes.BlobProofArgs{
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

	err = verifier.VerifyBlobProofBatch(args)
	require.NoError(t, err)
}
