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

	ctypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	ssz "github.com/ferranbt/fastssz"
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

func TestMarshalSSZUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name           string
		blobSidecars   types.BlobSidecars
		expectedResult types.BlobSidecars
		expectError    bool
	}{
		{
			name: "Valid BlobSidecars with two sidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{
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
					},
					{
						Index: 1,
						Blob:  eip4844.Blob{},
						BeaconBlockHeader: &ctypes.BeaconBlockHeader{
							BeaconBlockHeaderBase: ctypes.BeaconBlockHeaderBase{
								StateRoot: [32]byte{3},
							},
							BodyRoot: [32]byte{4},
						},
						InclusionProof: [][32]byte{
							byteslib.ToBytes32([]byte("9")),
							byteslib.ToBytes32([]byte("10")),
							byteslib.ToBytes32([]byte("11")),
							byteslib.ToBytes32([]byte("12")),
							byteslib.ToBytes32([]byte("13")),
							byteslib.ToBytes32([]byte("14")),
							byteslib.ToBytes32([]byte("15")),
							byteslib.ToBytes32([]byte("16")),
						},
					},
				},
			},
			expectedResult: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{
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
					},
					{
						Index: 1,
						Blob:  eip4844.Blob{},
						BeaconBlockHeader: &ctypes.BeaconBlockHeader{
							BeaconBlockHeaderBase: ctypes.BeaconBlockHeaderBase{
								StateRoot: [32]byte{3},
							},
							BodyRoot: [32]byte{4},
						},
						InclusionProof: [][32]byte{
							byteslib.ToBytes32([]byte("9")),
							byteslib.ToBytes32([]byte("10")),
							byteslib.ToBytes32([]byte("11")),
							byteslib.ToBytes32([]byte("12")),
							byteslib.ToBytes32([]byte("13")),
							byteslib.ToBytes32([]byte("14")),
							byteslib.ToBytes32([]byte("15")),
							byteslib.ToBytes32([]byte("16")),
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Empty BlobSidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{},
			},
			expectedResult: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{},
			},
			expectError: false,
		},
		{
			name: "BlobSidecars with more than 6 sidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{Index: 0}, {Index: 1}, {Index: 2}, {Index: 3},
					{Index: 4}, {Index: 5}, {Index: 6},
				},
			},
			expectedResult: types.BlobSidecars{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the BlobSidecars object
			marshalled, err := tt.blobSidecars.MarshalSSZ()
			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
				return
			}
			require.NoError(t, err,
				"Marshalling BlobSidecars should not produce an error")
			require.NotNil(t, marshalled,
				"Marshalling BlobSidecars should produce a result")

			// Unmarshal the BlobSidecars object
			unmarshalled := types.BlobSidecars{}
			err = unmarshalled.UnmarshalSSZ(marshalled)
			require.NoError(t,
				err,
				"Unmarshalling BlobSidecars should not produce an error")

			// Compare the original and unmarshalled BlobSidecars objects
			require.Equal(t,
				tt.expectedResult,
				unmarshalled,
				"The original and unmarshalled BlobSidecars should be equal")
		})
	}
}

func TestUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name         string
		blobSidecars types.BlobSidecars
		inputBuffer  []byte
		expectedErr  error
	}{
		{
			name:         "BlobSidecar fails to unmarshal due to invalid offset",
			blobSidecars: types.BlobSidecars{},
			inputBuffer:  []byte{0x00, 0x00, 0x00, 0x08},
			expectedErr:  ssz.ErrOffset,
		},
		{
			name:         "Fails to unmarshal due to invalid variable offset",
			blobSidecars: types.BlobSidecars{},
			inputBuffer:  []byte{0x01, 0x00, 0x00, 0x00},
			expectedErr:  ssz.ErrInvalidVariableOffset,
		},
		{
			name:         "BlobSidecar fails to unmarshal due to size less than 4",
			blobSidecars: types.BlobSidecars{},
			inputBuffer:  []byte{0x01, 0x02, 0x03},
			expectedErr:  ssz.ErrSize,
		},
		{
			name:         "BlobSidecar successfully unmarshals with valid buffer",
			blobSidecars: types.BlobSidecars{},
			inputBuffer:  []byte{0x04, 0x00, 0x00, 0x00},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.blobSidecars.UnmarshalSSZ(tt.inputBuffer)
			if tt.expectedErr != nil {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedErr, err,
					"Expected error %v, got %v", tt.expectedErr, err)
			} else {
				require.NoError(t, err,
					"Unmarshalling BlobSidecars should not produce an error")
			}
		})
	}
}

