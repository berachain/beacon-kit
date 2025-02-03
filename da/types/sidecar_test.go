// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"time"

	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	byteslib "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/math/log"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSidecarMarshalling(t *testing.T) {
	t.Parallel()
	// Create a sample BlobSidecar
	blob := eip4844.Blob{}
	for i := range blob {
		blob[i] = byte(i % 256)
	}
	inclusionProof := make([]common.Root, 0, ctypes.KZGInclusionProofDepth)
	for i := 1; i <= ctypes.KZGInclusionProofDepth; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, errBytes := byteslib.ToBytes32(it)
		require.NoError(t, errBytes)
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

func generateValidBeaconBlock(t *testing.T) *ctypes.BeaconBlock {
	t.Helper()

	// Initialize your block here
	deneb1 := version.Deneb1()
	beaconBlock, err := ctypes.NewBeaconBlockWithVersion(
		math.Slot(10),
		math.ValidatorIndex(5),
		common.Root{1, 2, 3, 4, 5}, // parent root
		deneb1,
	)
	require.NoError(t, err)

	beaconBlock.StateRoot = common.Root{5, 4, 3, 2, 1}
	beaconBlock.Body = &ctypes.BeaconBlockBody{
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
			EpVersion:     deneb1,
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
	}

	body := beaconBlock.GetBody()
	body.SetProposerSlashings(ctypes.ProposerSlashings{})
	body.SetAttesterSlashings(ctypes.AttesterSlashings{})
	body.SetAttestations(ctypes.Attestations{})
	body.SetSyncAggregate(&ctypes.SyncAggregate{})
	body.SetVoluntaryExits(ctypes.VoluntaryExits{})
	body.SetBlsToExecutionChanges(ctypes.BlsToExecutionChanges{})
	return beaconBlock
}

type InclusionSink struct{}

func (is InclusionSink) MeasureSince(_ string, _ time.Time, _ ...string) {}

func TestHasValidInclusionProof(t *testing.T) {
	t.Parallel()
	spec, err := spec.DevnetChainSpec()
	require.NoError(t, err)

	sink := InclusionSink{}
	tests := []struct {
		name           string
		sidecars       func(t *testing.T) types.BlobSidecars
		expectedResult bool
	}{
		{
			name: "Invalid inclusion proof",
			sidecars: func(t *testing.T) types.BlobSidecars {
				t.Helper()
				inclusionProof := make([]common.Root, 0)
				for i := 1; i <= ctypes.KZGInclusionProofDepth; i++ {
					it := byteslib.ExtendToSize(
						[]byte(strconv.Itoa(i)),
						byteslib.B32Size,
					)
					proof, err2 := byteslib.ToBytes32(it)
					require.NoError(t, err2)
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
			expectedResult: false,
		},
		{
			name: "Valid inclusion proof",
			sidecars: func(t *testing.T) types.BlobSidecars {
				t.Helper()
				block := generateValidBeaconBlock(t)

				sidecarFactory := blob.NewSidecarFactory(spec, sink)
				numBlobs := len(block.GetBody().GetBlobKzgCommitments())
				sidecars := make(types.BlobSidecars, numBlobs)
				for i := range numBlobs {
					inclusionProof, incErr := sidecarFactory.BuildKZGInclusionProof(
						block.GetBody(), math.U64(i),
					)
					require.NoError(t, incErr)
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
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sidecars := tt.sidecars(t)
			for _, sidecar := range sidecars {
				result := sidecar.HasValidInclusionProof()
				require.Equal(t, tt.expectedResult, result,
					"Result should match expected value")
			}
		})
	}
}

// Test taken from Prysm:
// https://github.com/prysmaticlabs/prysm/blob/6ce6b869e54c2f98fab5cc836a24e493df19ec49/consensus-types/blocks/kzg_test.go#L107-L120
// This test explains the calculation of the KZG commitment root's Merkle index
// in the Body's Merkle tree based on the index of the KZG commitment list in the Body.
func Test_KZGRootIndex(t *testing.T) {
	t.Parallel()
	// Level of the KZG commitment root's parent.
	kzgParentRootLevel := log.ILog2Ceil(ctypes.KZGPositionDeneb)
	require.NotEqual(t, 0, kzgParentRootLevel)
	// Merkle index of the KZG commitment root's parent.
	// The parent's left child is the KZG commitment root,
	// and its right child is the KZG commitment size.
	kzgParentRootIndex := ctypes.KZGPositionDeneb + (1 << kzgParentRootLevel)
	require.Equal(t, uint64(ctypes.KZGGeneralizedIndex), kzgParentRootIndex)
	// The KZG commitment root is the left child of its parent.
	// Its Merkle index is the double of its parent's Merkle index.
	require.Equal(t, 2*kzgParentRootIndex, uint64(ctypes.KZGRootIndexDeneb))
}

func TestHashTreeRoot(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
