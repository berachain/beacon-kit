// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//go:build !ckzg

package ckzg_test

import (
	"path/filepath"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"

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
				require.Error(
					t,
					err,
					//nolint:lll
					"github.com/ethereum/c-kzg-4844 requires an executable built with CGO_ENABLED=1",
				)
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
	//nolint:lll
	require.Error(
		t,
		err,
		"github.com/ethereum/c-kzg-4844 requires an executable built with CGO_ENABLED=1",
	)
}

// TestVerifyBlobKZGInvalidProof tests the VerifyBlobProof function for invalid
// proofs.
func TestVerifyBlobKZGInvalidProof(t *testing.T) {
	validBlob, invalidProof, validCommitment := setupTestData(
		t, "test_data_incorrect_proof.json")
	testCases := []struct {
		name        string
		blob        *eip4844.Blob
		proof       eip4844.KZGProof
		commitment  eip4844.KZGCommitment
		expectError bool
	}{
		{
			name:        "Invalid Proof",
			blob:        validBlob,
			proof:       invalidProof,
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
				require.Error(t, err, "invalid proof")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetImplementation(t *testing.T) {
	require.NotNil(t, globalVerifier)
	require.Equal(t, "ethereum/c-kzg-4844", globalVerifier.GetImplementation())
}