func TestHashTreeRootFunctions(t *testing.T) {
	tests := []struct {
		name          string
		blobSidecars  types.BlobSidecars
		expectError   bool
		expectedError error
	}{
		{
			name: "Valid BlobSidecars with two sidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{
						Index: 0,
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
					{
						Index: 1,
						InclusionProof: [][32]byte{
							byteslib.ToBytes32([]byte("9")),
							byteslib.ToBytes32([]byte("10")),
							byteslib.ToBytes32([]byte("11")),
							byteslib.ToBytes32([]byte("12")),
							byteslib.ToBytes32([]byte("13")),
							byteslib.ToBytes32([]byte("14")),
							byteslib.ToBytes32([]byte("15")),
							byteslib.ToBytes32([]byte("16")),
						},
					},
				},
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "Empty BlobSidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{},
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "BlobSidecars with more than 8 sidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{Index: 0}, {Index: 1}, {Index: 2}, {Index: 3}, {Index: 4},
					{Index: 5}, {Index: 6}, {Index: 7}, {Index: 8},
				},
			},
			expectError:   true,
			expectedError: ssz.ErrIncorrectListSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test HashTreeRoot
			root, err := tt.blobSidecars.HashTreeRoot()
			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedError, err,
					"Expected error %v, got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err,
					"HashTreeRoot should not produce an error")
				require.NotNil(t, root,
					"HashTreeRoot should produce a result")
			}

			// Test HashTreeRootWith
			hh := ssz.NewHasher()
			err = tt.blobSidecars.HashTreeRootWith(hh)
			if tt.expectError {
				require.Error(t,
					err,
					"Expected an error but got none")
				require.Equal(t, tt.expectedError, err,
					"Expected error %v, got %v",
					tt.expectedError, err)
			} else {
				require.NoError(t, err,
					"HashTreeRootWith should not produce an error")
				require.NotNil(t, hh.Hash(),
					"HashTreeRootWith should produce a result")
			}

			// Test GetTree
			tree, err := tt.blobSidecars.GetTree()
			if tt.expectError {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedError, err,
					"Expected error %v, got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err,
					"GetTree should not produce an error")
				require.NotNil(t, tree,
					"GetTree should produce a result")
			}
		})
	}
}

func TestVerifyInclusionProofs(t *testing.T) {
	tests := []struct {
		name          string
		blobSidecars  types.BlobSidecars
		kzgOffset     uint64
		expectedError error
	}{
		{
			name: "Invalid inclusion proof",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{
						Index:          0,
						InclusionProof: [][32]byte{},
						BeaconBlockHeader: &ctypes.BeaconBlockHeader{
							BodyRoot: [32]byte{7, 8, 9},
						},
					},
				},
			},
			kzgOffset:     0,
			expectedError: types.ErrInvalidInclusionProof,
		},
		{
			name: "Nil sidecar",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					nil,
				},
			},
			kzgOffset:     0,
			expectedError: types.ErrAttemptedToVerifyNilSidecar,
		},
		{
			name: "Empty sidecars",
			blobSidecars: types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{},
			},
			kzgOffset:     0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.blobSidecars.VerifyInclusionProofs(tt.kzgOffset)
			if tt.expectedError != nil {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedError.Error(), err.Error(),
					"Expected error %v, "+
						"got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err,
					"VerifyInclusionProofs should not produce an error")
			}
		})
	}
}

func TestBlobSidecars(t *testing.T) {
	tests := []struct {
		name         string
		blobSidecars *types.BlobSidecars
		expectedNil  bool
		expectedLen  int
	}{
		{
			name: "Nil Sidecars slice",
			blobSidecars: &types.BlobSidecars{
				Sidecars: nil,
			},
			expectedNil: true,
			expectedLen: 0,
		},
		{
			name: "Empty Sidecars slice",
			blobSidecars: &types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{},
			},
			expectedNil: false,
			expectedLen: 0,
		},
		{
			name: "Non-empty Sidecars slice",
			blobSidecars: &types.BlobSidecars{
				Sidecars: []*types.BlobSidecar{
					{
						Index: 0,
					},
				},
			},
			expectedNil: false,
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test IsNil
			resultNil := tt.blobSidecars.IsNil()
			require.Equal(t, tt.expectedNil, resultNil,
				"Expected IsNil to be %v, got %v",
				tt.expectedNil, resultNil)

			// Test Len
			resultLen := tt.blobSidecars.Len()
			require.Equal(t, tt.expectedLen, resultLen,
				"Expected Len to be %d, got %d",
				tt.expectedLen, resultLen)
		})
	}
}
