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

package types_test

import (
	"testing"

	ctypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types/v2"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSidecarMarshalling(t *testing.T) {
	// Create a sample BlobSidecar
	blob := eip4844.Blob{}
	for i := range blob {
		blob[i] = byte(i % 256)
	}
	sidecar := types.BlobSidecar{
		Index:             1,
		Blob:              blob,
		KzgCommitment:     [48]byte{},
		KzgProof:          [48]byte{},
		BeaconBlockHeader: &ctypes.BeaconBlockHeader{},
		InclusionProof: [][32]byte{
			byteslib.ToBytes32([]byte("1")),
			byteslib.ToBytes32([]byte("2")),
			byteslib.ToBytes32([]byte("3")),
			byteslib.ToBytes32([]byte("4")),
			byteslib.ToBytes32([]byte("5")),
			byteslib.ToBytes32([]byte("6")),
			byteslib.ToBytes32([]byte("7")),
			byteslib.ToBytes32([]byte("8")),
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

func TestHasValidInclusionProof(t *testing.T) {
	tests := []struct {
		name           string
		sidecar        *types.BlobSidecar
		kzgOffset      uint64
		expectedResult bool
	}{
		{
			name: "Invalid inclusion proof",
			sidecar: &types.BlobSidecar{
				Index:         0,
				KzgCommitment: eip4844.KZGCommitment{},
				InclusionProof: [][32]byte{
					byteslib.ToBytes32([]byte("4")),
					byteslib.ToBytes32([]byte("5")),
					byteslib.ToBytes32([]byte("6")),
				},
				BeaconBlockHeader: &ctypes.BeaconBlockHeader{
					BodyRoot: [32]byte{3},
				},
			},
			kzgOffset:      0,
			expectedResult: false,
		},
		{
			name: "Empty inclusion proof",
			sidecar: &types.BlobSidecar{
				Index:             0,
				KzgCommitment:     eip4844.KZGCommitment{},
				InclusionProof:    [][32]byte{},
				BeaconBlockHeader: &ctypes.BeaconBlockHeader{},
			},
			kzgOffset:      0,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sidecar.HasValidInclusionProof(tt.kzgOffset)
			require.Equal(t, tt.expectedResult, result,
				"Result should match expected value")
		})
	}
}

func TestHashTreeRoot(t *testing.T) {
	tests := []struct {
		name           string
		sidecar        *types.BlobSidecar
		expectedResult [32]byte
		expectError    bool
	}{
		{
			name: "Valid BlobSidecar",
			sidecar: &types.BlobSidecar{
				Index:         1,
				Blob:          eip4844.Blob{0, 1, 2, 3, 4, 5, 6, 7},
				KzgCommitment: [48]byte{1, 2, 3},
				KzgProof:      [48]byte{4, 5, 6},
				BeaconBlockHeader: &ctypes.BeaconBlockHeader{
					BodyRoot: [32]byte{7, 8, 9},
				},
				InclusionProof: [][32]byte{
					byteslib.ToBytes32([]byte("1")),
					byteslib.ToBytes32([]byte("2")),
					byteslib.ToBytes32([]byte("3")),
					byteslib.ToBytes32([]byte("4")),
					byteslib.ToBytes32([]byte("5")),
					byteslib.ToBytes32([]byte("6")),
					byteslib.ToBytes32([]byte("7")),
					byteslib.ToBytes32([]byte("8")),
				},
			},
			expectedResult: [32]uint8{
				0xce, 0x75, 0x41, 0x87, 0x48, 0x46, 0x6d, 0x26, 0x9e, 0x72, 0x5d,
				0xac, 0x5a, 0x6e, 0x36, 0xed, 0x8c, 0x2a, 0x98, 0x19, 0x6b, 0xe1,
				0xf1, 0xf7, 0xfa, 0xe1, 0x20, 0x5d, 0x2b, 0x3c, 0x57, 0x6a},
			expectError: false,
		},
		{
			name: "Invalid InclusionProof length",
			sidecar: &types.BlobSidecar{
				Index:         1,
				Blob:          eip4844.Blob{0, 1, 2, 3, 4, 5, 6, 7},
				KzgCommitment: [48]byte{1, 2, 3},
				KzgProof:      [48]byte{4, 5, 6},
				BeaconBlockHeader: &ctypes.BeaconBlockHeader{
					BodyRoot: [32]byte{7, 8, 9},
				},
				InclusionProof: [][32]byte{
					byteslib.ToBytes32([]byte("1")),
				},
			},
			expectedResult: [32]byte{},
			expectError:    true,
		},
		{
			name: "Nil BeaconBlockHeader",
			sidecar: &types.BlobSidecar{
				Index:             1,
				Blob:              eip4844.Blob{0, 1, 2, 3, 4, 5, 6, 7},
				KzgCommitment:     [48]byte{1, 2, 3},
				KzgProof:          [48]byte{4, 5, 6},
				BeaconBlockHeader: nil,
				InclusionProof: [][32]byte{
					byteslib.ToBytes32([]byte("1")),
					byteslib.ToBytes32([]byte("2")),
					byteslib.ToBytes32([]byte("3")),
					byteslib.ToBytes32([]byte("4")),
					byteslib.ToBytes32([]byte("5")),
					byteslib.ToBytes32([]byte("6")),
					byteslib.ToBytes32([]byte("7")),
					byteslib.ToBytes32([]byte("8")),
				},
			},
			expectedResult: [32]uint8{
				0xb8, 0x3d, 0x3b, 0xfb, 0x39, 0xd4, 0xce, 0x2a, 0x9e, 0x4c, 0xa1,
				0x40, 0xd2, 0x94, 0xeb, 0xaf, 0xdf, 0xbd, 0x85, 0x3d, 0xe0, 0x87,
				0xa4, 0xf3, 0x6, 0xf7, 0xe2, 0x9c, 0x27, 0x41, 0x27, 0x71},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.sidecar.HashTreeRoot()
			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err,
					"Did not expect an error but got one")
				assert.Equal(t, tt.expectedResult, result,
					"Hash result should match expected value")
			}
		})
	}
}
