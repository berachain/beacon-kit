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
	"fmt"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"testing"

	ctypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestEmptySidecarMarshalling(t *testing.T) {
	// Create an empty BlobSidecar
	sidecar := types.BlobSidecar{
		Index:             0,
		Blob:              eip4844.Blob{},
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

	// Marshal the empty sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(
		t,
		err,
		"Marshalling empty sidecar should not produce an error",
	)
	require.NotNil(
		t,
		marshalled,
		"Marshalling empty sidecar should produce a result",
	)

	// Unmarshal the empty sidecar
	unmarshalled := types.BlobSidecar{}
	err = unmarshalled.UnmarshalSSZ(marshalled)
	require.NoError(
		t,
		err,
		"Unmarshalling empty sidecar should not produce an error",
	)

	// Compare the original and unmarshalled empty sidecars
	require.Equal(
		t,
		sidecar,
		unmarshalled,
		"The original and unmarshalled empty sidecars should be equal",
	)
}
func TestValidateBlockRoots(t *testing.T) {
	// Create a sample BlobSidecar with valid roots
	validSidecar := types.BlobSidecar{
		Index: 0,
		Blob:  eip4844.Blob{},
		BeaconBlockHeader: &ctypes.BeaconBlockHeader{
			BeaconBlockHeaderBase: ctypes.BeaconBlockHeaderBase{
				StateRoot: [32]byte{1},
			},
			BodyRoot: [32]byte{2},
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
	}

	// Validate the sidecar with valid roots
	sidecars := types.BlobSidecars{
		Sidecars: []*types.BlobSidecar{&validSidecar},
	}
	err := sidecars.ValidateBlockRoots()
	require.NoError(
		t,
		err,
		"Validating sidecar with valid roots should not produce an error",
	)

	// Create a sample BlobSidecar with invalid roots
	differentBlockRootSidecar := types.BlobSidecar{
		Index: 0,
		Blob:  eip4844.Blob{},
		BeaconBlockHeader: &ctypes.BeaconBlockHeader{
			BeaconBlockHeaderBase: ctypes.BeaconBlockHeaderBase{
				StateRoot: [32]byte{1},
			},
			BodyRoot: [32]byte{3},
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
	}
	// Validate the sidecar with invalid roots
	sidecarsInvalid := types.BlobSidecars{
		Sidecars: []*types.BlobSidecar{
			&validSidecar,
			&differentBlockRootSidecar,
		},
	}
	err = sidecarsInvalid.ValidateBlockRoots()
	require.Error(
		t,
		err,
		"Validating sidecar with invalid roots should produce an error",
	)
}

func TestBuildBlobSidecar(t *testing.T) {
	tests := []struct {
		name           string
		index          math.U64
		header         *ctypes.BeaconBlockHeader
		blob           *eip4844.Blob
		commitment     eip4844.KZGCommitment
		proof          eip4844.KZGProof
		inclusionProof [][32]byte
	}{
		{
			name:   "Basic test",
			index:  math.U64(1),
			header: &ctypes.BeaconBlockHeader{
				// Initialize with sample data
			},
			blob: func() *eip4844.Blob {
				var b eip4844.Blob
				for i := range b {
					b[i] = byte(i % 256)
				}
				return &b
			}(),
			commitment: eip4844.KZGCommitment{},
			proof:      eip4844.KZGProof{},
			inclusionProof: [][32]byte{
				byteslib.ToBytes32([]byte("1")),
				byteslib.ToBytes32([]byte("2")),
				byteslib.ToBytes32([]byte("3")),
			},
		},
		{
			name:   "Empty inclusion proof",
			index:  math.U64(2),
			header: &ctypes.BeaconBlockHeader{
				// Initialize with sample data
			},
			blob: func() *eip4844.Blob {
				var b eip4844.Blob
				for i := range b {
					b[i] = byte((i + 1) % 256)
				}
				return &b
			}(),
			commitment:     eip4844.KZGCommitment{},
			proof:          eip4844.KZGProof{},
			inclusionProof: [][32]byte{},
		},
		{
			name:   "Different index",
			index:  math.U64(3),
			header: &ctypes.BeaconBlockHeader{
				// Initialize with sample data
			},
			blob: func() *eip4844.Blob {
				var b eip4844.Blob
				for i := range b {
					b[i] = byte((i + 2) % 256)
				}
				return &b
			}(),
			commitment: eip4844.KZGCommitment{},
			proof:      eip4844.KZGProof{},
			inclusionProof: [][32]byte{
				byteslib.ToBytes32([]byte("4")),
				byteslib.ToBytes32([]byte("5")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sidecar := types.BuildBlobSidecar(
				tt.index,
				tt.header,
				tt.blob,
				tt.commitment,
				tt.proof,
				tt.inclusionProof,
			)

			require.Equal(t, tt.index.Unwrap(), sidecar.Index, "Index should match")
			require.Equal(t, *tt.blob, sidecar.Blob, "Blob should match")
			require.Equal(t, tt.commitment, sidecar.KzgCommitment, "KzgCommitment should match")
			require.Equal(t, tt.proof, sidecar.KzgProof, "KzgProof should match")
			require.Equal(t, tt.header, sidecar.BeaconBlockHeader, "BeaconBlockHeader should match")
			require.Equal(t, tt.inclusionProof, sidecar.InclusionProof, "InclusionProof should match")
		})
	}
}

func TestHasValidInclusionProof(t *testing.T) {
	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("B")),
		byteslib.ToBytes32([]byte("C")),
		byteslib.ToBytes32([]byte("D")),
		byteslib.ToBytes32([]byte("E")),
		byteslib.ToBytes32([]byte("F")),
		byteslib.ToBytes32([]byte("G")),
		byteslib.ToBytes32([]byte("H")),
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte, [32]byte](
		items,
		32,
	)
	proof, _ := m.MerkleProofWithMixin(0)
	root, err := m.HashTreeRoot()
	require.True(t, merkle.VerifyProof(
		root, items[0], 0, proof,
	), "First Merkle proof did not verify")

	tests := []struct {
		name           string
		sidecar        *types.BlobSidecar
		kzgOffset      uint64
		expectedResult bool
	}{
		{
			name: "Valid inclusion proof",
			sidecar: &types.BlobSidecar{
				Index:          0,
				KzgCommitment:  m,
				InclusionProof: items,
				BeaconBlockHeader: &ctypes.BeaconBlockHeader{
					BodyRoot: [32]byte{3},
				},
			},
			kzgOffset:      0,
			expectedResult: true,
		},
		{
			name: "Invalid inclusion proof",
			sidecar: &types.BlobSidecar{
				Index:         0,
				KzgCommitment: eip4844.KZGCommitment{
					// Initialize with sample data
				},
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
				Index:         0,
				KzgCommitment: eip4844.KZGCommitment{
					// Initialize with sample data
				},
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
			fmt.Println(result)
			require.Equal(t, tt.expectedResult, result, "Result should match expected value")
		})
	}
}
