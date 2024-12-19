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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"strconv"
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/types"
	byteslib "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSidecarMarshalling(t *testing.T) {
	// Create a sample BlobSidecar
	blob := eip4844.Blob{}
	for i := range blob {
		blob[i] = byte(i % 256)
	}
	inclusionProof := make([]common.Root, 0)
	for i := int(1); i <= 8; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, err := byteslib.ToBytes32(it)
		require.NoError(t, err)
		inclusionProof = append(inclusionProof, common.Root(proof))
	}
	sidecar := types.BuildBlobSidecar(
		1,
		&ctypes.SignedBeaconBlockHeader{
			Header:    &ctypes.BeaconBlockHeader{},
			Signature: crypto.BLSSignature{},
		},
		&blob,
		eip4844.KZGCommitment{},
		eip4844.KZGProof{},
		inclusionProof,
	)

	// Marshal the sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(t, err, "Marshalling should not produce an error")
	require.NotNil(t, marshalled, "Marshalling should produce a result")

	// Unmarshal the sidecar
	unmarshalled := &types.BlobSidecar{}
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
	// Equates to KZG_COMMITMENT_INCLUSION_PROOF_DEPTH
	const inclusionProofDepth = 17
	tests := []struct {
		name           string
		sidecar        func(t *testing.T) *types.BlobSidecar
		kzgOffset      uint64
		expectedResult bool
	}{
		{
			name: "Invalid inclusion proof",
			sidecar: func(t *testing.T) *types.BlobSidecar {
				t.Helper()
				inclusionProof := make([]common.Root, 0)
				for i := int(1); i <= 8; i++ {
					it := byteslib.ExtendToSize(
						[]byte(strconv.Itoa(i)),
						byteslib.B32Size,
					)
					proof, err := byteslib.ToBytes32(it)
					require.NoError(t, err)
					inclusionProof = append(inclusionProof, common.Root(proof))
				}
				return types.BuildBlobSidecar(
					math.U64(0),
					&ctypes.SignedBeaconBlockHeader{
						Header: &ctypes.BeaconBlockHeader{
							BodyRoot: [32]byte{3},
						},
						Signature: crypto.BLSSignature{},
					},
					&eip4844.Blob{},
					eip4844.KZGCommitment{},
					eip4844.KZGProof{},
					inclusionProof,
				)
			},
			kzgOffset:      0,
			expectedResult: false,
		},
		{
			name: "Empty inclusion proof",
			sidecar: func(*testing.T) *types.BlobSidecar {
				return types.BuildBlobSidecar(
					math.U64(0),
					&ctypes.SignedBeaconBlockHeader{},
					&eip4844.Blob{},
					eip4844.KZGCommitment{},
					eip4844.KZGProof{},
					[]common.Root{},
				)
			},
			kzgOffset:      0,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sidecar := tt.sidecar(t)
			result := sidecar.HasValidInclusionProof(tt.kzgOffset, inclusionProofDepth)
			require.Equal(t, tt.expectedResult, result,
				"Result should match expected value")
		})
	}
}

func TestHashTreeRoot(t *testing.T) {
	tests := []struct {
		name           string
		sidecar        func(t *testing.T) *types.BlobSidecar
		expectedResult common.Root
		expectError    bool
	}{
		{
			name: "Valid BlobSidecar",
			sidecar: func(t *testing.T) *types.BlobSidecar {
				t.Helper()
				inclusionProof := make([]common.Root, 0)
				for i := int(1); i <= 8; i++ {
					it := byteslib.ExtendToSize(
						[]byte(strconv.Itoa(i)),
						byteslib.B32Size,
					)
					proof, err := byteslib.ToBytes32(it)
					require.NoError(t, err)
					inclusionProof = append(inclusionProof, common.Root(proof))
				}
				return types.BuildBlobSidecar(
					math.U64(1),
					&ctypes.SignedBeaconBlockHeader{
						Header: &ctypes.BeaconBlockHeader{
							BodyRoot: [32]byte{7, 8, 9},
						},
						Signature: crypto.BLSSignature{0xde, 0xad},
					},
					&eip4844.Blob{0, 1, 2, 3, 4, 5, 6, 7},
					eip4844.KZGCommitment{1, 2, 3},
					eip4844.KZGProof{4, 5, 6},
					inclusionProof,
				)
			},
			expectedResult: [32]uint8{
				0xd8, 0xb2, 0x91, 0x39, 0x93, 0x75, 0x38, 0x1f,
				0xd4, 0xdf, 0xef, 0xa7, 0x16, 0x91, 0xd9, 0x9,
				0x3, 0x62, 0xee, 0x3a, 0x79, 0x96, 0x57, 0xc4,
				0xc4, 0x6d, 0x86, 0x79, 0x78, 0x1b, 0xb4, 0xe3,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				sidecar := tt.sidecar(t)
				result := sidecar.HashTreeRoot()
				require.Equal(
					t,
					tt.expectedResult,
					result,
					"HashTreeRoot result should match expected value",
				)
			}, "HashTreeRoot should not panic")
		})
	}
}
