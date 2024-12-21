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
	"github.com/berachain/beacon-kit/chain-spec/chain"
	spec2 "github.com/berachain/beacon-kit/config/spec"
	"github.com/berachain/beacon-kit/da/blob"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"strconv"
	"testing"
	"time"

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

func generateValidBeaconBlock() *ctypes.BeaconBlock {
	// Initialize your block here
	return &ctypes.BeaconBlock{
		Slot:          10,
		ProposerIndex: 5,
		ParentRoot:    common.Root{1, 2, 3, 4, 5},
		StateRoot:     common.Root{5, 4, 3, 2, 1},
		Body: &ctypes.BeaconBlockBody{
			ExecutionPayload: &ctypes.ExecutionPayload{
				Timestamp: 10,
				ExtraData: []byte("dummy extra data for testing"),
				Transactions: [][]byte{
					[]byte("tx1"),
					[]byte("tx2"),
					[]byte("tx3"),
				},
				Withdrawals: engineprimitives.Withdrawals{
					{Index: 0, Amount: 100},
					{Index: 1, Amount: 200},
				},
				BaseFeePerGas: math.NewU256(0),
			},
			Eth1Data: &ctypes.Eth1Data{},
			Deposits: []*ctypes.Deposit{
				{
					Index: 1,
				},
			},
			BlobKzgCommitments: []eip4844.KZGCommitment{
				{0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab}, {2}, {0x69},
			},
		},
	}
}

type InclusionSink struct{}

func (is InclusionSink) MeasureSince(key string, start time.Time, args ...string) {}

func TestHasValidInclusionProof(t *testing.T) {
	// TODO: Get updated spec?
	specVals := spec2.BaseSpec()
	// TODO: we currently cannot change this MaxBlobCommitmentsPerBlock value.
	// Will likely have to change inclusionProofDepth to make any different values work.
	specVals.MaxBlobCommitmentsPerBlock = 16
	spec, err := chain.NewChainSpec(specVals)
	require.NoError(t, err)
	// TODO: get a good slot number that is current fork
	inclusionProofDepth, err := ctypes.KZGCommitmentInclusionProofDepth(0, spec)
	require.NoError(t, err)

	sink := InclusionSink{}
	tests := []struct {
		name           string
		sidecars       func(t *testing.T) types.BlobSidecars
		kzgOffset      uint64
		expectedResult bool
	}{
		{
			name: "Invalid inclusion proof",
			sidecars: func(t *testing.T) types.BlobSidecars {
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
				return types.BlobSidecars{types.BuildBlobSidecar(
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
				)}
			},
			kzgOffset:      0,
			expectedResult: false,
		},
		{
			name: "Empty inclusion proof",
			sidecars: func(*testing.T) types.BlobSidecars {
				return types.BlobSidecars{types.BuildBlobSidecar(
					math.U64(0),
					&ctypes.SignedBeaconBlockHeader{},
					&eip4844.Blob{},
					eip4844.KZGCommitment{},
					eip4844.KZGProof{},
					[]common.Root{},
				)}
			},
			kzgOffset:      0,
			expectedResult: false,
		},
		{
			name: "Valid inclusion proof",
			sidecars: func(t *testing.T) types.BlobSidecars {
				block := generateValidBeaconBlock()

				sidecarFactory := blob.NewSidecarFactory(
					spec,
					sink,
				)
				numBlobs := len(block.GetBody().GetBlobKzgCommitments())
				sidecars := make(types.BlobSidecars, numBlobs)
				for i := range numBlobs {
					inclusionProof, err := sidecarFactory.BuildKZGInclusionProof(
						block.GetBody(), math.U64(i), ctypes.KZGPositionDeneb,
					)
					require.NoError(t, err)
					sigHeader := ctypes.NewSignedBeaconBlockHeader(block.GetHeader(), crypto.BLSSignature{})
					sidecars[i] = types.BuildBlobSidecar(
						math.U64(i),
						sigHeader,
						&eip4844.Blob{},
						block.GetBody().BlobKzgCommitments[i],
						eip4844.KZGProof{},
						inclusionProof,
					)
				}
				return sidecars
			},
			kzgOffset:      ctypes.KZGMerkleIndexDeneb * spec.MaxBlobCommitmentsPerBlock(),
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sidecars := tt.sidecars(t)
			for _, sidecar := range sidecars {
				result := sidecar.HasValidInclusionProof(tt.kzgOffset, inclusionProofDepth)
				require.Equal(t, tt.expectedResult, result,
					"Result should match expected value")
			}
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
